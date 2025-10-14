package graph_test

import (
	"testing"

	"github.com/morphy76/ggraph/pkg/graph"
)

func TestCreateRuntime(t *testing.T) {
	endNode := graph.CreateEndNode[TestSharedState]()
	startEdge := graph.CreateStartEdge(endNode)

	runtime := graph.CreateRuntime(startEdge)

	if runtime == nil {
		t.Fatal("CreateRuntime returned nil")
	}
}

func TestRuntimeAddEdge(t *testing.T) {
	endNode := graph.CreateEndNode[TestSharedState]()
	startEdge := graph.CreateStartEdge(endNode)
	runtime := graph.CreateRuntime(startEdge)

	// Create an edge to add
	fromNode := graph.CreateStartNode[TestSharedState]()
	toNode := graph.CreateEndNode[TestSharedState]()
	edge := graph.CreateEdge(fromNode, toNode)

	err := runtime.AddEdge(edge)
	if err != nil {
		t.Errorf("AddEdge returned error: %v", err)
	}

	// Verify the edge was added by checking EdgesFrom
	edges := runtime.EdgesFrom(fromNode)
	if len(edges) != 1 {
		t.Errorf("Expected 1 edge from node, got %d", len(edges))
	}

	// Check that the edge has the correct target node
	if edges[0].To() != toNode {
		t.Error("AddEdge did not add an edge with the correct target node")
	}
}

func TestRuntimeEdgesFrom(t *testing.T) {
	endNode := graph.CreateEndNode[TestSharedState]()
	startNode := graph.CreateStartNode[TestSharedState]()
	startEdge := graph.CreateStartEdge(endNode)
	runtime := graph.CreateRuntime(startEdge)

	// Test EdgesFrom start node (should return the start edge)
	startEdges := runtime.EdgesFrom(startNode)
	if len(startEdges) != 1 {
		t.Errorf("Expected 1 edge from start node, got %d", len(startEdges))
	}

	// Test EdgesFrom a regular node
	node1, _ := graph.CreateNode("node1", func(state TestSharedState) (TestSharedState, error) {
		return state, nil
	})
	node2, _ := graph.CreateNode("node2", func(state TestSharedState) (TestSharedState, error) {
		return state, nil
	})

	edge1 := graph.CreateEdge(node1, node2)
	edge2 := graph.CreateEdge(node1, endNode)

	runtime.AddEdge(edge1)
	runtime.AddEdge(edge2)

	edges := runtime.EdgesFrom(node1)
	if len(edges) != 2 {
		t.Errorf("Expected 2 edges from node1, got %d", len(edges))
	}

	// Test EdgesFrom a node with no outbound edges
	node3, _ := graph.CreateNode("node3", func(state TestSharedState) (TestSharedState, error) {
		return state, nil
	})
	emptyEdges := runtime.EdgesFrom(node3)
	if len(emptyEdges) != 0 {
		t.Errorf("Expected 0 edges from node3, got %d", len(emptyEdges))
	}
}

func TestRuntimeInvoke(t *testing.T) {
	endNode := graph.CreateEndNode[TestSharedState]()
	startEdge := graph.CreateStartEdge(endNode)
	runtime := graph.CreateRuntime(startEdge)

	state := TestSharedState{Value: "initial"}

	result, err := runtime.Invoke(state)
	if err != nil {
		t.Errorf("Invoke returned error: %v", err)
	}

	// Since we go directly from start to end, state should be unchanged
	if result.Value != "initial" {
		t.Errorf("Expected result value 'initial', got '%s'", result.Value)
	}
}

func TestRuntimeInvokeWithProcessingNodes(t *testing.T) {
	// Create nodes that modify state
	node1, _ := graph.CreateNode("node1", func(state TestSharedState) (TestSharedState, error) {
		state.Value = "processed_by_node1"
		return state, nil
	})

	startEdge := graph.CreateStartEdge(node1)
	runtime := graph.CreateRuntime(startEdge)

	// Add edge from node1 to end
	edge1 := graph.CreateEndEdge(node1)
	runtime.AddEdge(edge1)

	state := TestSharedState{Value: "initial"}

	result, err := runtime.Invoke(state)
	if err != nil {
		t.Errorf("Invoke returned error: %v", err)
	}

	if result.Value != "processed_by_node1" {
		t.Errorf("Expected result value 'processed_by_node1', got '%s'", result.Value)
	}
}

func TestRuntimeValidateValid(t *testing.T) {
	// Create a valid graph: start -> node1 -> end
	node1, _ := graph.CreateNode("node1", func(state TestSharedState) (TestSharedState, error) {
		return state, nil
	})

	startEdge := graph.CreateStartEdge(node1)
	runtime := graph.CreateRuntime(startEdge)

	endEdge := graph.CreateEndEdge(node1)
	runtime.AddEdge(endEdge)

	err := runtime.Validate()
	if err != nil {
		t.Errorf("Validate returned error for valid graph: %v", err)
	}
}

func TestRuntimeValidateNoPathToEnd(t *testing.T) {
	// Create an invalid graph: start -> node1 (no end)
	node1, _ := graph.CreateNode("node1", func(state TestSharedState) (TestSharedState, error) {
		return state, nil
	})

	startEdge := graph.CreateStartEdge(node1)
	runtime := graph.CreateRuntime(startEdge)

	// Don't add any end edge - this should make validation fail
	err := runtime.Validate()
	if err == nil {
		t.Error("Validate should return error for graph with no path to end")
	}

	expectedError := "graph validation failed: no path from start edge to any end edge"
	if err.Error() != expectedError {
		t.Errorf("Expected error message '%s', got '%s'", expectedError, err.Error())
	}
}

func TestRuntimeValidateWithEndNode(t *testing.T) {
	// Create a valid graph: start -> endNode
	endNode := graph.CreateEndNode[TestSharedState]()
	startEdge := graph.CreateStartEdge(endNode)
	runtime := graph.CreateRuntime(startEdge)

	err := runtime.Validate()
	if err != nil {
		t.Errorf("Validate returned error for valid graph ending with EndNode: %v", err)
	}
}

func TestRuntimeValidateComplexGraph(t *testing.T) {
	// Create a more complex valid graph: start -> node1 -> node2 -> end
	node1, _ := graph.CreateNode("node1", func(state TestSharedState) (TestSharedState, error) {
		return state, nil
	})
	node2, _ := graph.CreateNode("node2", func(state TestSharedState) (TestSharedState, error) {
		return state, nil
	})

	startEdge := graph.CreateStartEdge(node1)
	runtime := graph.CreateRuntime(startEdge)

	edge1 := graph.CreateEdge(node1, node2)
	runtime.AddEdge(edge1)

	endEdge := graph.CreateEndEdge(node2)
	runtime.AddEdge(endEdge)

	err := runtime.Validate()
	if err != nil {
		t.Errorf("Validate returned error for valid complex graph: %v", err)
	}
}

func TestRuntimeCircularReference(t *testing.T) {
	// Create nodes for a circular graph
	node1, _ := graph.CreateNode("node1", func(state TestSharedState) (TestSharedState, error) {
		return state, nil
	})
	node2, _ := graph.CreateNode("node2", func(state TestSharedState) (TestSharedState, error) {
		return state, nil
	})

	startEdge := graph.CreateStartEdge(node1)
	runtime := graph.CreateRuntime(startEdge)

	// Create circular reference: node1 -> node2 -> node1
	edge1 := graph.CreateEdge(node1, node2)
	edge2 := graph.CreateEdge(node2, node1)
	runtime.AddEdge(edge1)
	runtime.AddEdge(edge2)

	// This should fail validation because there's no path to end
	err := runtime.Validate()
	if err == nil {
		t.Error("Validate should return error for circular graph with no end")
	}
}
