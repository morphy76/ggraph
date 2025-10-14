package graph

import (
	"context"
	"fmt"
	"time"
)

// NodeFunc defines a function type that processes a node with the given SharedState type.
type NodeFunc[T SharedState] func(state T, notify func(T)) (T, error)

// CreateRouter creates a new instance of Node with the specified SharedState type and routing policy.
func CreateRouter[T SharedState](name string, policy RoutePolicy[T]) (Node[T], error) {
	passthrough := func(state T, notify func(T)) (T, error) { return state, nil }
	return CreateNodeWithRoutingPolicy(name, passthrough, policy)
}

// CreateNodeWithRoutingPolicy creates a new instance of Node with the specified SharedState type and routing policy.
func CreateNodeWithRoutingPolicy[T SharedState](name string, fn NodeFunc[T], policy RoutePolicy[T]) (Node[T], error) {
	if name == "" {
		return nil, fmt.Errorf("node creation failed: name cannot be empty")
	}
	if fn == nil {
		return nil, fmt.Errorf("node creation failed: function cannot be nil")
	}
	if policy == nil {
		return nil, fmt.Errorf("node creation failed: route policy cannot be nil")
	}
	return &nodeImpl[T]{
		mailbox: make(chan T),

		name: name,
		fn:   fn,

		routePolicy: policy,
	}, nil
}

// CreateNode creates a new instance of Node with the specified SharedState type.
func CreateNode[T SharedState](name string, fn NodeFunc[T]) (Node[T], error) {
	policy, err := CreateAnyRoutePolicy[T]()
	if err != nil {
		return nil, err
	}
	return CreateNodeWithRoutingPolicy(name, fn, policy)
}

// Node represents a node in the graph.
type Node[T SharedState] interface {
	// Accept processes the node with the given state and returns the updated state.
	Accept(state T, runtime StateObserver[T])
	// Name returns the name of the node.
	Name() string
	// RoutePolicy returns the routing policy associated with the node.
	RoutePolicy() RoutePolicy[T]
}

var _ Node[SharedState] = (*nodeImpl[SharedState])(nil)

type nodeImpl[T SharedState] struct {
	mailbox chan T

	name string
	fn   NodeFunc[T]

	routePolicy RoutePolicy[T]
}

func (n *nodeImpl[T]) Name() string {
	return n.name
}

func (n *nodeImpl[T]) Accept(state T, runtime StateObserver[T]) {
	go func() {
		ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
		defer cancel()

		partialStateChange := func(state T) {
			runtime.NotifyStateChange(n, state, nil, true)
		}

		select {
		case state := <-n.mailbox:
			updatedState, err := n.fn(state, partialStateChange)
			if err != nil {
				runtime.NotifyStateChange(n, updatedState, fmt.Errorf("error executing node %s: %w", n.name, err), false)
				return
			}
			runtime.NotifyStateChange(n, updatedState, nil, false)
		case <-ctx.Done():
			runtime.NotifyStateChange(n, state, fmt.Errorf("timeout executing node %s: %w", n.name, ctx.Err()), false)
			return
		}
	}()

	n.mailbox <- state
}

func (n *nodeImpl[T]) RoutePolicy() RoutePolicy[T] {
	return n.routePolicy
}
