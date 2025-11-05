package agent

import (
	i "github.com/morphy76/ggraph/internal/graph"
	g "github.com/morphy76/ggraph/pkg/graph"
)

const (
	// RouteTagToolExecutor is the label key used to identify edges meant for tool execution.
	RouteTagToolExecutor = "tool_executor"
)

// ToolProcessorRoutingFn is a routing function that directs the graph execution
// based on whether the current conversation state includes tool calls.
//
// If there are no tool calls in the current state, it allows any available edge.
// If there are tool calls, it specifically routes to the edge labeled for tool execution.
//
// Parameters:
//   - userInput: The input provided by the user.
//   - currentState: The current state of the conversation.
//   - edges: The available edges to choose from.
//
// Returns:
//   - The selected edge based on the routing logic.
func ToolProcessorRoutingFn(userInput, currentState Conversation, edges []g.Edge[Conversation]) g.Edge[Conversation] {
	if len(currentState.ToolCalls) == 0 {
		return i.AnyRoute(userInput, currentState, edges)
	}
	for _, edge := range edges {
		if _, ok := edge.LabelByKey(RouteTagToolExecutor); ok {
			return edge
		}
	}
	return nil
}
