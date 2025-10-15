package graph

// Edge represents an edge in the graph.
type Edge[T SharedState] interface {
	// From returns the source node of the edge.
	From() Node[T]
	// To returns the destination node of the edge.
	To() Node[T]
	// LabelByKey returns the label value associated with the specified key and a boolean indicating if the key exists.
	LabelByKey(key string) (string, bool)
}
