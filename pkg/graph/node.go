package graph

import (
	"context"
	"fmt"
	"time"
)

// NodeFunc defines a function type that processes a node with the given SharedState type.
type NodeFunc[T SharedState] func(state T) (T, error)

// CreateNode creates a new instance of Node with the specified SharedState type.
func CreateNode[T SharedState](name string, fn NodeFunc[T]) (Node[T], error) {
	if name == "" {
		return nil, fmt.Errorf("node creation failed: name cannot be empty")
	}
	if fn == nil {
		return nil, fmt.Errorf("node creation failed: function cannot be nil")
	}
	return &nodeImpl[T]{
		mailbox: make(chan T),

		name: name,
		fn:   fn,
	}, nil
}

// Node represents a node in the graph.
type Node[T SharedState] interface {
	// Accept processes the node with the given state and returns the updated state.
	Accept(state T, runtime StateObserver[T])
	// Name returns the name of the node.
	Name() string
}

var _ Node[SharedState] = (*nodeImpl[SharedState])(nil)

type nodeImpl[T SharedState] struct {
	mailbox chan T

	name string
	fn   NodeFunc[T]
}

func (n *nodeImpl[T]) Name() string {
	return n.name
}

func (n *nodeImpl[T]) Accept(state T, runtime StateObserver[T]) {
	go func() {
		ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
		defer cancel()

		select {
		case state := <-n.mailbox:
			updatedState, err := n.fn(state)
			if err != nil {
				runtime.NotifyStateChange(n, updatedState, fmt.Errorf("error executing node %s: %w", n.name, err))
				return
			}
			runtime.NotifyStateChange(n, updatedState, nil)
		case <-ctx.Done():
			runtime.NotifyStateChange(n, state, fmt.Errorf("timeout executing node %s: %w", n.name, ctx.Err()))
			return
		}
	}()

	n.mailbox <- state
}
