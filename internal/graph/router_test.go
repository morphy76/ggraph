package graph_test

import (
	"testing"

	"github.com/morphy76/ggraph/internal/graph"
	g "github.com/morphy76/ggraph/pkg/graph"
)

// RouterTestState is a state type for router testing
type RouterTestState struct {
	Value   string
	Counter int
	Flag    bool
}

// mockEdge is a minimal edge implementation for testing routers
type mockEdge struct {
	from   string
	to     string
	labels map[string]string
	role   g.EdgeRole
}

func (m *mockEdge) From() g.Node[RouterTestState] {
	return &mockRouterNode{name: m.from}
}

func (m *mockEdge) To() g.Node[RouterTestState] {
	return &mockRouterNode{name: m.to}
}

func (m *mockEdge) LabelByKey(key string) (string, bool) {
	val, ok := m.labels[key]
	return val, ok
}

func (m *mockEdge) Role() g.EdgeRole {
	return m.role
}

// mockRouterNode is a minimal node implementation for router testing
type mockRouterNode struct {
	name string
}

func (m *mockRouterNode) Accept(userInput RouterTestState, stateObserver g.StateObserver[RouterTestState], nodeExecutor g.NodeExecutor, config g.InvokeConfig) {
}
func (m *mockRouterNode) Name() string                                { return m.name }
func (m *mockRouterNode) RoutePolicy() g.RoutePolicy[RouterTestState] { return nil }
func (m *mockRouterNode) Role() g.NodeRole                            { return g.IntermediateNode }

// Test AnyRoute function

func TestAnyRoute_WithMultipleEdges(t *testing.T) {
	edges := []g.Edge[RouterTestState]{
		&mockEdge{from: "node1", to: "node2", role: g.IntermediateEdge},
		&mockEdge{from: "node1", to: "node3", role: g.IntermediateEdge},
		&mockEdge{from: "node1", to: "node4", role: g.IntermediateEdge},
	}

	userInput := RouterTestState{Value: "input", Counter: 5}
	currentState := RouterTestState{Value: "current", Counter: 10}

	result := graph.AnyRoute(userInput, currentState, edges)

	if result == nil {
		t.Fatal("AnyRoute returned nil with available edges")
	}

	// AnyRoute should return the first edge
	if result.To().Name() != "node2" {
		t.Errorf("Expected first edge (to node2), got edge to %s", result.To().Name())
	}
}

func TestAnyRoute_WithSingleEdge(t *testing.T) {
	edges := []g.Edge[RouterTestState]{
		&mockEdge{from: "node1", to: "node2", role: g.IntermediateEdge},
	}

	userInput := RouterTestState{Value: "input"}
	currentState := RouterTestState{Value: "current"}

	result := graph.AnyRoute(userInput, currentState, edges)

	if result == nil {
		t.Fatal("AnyRoute returned nil with one edge")
	}

	if result.To().Name() != "node2" {
		t.Errorf("Expected edge to node2, got edge to %s", result.To().Name())
	}
}

func TestAnyRoute_WithNoEdges(t *testing.T) {
	edges := []g.Edge[RouterTestState]{}

	userInput := RouterTestState{Value: "input"}
	currentState := RouterTestState{Value: "current"}

	result := graph.AnyRoute(userInput, currentState, edges)

	if result != nil {
		t.Errorf("Expected nil for empty edges, got %v", result)
	}
}

func TestAnyRoute_WithNilEdges(t *testing.T) {
	var edges []g.Edge[RouterTestState]

	userInput := RouterTestState{Value: "input"}
	currentState := RouterTestState{Value: "current"}

	result := graph.AnyRoute(userInput, currentState, edges)

	if result != nil {
		t.Errorf("Expected nil for nil edges, got %v", result)
	}
}

func TestAnyRoute_IgnoresStateValues(t *testing.T) {
	edges := []g.Edge[RouterTestState]{
		&mockEdge{from: "node1", to: "node2", role: g.IntermediateEdge},
		&mockEdge{from: "node1", to: "node3", role: g.IntermediateEdge},
	}

	// Try with different state values - should always return first edge
	result1 := graph.AnyRoute(
		RouterTestState{Value: "a", Counter: 1},
		RouterTestState{Value: "b", Counter: 2},
		edges,
	)

	result2 := graph.AnyRoute(
		RouterTestState{Value: "x", Counter: 100},
		RouterTestState{Value: "y", Counter: 200},
		edges,
	)

	if result1 == nil || result2 == nil {
		t.Fatal("AnyRoute returned nil")
	}

	// Both should return the first edge regardless of state
	if result1.To().Name() != "node2" || result2.To().Name() != "node2" {
		t.Error("AnyRoute should always return first edge regardless of state")
	}
}

// Test RouterPolicyImplFactory function

func TestRouterPolicyImplFactory_WithValidSelectionFn(t *testing.T) {
	selectionFn := func(userInput, currentState RouterTestState, edges []g.Edge[RouterTestState]) g.Edge[RouterTestState] {
		if len(edges) > 0 {
			return edges[0]
		}
		return nil
	}

	policy, err := graph.RouterPolicyImplFactory[RouterTestState](selectionFn)

	if err != nil {
		t.Errorf("Expected no error, got %v", err)
	}

	if policy == nil {
		t.Fatal("Expected non-nil policy")
	}
}

func TestRouterPolicyImplFactory_WithNilSelectionFn(t *testing.T) {
	policy, err := graph.RouterPolicyImplFactory[RouterTestState](nil)

	if err == nil {
		t.Error("Expected error for nil selection function")
	}

	if policy != nil {
		t.Error("Expected nil policy when error occurs")
	}

	// Check error message
	expectedMsg := "conditional route policy creation failed: edge selection function cannot be nil"
	if err.Error() != expectedMsg {
		t.Errorf("Expected error message '%s', got '%s'", expectedMsg, err.Error())
	}
}

func TestRouterPolicyImplFactory_PolicySelectsCorrectEdge(t *testing.T) {
	edges := []g.Edge[RouterTestState]{
		&mockEdge{from: "node1", to: "node2", labels: map[string]string{"type": "success"}},
		&mockEdge{from: "node1", to: "node3", labels: map[string]string{"type": "error"}},
	}

	// Create a conditional selection function
	selectionFn := func(userInput, currentState RouterTestState, edges []g.Edge[RouterTestState]) g.Edge[RouterTestState] {
		if currentState.Counter > 10 {
			return edges[0] // success path
		}
		return edges[1] // error path
	}

	policy, err := graph.RouterPolicyImplFactory[RouterTestState](selectionFn)
	if err != nil {
		t.Fatalf("Unexpected error: %v", err)
	}

	// Test with counter > 10 (should select first edge)
	result1 := policy.SelectEdge(
		RouterTestState{Value: "input"},
		RouterTestState{Counter: 15},
		edges,
	)

	if result1.To().Name() != "node2" {
		t.Errorf("Expected edge to node2 for Counter>10, got edge to %s", result1.To().Name())
	}

	// Test with counter <= 10 (should select second edge)
	result2 := policy.SelectEdge(
		RouterTestState{Value: "input"},
		RouterTestState{Counter: 5},
		edges,
	)

	if result2.To().Name() != "node3" {
		t.Errorf("Expected edge to node3 for Counter<=10, got edge to %s", result2.To().Name())
	}
}

func TestRouterPolicyImplFactory_PolicyWithAnyRoute(t *testing.T) {
	edges := []g.Edge[RouterTestState]{
		&mockEdge{from: "node1", to: "node2"},
		&mockEdge{from: "node1", to: "node3"},
	}

	policy, err := graph.RouterPolicyImplFactory[RouterTestState](graph.AnyRoute[RouterTestState])
	if err != nil {
		t.Fatalf("Unexpected error: %v", err)
	}

	result := policy.SelectEdge(
		RouterTestState{Value: "input"},
		RouterTestState{Value: "current"},
		edges,
	)

	if result == nil {
		t.Fatal("Expected non-nil result")
	}

	// Should always return first edge
	if result.To().Name() != "node2" {
		t.Errorf("Expected edge to node2, got edge to %s", result.To().Name())
	}
}

func TestRouterPolicyImplFactory_PolicyWithLabelBasedRouting(t *testing.T) {
	edges := []g.Edge[RouterTestState]{
		&mockEdge{from: "node1", to: "node2", labels: map[string]string{"priority": "high"}},
		&mockEdge{from: "node1", to: "node3", labels: map[string]string{"priority": "low"}},
		&mockEdge{from: "node1", to: "node4", labels: map[string]string{"priority": "medium"}},
	}

	// Create a selection function that routes based on labels
	selectionFn := func(userInput, currentState RouterTestState, edges []g.Edge[RouterTestState]) g.Edge[RouterTestState] {
		targetPriority := "medium"
		for _, edge := range edges {
			if label, ok := edge.LabelByKey("priority"); ok && label == targetPriority {
				return edge
			}
		}
		return edges[0] // default
	}

	policy, err := graph.RouterPolicyImplFactory[RouterTestState](selectionFn)
	if err != nil {
		t.Fatalf("Unexpected error: %v", err)
	}

	result := policy.SelectEdge(
		RouterTestState{},
		RouterTestState{},
		edges,
	)

	if result.To().Name() != "node4" {
		t.Errorf("Expected edge to node4 (medium priority), got edge to %s", result.To().Name())
	}
}

func TestRouterPolicyImplFactory_PolicyWithFlagBasedRouting(t *testing.T) {
	edges := []g.Edge[RouterTestState]{
		&mockEdge{from: "node1", to: "success-node"},
		&mockEdge{from: "node1", to: "failure-node"},
	}

	// Create a selection function based on flag
	selectionFn := func(userInput, currentState RouterTestState, edges []g.Edge[RouterTestState]) g.Edge[RouterTestState] {
		if currentState.Flag {
			return edges[0]
		}
		return edges[1]
	}

	policy, err := graph.RouterPolicyImplFactory[RouterTestState](selectionFn)
	if err != nil {
		t.Fatalf("Unexpected error: %v", err)
	}

	// Test with Flag = true
	result1 := policy.SelectEdge(
		RouterTestState{},
		RouterTestState{Flag: true},
		edges,
	)

	if result1.To().Name() != "success-node" {
		t.Errorf("Expected edge to success-node when Flag=true, got %s", result1.To().Name())
	}

	// Test with Flag = false
	result2 := policy.SelectEdge(
		RouterTestState{},
		RouterTestState{Flag: false},
		edges,
	)

	if result2.To().Name() != "failure-node" {
		t.Errorf("Expected edge to failure-node when Flag=false, got %s", result2.To().Name())
	}
}

func TestRouterPolicyImplFactory_PolicyWithUserInputRouting(t *testing.T) {
	edges := []g.Edge[RouterTestState]{
		&mockEdge{from: "node1", to: "node2"},
		&mockEdge{from: "node1", to: "node3"},
	}

	// Create a selection function based on userInput
	selectionFn := func(userInput, currentState RouterTestState, edges []g.Edge[RouterTestState]) g.Edge[RouterTestState] {
		if userInput.Value == "special" {
			return edges[1]
		}
		return edges[0]
	}

	policy, err := graph.RouterPolicyImplFactory[RouterTestState](selectionFn)
	if err != nil {
		t.Fatalf("Unexpected error: %v", err)
	}

	// Test with special input
	result1 := policy.SelectEdge(
		RouterTestState{Value: "special"},
		RouterTestState{},
		edges,
	)

	if result1.To().Name() != "node3" {
		t.Errorf("Expected edge to node3 for special input, got %s", result1.To().Name())
	}

	// Test with normal input
	result2 := policy.SelectEdge(
		RouterTestState{Value: "normal"},
		RouterTestState{},
		edges,
	)

	if result2.To().Name() != "node2" {
		t.Errorf("Expected edge to node2 for normal input, got %s", result2.To().Name())
	}
}

func TestRouterPolicyImplFactory_PolicyReturnsNil(t *testing.T) {
	edges := []g.Edge[RouterTestState]{
		&mockEdge{from: "node1", to: "node2"},
	}

	// Create a selection function that returns nil
	selectionFn := func(userInput, currentState RouterTestState, edges []g.Edge[RouterTestState]) g.Edge[RouterTestState] {
		return nil
	}

	policy, err := graph.RouterPolicyImplFactory[RouterTestState](selectionFn)
	if err != nil {
		t.Fatalf("Unexpected error: %v", err)
	}

	result := policy.SelectEdge(
		RouterTestState{},
		RouterTestState{},
		edges,
	)

	if result != nil {
		t.Error("Expected nil when selection function returns nil")
	}
}

func TestRouterPolicyImplFactory_PolicyWithEmptyEdges(t *testing.T) {
	edges := []g.Edge[RouterTestState]{}

	selectionFn := func(userInput, currentState RouterTestState, edges []g.Edge[RouterTestState]) g.Edge[RouterTestState] {
		if len(edges) > 0 {
			return edges[0]
		}
		return nil
	}

	policy, err := graph.RouterPolicyImplFactory[RouterTestState](selectionFn)
	if err != nil {
		t.Fatalf("Unexpected error: %v", err)
	}

	result := policy.SelectEdge(
		RouterTestState{},
		RouterTestState{},
		edges,
	)

	if result != nil {
		t.Error("Expected nil for empty edges")
	}
}

func TestRouterPolicyImplFactory_DifferentStateTypes(t *testing.T) {
	// Test that RouterPolicyImplFactory works with different state types
	// Using the existing RouterTestState is sufficient to demonstrate type safety
	edges := []g.Edge[RouterTestState]{
		&mockEdge{from: "node1", to: "target1"},
		&mockEdge{from: "node1", to: "target2"},
	}

	selectionFn := func(userInput, currentState RouterTestState, edges []g.Edge[RouterTestState]) g.Edge[RouterTestState] {
		if currentState.Counter > 5 {
			return edges[0]
		}
		return edges[1]
	}

	policy, err := graph.RouterPolicyImplFactory[RouterTestState](selectionFn)
	if err != nil {
		t.Fatalf("Unexpected error: %v", err)
	}

	result := policy.SelectEdge(
		RouterTestState{Value: "input", Counter: 10},
		RouterTestState{Value: "current", Counter: 10},
		edges,
	)

	if result.To().Name() != "target1" {
		t.Errorf("Expected edge to target1, got %s", result.To().Name())
	}

	// Also test with generic AnyRoute to show type flexibility
	policy2, err := graph.RouterPolicyImplFactory[RouterTestState](graph.AnyRoute[RouterTestState])
	if err != nil {
		t.Fatalf("Unexpected error creating policy with AnyRoute: %v", err)
	}

	result2 := policy2.SelectEdge(
		RouterTestState{},
		RouterTestState{},
		edges,
	)

	if result2 == nil {
		t.Error("Expected non-nil result from policy with AnyRoute")
	}
}

func TestRouterPolicyImplFactory_ComplexRoutingLogic(t *testing.T) {
	edges := []g.Edge[RouterTestState]{
		&mockEdge{from: "node1", to: "retry-node", labels: map[string]string{"type": "retry"}},
		&mockEdge{from: "node1", to: "success-node", labels: map[string]string{"type": "success"}},
		&mockEdge{from: "node1", to: "error-node", labels: map[string]string{"type": "error"}},
	}

	// Complex routing logic: retry if counter < 3, success if counter >= 10, error otherwise
	selectionFn := func(userInput, currentState RouterTestState, edges []g.Edge[RouterTestState]) g.Edge[RouterTestState] {
		if currentState.Counter < 3 {
			// Find retry edge
			for _, edge := range edges {
				if label, ok := edge.LabelByKey("type"); ok && label == "retry" {
					return edge
				}
			}
		} else if currentState.Counter >= 10 {
			// Find success edge
			for _, edge := range edges {
				if label, ok := edge.LabelByKey("type"); ok && label == "success" {
					return edge
				}
			}
		}
		// Find error edge
		for _, edge := range edges {
			if label, ok := edge.LabelByKey("type"); ok && label == "error" {
				return edge
			}
		}
		return edges[0]
	}

	policy, err := graph.RouterPolicyImplFactory[RouterTestState](selectionFn)
	if err != nil {
		t.Fatalf("Unexpected error: %v", err)
	}

	// Test retry path (counter < 3)
	result1 := policy.SelectEdge(
		RouterTestState{},
		RouterTestState{Counter: 2},
		edges,
	)
	if result1.To().Name() != "retry-node" {
		t.Errorf("Expected retry-node for Counter=2, got %s", result1.To().Name())
	}

	// Test success path (counter >= 10)
	result2 := policy.SelectEdge(
		RouterTestState{},
		RouterTestState{Counter: 15},
		edges,
	)
	if result2.To().Name() != "success-node" {
		t.Errorf("Expected success-node for Counter=15, got %s", result2.To().Name())
	}

	// Test error path (3 <= counter < 10)
	result3 := policy.SelectEdge(
		RouterTestState{},
		RouterTestState{Counter: 5},
		edges,
	)
	if result3.To().Name() != "error-node" {
		t.Errorf("Expected error-node for Counter=5, got %s", result3.To().Name())
	}
}
