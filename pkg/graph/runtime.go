package graph

import (
	"context"
	"fmt"
	"sync"
)

// CreateRuntime creates a new instance of Runtime with the specified SharedState type.
func CreateRuntime[T SharedState](
	startEdge *StartEdge[T],
	stateMonitorCh chan StateMonitorEntry[T],
) (Runtime[T], error) {
	return CreateRuntimeWithMerger(startEdge, stateMonitorCh, nil)
}

// CreateRuntimeWithMerger creates a new instance of Runtime with the specified SharedState type and state merger function.
func CreateRuntimeWithMerger[T SharedState](
	startEdge *StartEdge[T],
	stateMonitorCh chan StateMonitorEntry[T],
	merger StateMergeFn[T],
) (Runtime[T], error) {
	var zero T
	return CreateRuntimeWithMergerAndInitialState(startEdge, stateMonitorCh, merger, zero)
}

// CreateRuntimeWithMergerAndInitialState creates a new instance of Runtime with the specified SharedState type, state merger function, and initial state.
func CreateRuntimeWithMergerAndInitialState[T SharedState](
	startEdge *StartEdge[T],
	stateMonitorCh chan StateMonitorEntry[T],
	merger StateMergeFn[T],
	initialState T,
) (Runtime[T], error) {
	if startEdge == nil {
		return nil, fmt.Errorf("runtime creation failed: start edge cannot be nil")
	}
	ctx, cancelFn := context.WithCancel(context.Background())
	rv := &runtimeImpl[T]{
		ctx:    ctx,
		cancel: cancelFn,

		outcomeCh:      make(chan nodeFnReturnStruct[T]),
		stateMonitorCh: stateMonitorCh,

		startEdge: *startEdge,
		edges:     []Edge[T]{},

		state:  initialState,
		merger: merger,
		lock:   &sync.RWMutex{},
	}
	rv.start()
	return rv, nil
}

// Connected provides access to the connected graph components.
type Connected[T SharedState] interface {
	// AddEdge adds an edge to the runtime's graph.
	AddEdge(edge ...Edge[T])
	// Validate checks the integrity of the graph structure.
	Validate() error
}

// Runtime represents the runtime environment for graph processing.
type Runtime[T SharedState] interface {
	Connected[T]
	// Invoke executes the graph processing with the given entry state.
	Invoke(entryState T)
	// Shutdown gracefully stops the runtime processing.
	Shutdown()
}

var _ Runtime[SharedState] = (*runtimeImpl[SharedState])(nil)

type nodeFnReturnStruct[T SharedState] struct {
	node  Node[T]
	state T
	err   error
}

type runtimeImpl[T SharedState] struct {
	ctx    context.Context
	cancel context.CancelFunc

	outcomeCh      chan nodeFnReturnStruct[T]
	stateMonitorCh chan StateMonitorEntry[T]

	startEdge StartEdge[T]
	edges     []Edge[T]

	state  T
	merger StateMergeFn[T]
	lock   *sync.RWMutex
}

func (r *runtimeImpl[T]) Invoke(entryState T) {
	r.startEdge.from.Accept(entryState, r)
}

func (r *runtimeImpl[T]) AddEdge(edge ...Edge[T]) {
	r.edges = append(r.edges, edge...)
}

func (r *runtimeImpl[T]) Validate() error {
	if r.startEdge.from == nil {
		return fmt.Errorf("graph validation failed: start edge 'from' node is nil")
	}

	// Check if there's at least one path from start to an end edge
	visited := make(map[string]bool)
	// Include the start edge in the traversal by starting from its target node
	hasPathToEnd := r.hasPathToEndEdge(r.startEdge.to, visited)
	if !hasPathToEnd {
		return fmt.Errorf("graph validation failed: no path from start edge to any end edge")
	}

	return nil
}

func (r *runtimeImpl[T]) Shutdown() {
	r.cancel()
}

func (r *runtimeImpl[T]) NotifyStateChange(node Node[T], state T, err error) {
	r.outcomeCh <- nodeFnReturnStruct[T]{node: node, state: state, err: err}
}

func (r *runtimeImpl[T]) start() {
	go r.onStatusChange()
}

func (r *runtimeImpl[T]) onStatusChange() {
	for {
		select {
		case <-r.ctx.Done():
			return
		case result := <-r.outcomeCh:
			if result.err != nil {
				if r.stateMonitorCh != nil {
					r.stateMonitorCh <- GraphError(result.node.Name(), r.state, result.err)
				}
				continue
			}

			previous := r.merge(r.state, result.state)
			if r.stateMonitorCh != nil {
				r.stateMonitorCh <- GraphRunning(result.node.Name(), previous, r.state)
			}

			_, isEndNode := result.node.(*EndNode[T])
			if isEndNode {
				if r.stateMonitorCh != nil {
					r.stateMonitorCh <- GraphCompleted(r.state)
				}
				continue
			}

			outboundEdges := r.edgesFrom(result.node)
			if len(outboundEdges) == 0 {
				r.stateMonitorCh <- GraphError(result.node.Name(), r.state, fmt.Errorf("no outbound edges from node %s", result.node.Name()))
				continue
			}

			policy := result.node.RoutePolicy()
			if policy == nil {
				r.stateMonitorCh <- GraphError(result.node.Name(), r.state, fmt.Errorf("node %s has no routing policy", result.node.Name()))
				continue
			}

			nextEdge := policy.SelectEdge(r.state, outboundEdges)
			if nextEdge == nil {
				r.stateMonitorCh <- GraphError(result.node.Name(), r.state, fmt.Errorf("routing policy for node %s returned nil edge", result.node.Name()))
				continue
			}

			nextNode := nextEdge.To()
			if nextNode == nil {
				r.stateMonitorCh <- GraphError(result.node.Name(), r.state, fmt.Errorf("next edge from node %s has nil target node", result.node.Name()))
				continue
			}

			nextNode.Accept(r.state, r)
		}
	}
}

func (r *runtimeImpl[T]) merge(current, other T) T {
	r.lock.Lock()
	defer r.lock.Unlock()

	previous := r.state

	if r.merger == nil || any(current) == nil {
		r.state = other
	} else {
		r.state = r.merger(current, other)
	}

	return previous
}

func (r *runtimeImpl[T]) edgesFrom(node Node[T]) []Edge[T] {
	if r.startEdge.from == node {
		return []Edge[T]{&r.startEdge}
	}
	var outboundEdges []Edge[T]
	for _, edge := range r.edges {
		if edgeFrom(edge) == node {
			outboundEdges = append(outboundEdges, edge)
		}
	}
	return outboundEdges
}

func (r *runtimeImpl[T]) hasPathToEndEdge(node Node[T], visited map[string]bool) bool {
	// Check if the node is an EndNode
	if _, ok := node.(*EndNode[T]); ok {
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
		if endEdge, ok := edge.(*EndEdge[T]); ok {
			if edgeFrom[T](endEdge) == node {
				return true
			}
		}
	}

	// Explore all edges to find connected nodes
	for _, edge := range r.edges {
		if edgeFrom[T](edge) == node {
			if r.hasPathToEndEdge(edgeTo[T](edge), visited) {
				return true
			}
		}
	}

	return false
}

func edgeFrom[T SharedState](edge Edge[T]) Node[T] {
	switch e := edge.(type) {
	case *edgeImpl[T]:
		return e.from
	case *StartEdge[T]:
		return e.from
	case *EndEdge[T]:
		return e.from
	default:
		return nil
	}
}

func edgeTo[T SharedState](edge Edge[T]) Node[T] {
	switch e := edge.(type) {
	case *edgeImpl[T]:
		return e.to
	case *StartEdge[T]:
		return e.to
	case *EndEdge[T]:
		return e.to
	default:
		return nil
	}
}
