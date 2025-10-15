package builders

import (
	"fmt"

	i "github.com/morphy76/ggraph/internal/graph"
	g "github.com/morphy76/ggraph/pkg/graph"
)

// CreateStartNode creates a new instance of StartNode with the specified SharedState type.
func CreateStartNode[T g.SharedState]() g.Node[T] {
	policy, _ := CreateAnyRoutePolicy[T]()
	return i.NodeImplFactory("StartNode", nil, policy)
}

// CreateEndNode creates a new instance of EndNode with the specified SharedState type.
func CreateEndNode[T g.SharedState]() g.Node[T] {
	return i.NodeImplFactory[T]("EndNode", nil, nil)
}

// CreateRouter creates a new instance of Node with the specified SharedState type and routing policy.
func CreateRouter[T g.SharedState](name string, policy g.RoutePolicy[T]) (g.Node[T], error) {
	passthrough := func(state T, notify func(T)) (T, error) { return state, nil }
	return CreateNodeWithRoutingPolicy(name, passthrough, policy)
}

// CreateNodeWithRoutingPolicy creates a new instance of Node with the specified SharedState type and routing policy.
func CreateNodeWithRoutingPolicy[T g.SharedState](name string, fn g.NodeFunc[T], policy g.RoutePolicy[T]) (g.Node[T], error) {
	if name == "" {
		return nil, fmt.Errorf("node creation failed: name cannot be empty")
	}
	if fn == nil {
		return nil, fmt.Errorf("node creation failed: function cannot be nil")
	}
	if policy == nil {
		return nil, fmt.Errorf("node creation failed: route policy cannot be nil")
	}
	return i.NodeImplFactory(name, fn, policy), nil
}

// CreateNode creates a new instance of Node with the specified SharedState type.
func CreateNode[T g.SharedState](name string, fn g.NodeFunc[T]) (g.Node[T], error) {
	policy, err := CreateAnyRoutePolicy[T]()
	if err != nil {
		return nil, err
	}
	return CreateNodeWithRoutingPolicy(name, fn, policy)
}
