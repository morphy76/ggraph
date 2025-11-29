package graph

// NodeOptions holds the configuration for a node.
type NodeOptions[T SharedState] struct {
	RoutingPolicy RoutePolicy[T]
	Reducer       ReducerFn[T]
}

// NodeOption is a functional option for configuring a node.
type NodeOption[T SharedState] interface {
	// Apply applies the option to the NodeOptions.
	//
	// Parameters:
	//   - r: A pointer to NodeOptions to modify.
	//
	// Returns:
	//   - An error if the application of the option fails, otherwise nil.
	Apply(r *NodeOptions[T]) error
}

// NodeOptionFunc is a function type that implements the NodeOption interface.
type NodeOptionFunc[T SharedState] func(*NodeOptions[T]) error

// Apply applies the NodeOptionFunc to the given NodeOptions.
//
// Parameters:
//   - r: A pointer to NodeOptions to modify.
//
// Returns:
//   - An error if the application of the option fails, otherwise nil.
func (s NodeOptionFunc[T]) Apply(r *NodeOptions[T]) error { return s(r) }

// WithRoutingPolicy sets a custom routing policy for the node.
//
// Parameters:
//   - policy: The RoutePolicy to use for routing decisions.
//
// Returns:
//   - A NodeOption that sets the routing policy.
//
// Example:
//
//	node, err := builders.NewNode("MyNode", myNodeFunction,
//	    builders.WithRoutingPolicy(myRoutingPolicy))
func WithRoutingPolicy[T SharedState](policy RoutePolicy[T]) NodeOption[T] {
	return NodeOptionFunc[T](func(r *NodeOptions[T]) error {
		r.RoutingPolicy = policy
		return nil
	})
}

// WithReducer sets a custom state reducer function for the node.
//
// Parameters:
//   - reducer: The ReducerFn to use for combining state updates.
//
// Returns:
//   - A NodeOption that sets the reducer function.
//
// Example:
//
//	node, err := builders.NewNode("MyNode", myNodeFunction,
//	    builders.WithReducer(myReducerFunction))
func WithReducer[T SharedState](reducer ReducerFn[T]) NodeOption[T] {
	return NodeOptionFunc[T](func(r *NodeOptions[T]) error {
		r.Reducer = reducer
		return nil
	})
}
