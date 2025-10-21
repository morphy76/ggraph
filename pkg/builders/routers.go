package builders

import (
	i "github.com/morphy76/ggraph/internal/graph"
	g "github.com/morphy76/ggraph/pkg/graph"
)

// CreateRouter creates a new router node that directs graph execution flow.
//
// A router node is a special type of node without an execution function that solely
// determines which path(s) the graph execution should take based on the provided
// routing policy. Routers are essential for implementing branching logic, conditional
// workflows, and dynamic path selection in complex graph structures.
//
// Type Parameters:
//   - T: The SharedState type that will be passed through the graph execution.
//
// Parameters:
//   - name: A unique identifier for the router node. Must be non-empty.
//   - policy: The routing policy that determines which outgoing edge(s) to follow.
//     This policy is invoked each time the router is reached during execution.
//
// Returns:
//   - A new router Node instance that uses the specified routing policy.
//   - An error if the name is empty.
//
// Example:
//
//	policy, _ := CreateConditionalRoutePolicy(func(userInput, state MyState, edges []Edge[MyState]) Edge[MyState] {
//	    if state.Score > 100 {
//	        return edges[0] // High score path
//	    }
//	    return edges[1] // Normal path
//	})
//	router, err := CreateRouter[MyState]("scoreRouter", policy)
func CreateRouter[T g.SharedState](name string, policy g.RoutePolicy[T]) (g.Node[T], error) {
	return CreateNodeWithRoutingPolicy(name, nil, policy)
}

// CreateAnyRoutePolicy creates a default routing policy that allows any available edge.
//
// This policy does not impose any restrictions on edge selection and allows the graph
// runtime to choose any available outgoing edge. It is the default policy used by
// CreateNode and is suitable for simple linear workflows where branching logic is not
// required. When multiple edges are available, the runtime will typically select the
// first available edge.
//
// Type Parameters:
//   - T: The SharedState type that will be passed through the graph execution.
//
// Returns:
//   - A new RoutePolicy instance that accepts any available edge.
//   - An error if the policy cannot be created (typically never fails).
//
// Example:
//
//	policy, err := CreateAnyRoutePolicy[MyState]()
//	node, _ := CreateNodeWithRoutingPolicy[MyState]("simpleNode", myFunction, policy)
func CreateAnyRoutePolicy[T g.SharedState]() (g.RoutePolicy[T], error) {
	return CreateConditionalRoutePolicy(i.AnyRoute[T])
}

// CreateConditionalRoutePolicy creates a custom routing policy with conditional logic.
//
// This function allows you to define sophisticated routing behavior by providing a
// selection function that inspects the current state and available edges to determine
// which edge should be followed. This is the foundation for implementing conditional
// branching, loops, state-based routing, and other dynamic workflow patterns.
//
// The selection function receives three parameters:
//   - userInput: The original input state provided to the graph execution
//   - currentState: The current state at the time of routing decision
//   - edges: All available outgoing edges from the current node
//
// Type Parameters:
//   - T: The SharedState type that will be passed through the graph execution.
//
// Parameters:
//   - selectionFn: A function that examines the state and edges to select which edge
//     to follow. This function is called each time the routing decision is made.
//
// Returns:
//   - A new RoutePolicy instance that uses the provided selection logic.
//   - An error if the policy cannot be created.
//
// Example:
//
//	policy, err := CreateConditionalRoutePolicy(func(userInput, currentState GameState, edges []Edge[GameState]) Edge[GameState] {
//	    switch {
//	    case currentState.Lives <= 0:
//	        return edges[0] // Game over edge
//	    case currentState.Level >= 10:
//	        return edges[1] // Victory edge
//	    default:
//	        return edges[2] // Continue playing edge
//	    }
//	})
func CreateConditionalRoutePolicy[T g.SharedState](selectionFn g.EdgeSelectionFn[T]) (g.RoutePolicy[T], error) {
	return i.RouterPolicyImplFactory(selectionFn)
}
