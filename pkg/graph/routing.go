package graph

import "errors"

var (
	// ErrEdgeSelectionFnNil indicates that the edge selection function is nil.
	ErrEdgeSelectionFnNil = errors.New("edge selection function cannot be nil")
	// ErrNoOutboundEdges indicates that there are no outbound edges from a node.
	ErrNoOutboundEdges = errors.New("no outbound edges from node")
	// ErrNoRoutingPolicy indicates that a node has no routing policy defined.
	ErrNoRoutingPolicy = errors.New("no routing policy defined for node")
	// ErrNilEdge indicates that a routing policy returned a nil edge.
	ErrNilEdge = errors.New("routing policy returned nil edge")
	// ErrNextEdgeNil indicates that the next edge from a node has a nil target node.
	ErrNextEdgeNil = errors.New("next edge from node has nil target node")
)

// RoutePolicy defines the strategy for selecting which edge to follow after node execution.
//
// Routing policies are the mechanism that enables dynamic workflow behavior in graphs.
// They determine how execution flows through the graph by selecting which outgoing edge(s)
// to traverse based on the current state, user input, and available edges.
//
// Different routing policies enable different workflow patterns:
//   - "Any route" policy: Allows the runtime to select any available edge (default)
//   - Conditional policy: Implements if/else branching based on state conditions
//   - Loop policy: Creates cycles by routing back to previous nodes
//   - Custom policy: Implements domain-specific routing logic
//
// Routing policies are created using builder functions:
//   - builders.CreateAnyRoutePolicy() for default behavior
//   - builders.CreateConditionalRoutePolicy() for custom logic
//
// Example conditional policy:
//
//	policy, _ := builders.CreateConditionalRoutePolicy(
//	    func(userInput, state MyState, edges []Edge[MyState]) Edge[MyState] {
//	        if state.Value > 100 {
//	            return edges[0] // Success path
//	        }
//	        return edges[1] // Retry path
//	    },
//	)
type RoutePolicy[T SharedState] interface {
	// SelectEdge determines which outgoing edge to follow after the current node.
	//
	// This method is called by the runtime after a node completes execution. It
	// examines the user input, current state, and available edges to make a routing
	// decision. The selected edge determines which node executes next.
	//
	// The implementation must return one of the edges from the provided edges slice.
	// Returning an edge not in the slice will cause a runtime error.
	//
	// Parameters:
	//   - userInput: The original input provided to Runtime.Invoke(), unchanged
	//     throughout execution.
	//   - currentState: The current state after the node's execution, potentially
	//     modified by the node's processing function.
	//   - edges: All available outgoing edges from the current node. This slice
	//     contains only edges where From() matches the current node.
	//
	// Returns:
	//   - The Edge to traverse next. Must be one of the edges from the provided slice.
	//
	// Example implementation:
	//
	//	func (p *MyPolicy) SelectEdge(userInput, currentState MyState, edges []Edge[MyState]) Edge[MyState] {
	//	    // Conditional routing based on state
	//	    if currentState.Counter >= currentState.MaxIterations {
	//	        // Find and return the exit edge
	//	        for _, edge := range edges {
	//	            if label, ok := edge.LabelByKey("type"); ok && label == "exit" {
	//	                return edge
	//	            }
	//	        }
	//	    }
	//	    // Default: return the first edge
	//	    return edges[0]
	//	}
	SelectEdge(userInput T, currentState T, edges []Edge[T]) Edge[T]
}
