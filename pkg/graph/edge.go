package graph

// CreateEdge creates a new instance of Edge with the specified SharedState type.
func CreateEdge[T SharedState](from, to Node[T]) Edge[T] {
	return &edgeImpl[T]{from: from, to: to}
}

// TODO predicates

// Edge represents an edge in the graph.
type Edge[T SharedState] interface {
	To() Node[T]
}

type edgeImpl[T SharedState] struct {
	from Node[T]
	to   Node[T]
}

func (e *edgeImpl[T]) To() Node[T] {
	return e.to
}
