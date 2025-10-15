package graph

// Node represents a node in the graph.
type Node[T SharedState] interface {
	// Accept processes the node with the given state and returns the updated state.
	Accept(deltaState T, runtime StateObserver[T])
	// Name returns the name of the node.
	Name() string
	// RoutePolicy returns the routing policy associated with the node.
	RoutePolicy() RoutePolicy[T]
}
