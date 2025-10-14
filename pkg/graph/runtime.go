package graph

import "fmt"

// CreateRuntime creates a new instance of Runtime with the specified SharedState type.
func CreateRuntime[T SharedState](startEdge *StartEdge[T]) Runtime[T] {
	return &runtimeImpl[T]{
		startEdge: *startEdge,
		edges:     []Edge[T]{},
	}
}

// Connected provides access to the connected graph components.
type Connected[T SharedState] interface {
	// AddEdge adds an edge to the runtime's graph.
	AddEdge(edge Edge[T]) error
	// Validate checks the integrity of the graph structure.
	Validate() error
	// EdgesFrom returns all edges originating from the given node.
	EdgesFrom(node Node[T]) []Edge[T]
}

// Runtime represents the runtime environment for graph processing.
type Runtime[T SharedState] interface {
	Connected[T]
	// Invoke executes the graph processing with the given entry state.
	Invoke(entryState T) (T, error)
}

var _ Runtime[SharedState] = (*runtimeImpl[SharedState])(nil)

type runtimeImpl[T SharedState] struct {
	startEdge StartEdge[T]
	edges     []Edge[T]
}

func (r *runtimeImpl[T]) Invoke(entryState T) (T, error) {
	rv, err := r.startEdge.from.Accept(entryState, r)
	if err != nil {
		return rv, fmt.Errorf("error invoking start edge: %w", err)
	}

	return rv, nil
}

func (r *runtimeImpl[T]) AddEdge(edge Edge[T]) error {
	r.edges = append(r.edges, edge)
	return nil
}

func (r *runtimeImpl[T]) Validate() error {
	var zeroStartEdge StartEdge[T]
	if r.startEdge == zeroStartEdge {
		return fmt.Errorf("graph validation failed: start edge is nil")
	}
	if r.startEdge.from == nil {
		return fmt.Errorf("graph validation failed: start edge 'from' node is nil")
	}

	// Check if there's at least one path from start to an end edge
	visited := make(map[string]bool)
	// Include the start edge in the traversal by starting from its target node
	hasPathToEnd := r.hasPathToEndEdge(r.startEdge.to, visited)
	if !hasPathToEnd {
		return fmt.Errorf("graph validation failed: no path from start edge to any end edge")
	}

	return nil
}

func (r *runtimeImpl[T]) EdgesFrom(node Node[T]) []Edge[T] {
	if r.startEdge.from == node {
		return []Edge[T]{&r.startEdge}
	}
	var outboundEdges []Edge[T]
	for _, edge := range r.edges {
		if edgeFrom(edge) == node {
			outboundEdges = append(outboundEdges, edge)
		}
	}
	return outboundEdges
}

func (r *runtimeImpl[T]) hasPathToEndEdge(node Node[T], visited map[string]bool) bool {
	// Check if the node is an EndNode
	if _, ok := node.(*EndNode[T]); ok {
		return true
	}

	// Mark the node as visited
	nodeKey := fmt.Sprintf("%p", node)
	if visited[nodeKey] {
		return false
	}
	visited[nodeKey] = true

	// Check if any EndEdge starts from this node
	for _, edge := range r.edges {
		if endEdge, ok := edge.(*EndEdge[T]); ok {
			if edgeFrom[T](endEdge) == node {
				return true
			}
		}
	}

	// Explore all edges to find connected nodes
	for _, edge := range r.edges {
		if edgeFrom[T](edge) == node {
			if r.hasPathToEndEdge(edgeTo[T](edge), visited) {
				return true
			}
		}
	}

	return false
}

func edgeFrom[T SharedState](edge Edge[T]) Node[T] {
	switch e := edge.(type) {
	case *edgeImpl[T]:
		return e.from
	case *StartEdge[T]:
		return e.from
	case *EndEdge[T]:
		return e.from
	default:
		return nil
	}
}

func edgeTo[T SharedState](edge Edge[T]) Node[T] {
	switch e := edge.(type) {
	case *edgeImpl[T]:
		return e.to
	case *StartEdge[T]:
		return e.to
	case *EndEdge[T]:
		return e.to
	default:
		return nil
	}
}
