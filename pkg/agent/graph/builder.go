package graph

import (
	t "github.com/morphy76/ggraph/internal/agent/tool"
	a "github.com/morphy76/ggraph/pkg/agent"
	pt "github.com/morphy76/ggraph/pkg/agent/tool"
	g "github.com/morphy76/ggraph/pkg/graph"
)

// CreateToolNode creates a new Node capable of processing tool calls within an agent conversation.
//
// Parameters:
//   - name: The unique name for the tool node.
//   - tools: A variadic list of tools that the node can utilize.
//
// Returns:
//   - An instance of g.Node[a.Conversation] configured for tool processing.
//   - An error if the node creation fails.
//
// Example usage:
//
//	toolNode, err := CreateToolNode("ToolProcessorNode")
func CreateToolNode(name string, tools ...*pt.Tool) (g.Node[a.Conversation], error) {
	return t.NodeToolFactory(name, tools...)
}
