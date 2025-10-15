package graph

// NodeRole represents the role of a node in the graph.
type NodeRole int

const (
	// StartNode represents the starting point of the graph.
	StartNode NodeRole = iota
	// IntermediateNode represents a node that is neither a start nor an end node.
	IntermediateNode
	// EndNode represents the endpoint of the graph.
	EndNode
)

// Node represents a node in the graph.
type Node[T SharedState] interface {
	// Accept processes the node with the given state and returns the updated state.
	Accept(deltaState T, runtime StateObserver[T])
	// Name returns the name of the node.
	Name() string
	// RoutePolicy returns the routing policy associated with the node.
	RoutePolicy() RoutePolicy[T]
	// Role returns the role of the node in the graph.
	Role() NodeRole
}
