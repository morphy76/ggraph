package builders

import (
	"fmt"

	i "github.com/morphy76/ggraph/internal/graph"
	g "github.com/morphy76/ggraph/pkg/graph"
)

// CreateNodeWithRoutingPolicy creates a new operational node with a custom routing policy.
//
// This function constructs a node that can execute a function and determine the next
// node(s) to execute based on the provided routing policy. The routing policy controls
// how the node selects which outgoing edge(s) to follow after execution. This is useful
// for implementing conditional branching, parallel execution, or custom routing logic.
//
// Type Parameters:
//   - T: The SharedState type that will be passed through the graph execution.
//
// Parameters:
//   - name: A unique identifier for the node. Must be non-empty.
//   - fn: The function to execute when the node is reached. Can be nil for router-only nodes.
//   - policy: The routing policy that determines which edge(s) to follow after node execution.
//
// Returns:
//   - A new Node instance with the specified configuration.
//   - An error if the name is empty.
//
// Example:
//
//	policy, _ := CreateConditionalRoutePolicy(func(state MyState, edges []Edge[MyState]) []Edge[MyState] {
//	    if state.Value > 10 {
//	        return []Edge[MyState]{edges[0]} // Take first edge
//	    }
//	    return []Edge[MyState]{edges[1]} // Take second edge
//	})
//	node, err := CreateNodeWithRoutingPolicy[MyState]("decisionNode", myFunction, policy)
func CreateNodeWithRoutingPolicy[T g.SharedState](name string, fn g.NodeFn[T], policy g.RoutePolicy[T]) (g.Node[T], error) {
	if name == "" {
		return nil, fmt.Errorf("node creation failed: name cannot be empty")
	}
	return i.NodeImplFactory(name, fn, policy, g.IntermediateNode), nil
}

// CreateNode creates a new operational node with default routing behavior.
//
// This is the primary function for creating standard nodes in a graph workflow.
// The node will execute the provided function when reached during graph execution,
// and will use the default "any route" policy, which allows the graph runtime to
// select any available outgoing edge. This is suitable for simple linear workflows
// or when custom routing logic is not required.
//
// Type Parameters:
//   - T: The SharedState type that will be passed through the graph execution.
//
// Parameters:
//   - name: A unique identifier for the node. Must be non-empty.
//   - fn: The function to execute when the node is reached. This function receives
//     the current state and should return the updated state along with any error.
//
// Returns:
//   - A new Node instance with default routing policy.
//   - An error if the name is empty or if the default routing policy cannot be created.
//
// Example:
//
//	node, err := CreateNode[MyState]("processData", func(state MyState) (MyState, error) {
//	    state.Value += 1
//	    state.Message = "Processed"
//	    return state, nil
//	})
//	if err != nil {
//	    log.Fatal(err)
//	}
func CreateNode[T g.SharedState](name string, fn g.NodeFn[T]) (g.Node[T], error) {
	policy, err := CreateAnyRoutePolicy[T]()
	if err != nil {
		return nil, err
	}
	return CreateNodeWithRoutingPolicy(name, fn, policy)
}

func createStartNode[T g.SharedState]() g.Node[T] {
	policy, _ := CreateAnyRoutePolicy[T]()
	return i.NodeImplFactory("StartNode", nil, policy, g.StartNode)
}

func createEndNode[T g.SharedState]() g.Node[T] {
	return i.NodeImplFactory[T]("EndNode", nil, nil, g.EndNode)
}
