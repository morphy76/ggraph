package graph

// CreateEdge creates a new instance of Edge with the specified SharedState type.
func CreateEdge[T SharedState](from, to Node[T]) Edge[T] {
	return &edgeImpl[T]{from: from, to: to}
}

// CreateStartEdge creates a new instance of StartEdge with the specified SharedState type.
func CreateStartEdge[T SharedState](to Node[T]) *StartEdge[T] {
	return &StartEdge[T]{edgeImpl: edgeImpl[T]{from: CreateStartNode[T](), to: to}}
}

// CreateEndEdge creates a new instance of EndEdge with the specified SharedState type.
func CreateEndEdge[T SharedState](from Node[T]) *EndEdge[T] {
	return &EndEdge[T]{edgeImpl: edgeImpl[T]{from: from, to: CreateEndNode[T]()}}
}

// Edge represents an edge in the graph.
type Edge[T SharedState] interface {
	To() Node[T]
}

var _ Edge[SharedState] = (*StartEdge[SharedState])(nil)

// StartEdge represents the starting edge of a graph.
type StartEdge[T SharedState] struct {
	edgeImpl[T]
}

func (e *StartEdge[T]) To() Node[T] {
	return e.to
}

var _ Edge[SharedState] = (*EndEdge[SharedState])(nil)

// EndEdge represents the ending edge of a graph.
type EndEdge[T SharedState] struct {
	edgeImpl[T]
}

func (e *EndEdge[T]) To() Node[T] {
	return e.to
}

type edgeImpl[T SharedState] struct {
	from Node[T]
	to   Node[T]
}

func (e *edgeImpl[T]) To() Node[T] {
	return e.to
}
