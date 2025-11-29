package builders

import (
	"fmt"

	i "github.com/morphy76/ggraph/internal/graph"
	g "github.com/morphy76/ggraph/pkg/graph"
)

const (
	// ReservedNodeNameStart is the reserved name for the start node.
	ReservedNodeNameStart string = "StartNode"
	// ReservedNodeNameEnd is the reserved name for the end node.
	ReservedNodeNameEnd string = "EndNode"
)

// NewNode creates a new node with the specified name, processing function, and options.
//
// It initializes the node with default values for any options not provided.
// Parameters:
//   - name: The unique name for the node.
//   - fn: The processing function (NodeFn) for the node.
//   - opts: Optional configuration options for the node.
//
// Returns:
//   - The constructed Node[T] instance.
//   - An error if the node could not be created.
//
// Example:
//
//	node, err := builders.NewNode("MyNode", myNodeFunction,
//	    builders.WithRoutingPolicy(myRoutingPolicy),
//	    builders.WithReducer(myReducerFunction))
func NewNode[T g.SharedState](name string, fn g.NodeFn[T], opts ...g.NodeOption[T]) (g.Node[T], error) {
	if name == ReservedNodeNameStart || name == ReservedNodeNameEnd {
		return nil, fmt.Errorf("node creation error for name %s: %w", name, g.ErrReservedNodeName)
	}

	useOpts := &g.NodeOptions[T]{
		Reducer: i.Replacer[T],
	}
	for _, opt := range opts {
		opt.Apply(useOpts)
	}

	if useOpts.RoutingPolicy == nil {
		var err error
		useOpts.RoutingPolicy, err = CreateAnyRoutePolicy[T]()
		if err != nil {
			return nil, err
		}
	}

	return i.NodeImplFactory(g.IntermediateNode, name, fn, useOpts), nil
}

func createStartNode[T g.SharedState]() (g.Node[T], error) {
	policy, _ := CreateAnyRoutePolicy[T]()
	useOpts := &g.NodeOptions[T]{
		RoutingPolicy: policy,
		Reducer:       i.Replacer[T],
	}
	return i.NodeImplFactory(g.StartNode, ReservedNodeNameStart, nil, useOpts), nil
}

func createEndNode[T g.SharedState]() (g.Node[T], error) {
	useOpts := &g.NodeOptions[T]{
		Reducer: i.Replacer[T],
	}
	return i.NodeImplFactory(g.EndNode, ReservedNodeNameEnd, nil, useOpts), nil
}
