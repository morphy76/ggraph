package graph

import (
	"context"
	"fmt"

	g "github.com/morphy76/ggraph/pkg/graph"
)

// NodeImplFactory creates a new instance of Node with the specified SharedState type.
func NodeImplFactory[T g.SharedState](role g.NodeRole, name string, fn g.NodeFn[T], opt *g.NodeOptions[T]) (g.Node[T], error) {
	if name == "" {
		return nil, fmt.Errorf("node creation failed: %w", g.ErrNodeNameEmpty)
	}
	if name == "StartNode" || name == "EndNode" {
		if role != g.StartNode && role != g.EndNode {
			return nil, fmt.Errorf("node creation failed: %w", g.ErrReservedNodeName)
		}
	}
	if opt == nil {
		return nil, fmt.Errorf("node creation failed: %w", g.ErrNodeOptionsNil)
	}
	if role < g.StartNode || role > g.EndNode {
		return nil, fmt.Errorf("node creation failed: %w", g.ErrInvalidNodeRole)
	}

	opt.NodeSettings = g.FillNodeSettingsWithDefaults(opt.NodeSettings)

	useFn := fn
	if useFn == nil {
		useFn = func(userInput T, currentState T, notifyPartial g.NotifyPartialFn[T]) (T, error) {
			return currentState, nil
		}
	}
	usePolicy := opt.RoutingPolicy
	if usePolicy == nil {
		usePolicy, _ = RouterPolicyImplFactory[T](AnyRoute)
	}
	return &nodeImpl[T]{
		mailbox:     make(chan T, opt.NodeSettings.MailboxSize),
		name:        name,
		fn:          useFn,
		routePolicy: usePolicy,
		role:        role,
		reducer:     opt.Reducer,
		settings:    opt.NodeSettings,
	}, nil
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

	settings g.NodeSettings
}

func (n *nodeImpl[T]) Name() string {
	return n.name
}

func (n *nodeImpl[T]) Accept(
	userInput T,
	stateObserver g.StateObserver[T],
	nodeExecutor g.NodeExecutor,
	config g.InvokeConfig,
) {
	useThreadID := config.ThreadID

	task := func() {
		ctx, cancel := context.WithTimeout(context.Background(), n.settings.AcceptTimeout)
		defer cancel()

		partialStateChange := func(state T) {
			stateObserver.NotifyStateChange(n, config, userInput, state, n.reducer, nil, true)
		}

		select {
		case asyncDeltaState := <-n.mailbox:
			stateChange, err := n.fn(asyncDeltaState, stateObserver.CurrentState(useThreadID), partialStateChange)
			if err != nil {
				stateObserver.NotifyStateChange(n, config, userInput, stateChange, n.reducer, fmt.Errorf("error executing node %s: %w", n.name, err), false)
				return
			}
			stateObserver.NotifyStateChange(n, config, userInput, stateChange, n.reducer, nil, false)
		case <-ctx.Done():
			stateObserver.NotifyStateChange(n, config, userInput, stateObserver.CurrentState(useThreadID), n.reducer, fmt.Errorf("error executing node %s: %w", n.name, ctx.Err()), false)
			return
		}
	}

	nodeExecutor.Submit(task)

	n.mailbox <- userInput
}

func (n *nodeImpl[T]) RoutePolicy() g.RoutePolicy[T] {
	return n.routePolicy
}

func (n *nodeImpl[T]) Role() g.NodeRole {
	return n.role
}
