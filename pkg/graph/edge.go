package graph

// CreateEdge creates a new instance of Edge with the specified SharedState type.
func CreateEdge[T SharedState](from, to Node[T], labels ...map[string]string) Edge[T] {
	useLabels := make(map[string]string)
	for _, lbls := range labels {
		for k, v := range lbls {
			useLabels[k] = v
		}
	}
	return &edgeImpl[T]{labels: useLabels, from: from, to: to}
}

// TODO predicates

// Edge represents an edge in the graph.
type Edge[T SharedState] interface {
	// To returns the destination node of the edge.
	To() Node[T]
	// LabelByKey returns the label value associated with the specified key and a boolean indicating if the key exists.
	LabelByKey(key string) (string, bool)
}

type edgeImpl[T SharedState] struct {
	labels map[string]string
	from   Node[T]
	to     Node[T]
}

func (e *edgeImpl[T]) To() Node[T] {
	return e.to
}

func (e *edgeImpl[T]) LabelByKey(key string) (string, bool) {
	val, ok := e.labels[key]
	return val, ok
}
