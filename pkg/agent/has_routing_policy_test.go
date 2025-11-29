package agent_test

import (
	"testing"

	a "github.com/morphy76/ggraph/pkg/agent"
	tool "github.com/morphy76/ggraph/pkg/agent/tool"
	b "github.com/morphy76/ggraph/pkg/builders"
	g "github.com/morphy76/ggraph/pkg/graph"
)

// mockNodeFn is a simple mock function for testing nodes
func mockNodeFn(userInput, currentState a.Conversation, notify g.NotifyPartialFn[a.Conversation]) (a.Conversation, error) {
	return currentState, nil
}

func TestToolProcessorRoutingFn(t *testing.T) {
	t.Run("no_tool_calls_returns_any_edge", func(t *testing.T) {
		// Create test nodes
		node1, err := b.NewNode("node1", mockNodeFn)
		if err != nil {
			t.Fatalf("Failed to create node1: %v", err)
		}

		node2, err := b.NewNode("node2", mockNodeFn)
		if err != nil {
			t.Fatalf("Failed to create node2: %v", err)
		}

		node3, err := b.NewNode("node3", mockNodeFn)
		if err != nil {
			t.Fatalf("Failed to create node3: %v", err)
		}

		// Create test edges without tool executor label
		edge1 := b.CreateEdge(node1, node2)
		edge2 := b.CreateEdge(node1, node3, map[string]string{"type": "other"})

		edges := []g.Edge[a.Conversation]{edge1, edge2}

		// Create conversation state with no tool calls
		conversation := a.Conversation{
			Messages:         []a.Message{a.CreateMessage(a.User, "Hello")},
			CurrentToolCalls: []tool.FnCall{}, // Empty tool calls
		}

		// Call the routing function
		selectedEdge := a.ToolProcessorRoutingFn(conversation, conversation, edges)

		// Should return one of the available edges (AnyRoute behavior)
		if selectedEdge == nil {
			t.Error("Expected an edge to be selected, but got nil")
		}

		// Verify it's one of the provided edges
		found := false
		for _, edge := range edges {
			if edge == selectedEdge {
				found = true
				break
			}
		}
		if !found {
			t.Error("Selected edge is not one of the provided edges")
		}
	})

	t.Run("with_tool_calls_returns_tool_executor_edge", func(t *testing.T) {
		// Create test nodes
		node1, err := b.NewNode("node1", mockNodeFn)
		if err != nil {
			t.Fatalf("Failed to create node1: %v", err)
		}

		node2, err := b.NewNode("node2", mockNodeFn)
		if err != nil {
			t.Fatalf("Failed to create node2: %v", err)
		}

		node3, err := b.NewNode("toolExecutor", mockNodeFn)
		if err != nil {
			t.Fatalf("Failed to create toolExecutor node: %v", err)
		}

		// Create test edges - one with tool executor label, one without
		edge1 := b.CreateEdge(node1, node2, map[string]string{"type": "normal"})
		edge2 := b.CreateEdge(node1, node3, map[string]string{a.RouteTagToolKey: a.RouteTagToolRequest})

		edges := []g.Edge[a.Conversation]{edge1, edge2}

		// Create a mock tool call
		mockTool := createMockTool(t)
		toolCall := tool.FnCall{
			ID:        "call_123",
			ToolName:  mockTool.Name,
			Arguments: map[string]any{"key1": "value1"},
		}

		// Create conversation state with tool calls
		conversation := a.Conversation{
			Messages:         []a.Message{a.CreateMessage(a.User, "Use a tool")},
			CurrentToolCalls: []tool.FnCall{toolCall},
		}

		// Call the routing function
		selectedEdge := a.ToolProcessorRoutingFn(conversation, conversation, edges)

		// Should return the edge with tool_executor label
		if selectedEdge == nil {
			t.Fatal("Expected tool executor edge to be selected, but got nil")
		}

		if selectedEdge != edge2 {
			t.Error("Expected edge2 (tool executor) to be selected")
		}

		// Verify the selected edge has the tool_executor label
		if label, ok := selectedEdge.LabelByKey(a.RouteTagToolKey); !ok || label != a.RouteTagToolRequest {
			t.Errorf("Selected edge should have tool_executor label with value '%s', got: %v, %v", a.RouteTagToolRequest, label, ok)
		}
	})

	t.Run("with_tool_calls_but_no_executor_edge_returns_nil", func(t *testing.T) {
		// Create test nodes
		node1, err := b.NewNode("node1", mockNodeFn)
		if err != nil {
			t.Fatalf("Failed to create node1: %v", err)
		}

		node2, err := b.NewNode("node2", mockNodeFn)
		if err != nil {
			t.Fatalf("Failed to create node2: %v", err)
		}

		// Create test edges WITHOUT tool executor label
		edge1 := b.CreateEdge(node1, node2, map[string]string{"type": "normal"})
		edge2 := b.CreateEdge(node1, node2, map[string]string{"type": "other"})

		edges := []g.Edge[a.Conversation]{edge1, edge2}

		// Create a mock tool call
		mockTool := createMockTool(t)
		toolCall := tool.FnCall{
			ID:        "call_456",
			ToolName:  mockTool.Name,
			Arguments: map[string]any{"key1": "value1"},
		}

		// Create conversation state with tool calls
		conversation := a.Conversation{
			Messages:         []a.Message{a.CreateMessage(a.User, "Use a tool")},
			CurrentToolCalls: []tool.FnCall{toolCall},
		}

		// Call the routing function
		selectedEdge := a.ToolProcessorRoutingFn(conversation, conversation, edges)

		// Should return nil because no edge has the tool_executor label
		if selectedEdge != nil {
			t.Errorf("Expected nil when no tool executor edge exists, but got: %v", selectedEdge)
		}
	})

	t.Run("with_multiple_tool_calls_returns_executor_edge", func(t *testing.T) {
		// Create test nodes
		node1, err := b.NewNode("node1", mockNodeFn)
		if err != nil {
			t.Fatalf("Failed to create node1: %v", err)
		}

		node2, err := b.NewNode("toolExecutor", mockNodeFn)
		if err != nil {
			t.Fatalf("Failed to create toolExecutor node: %v", err)
		}

		// Create test edge with tool executor label
		executorEdge := b.CreateEdge(node1, node2, map[string]string{a.RouteTagToolKey: a.RouteTagToolRequest})

		edges := []g.Edge[a.Conversation]{executorEdge}

		// Create multiple mock tool calls
		mockTool1 := createMockTool(t)
		mockTool2 := createMockTool(t)

		toolCalls := []tool.FnCall{
			{
				ID:        "call_1",
				ToolName:  mockTool1.Name,
				Arguments: map[string]any{"key1": "value1"},
			},
			{
				ID:        "call_2",
				ToolName:  mockTool2.Name,
				Arguments: map[string]any{"key1": "value2"},
			},
		}

		// Create conversation state with multiple tool calls
		conversation := a.Conversation{
			Messages:         []a.Message{a.CreateMessage(a.User, "Use multiple tools")},
			CurrentToolCalls: toolCalls,
		}

		// Call the routing function
		selectedEdge := a.ToolProcessorRoutingFn(conversation, conversation, edges)

		// Should return the executor edge
		if selectedEdge == nil {
			t.Fatal("Expected tool executor edge to be selected, but got nil")
		}

		if selectedEdge != executorEdge {
			t.Error("Expected executorEdge to be selected")
		}
	})

	t.Run("empty_edges_with_no_tool_calls_returns_nil", func(t *testing.T) {
		// Empty edges list
		edges := []g.Edge[a.Conversation]{}

		// Create conversation state with no tool calls
		conversation := a.Conversation{
			Messages:         []a.Message{a.CreateMessage(a.User, "Hello")},
			CurrentToolCalls: []tool.FnCall{},
		}

		// Call the routing function
		selectedEdge := a.ToolProcessorRoutingFn(conversation, conversation, edges)

		// Should return nil (AnyRoute with empty edges returns nil)
		if selectedEdge != nil {
			t.Errorf("Expected nil with empty edges, but got: %v", selectedEdge)
		}
	})
}

// createMockTool creates a mock tool for testing purposes
func createMockTool(t *testing.T) *tool.Tool {
	mockToolFn := func(arg1 string) (string, error) {
		return "result", nil
	}

	mockTool, err := tool.CreateTool[string](mockToolFn, "prompt: Mock tool for testing")
	if err != nil {
		t.Fatalf("Failed to create mock tool: %v", err)
	}

	return mockTool
}
