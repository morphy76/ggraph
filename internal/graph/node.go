package graph

import (
	"context"
	"fmt"
	"time"

	g "github.com/morphy76/ggraph/pkg/graph"
)

// NodeImplFactory creates a new instance of Node with the specified SharedState type.
func NodeImplFactory[T g.SharedState](name string, fn g.NodeFunc[T], routePolicy g.RoutePolicy[T], role g.NodeRole) g.Node[T] {
	useFn := fn
	if useFn == nil {
		useFn = func(userInput T, currentState T, notify func(T)) (T, error) {
			return currentState, nil
		}
	}
	return &nodeImpl[T]{
		mailbox:     make(chan T, 100),
		name:        name,
		fn:          useFn,
		routePolicy: routePolicy,
		role:        role,
	}
}

// ------------------------------------------------------------------------------
// Node Implementation
// ------------------------------------------------------------------------------

var _ g.Node[g.SharedState] = (*nodeImpl[g.SharedState])(nil)

type nodeImpl[T g.SharedState] struct {
	mailbox chan T

	name        string
	fn          g.NodeFunc[T]
	routePolicy g.RoutePolicy[T]

	role g.NodeRole
}

func (n *nodeImpl[T]) Name() string {
	return n.name
}

func (n *nodeImpl[T]) Accept(deltaState T, runtime g.StateObserver[T]) {
	go func() {
		ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
		defer cancel()

		partialStateChange := func(state T) {
			runtime.NotifyStateChange(n, state, nil, true)
		}

		select {
		case asyncDeltaState := <-n.mailbox:
			updatedState, err := n.fn(asyncDeltaState, runtime.CurrentState(), partialStateChange)
			if err != nil {
				runtime.NotifyStateChange(n, updatedState, fmt.Errorf("error executing node %s: %w", n.name, err), false)
				return
			}
			runtime.NotifyStateChange(n, updatedState, nil, false)
		case <-ctx.Done():
			runtime.NotifyStateChange(n, deltaState, fmt.Errorf("timeout executing node %s: %w", n.name, ctx.Err()), false)
			return
		}
	}()

	n.mailbox <- deltaState
}

func (n *nodeImpl[T]) RoutePolicy() g.RoutePolicy[T] {
	return n.routePolicy
}

func (n *nodeImpl[T]) Role() g.NodeRole {
	return n.role
}
