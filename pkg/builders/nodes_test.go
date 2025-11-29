package builders_test

import (
	"errors"
	"testing"

	"github.com/morphy76/ggraph/pkg/builders"
	g "github.com/morphy76/ggraph/pkg/graph"
)

// TestState is a simple state type for testing
type TestState struct {
	Value   string
	Counter int
}

// mockNodeFn is a simple node function for testing
func mockNodeFn(userInput TestState, currentState TestState, notify g.NotifyPartialFn[TestState]) (TestState, error) {
	currentState.Counter++
	return currentState, nil
}

// mockNodeFnWithError returns an error
func mockNodeFnWithError(userInput TestState, currentState TestState, notify g.NotifyPartialFn[TestState]) (TestState, error) {
	return currentState, errors.New("mock error")
} // mockReducer is a simple reducer for testing
func mockReducer(current TestState, change TestState) TestState {
	current.Value = change.Value
	current.Counter += change.Counter
	return current
}

// mockEdgeSelectionFn selects the first edge
func mockEdgeSelectionFn(userInput TestState, currentState TestState, edges []g.Edge[TestState]) g.Edge[TestState] {
	if len(edges) > 0 {
		return edges[0]
	}
	return nil
}

// TestNewNode_BasicCreation tests creating a new node with basic parameters
func TestNewNode_BasicCreation(t *testing.T) {
	node, err := builders.NewNode("TestNode", mockNodeFn)
	if err != nil {
		t.Fatalf("NewNode() failed: %v", err)
	}

	if node == nil {
		t.Fatal("NewNode() returned nil node")
	}

	if node.Name() != "TestNode" {
		t.Errorf("Expected node name 'TestNode', got '%s'", node.Name())
	}

	if node.Role() != g.IntermediateNode {
		t.Errorf("Expected role IntermediateNode, got %v", node.Role())
	}
}

// TestNewNode_WithNilFunction tests creating a node with nil function (for routers)
func TestNewNode_WithNilFunction(t *testing.T) {
	node, err := builders.NewNode[TestState]("RouterNode", nil)
	if err != nil {
		t.Fatalf("NewNode() failed: %v", err)
	}

	if node == nil {
		t.Fatal("NewNode() returned nil node")
	}

	if node.Name() != "RouterNode" {
		t.Errorf("Expected node name 'RouterNode', got '%s'", node.Name())
	}
}

// TestNode_WithRoutingPolicy tests setting a custom routing policy
func TestNode_WithRoutingPolicy(t *testing.T) {
	policy, err := builders.CreateConditionalRoutePolicy(mockEdgeSelectionFn)
	if err != nil {
		t.Fatalf("CreateConditionalRoutePolicy() failed: %v", err)
	}

	node, err := builders.NewNode("TestNode", mockNodeFn,
		g.WithRoutingPolicy(policy))
	if err != nil {
		t.Fatalf("NewNode() failed: %v", err)
	}

	if node.RoutePolicy() == nil {
		t.Fatal("RoutePolicy() returned nil")
	}

	// Verify the policy is the one we set
	if node.RoutePolicy() != policy {
		t.Error("RoutePolicy() did not return the expected policy")
	}
}

// TestNode_WithReducer tests setting a custom reducer function
func TestNode_WithReducer(t *testing.T) {
	node, err := builders.NewNode("TestNode", mockNodeFn,
		g.WithReducer(mockReducer))
	if err != nil {
		t.Fatalf("NewNode() failed: %v", err)
	}

	if node == nil {
		t.Fatal("NewNode() returned nil node")
	}

	// We can't directly test the reducer, but we can verify the node was created successfully
	if node.Name() != "TestNode" {
		t.Errorf("Expected node name 'TestNode', got '%s'", node.Name())
	}
}

// TestNode_MultipleOptions tests using multiple configuration options
func TestNode_MultipleOptions(t *testing.T) {
	policy, _ := builders.CreateAnyRoutePolicy[TestState]()

	node, err := builders.NewNode("TestNode", mockNodeFn,
		g.WithRoutingPolicy(policy),
		g.WithReducer(mockReducer))
	if err != nil {
		t.Fatalf("NewNode() failed: %v", err)
	}

	if node == nil {
		t.Fatal("NewNode() returned nil node")
	}

	if node.Name() != "TestNode" {
		t.Errorf("Expected node name 'TestNode', got '%s'", node.Name())
	}

	if node.RoutePolicy() == nil {
		t.Fatal("RoutePolicy() returned nil")
	}
}

// TestNode_DefaultRoutingPolicy tests that a default routing policy is set when none is provided
func TestNode_DefaultRoutingPolicy(t *testing.T) {
	node, err := builders.NewNode("TestNode", mockNodeFn)
	if err != nil {
		t.Fatalf("NewNode() failed: %v", err)
	}

	if node.RoutePolicy() == nil {
		t.Fatal("Expected default routing policy to be set, got nil")
	}
}

// TestNode_ReservedNameStart tests that using the reserved StartNode name fails
func TestNode_ReservedNameStart(t *testing.T) {
	node, err := builders.NewNode(builders.ReservedNodeNameStart, mockNodeFn)
	if err == nil {
		t.Fatal("Expected error when using reserved name 'StartNode', got nil")
	}

	if node != nil {
		t.Error("Expected nil node when using reserved name, got non-nil node")
	}

	expectedErrMsg := "node creation error for name StartNode: node name is reserved and cannot be used"
	if err.Error() != expectedErrMsg {
		t.Errorf("Expected error message '%s', got '%s'", expectedErrMsg, err.Error())
	}
}

// TestNode_ReservedNameEnd tests that using the reserved EndNode name fails
func TestNode_ReservedNameEnd(t *testing.T) {
	node, err := builders.NewNode(builders.ReservedNodeNameEnd, mockNodeFn)
	if err == nil {
		t.Fatal("Expected error when using reserved name 'EndNode', got nil")
	}

	if node != nil {
		t.Error("Expected nil node when using reserved name, got non-nil node")
	}

	expectedErrMsg := "node creation error for name EndNode: node name is reserved and cannot be used"
	if err.Error() != expectedErrMsg {
		t.Errorf("Expected error message '%s', got '%s'", expectedErrMsg, err.Error())
	}
}

// TestNode_EmptyName tests creating a node with an empty name
func TestNode_EmptyName(t *testing.T) {
	node, err := builders.NewNode("", mockNodeFn)
	if err != nil {
		t.Fatalf("NewNode() failed: %v", err)
	}

	// Empty names are allowed (not reserved)
	if node == nil {
		t.Fatal("NewNode() returned nil node")
	}

	if node.Name() != "" {
		t.Errorf("Expected empty node name, got '%s'", node.Name())
	}
}

// TestNode_DifferentStateTypes tests NewNode with different state types
func TestNode_DifferentStateTypes(t *testing.T) {
	tests := []struct {
		name        string
		stateName   string
		buildFunc   func() (interface{}, error)
		expectedErr bool
	}{
		{
			name:      "IntState",
			stateName: "IntNode",
			buildFunc: func() (interface{}, error) {
				type IntState struct{ Value int }
				nodeFn := func(userInput IntState, currentState IntState, notify g.NotifyPartialFn[IntState]) (IntState, error) {
					return currentState, nil
				}
				return builders.NewNode("IntNode", nodeFn)
			},
			expectedErr: false,
		},
		{
			name:      "StringState",
			stateName: "StringNode",
			buildFunc: func() (interface{}, error) {
				type StringState struct{ Value string }
				nodeFn := func(userInput StringState, currentState StringState, notify g.NotifyPartialFn[StringState]) (StringState, error) {
					return currentState, nil
				}
				return builders.NewNode("StringNode", nodeFn)
			},
			expectedErr: false,
		},
		{
			name:      "ComplexState",
			stateName: "ComplexNode",
			buildFunc: func() (interface{}, error) {
				type ComplexState struct {
					ID    int
					Name  string
					Items []string
					Meta  map[string]interface{}
				}
				nodeFn := func(userInput ComplexState, currentState ComplexState, notify g.NotifyPartialFn[ComplexState]) (ComplexState, error) {
					return currentState, nil
				}
				return builders.NewNode("ComplexNode", nodeFn)
			},
			expectedErr: false,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			node, err := tt.buildFunc()

			if tt.expectedErr && err == nil {
				t.Fatal("Expected error but got none")
			}

			if !tt.expectedErr && err != nil {
				t.Fatalf("Unexpected error: %v", err)
			}

			if !tt.expectedErr && node == nil {
				t.Fatal("Expected non-nil node")
			}
		})
	}
}

// TestNode_MultipleCalls tests creating multiple nodes with the same options
func TestNode_MultipleCalls(t *testing.T) {
	// First creation
	node1, err1 := builders.NewNode("TestNode", mockNodeFn)
	if err1 != nil {
		t.Fatalf("First NewNode() failed: %v", err1)
	}

	// Second creation with same parameters
	node2, err2 := builders.NewNode("TestNode", mockNodeFn)
	if err2 != nil {
		t.Fatalf("Second NewNode() failed: %v", err2)
	}

	// Both should be valid but separate instances
	if node1 == nil || node2 == nil {
		t.Fatal("One or both nodes returned nil")
	}

	if node1.Name() != node2.Name() {
		t.Error("Expected both nodes to have the same name")
	}
}

// TestNode_DifferentConfigurations tests creating nodes with different configurations
func TestNode_DifferentConfigurations(t *testing.T) {
	node1, err := builders.NewNode("Node1", mockNodeFn)
	if err != nil {
		t.Fatalf("First NewNode() failed: %v", err)
	}

	// Create node with different configuration
	policy, _ := builders.CreateConditionalRoutePolicy(mockEdgeSelectionFn)
	node2, err := builders.NewNode("Node1", mockNodeFn,
		g.WithRoutingPolicy(policy),
		g.WithReducer(mockReducer))
	if err != nil {
		t.Fatalf("Second NewNode() failed: %v", err)
	}

	// Both nodes should be valid
	if node1 == nil || node2 == nil {
		t.Fatal("One or both nodes returned nil")
	}

	// Names should match
	if node1.Name() != node2.Name() {
		t.Error("Expected both nodes to have the same name")
	}
}

// TestNode_NilReducerUsesDefault tests that nil reducer falls back to default
func TestNode_NilReducerUsesDefault(t *testing.T) {
	node, err := builders.NewNode("TestNode", mockNodeFn,
		g.WithReducer[TestState](nil))
	if err != nil {
		t.Fatalf("NewNode() failed: %v", err)
	}

	if node == nil {
		t.Fatal("NewNode() returned nil node")
	}

	// The node should still be created with a default reducer
	if node.Name() != "TestNode" {
		t.Errorf("Expected node name 'TestNode', got '%s'", node.Name())
	}
}

// TestNode_NilRoutingPolicyUsesDefault tests that nil routing policy falls back to default
func TestNode_NilRoutingPolicyUsesDefault(t *testing.T) {
	node, err := builders.NewNode("TestNode", mockNodeFn,
		g.WithRoutingPolicy[TestState](nil))
	if err != nil {
		t.Fatalf("NewNode() failed: %v", err)
	}

	if node == nil {
		t.Fatal("NewNode() returned nil node")
	}

	// Should have a default routing policy
	if node.RoutePolicy() == nil {
		t.Fatal("Expected default routing policy, got nil")
	}
}

// TestNode_CompleteConfiguration tests a fully configured node
func TestNode_CompleteConfiguration(t *testing.T) {
	policy, err := builders.CreateConditionalRoutePolicy(mockEdgeSelectionFn)
	if err != nil {
		t.Fatalf("CreateConditionalRoutePolicy() failed: %v", err)
	}

	node, err := builders.NewNode("CompleteNode", mockNodeFn,
		g.WithRoutingPolicy(policy),
		g.WithReducer(mockReducer))
	if err != nil {
		t.Fatalf("NewNode() failed: %v", err)
	}

	// Verify all properties
	if node == nil {
		t.Fatal("NewNode() returned nil node")
	}

	if node.Name() != "CompleteNode" {
		t.Errorf("Expected node name 'CompleteNode', got '%s'", node.Name())
	}

	if node.Role() != g.IntermediateNode {
		t.Errorf("Expected role IntermediateNode, got %v", node.Role())
	}

	if node.RoutePolicy() == nil {
		t.Fatal("RoutePolicy() returned nil")
	}

	if node.RoutePolicy() != policy {
		t.Error("RoutePolicy() did not return the expected policy")
	}
}

// TestNode_LongNodeName tests creating a node with a very long name
func TestNode_LongNodeName(t *testing.T) {
	longName := "ThisIsAVeryLongNodeNameThatExceedsNormalExpectationsButShouldStillBeValidBecauseThereIsNoLengthRestrictionOnNodeNames"

	node, err := builders.NewNode(longName, mockNodeFn)
	if err != nil {
		t.Fatalf("NewNode() failed: %v", err)
	}

	if node.Name() != longName {
		t.Errorf("Expected node name '%s', got '%s'", longName, node.Name())
	}
}

// TestNode_SpecialCharactersInName tests node names with special characters
func TestNode_SpecialCharactersInName(t *testing.T) {
	names := []string{
		"Node-With-Dashes",
		"Node_With_Underscores",
		"Node.With.Dots",
		"Node123WithNumbers",
		"Node With Spaces",
		"Node/With/Slashes",
	}

	for _, name := range names {
		t.Run(name, func(t *testing.T) {
			node, err := builders.NewNode(name, mockNodeFn)
			if err != nil {
				t.Fatalf("NewNode() failed for name '%s': %v", name, err)
			}

			if node.Name() != name {
				t.Errorf("Expected node name '%s', got '%s'", name, node.Name())
			}
		})
	}
}

// TestNode_RouterPattern tests creating a router-style node (nil function)
func TestNode_RouterPattern(t *testing.T) {
	selectionFn := func(userInput TestState, currentState TestState, edges []g.Edge[TestState]) g.Edge[TestState] {
		// Select based on state value
		for _, edge := range edges {
			if label, ok := edge.LabelByKey("route"); ok && label == currentState.Value {
				return edge
			}
		}
		return nil
	}

	policy, err := builders.CreateConditionalRoutePolicy(selectionFn)
	if err != nil {
		t.Fatalf("CreateConditionalRoutePolicy() failed: %v", err)
	}

	// Router nodes have nil function
	node, err := builders.NewNode[TestState]("ConditionalRouter", nil,
		g.WithRoutingPolicy(policy))
	if err != nil {
		t.Fatalf("NewNode() failed: %v", err)
	}

	if node == nil {
		t.Fatal("NewNode() returned nil node")
	}

	if node.Name() != "ConditionalRouter" {
		t.Errorf("Expected node name 'ConditionalRouter', got '%s'", node.Name())
	}

	if node.RoutePolicy() != policy {
		t.Error("RoutePolicy() did not return the expected policy")
	}
}
