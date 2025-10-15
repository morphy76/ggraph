package graph

import (
	g "github.com/morphy76/ggraph/pkg/graph"
)

// EdgeImplFactory creates a new instance of Edge with the specified SharedState type.
func EdgeImplFactory[T g.SharedState](from, to g.Node[T], labels ...map[string]string) g.Edge[T] {
	useLabels := make(map[string]string)
	for _, lbls := range labels {
		for k, v := range lbls {
			useLabels[k] = v
		}
	}
	return &edgeImpl[T]{labels: useLabels, from: from, to: to}
}

// ------------------------------------------------------------------------------
// Edge Implementation
// ------------------------------------------------------------------------------

type edgeImpl[T g.SharedState] struct {
	labels map[string]string
	from   g.Node[T]
	to     g.Node[T]
}

func (e *edgeImpl[T]) To() g.Node[T] {
	return e.to
}

func (e *edgeImpl[T]) LabelByKey(key string) (string, bool) {
	val, ok := e.labels[key]
	return val, ok
}
