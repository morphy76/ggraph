package builders

import (
	i "github.com/morphy76/ggraph/internal/graph"
	g "github.com/morphy76/ggraph/pkg/graph"
)

// CreateEdge creates a new instance of Edge with the specified SharedState type.
func CreateEdge[T g.SharedState](from, to g.Node[T], labels ...map[string]string) g.Edge[T] {
	return i.EdgeImplFactory(from, to, g.IntermediateEdge, labels...)
}

// CreateStartEdge creates a new instance of StartEdge with the specified SharedState type.
func CreateStartEdge[T g.SharedState](to g.Node[T]) g.Edge[T] {
	return i.EdgeImplFactory(nil, to, g.StartEdge)
}

// CreateEndEdge creates a new instance of EndEdge with the specified SharedState type.
func CreateEndEdge[T g.SharedState](from g.Node[T]) g.Edge[T] {
	return i.EdgeImplFactory(from, nil, g.EndEdge)
}
