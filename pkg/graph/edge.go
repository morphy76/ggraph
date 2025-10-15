package graph

// EdgeRole defines the role of an edge in the graph.
type EdgeRole int

const (
	// StartEdge indicates the starting edge of the graph.
	StartEdge EdgeRole = iota
	// EndEdge indicates the ending edge of the graph.
	EndEdge
	// IntermediateEdge indicates an intermediate edge in the graph.
	IntermediateEdge
)

// Edge represents an edge in the graph.
type Edge[T SharedState] interface {
	// From returns the source node of the edge.
	From() Node[T]
	// To returns the destination node of the edge.
	To() Node[T]
	// LabelByKey returns the label value associated with the specified key and a boolean indicating if the key exists.
	LabelByKey(key string) (string, bool)
	// Role returns the role of the edge in the graph.
	Role() EdgeRole
}
