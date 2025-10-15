package graph

// RoutePolicy defines a policy for routing between nodes in the graph.
type RoutePolicy[T SharedState] interface {
	// SelectEdge selects an edge from the available edges based on the current state.
	SelectEdge(userInput T, currentState T, edges []Edge[T]) Edge[T]
}
