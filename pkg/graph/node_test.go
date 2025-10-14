package graph_test

import (
	"errors"
	"testing"

	"github.com/morphy76/ggraph/pkg/graph"
)

func TestCreateStartNode(t *testing.T) {
	node := graph.CreateStartNode[TestSharedState]()

	if node == nil {
		t.Fatal("CreateStartNode returned nil")
	}

	if _, ok := node.(*graph.StartNode[TestSharedState]); !ok {
		t.Error("CreateStartNode did not return a StartNode")
	}
}

func TestCreateEndNode(t *testing.T) {
	node := graph.CreateEndNode[TestSharedState]()

	if node == nil {
		t.Fatal("CreateEndNode returned nil")
	}

	if _, ok := node.(*graph.EndNode[TestSharedState]); !ok {
		t.Error("CreateEndNode did not return an EndNode")
	}
}

func TestCreateNode(t *testing.T) {
	testFunc := func(state TestSharedState) (TestSharedState, error) {
		state.Value = "processed"
		return state, nil
	}

	// Test successful node creation
	node, err := graph.CreateNode("testNode", testFunc)
	if err != nil {
		t.Fatalf("CreateNode returned error: %v", err)
	}

	if node == nil {
		t.Fatal("CreateNode returned nil node")
	}
}

func TestCreateNodeEmptyName(t *testing.T) {
	testFunc := func(state TestSharedState) (TestSharedState, error) {
		return state, nil
	}

	// Test with empty name
	node, err := graph.CreateNode("", testFunc)
	if err == nil {
		t.Error("CreateNode should return error for empty name")
	}

	if node != nil {
		t.Error("CreateNode should return nil node for empty name")
	}

	expectedError := "node creation failed: name cannot be empty"
	if err.Error() != expectedError {
		t.Errorf("Expected error message '%s', got '%s'", expectedError, err.Error())
	}
}

func TestCreateNodeNilFunction(t *testing.T) {
	// Test with nil function
	node, err := graph.CreateNode[TestSharedState]("testNode", nil)
	if err == nil {
		t.Error("CreateNode should return error for nil function")
	}

	if node != nil {
		t.Error("CreateNode should return nil node for nil function")
	}

	expectedError := "node creation failed: function cannot be nil"
	if err.Error() != expectedError {
		t.Errorf("Expected error message '%s', got '%s'", expectedError, err.Error())
	}
}

func TestStartNodeAccept(t *testing.T) {
	// Test that StartNode can be created - detailed testing requires integration with Runtime
	startNode := graph.CreateStartNode[TestSharedState]()
	if startNode == nil {
		t.Fatal("CreateStartNode returned nil")
	}

	// Note: Full testing of Accept method requires a proper Runtime setup
	// This is tested in integration tests with the actual runtime
}

func TestEndNodeAccept(t *testing.T) {
	endNode := graph.CreateEndNode[TestSharedState]()
	mockRuntime := graph.CreateRuntime(graph.CreateStartEdge(endNode))

	state := TestSharedState{Value: "test"}

	result, err := endNode.Accept(state, mockRuntime)
	if err != nil {
		t.Errorf("EndNode.Accept returned error: %v", err)
	}

	if result.Value != state.Value {
		t.Errorf("EndNode.Accept should return unchanged state, expected '%s', got '%s'", state.Value, result.Value)
	}
}

func TestNodeImplAccept(t *testing.T) {
	// Test successful processing
	testFunc := func(state TestSharedState) (TestSharedState, error) {
		state.Value = "processed"
		return state, nil
	}

	node, _ := graph.CreateNode("testNode", testFunc)

	// Create mock runtime with outbound edge
	mockRuntime := graph.CreateRuntime(graph.CreateStartEdge(node))
	endNode := graph.CreateEndNode[TestSharedState]()
	edge := graph.CreateEdge(node, endNode)
	mockRuntime.AddEdge(edge)

	state := TestSharedState{Value: "initial"}

	result, err := node.Accept(state, mockRuntime)
	if err != nil {
		t.Errorf("node.Accept returned error: %v", err)
	}

	if result.Value != "processed" {
		t.Errorf("Expected processed state value 'processed', got '%s'", result.Value)
	}
}

func TestNodeImplAcceptWithError(t *testing.T) {
	// Test function that returns error
	testFunc := func(state TestSharedState) (TestSharedState, error) {
		return state, errors.New("processing error")
	}

	node, _ := graph.CreateNode("testNode", testFunc)

	mockRuntime := graph.CreateRuntime(graph.CreateStartEdge(node))
	endNode := graph.CreateEndNode[TestSharedState]()
	edge := graph.CreateEdge(node, endNode)
	mockRuntime.AddEdge(edge)

	state := TestSharedState{Value: "initial"}

	_, err := node.Accept(state, mockRuntime)
	if err == nil {
		t.Error("node.Accept should return error when function fails")
	}

	expectedError := "error processing node 'testNode': processing error"
	if err.Error() != expectedError {
		t.Errorf("Expected error message '%s', got '%s'", expectedError, err.Error())
	}
}

func TestNodeImplAcceptNoOutboundEdges(t *testing.T) {
	testFunc := func(state TestSharedState) (TestSharedState, error) {
		state.Value = "processed"
		return state, nil
	}

	node, _ := graph.CreateNode("testNode", testFunc)

	// Create runtime with no outbound edges from our test node
	// We create a start edge to a different node to avoid runtime validation issues
	dummyNode := graph.CreateEndNode[TestSharedState]()
	mockRuntime := graph.CreateRuntime(graph.CreateStartEdge(dummyNode))

	// Don't add any edges from our test node - this should cause the error

	state := TestSharedState{Value: "initial"}

	_, err := node.Accept(state, mockRuntime)
	if err == nil {
		t.Error("node.Accept should return error when no outbound edges")
	}

	expectedError := "error browsing the graph from node 'testNode'"
	if err.Error() != expectedError {
		t.Errorf("Expected error message '%s', got '%s'", expectedError, err.Error())
	}
}

func TestNodeInterfaceCompliance(t *testing.T) {
	// Test that all node types implement the Node interface
	var node graph.Node[TestSharedState]

	// Test StartNode
	startNode := graph.CreateStartNode[TestSharedState]()
	node = startNode
	if node == nil {
		t.Error("StartNode does not implement Node interface")
	}

	// Test EndNode
	endNode := graph.CreateEndNode[TestSharedState]()
	node = endNode
	if node == nil {
		t.Error("EndNode does not implement Node interface")
	}

	// Test nodeImpl
	testFunc := func(state TestSharedState) (TestSharedState, error) {
		return state, nil
	}
	nodeImpl, _ := graph.CreateNode("test", testFunc)
	node = nodeImpl
	if node == nil {
		t.Error("nodeImpl does not implement Node interface")
	}
}
