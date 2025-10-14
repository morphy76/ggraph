package graph

import (
	"fmt"
)

// NodeFunc defines a function type that processes a node with the given SharedState type.
type NodeFunc[T SharedState] func(state T) (T, error)

// CreateStartNode creates a new instance of StartNode with the specified SharedState type.
func CreateStartNode[T SharedState]() Node[T] {
	return &StartNode[T]{}
}

// CreateEndNode creates a new instance of EndNode with the specified SharedState type.
func CreateEndNode[T SharedState]() Node[T] {
	return &EndNode[T]{}
}

// CreateNode creates a new instance of Node with the specified SharedState type.
func CreateNode[T SharedState](name string, fn NodeFunc[T]) (Node[T], error) {
	if name == "" {
		return nil, fmt.Errorf("node creation failed: name cannot be empty")
	}
	if fn == nil {
		return nil, fmt.Errorf("node creation failed: function cannot be nil")
	}
	return &nodeImpl[T]{name: name, fn: fn}, nil
}

// Node represents a node in the graph.
type Node[T SharedState] interface {
	// Accept processes the node with the given state and returns the updated state.
	Accept(state T, runtime Runtime[T]) (T, error)
}

var _ Node[SharedState] = (*StartNode[SharedState])(nil)

// StartNode represents the starting node of a graph.
type StartNode[T SharedState] struct {
}

func (n *StartNode[T]) Accept(state T, runtime Runtime[T]) (T, error) {
	outboundEdges := runtime.EdgesFrom(n)
	if len(outboundEdges) == 0 {
		return state, fmt.Errorf("error browsing the graph from start node")
	}

	selectedEdge := outboundEdges[0]
	return selectedEdge.To().Accept(state, runtime)
}

var _ Node[SharedState] = (*EndNode[SharedState])(nil)

// EndNode represents the ending node of a graph.
type EndNode[T SharedState] struct {
}

func (n *EndNode[T]) Accept(state T, runtime Runtime[T]) (T, error) {
	return state, nil
}

var _ Node[SharedState] = (*nodeImpl[SharedState])(nil)

type nodeImpl[T SharedState] struct {
	name string
	fn   NodeFunc[T]
}

func (n *nodeImpl[T]) Accept(state T, runtime Runtime[T]) (T, error) {
	nextState, err := n.fn(state)
	if err != nil {
		return nextState, fmt.Errorf("error processing node '%s': %w", n.name, err)
	}

	outboundEdges := runtime.EdgesFrom(n)
	if len(outboundEdges) == 0 {
		return nextState, fmt.Errorf("error browsing the graph from node '%s'", n.name)
	}

	return outboundEdges[0].To().Accept(nextState, runtime)
}
