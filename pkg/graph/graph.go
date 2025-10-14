package graph

// StateObserver is an interface for observing state changes in nodes during graph processing.
type StateObserver[T SharedState] interface {
	// NotifyStateChange is called when a node changes state during processing.
	NotifyStateChange(node Node[T], state T, err error, partial bool)
}

// CreateStartNode creates a new instance of StartNode with the specified SharedState type.
func CreateStartNode[T SharedState]() Node[T] {
	policy, _ := CreateAnyRoutePolicy[T]()
	return &StartNode[T]{
		policy: policy,
	}
}

var _ Node[SharedState] = (*StartNode[SharedState])(nil)

// StartNode represents the starting node of a graph.
type StartNode[T SharedState] struct {
	policy RoutePolicy[T]
}

func (n *StartNode[T]) Name() string {
	return "StartNode"
}

func (n *StartNode[T]) Accept(state T, runtime StateObserver[T]) {
	go runtime.NotifyStateChange(n, state, nil, false)
}

func (n *StartNode[T]) RoutePolicy() RoutePolicy[T] {
	return n.policy
}

// CreateEndNode creates a new instance of EndNode with the specified SharedState type.
func CreateEndNode[T SharedState]() Node[T] {
	return &EndNode[T]{}
}

var _ Node[SharedState] = (*EndNode[SharedState])(nil)

// EndNode represents the ending node of a graph.
type EndNode[T SharedState] struct {
	policy RoutePolicy[T]
}

func (n *EndNode[T]) Name() string {
	return "EndNode"
}

func (n *EndNode[T]) Accept(state T, runtime StateObserver[T]) {
	go runtime.NotifyStateChange(n, state, nil, false)
}

func (n *EndNode[T]) RoutePolicy() RoutePolicy[T] {
	return nil
}

// CreateStartEdge creates a new instance of StartEdge with the specified SharedState type.
func CreateStartEdge[T SharedState](to Node[T]) *StartEdge[T] {
	return &StartEdge[T]{edgeImpl: edgeImpl[T]{from: CreateStartNode[T](), to: to, labels: map[string]string{"type": "start"}}}
}

var _ Edge[SharedState] = (*StartEdge[SharedState])(nil)

// StartEdge represents the starting edge of a graph.
type StartEdge[T SharedState] struct {
	edgeImpl[T]
}

func (e *StartEdge[T]) To() Node[T] {
	return e.to
}

// CreateEndEdge creates a new instance of EndEdge with the specified SharedState type.
func CreateEndEdge[T SharedState](from Node[T]) *EndEdge[T] {
	return &EndEdge[T]{edgeImpl: edgeImpl[T]{from: from, to: CreateEndNode[T](), labels: map[string]string{"type": "end"}}}
}

var _ Edge[SharedState] = (*EndEdge[SharedState])(nil)

// EndEdge represents the ending edge of a graph.
type EndEdge[T SharedState] struct {
	edgeImpl[T]
}

func (e *EndEdge[T]) To() Node[T] {
	return e.to
}
