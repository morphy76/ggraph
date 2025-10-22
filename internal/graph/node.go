package graph

import (
	"context"
	"fmt"
	"time"

	g "github.com/morphy76/ggraph/pkg/graph"
)

// NodeImplFactory creates a new instance of Node with the specified SharedState type.
func NodeImplFactory[T g.SharedState](role g.NodeRole, name string, fn g.NodeFn[T], routePolicy g.RoutePolicy[T], reducer g.ReducerFn[T]) g.Node[T] {
	useFn := fn
	if useFn == nil {
		useFn = func(userInput T, currentState T, notifyPartial g.NotifyPartialFn[T]) (T, error) {
			return currentState, nil
		}
	}
	usePolicy := routePolicy
	if usePolicy == nil {
		usePolicy, _ = RouterPolicyImplFactory[T](AnyRoute)
	}
	return &nodeImpl[T]{
		mailbox:     make(chan T, 100),
		name:        name,
		fn:          useFn,
		routePolicy: usePolicy,
		role:        role,
		reducer:     reducer,
	}
}

// ------------------------------------------------------------------------------
// Node Implementation
// ------------------------------------------------------------------------------

var _ g.Node[g.SharedState] = (*nodeImpl[g.SharedState])(nil)

type nodeImpl[T g.SharedState] struct {
	mailbox chan T

	name        string
	fn          g.NodeFn[T]
	routePolicy g.RoutePolicy[T]

	role g.NodeRole

	reducer g.ReducerFn[T]
}

func (n *nodeImpl[T]) Name() string {
	return n.name
}

func (n *nodeImpl[T]) Accept(userInput T, runtime g.StateObserver[T]) {
	go func() {
		ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
		defer cancel()

		partialStateChange := func(state T) {
			runtime.NotifyStateChange(n, userInput, state, n.reducer, nil, true)
		}

		select {
		case asyncDeltaState := <-n.mailbox:
			stateChange, err := n.fn(asyncDeltaState, runtime.CurrentState(), partialStateChange)
			if err != nil {
				runtime.NotifyStateChange(n, userInput, stateChange, n.reducer, fmt.Errorf("error executing node %s: %w", n.name, err), false)
				return
			}
			runtime.NotifyStateChange(n, userInput, stateChange, n.reducer, nil, false)
		case <-ctx.Done():
			runtime.NotifyStateChange(n, userInput, runtime.CurrentState(), n.reducer, fmt.Errorf("timeout executing node %s: %w", n.name, ctx.Err()), false)
			return
		}
	}()

	n.mailbox <- userInput
}

func (n *nodeImpl[T]) RoutePolicy() g.RoutePolicy[T] {
	return n.routePolicy
}

func (n *nodeImpl[T]) Role() g.NodeRole {
	return n.role
}
