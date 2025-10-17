package graph

import (
	"context"
	"fmt"
	"sync"
	"sync/atomic"

	"github.com/google/uuid"
	g "github.com/morphy76/ggraph/pkg/graph"
)

// RuntimeFactory creates a new instance of Runtime with the specified SharedState type, state merger function, and initial state.
func RuntimeFactory[T g.SharedState](
	startEdge g.Edge[T],
	stateMonitorCh chan g.StateMonitorEntry[T],
	initialState T,
) (g.Runtime[T], error) {
	if startEdge == nil {
		return nil, fmt.Errorf("runtime creation failed: start edge cannot be nil")
	}
	ctx, cancelFn := context.WithCancel(context.Background())
	rv := &runtimeImpl[T]{
		ctx:    ctx,
		cancel: cancelFn,

		outcomeCh:      make(chan nodeFnReturnStruct[T]),
		stateMonitorCh: stateMonitorCh,

		startEdge: startEdge,
		edges:     []g.Edge[T]{},

		state:          initialState,
		stateMergeLock: &sync.Mutex{},

		pendingPersist: make(chan T, 10),
	}
	rv.executing.Store(false)
	rv.start()
	rv.startPersistenceWorker()
	return rv, nil
}

// ------------------------------------------------------------------------------
// Runtime Implementation
// ------------------------------------------------------------------------------

var _ g.Runtime[g.SharedState] = (*runtimeImpl[g.SharedState])(nil)
var _ g.StateObserver[g.SharedState] = (*runtimeImpl[g.SharedState])(nil)

type nodeFnReturnStruct[T g.SharedState] struct {
	node      g.Node[T]
	userInput T
	newState  T
	err       error
	partial   bool
}

type runtimeImpl[T g.SharedState] struct {
	ctx    context.Context
	cancel context.CancelFunc

	outcomeCh      chan nodeFnReturnStruct[T]
	stateMonitorCh chan g.StateMonitorEntry[T]

	startEdge g.Edge[T]
	edges     []g.Edge[T]

	state          T
	stateMergeLock *sync.Mutex

	executing atomic.Bool

	identity       uuid.UUID
	persistFn      g.PersistFn[T]
	restoreFn      g.RestoreFn[T]
	lastPersisted  T
	persistLock    sync.RWMutex
	pendingPersist chan T
	persistWg      sync.WaitGroup
}

func (r *runtimeImpl[T]) Invoke(userInput T) {
	// Use atomic compare-and-swap to prevent concurrent invocations
	if !r.executing.CompareAndSwap(false, true) {
		// Already executing, send error to monitor channel
		if r.stateMonitorCh != nil {
			r.stateMonitorCh <- GraphError("Runtime", r.CurrentState(), fmt.Errorf("runtime is already executing, concurrent invocations not allowed"))
		}
		return
	}

	r.startEdge.From().Accept(userInput, r)
}

func (r *runtimeImpl[T]) AddEdge(edge ...g.Edge[T]) {
	r.edges = append(r.edges, edge...)
}

func (r *runtimeImpl[T]) Validate() error {
	if r.startEdge.From() == nil {
		return fmt.Errorf("graph validation failed: start edge 'from' node is nil")
	}

	// Check if there's at least one path from start to an end edge
	visited := make(map[string]bool)
	// Include the start edge in the traversal by starting from its target node
	hasPathToEnd := r.hasPathToEndEdge(r.startEdge.To(), visited)
	if !hasPathToEnd {
		return fmt.Errorf("graph validation failed: no path from start edge to any end edge")
	}

	return nil
}

func (r *runtimeImpl[T]) Shutdown() {
	r.cancel()
	r.persistWg.Wait()
}

func (r *runtimeImpl[T]) NotifyStateChange(
	node g.Node[T],
	userInput T,
	newState T,
	err error,
	partial bool,
) {
	r.outcomeCh <- nodeFnReturnStruct[T]{node: node, userInput: userInput, newState: newState, err: err, partial: partial}
}

func (r *runtimeImpl[T]) StartEdge() g.Edge[T] {
	return r.startEdge
}

func (r *runtimeImpl[T]) CurrentState() T {
	r.stateMergeLock.Lock()
	defer r.stateMergeLock.Unlock()
	return r.state
}

func (r *runtimeImpl[T]) SetPersistentState(
	persist g.PersistFn[T],
	restore g.RestoreFn[T],
	runtimeID uuid.UUID,
) {
	r.persistFn = persist
	r.restoreFn = restore
	r.identity = runtimeID
}

func (r *runtimeImpl[T]) Restore() error {
	if r.restoreFn == nil {
		return fmt.Errorf("restore function is not set")
	}
	if r.identity == uuid.Nil {
		return fmt.Errorf("runtime identity is not set")
	}
	restoredState, err := r.restoreFn(r.ctx, r.identity)
	if err != nil {
		return fmt.Errorf("state restoration failed: %w", err)
	}
	r.stateMergeLock.Lock()
	r.state = restoredState
	r.stateMergeLock.Unlock()
	r.persistLock.Lock()
	r.lastPersisted = restoredState
	r.persistLock.Unlock()
	return nil
}

func (r *runtimeImpl[T]) persistState() error {
	if r.persistFn == nil {
		return nil
	}
	if r.identity == uuid.Nil {
		return fmt.Errorf("runtime identity is not set")
	}
	r.stateMergeLock.Lock()
	currentState := r.state
	r.stateMergeLock.Unlock()

	r.persistLock.RLock()
	lastPersisted := r.lastPersisted
	r.persistLock.RUnlock()

	if r.statesEqual(currentState, lastPersisted) {
		return nil
	}

	select {
	case r.pendingPersist <- currentState:
	default:
	}

	return nil
}

func (r *runtimeImpl[T]) start() {
	go r.onStateChange()
}

func (r *runtimeImpl[T]) onStateChange() {
	for {
		select {
		case <-r.ctx.Done():
			r.executing.Store(false)
			return
		case result := <-r.outcomeCh:
			if result.err != nil {
				if r.stateMonitorCh != nil {
					r.stateMonitorCh <- GraphError(result.node.Name(), r.state, result.err)
					r.executing.Store(false)
				}
				continue
			}

			if result.partial {
				if r.stateMonitorCh != nil {
					r.stateMonitorCh <- GraphPartial(result.node.Name(), result.newState)
				}
				continue
			}

			previous := r.replace(result.newState)
			err := r.persistState()
			if err != nil {
				if r.stateMonitorCh != nil {
					r.stateMonitorCh <- GraphNonFatalError[T](result.node.Name(), fmt.Errorf("state persistence error: %w", err))
				}
			}

			if result.node.Role() == g.EndNode {
				if r.stateMonitorCh != nil {
					r.stateMonitorCh <- GraphCompleted(result.node.Name(), r.state)
					r.executing.Store(false)
				}
				continue
			} else {
				if r.stateMonitorCh != nil {
					r.stateMonitorCh <- GraphRunning(result.node.Name(), previous, r.state)
				}
			}

			outboundEdges := r.edgesFrom(result.node)
			if len(outboundEdges) == 0 {
				r.stateMonitorCh <- GraphError(result.node.Name(), r.state, fmt.Errorf("no outbound edges from node %s", result.node.Name()))
				r.executing.Store(false)
				continue
			}

			policy := result.node.RoutePolicy()
			if policy == nil {
				r.stateMonitorCh <- GraphError(result.node.Name(), r.state, fmt.Errorf("node %s has no routing policy", result.node.Name()))
				r.executing.Store(false)
				continue
			}

			nextEdge := policy.SelectEdge(result.userInput, r.state, outboundEdges)
			if nextEdge == nil {
				r.stateMonitorCh <- GraphError(result.node.Name(), r.state, fmt.Errorf("routing policy for node %s returned nil edge", result.node.Name()))
				r.executing.Store(false)
				continue
			}

			nextNode := nextEdge.To()
			if nextNode == nil {
				r.stateMonitorCh <- GraphError(result.node.Name(), r.state, fmt.Errorf("next edge from node %s has nil target node", result.node.Name()))
				r.executing.Store(false)
				continue
			}

			nextNode.Accept(result.userInput, r)
		}
	}
}

func (r *runtimeImpl[T]) replace(newState T) T {
	r.stateMergeLock.Lock()
	defer r.stateMergeLock.Unlock()

	previous := r.state
	r.state = newState
	return previous
}

func (r *runtimeImpl[T]) edgesFrom(node g.Node[T]) []g.Edge[T] {
	if r.startEdge.From() == node {
		return []g.Edge[T]{r.StartEdge()}
	}
	var outboundEdges []g.Edge[T]
	for _, edge := range r.edges {
		if edge.From() == node {
			outboundEdges = append(outboundEdges, edge)
		}
	}
	return outboundEdges
}

func (r *runtimeImpl[T]) hasPathToEndEdge(node g.Node[T], visited map[string]bool) bool {
	// Check if the node is an EndNode
	if node.Role() == g.EndNode {
		return true
	}

	// Mark the node as visited
	nodeKey := fmt.Sprintf("%p", node)
	if visited[nodeKey] {
		return false
	}
	visited[nodeKey] = true

	// Check if any EndEdge starts from this node
	for _, edge := range r.edges {
		if edge.Role() == g.EndEdge {
			if edge.From() == node {
				return true
			}
		}
	}

	// Explore all edges to find connected nodes
	for _, edge := range r.edges {
		if edge.From() == node {
			if r.hasPathToEndEdge(edge.To(), visited) {
				return true
			}
		}
	}

	return false
}

func (r *runtimeImpl[T]) startPersistenceWorker() {
	r.persistWg.Add(1)
	go r.persistenceWorker()
}

func (r *runtimeImpl[T]) persistenceWorker() {
	defer r.persistWg.Done()

	for {
		select {
		case <-r.ctx.Done():
			r.flushPendingStates()
			return
		case state := <-r.pendingPersist:
			if err := r.persistFn(r.ctx, r.identity, state); err != nil {
				if r.stateMonitorCh != nil {
					r.stateMonitorCh <- GraphNonFatalError[T]("Persistence", fmt.Errorf("state persistence error: %w", err))
				}
			} else {
				r.persistLock.Lock()
				r.lastPersisted = state
				r.persistLock.Unlock()
			}
		}
	}
}

func (r *runtimeImpl[T]) flushPendingStates() {
	for {
		select {
		case state := <-r.pendingPersist:
			if err := r.persistFn(r.ctx, r.identity, state); err != nil {
				if r.stateMonitorCh != nil {
					r.stateMonitorCh <- GraphNonFatalError[T]("Persistence", fmt.Errorf("state persistence error during flush: %w", err))
				}
			}
		default:
			return
		}
	}
}

func (r *runtimeImpl[T]) statesEqual(a, b T) bool {
	return fmt.Sprintf("%v", a) == fmt.Sprintf("%v", b)
}
