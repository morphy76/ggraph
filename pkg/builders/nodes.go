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

// NodeBuilder is a builder for creating nodes with customizable properties.
type NodeBuilder[T g.SharedState] struct {
	name          string
	fn            g.NodeFn[T]
	routingPolicy g.RoutePolicy[T]
	reducer       g.ReducerFn[T]
}

// NewNodeBuilder creates a new NodeBuilder with the specified name and processing function.
//
// It initializes the NodeBuilder with default values for routing policy and reducer.
// Parameters:
//   - name: The unique name for the node.
//   - fn: The processing function (NodeFn) for the node.
//
// Returns:
//   - A new instance of NodeBuilder[T].
//
// Example:
//
//	builder := builders.NewNodeBuilder("MyNode", myNodeFunction)
func NewNodeBuilder[T g.SharedState](name string, fn g.NodeFn[T]) NodeBuilder[T] {
	return NodeBuilder[T]{
		name: name,
		fn:   fn,
	}
}

// WithRoutingPolicy sets a custom routing policy for the node.
//
// Parameters:
//   - policy: The RoutePolicy to use for routing decisions.
//
// Returns:
//   - The updated NodeBuilder[T] with the specified routing policy.
//
// Example:
//
//	builder := builders.NewNodeBuilder("MyNode", myNodeFunction).
//	              WithRoutingPolicy(myRoutingPolicy)
func (b NodeBuilder[T]) WithRoutingPolicy(policy g.RoutePolicy[T]) NodeBuilder[T] {
	b.routingPolicy = policy
	return b
}

// WithReducer sets a custom state reducer function for the node.
//
// Parameters:
//   - reducer: The ReducerFn to use for combining state updates.
//
// Returns:
//   - The updated NodeBuilder[T] with the specified reducer function.
//
// Example:
//
//	builder := builders.NewNodeBuilder("MyNode", myNodeFunction).
//	              WithReducer(myReducerFunction)
func (b NodeBuilder[T]) WithReducer(reducer g.ReducerFn[T]) NodeBuilder[T] {
	b.reducer = reducer
	return b
}

// Build constructs the Node[T] instance based on the configured properties.
//
// It uses default values for any properties that were not explicitly set.
//
// Returns:
//   - The constructed Node[T] instance.
//   - An error if the node could not be created.
//
// Example:
//
//	node := builders.NewNodeBuilder("MyNode", myNodeFunction).
//	            WithRoutingPolicy(myRoutingPolicy).
//	            WithReducer(myReducerFunction).
//	            Build()
func (b NodeBuilder[T]) Build() (g.Node[T], error) {
	if b.name == ReservedNodeNameStart || b.name == ReservedNodeNameEnd {
		return nil, fmt.Errorf("node creation error for name %s: %w", b.name, g.ErrReservedNodeName)
	}
	policy := b.routingPolicy
	if policy == nil {
		var err error
		policy, err = CreateAnyRoutePolicy[T]()
		if err != nil {
			return nil, err
		}
	}
	reducer := b.reducer
	if reducer == nil {
		reducer = i.Replacer[T]
	}
	return i.NodeImplFactory(g.IntermediateNode, b.name, b.fn, policy, reducer), nil
}

func createStartNode[T g.SharedState]() (g.Node[T], error) {
	policy, _ := CreateAnyRoutePolicy[T]()
	return i.NodeImplFactory(g.StartNode, ReservedNodeNameStart, nil, policy, i.Replacer[T]), nil
}

func createEndNode[T g.SharedState]() (g.Node[T], error) {
	return i.NodeImplFactory[T](g.EndNode, ReservedNodeNameEnd, nil, nil, i.Replacer[T]), nil
}
