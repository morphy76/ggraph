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

// TestNewNodeBuilder_BasicCreation tests creating a new NodeBuilder with basic parameters
func TestNewNodeBuilder_BasicCreation(t *testing.T) {
	builder := builders.NewNodeBuilder("TestNode", mockNodeFn)

	node, err := builder.Build()
	if err != nil {
		t.Fatalf("Build() failed: %v", err)
	}

	if node == nil {
		t.Fatal("Build() returned nil node")
	}

	if node.Name() != "TestNode" {
		t.Errorf("Expected node name 'TestNode', got '%s'", node.Name())
	}

	if node.Role() != g.IntermediateNode {
		t.Errorf("Expected role IntermediateNode, got %v", node.Role())
	}
}

// TestNewNodeBuilder_WithNilFunction tests creating a NodeBuilder with nil function (for routers)
func TestNewNodeBuilder_WithNilFunction(t *testing.T) {
	builder := builders.NewNodeBuilder[TestState]("RouterNode", nil)

	node, err := builder.Build()
	if err != nil {
		t.Fatalf("Build() failed: %v", err)
	}

	if node == nil {
		t.Fatal("Build() returned nil node")
	}

	if node.Name() != "RouterNode" {
		t.Errorf("Expected node name 'RouterNode', got '%s'", node.Name())
	}
}

// TestNodeBuilder_WithRoutingPolicy tests setting a custom routing policy
func TestNodeBuilder_WithRoutingPolicy(t *testing.T) {
	policy, err := builders.CreateConditionalRoutePolicy(mockEdgeSelectionFn)
	if err != nil {
		t.Fatalf("CreateConditionalRoutePolicy() failed: %v", err)
	}

	builder := builders.NewNodeBuilder("TestNode", mockNodeFn).
		WithRoutingPolicy(policy)

	node, err := builder.Build()
	if err != nil {
		t.Fatalf("Build() failed: %v", err)
	}

	if node.RoutePolicy() == nil {
		t.Fatal("RoutePolicy() returned nil")
	}

	// Verify the policy is the one we set
	if node.RoutePolicy() != policy {
		t.Error("RoutePolicy() did not return the expected policy")
	}
}

// TestNodeBuilder_WithReducer tests setting a custom reducer function
func TestNodeBuilder_WithReducer(t *testing.T) {
	builder := builders.NewNodeBuilder("TestNode", mockNodeFn).
		WithReducer(mockReducer)

	node, err := builder.Build()
	if err != nil {
		t.Fatalf("Build() failed: %v", err)
	}

	if node == nil {
		t.Fatal("Build() returned nil node")
	}

	// We can't directly test the reducer, but we can verify the node was created successfully
	if node.Name() != "TestNode" {
		t.Errorf("Expected node name 'TestNode', got '%s'", node.Name())
	}
}

// TestNodeBuilder_ChainedConfiguration tests chaining multiple configuration methods
func TestNodeBuilder_ChainedConfiguration(t *testing.T) {
	policy, _ := builders.CreateAnyRoutePolicy[TestState]()

	builder := builders.NewNodeBuilder("TestNode", mockNodeFn).
		WithRoutingPolicy(policy).
		WithReducer(mockReducer)

	node, err := builder.Build()
	if err != nil {
		t.Fatalf("Build() failed: %v", err)
	}

	if node == nil {
		t.Fatal("Build() returned nil node")
	}

	if node.Name() != "TestNode" {
		t.Errorf("Expected node name 'TestNode', got '%s'", node.Name())
	}

	if node.RoutePolicy() == nil {
		t.Fatal("RoutePolicy() returned nil")
	}
}

// TestNodeBuilder_DefaultRoutingPolicy tests that a default routing policy is set when none is provided
func TestNodeBuilder_DefaultRoutingPolicy(t *testing.T) {
	builder := builders.NewNodeBuilder("TestNode", mockNodeFn)

	node, err := builder.Build()
	if err != nil {
		t.Fatalf("Build() failed: %v", err)
	}

	if node.RoutePolicy() == nil {
		t.Fatal("Expected default routing policy to be set, got nil")
	}
}

// TestNodeBuilder_ReservedNameStart tests that using the reserved StartNode name fails
func TestNodeBuilder_ReservedNameStart(t *testing.T) {
	builder := builders.NewNodeBuilder(builders.ReservedNodeNameStart, mockNodeFn)

	node, err := builder.Build()
	if err == nil {
		t.Fatal("Expected error when using reserved name 'StartNode', got nil")
	}

	if node != nil {
		t.Error("Expected nil node when using reserved name, got non-nil node")
	}

	expectedErrMsg := "node name StartNode is reserved and cannot be used"
	if err.Error() != expectedErrMsg {
		t.Errorf("Expected error message '%s', got '%s'", expectedErrMsg, err.Error())
	}
}

// TestNodeBuilder_ReservedNameEnd tests that using the reserved EndNode name fails
func TestNodeBuilder_ReservedNameEnd(t *testing.T) {
	builder := builders.NewNodeBuilder(builders.ReservedNodeNameEnd, mockNodeFn)

	node, err := builder.Build()
	if err == nil {
		t.Fatal("Expected error when using reserved name 'EndNode', got nil")
	}

	if node != nil {
		t.Error("Expected nil node when using reserved name, got non-nil node")
	}

	expectedErrMsg := "node name EndNode is reserved and cannot be used"
	if err.Error() != expectedErrMsg {
		t.Errorf("Expected error message '%s', got '%s'", expectedErrMsg, err.Error())
	}
}

// TestNodeBuilder_EmptyName tests creating a node with an empty name
func TestNodeBuilder_EmptyName(t *testing.T) {
	builder := builders.NewNodeBuilder("", mockNodeFn)

	node, err := builder.Build()
	if err != nil {
		t.Fatalf("Build() failed: %v", err)
	}

	// Empty names are allowed (not reserved)
	if node == nil {
		t.Fatal("Build() returned nil node")
	}

	if node.Name() != "" {
		t.Errorf("Expected empty node name, got '%s'", node.Name())
	}
}

// TestNodeBuilder_DifferentStateTypes tests NodeBuilder with different state types
func TestNodeBuilder_DifferentStateTypes(t *testing.T) {
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
				return builders.NewNodeBuilder("IntNode", nodeFn).Build()
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
				return builders.NewNodeBuilder("StringNode", nodeFn).Build()
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
				return builders.NewNodeBuilder("ComplexNode", nodeFn).Build()
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

// TestNodeBuilder_MultipleBuildsFromSameBuilder tests that the builder pattern works correctly
func TestNodeBuilder_MultipleBuildsFromSameBuilder(t *testing.T) {
	baseBuilder := builders.NewNodeBuilder("TestNode", mockNodeFn)

	// First build
	node1, err1 := baseBuilder.Build()
	if err1 != nil {
		t.Fatalf("First Build() failed: %v", err1)
	}

	// Second build from same builder
	node2, err2 := baseBuilder.Build()
	if err2 != nil {
		t.Fatalf("Second Build() failed: %v", err2)
	}

	// Both should be valid but separate instances
	if node1 == nil || node2 == nil {
		t.Fatal("One or both builds returned nil")
	}

	if node1.Name() != node2.Name() {
		t.Error("Expected both nodes to have the same name")
	}
}

// TestNodeBuilder_ModifyingBuilderAfterBuild tests that modifying builder after Build doesn't affect previous builds
func TestNodeBuilder_ModifyingBuilderAfterBuild(t *testing.T) {
	builder := builders.NewNodeBuilder("Node1", mockNodeFn)

	node1, err := builder.Build()
	if err != nil {
		t.Fatalf("First Build() failed: %v", err)
	}

	// Modify builder with different reducer
	policy, _ := builders.CreateConditionalRoutePolicy(mockEdgeSelectionFn)
	builder = builder.WithRoutingPolicy(policy).WithReducer(mockReducer)

	node2, err := builder.Build()
	if err != nil {
		t.Fatalf("Second Build() failed: %v", err)
	}

	// Both nodes should be valid
	if node1 == nil || node2 == nil {
		t.Fatal("One or both builds returned nil")
	}

	// Names should match
	if node1.Name() != node2.Name() {
		t.Error("Expected both nodes to have the same name")
	}
}

// TestNodeBuilder_NilReducerUsesDefault tests that nil reducer falls back to default
func TestNodeBuilder_NilReducerUsesDefault(t *testing.T) {
	builder := builders.NewNodeBuilder("TestNode", mockNodeFn).
		WithReducer(nil)

	node, err := builder.Build()
	if err != nil {
		t.Fatalf("Build() failed: %v", err)
	}

	if node == nil {
		t.Fatal("Build() returned nil node")
	}

	// The node should still be created with a default reducer
	if node.Name() != "TestNode" {
		t.Errorf("Expected node name 'TestNode', got '%s'", node.Name())
	}
}

// TestNodeBuilder_NilRoutingPolicyUsesDefault tests that nil routing policy falls back to default
func TestNodeBuilder_NilRoutingPolicyUsesDefault(t *testing.T) {
	builder := builders.NewNodeBuilder("TestNode", mockNodeFn).
		WithRoutingPolicy(nil)

	node, err := builder.Build()
	if err != nil {
		t.Fatalf("Build() failed: %v", err)
	}

	if node == nil {
		t.Fatal("Build() returned nil node")
	}

	// Should have a default routing policy
	if node.RoutePolicy() == nil {
		t.Fatal("Expected default routing policy, got nil")
	}
}

// TestNodeBuilder_CompleteConfiguration tests a fully configured node
func TestNodeBuilder_CompleteConfiguration(t *testing.T) {
	policy, err := builders.CreateConditionalRoutePolicy(mockEdgeSelectionFn)
	if err != nil {
		t.Fatalf("CreateConditionalRoutePolicy() failed: %v", err)
	}

	builder := builders.NewNodeBuilder("CompleteNode", mockNodeFn).
		WithRoutingPolicy(policy).
		WithReducer(mockReducer)

	node, err := builder.Build()
	if err != nil {
		t.Fatalf("Build() failed: %v", err)
	}

	// Verify all properties
	if node == nil {
		t.Fatal("Build() returned nil node")
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

// TestNodeBuilder_LongNodeName tests creating a node with a very long name
func TestNodeBuilder_LongNodeName(t *testing.T) {
	longName := "ThisIsAVeryLongNodeNameThatExceedsNormalExpectationsButShouldStillBeValidBecauseThereIsNoLengthRestrictionOnNodeNames"

	builder := builders.NewNodeBuilder(longName, mockNodeFn)

	node, err := builder.Build()
	if err != nil {
		t.Fatalf("Build() failed: %v", err)
	}

	if node.Name() != longName {
		t.Errorf("Expected node name '%s', got '%s'", longName, node.Name())
	}
}

// TestNodeBuilder_SpecialCharactersInName tests node names with special characters
func TestNodeBuilder_SpecialCharactersInName(t *testing.T) {
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
			builder := builders.NewNodeBuilder(name, mockNodeFn)

			node, err := builder.Build()
			if err != nil {
				t.Fatalf("Build() failed for name '%s': %v", name, err)
			}

			if node.Name() != name {
				t.Errorf("Expected node name '%s', got '%s'", name, node.Name())
			}
		})
	}
}

// TestNodeBuilder_RouterPattern tests creating a router-style node (nil function)
func TestNodeBuilder_RouterPattern(t *testing.T) {
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
	builder := builders.NewNodeBuilder[TestState]("ConditionalRouter", nil).
		WithRoutingPolicy(policy)

	node, err := builder.Build()
	if err != nil {
		t.Fatalf("Build() failed: %v", err)
	}

	if node == nil {
		t.Fatal("Build() returned nil node")
	}

	if node.Name() != "ConditionalRouter" {
		t.Errorf("Expected node name 'ConditionalRouter', got '%s'", node.Name())
	}

	if node.RoutePolicy() != policy {
		t.Error("RoutePolicy() did not return the expected policy")
	}
}
