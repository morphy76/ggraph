package graph

import (
	"context"
	"fmt"
	"reflect"
	"sync"
	"sync/atomic"
	"time"

	g "github.com/morphy76/ggraph/pkg/graph"
)

type pendingPersistEntry[T g.SharedState] struct {
	threadID string
	state    T
}

// RuntimeFactory creates a new instance of Runtime with the specified SharedState type, state merger function, and initial state.
func RuntimeFactory[T g.SharedState](
	startEdge g.Edge[T],
	stateMonitorCh chan g.StateMonitorEntry[T],
	initialState T,
) (g.Runtime[T], error) {
	if startEdge == nil {
		return nil, fmt.Errorf("runtime creation failed: %w", g.ErrStartEdgeNil)
	}
	ctx, cancelFn := context.WithCancel(context.Background())
	rv := &runtimeImpl[T]{
		ctx:    ctx,
		cancel: cancelFn,

		outcomeCh:      make(chan nodeFnReturnStruct[T], 1000),
		stateMonitorCh: stateMonitorCh,

		startEdge: startEdge,
		edges:     []g.Edge[T]{},

		initialState:    initialState,
		state:           make(map[string]T),
		stateChangeLock: make(map[string]*sync.Mutex),

		executing: make(map[string]*atomic.Bool),

		lastPersisted: make(map[string]T),

		pendingPersist: make(chan pendingPersistEntry[T], 10),

		threadTTL: make(map[string]time.Time),
	}
	rv.start()
	rv.startPersistenceWorker()
	rv.startThreadEvictor()
	return rv, nil
}

// ------------------------------------------------------------------------------
// Runtime Implementation
// ------------------------------------------------------------------------------

var _ g.Runtime[g.SharedState] = (*runtimeImpl[g.SharedState])(nil)
var _ g.StateObserver[g.SharedState] = (*runtimeImpl[g.SharedState])(nil)

type nodeFnReturnStruct[T g.SharedState] struct {
	node        g.Node[T]
	userInput   T
	stateChange T
	err         error
	partial     bool
	reducer     g.ReducerFn[T]
	config      g.InvokeConfig
}

type runtimeImpl[T g.SharedState] struct {
	ctx    context.Context
	cancel context.CancelFunc

	outcomeCh      chan nodeFnReturnStruct[T]
	stateMonitorCh chan g.StateMonitorEntry[T]

	startEdge g.Edge[T]
	edges     []g.Edge[T]

	initialState    T
	state           map[string]T
	stateChangeLock map[string]*sync.Mutex

	executing map[string]*atomic.Bool

	persistFn     g.PersistFn[T]
	restoreFn     g.RestoreFn[T]
	lastPersisted map[string]T

	pendingPersist chan pendingPersistEntry[T]

	threadTTL map[string]time.Time

	backgroundWorkers sync.WaitGroup
}

func (r *runtimeImpl[T]) Invoke(userInput T, configs ...g.InvokeConfig) string {
	requestedConfig := g.MergeInvokeConfig(configs...)
	useConfig := g.MergeInvokeConfig(g.DefaultInvokeConfig(), requestedConfig)

	if !r.threadExistsWithinTTL(useConfig.ThreadID) {
		_ = r.Restore(useConfig.ThreadID)
	}

	r.threadTTL[useConfig.ThreadID] = time.Now().Add(1 * time.Hour)

	if !r.executingByThreadID(useConfig).CompareAndSwap(false, true) {
		r.sendMonitorEntry(monitorError[T]("Runtime", useConfig.ThreadID, fmt.Errorf("cannot invoke graph for thread %s: %w", useConfig.ThreadID, g.ErrRuntimeExecuting)))
		return useConfig.ThreadID
	}

	r.startEdge.From().Accept(userInput, r, useConfig)
	return useConfig.ThreadID
}

func (r *runtimeImpl[T]) AddEdge(edge ...g.Edge[T]) {
	r.edges = append(r.edges, edge...)
}

func (r *runtimeImpl[T]) Validate() error {
	if r.startEdge.From() == nil {
		return fmt.Errorf("graph validation failed: %w", g.ErrStartNodeNil)
	}

	// Check if there's at least one path from start to an end edge
	visited := make(map[string]bool)
	// Include the start edge in the traversal by starting from its target node
	hasPathToEnd := r.hasPathToEndEdge(r.startEdge.To(), visited)
	if !hasPathToEnd {
		return fmt.Errorf("graph validation failed: %w", g.ErrNoPathToEnd)
	}

	return nil
}

func (r *runtimeImpl[T]) Shutdown() {
	r.cancel()

	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()

	done := make(chan struct{})
	go func() {
		r.backgroundWorkers.Wait()
		close(done)
	}()

	select {
	case <-done:
	case <-ctx.Done():
	}
}

func (r *runtimeImpl[T]) NotifyStateChange(
	node g.Node[T],
	config g.InvokeConfig,
	userInput T,
	stateChange T,
	reducer g.ReducerFn[T],
	err error,
	partial bool,
) {
	r.outcomeCh <- nodeFnReturnStruct[T]{node: node, userInput: userInput, stateChange: stateChange, err: err, partial: partial, reducer: reducer, config: config}
}

func (r *runtimeImpl[T]) CurrentState(threadID string) T {
	useLock := r.lockByThreadID(threadID)
	useLock.Lock()
	defer useLock.Unlock()

	useState := r.initialState
	if state, exists := r.state[threadID]; exists {
		useState = state
	}
	return useState
}

func (r *runtimeImpl[T]) InitialState() T {
	return r.initialState
}

func (r *runtimeImpl[T]) StartEdge() g.Edge[T] {
	return r.startEdge
}

func (r *runtimeImpl[T]) SetPersistentState(
	persist g.PersistFn[T],
	restore g.RestoreFn[T],
) {
	r.persistFn = persist
	r.restoreFn = restore
}

func (r *runtimeImpl[T]) Restore(threadID string) error {
	if r.restoreFn == nil {
		return nil
	}
	restoredState, err := r.restoreFn(r.ctx, threadID)
	if err != nil {
		return fmt.Errorf("state restoration failed: %w", err)
	}

	useLock := r.lockByThreadID(threadID)
	useLock.Lock()
	r.state[threadID] = restoredState
	r.lastPersisted[threadID] = restoredState
	useLock.Unlock()

	return nil
}

func (r *runtimeImpl[T]) persistState(threadID string) error {
	if r.persistFn == nil {
		return nil
	}

	useLock := r.lockByThreadID(threadID)
	useLock.Lock()
	defer useLock.Unlock()

	currentState := r.state[threadID]
	lastPersisted := r.lastPersisted[threadID]

	if r.statesEqual(currentState, lastPersisted) {
		return nil
	}

	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	select {
	case r.pendingPersist <- pendingPersistEntry[T]{threadID: threadID, state: currentState}:
	case <-ctx.Done():
		r.sendMonitorEntry(monitorNonFatalError[T]("Persistence", threadID, fmt.Errorf("persistence timed out: %w", ctx.Err())))
	default:
		r.sendMonitorEntry(monitorNonFatalError[T]("Persistence", threadID, fmt.Errorf("cannot persist state: %w", g.ErrPersistenceQueueFull)))
	}

	return nil
}

func (r *runtimeImpl[T]) start() {
	go r.onNodeOutcome()
}

func (r *runtimeImpl[T]) onNodeOutcome() {
	for {
		select {
		case <-r.ctx.Done():
			return
		case result := <-r.outcomeCh:
			useThreadID := result.config.ThreadID
			useInvocationContext := result.config.Context
			useExecuting := r.executingByThreadID(result.config)

			if result.err != nil {
				r.sendMonitorEntry(monitorError[T](result.node.Name(), useThreadID, result.err))
				useExecuting.Store(false)
				r.clearThread(useThreadID)
				continue
			}

			select {
			case <-useInvocationContext.Done():
				err := r.persistState(useThreadID)
				if err != nil {
					r.sendMonitorEntry(monitorNonFatalError[T](result.node.Name(), useThreadID, fmt.Errorf("state persistence error: %w", err)))
				}
				r.sendMonitorEntry(monitorError[T](result.node.Name(), useThreadID, fmt.Errorf("invocation context done: %w", useInvocationContext.Err())))
				useExecuting.Store(false)
				r.clearThread(useThreadID)
				continue
			default:
				if result.partial {
					r.sendMonitorEntry(monitorPartial(result.node.Name(), useThreadID, result.stateChange))
					continue
				}

				newState := r.replace(useThreadID, result.stateChange, result.reducer)
				err := r.persistState(useThreadID)
				if err != nil {
					r.sendMonitorEntry(monitorNonFatalError[T](result.node.Name(), useThreadID, fmt.Errorf("state persistence error: %w", err)))
				}

				if result.node.Role() == g.EndNode {
					if r.stateMonitorCh != nil {
						r.sendMonitorEntry(monitorCompleted(result.node.Name(), useThreadID, newState))
					}
					useExecuting.Store(false)
					// Don't clear thread state immediately if there's no persistence
					// This allows CurrentState() to return the final state
					if r.persistFn != nil {
						r.clearThread(useThreadID)
					}
					continue
				} else {
					if r.stateMonitorCh != nil {
						r.sendMonitorEntry(monitorRunning(result.node.Name(), useThreadID, newState))
					}
				}

				outboundEdges := r.edgesFrom(result.node)
				if len(outboundEdges) == 0 {
					r.sendMonitorEntry(monitorError[T](result.node.Name(), useThreadID, fmt.Errorf("routing error for node %s: %w", result.node.Name(), g.ErrNoOutboundEdges)))
					useExecuting.Store(false)
					r.clearThread(useThreadID)
					continue
				}

				policy := result.node.RoutePolicy()
				if policy == nil {
					r.sendMonitorEntry(monitorError[T](result.node.Name(), useThreadID, fmt.Errorf("routing error for node %s: %w", result.node.Name(), g.ErrNoRoutingPolicy)))
					useExecuting.Store(false)
					r.clearThread(useThreadID)
					continue
				}

				nextEdge := policy.SelectEdge(result.userInput, r.state[useThreadID], outboundEdges)
				if nextEdge == nil {
					r.sendMonitorEntry(monitorError[T](result.node.Name(), useThreadID, fmt.Errorf("routing error for node %s: %w", result.node.Name(), g.ErrNilEdge)))
					useExecuting.Store(false)
					r.clearThread(useThreadID)
					continue
				}

				nextNode := nextEdge.To()
				if nextNode == nil {
					r.sendMonitorEntry(monitorError[T](result.node.Name(), useThreadID, fmt.Errorf("routing error for node %s: %w", result.node.Name(), g.ErrNextEdgeNil)))
					useExecuting.Store(false)
					r.clearThread(useThreadID)
					continue
				}

				nextNode.Accept(result.userInput, r, result.config)
			}
		}
	}
}

func (r *runtimeImpl[T]) sendMonitorEntry(entry g.StateMonitorEntry[T]) {
	if r.stateMonitorCh == nil {
		return
	}

	select {
	case r.stateMonitorCh <- entry:
	case <-time.After(100 * time.Millisecond):
	case <-r.ctx.Done():
	}
}

func (r *runtimeImpl[T]) replace(threadID string, stateChange T, reducer g.ReducerFn[T]) T {
	useLock := r.lockByThreadID(threadID)
	useLock.Lock()
	defer useLock.Unlock()

	// Get current state without calling CurrentState() to avoid deadlock
	useState := r.initialState
	if state, exists := r.state[threadID]; exists {
		useState = state
	}

	r.state[threadID] = reducer(useState, stateChange)
	return r.state[threadID]
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
	r.backgroundWorkers.Add(1)
	go r.persistenceWorker()
}

func (r *runtimeImpl[T]) persistenceWorker() {
	defer r.backgroundWorkers.Done()

	for {
		select {
		case <-r.ctx.Done():
			r.flushPendingStates()
			return
		case state := <-r.pendingPersist:
			if err := r.persistFn(r.ctx, state.threadID, state.state); err != nil {
				r.sendMonitorEntry(monitorNonFatalError[T]("Persistence", state.threadID, fmt.Errorf("state persistence error: %w", err)))
			}
		}
	}
}

func (r *runtimeImpl[T]) startThreadEvictor() {
	r.backgroundWorkers.Add(1)
	go r.threadEvictor()
}

func (r *runtimeImpl[T]) threadEvictor() {
	defer r.backgroundWorkers.Done()

	ticker := time.NewTicker(10 * time.Minute)
	defer ticker.Stop()
	for {
		select {
		case <-r.ctx.Done():
			return
		case <-ticker.C:
			now := time.Now()
			for threadID, expiry := range r.threadTTL {
				if now.After(expiry) {
					err := r.persistState(threadID)
					if err != nil {
						r.sendMonitorEntry(monitorNonFatalError[T]("ThreadEvictor", threadID, fmt.Errorf("state persistence error during eviction: %w", err)))
					}

					r.clearThread(threadID)

					r.sendMonitorEntry(monitorNonFatalError[T]("ThreadEvictor", threadID, fmt.Errorf("evicted thread %s: %w", threadID, g.ErrEvictionByInactivity)))
				}
			}
		}
	}
}

func (r *runtimeImpl[T]) flushPendingStates() {
	for {
		select {
		case state := <-r.pendingPersist:
			if err := r.persistFn(r.ctx, state.threadID, state.state); err != nil {
				r.sendMonitorEntry(monitorNonFatalError[T]("Persistence", state.threadID, fmt.Errorf("state persistence error during flush: %w", err)))
			}
		default:
			return
		}
	}
}

func (r *runtimeImpl[T]) statesEqual(a, b T) bool {
	return reflect.DeepEqual(a, b)
}

func (r *runtimeImpl[T]) executingByThreadID(config g.InvokeConfig) *atomic.Bool {
	exec, exists := r.executing[config.ThreadID]
	if !exists {
		exec = &atomic.Bool{}
		r.executing[config.ThreadID] = exec
	}
	return exec
}

func (r *runtimeImpl[T]) lockByThreadID(threadID string) *sync.Mutex {
	lock, exists := r.stateChangeLock[threadID]
	if !exists {
		lock = &sync.Mutex{}
		r.stateChangeLock[threadID] = lock
	}
	return lock
}

func (r *runtimeImpl[T]) threadExistsWithinTTL(threadID string) bool {
	ttl, exists := r.threadTTL[threadID]
	return exists && time.Now().Before(ttl)
}

func (r *runtimeImpl[T]) clearThread(threadID string) {
	delete(r.threadTTL, threadID)
	useLock := r.lockByThreadID(threadID)
	useLock.Lock()
	delete(r.state, threadID)
	useLock.Unlock()

	delete(r.executing, threadID)
	delete(r.lastPersisted, threadID)
	delete(r.stateChangeLock, threadID)
}
