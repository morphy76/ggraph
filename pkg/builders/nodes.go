package builders

import (
	"fmt"

	i "github.com/morphy76/ggraph/internal/graph"
	g "github.com/morphy76/ggraph/pkg/graph"
)

// CreateStartNode creates a new instance of StartNode with the specified SharedState type.
func CreateStartNode[T g.SharedState]() g.Node[T] {
	policy, _ := CreateAnyRoutePolicy[T]()
	return i.NodeImplFactory("StartNode", nil, policy, g.StartNode)
}

// CreateEndNode creates a new instance of EndNode with the specified SharedState type.
func CreateEndNode[T g.SharedState]() g.Node[T] {
	return i.NodeImplFactory[T]("EndNode", nil, nil, g.EndNode)
}

// CreateRouter creates a new instance of Node with the specified SharedState type and routing policy.
func CreateRouter[T g.SharedState](name string, policy g.RoutePolicy[T]) (g.Node[T], error) {
	return CreateNodeWithRoutingPolicy(name, nil, policy)
}

// CreateNodeWithRoutingPolicy creates a new instance of Node with the specified SharedState type and routing policy.
func CreateNodeWithRoutingPolicy[T g.SharedState](name string, fn g.NodeFn[T], policy g.RoutePolicy[T]) (g.Node[T], error) {
	if name == "" {
		return nil, fmt.Errorf("node creation failed: name cannot be empty")
	}
	return i.NodeImplFactory(name, fn, policy, g.IntermediateNode), nil
}

// CreateNode creates a new instance of Node with the specified SharedState type.
func CreateNode[T g.SharedState](name string, fn g.NodeFn[T]) (g.Node[T], error) {
	policy, err := CreateAnyRoutePolicy[T]()
	if err != nil {
		return nil, err
	}
	return CreateNodeWithRoutingPolicy(name, fn, policy)
}
