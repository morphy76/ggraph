package graph_test

import (
	"testing"

	"github.com/morphy76/ggraph/internal/graph"
	g "github.com/morphy76/ggraph/pkg/graph"
)

// TestState is a simple state type for testing
type TestState struct {
	Value string
}

// mockNode is a minimal node implementation for testing edges
type mockNode struct {
	name string
	role g.NodeRole
}

func (m *mockNode) Accept(userInput TestState, runtime g.StateObserver[TestState]) {}
func (m *mockNode) Name() string                                                   { return m.name }
func (m *mockNode) RoutePolicy() g.RoutePolicy[TestState]                          { return nil }
func (m *mockNode) Role() g.NodeRole                                               { return m.role }

// mockNodeGeneric is a generic mock node for testing with different state types
type mockNodeGeneric[T g.SharedState] struct {
	name string
	role g.NodeRole
}

func (m *mockNodeGeneric[T]) Accept(userInput T, runtime g.StateObserver[T]) {}
func (m *mockNodeGeneric[T]) Name() string                                   { return m.name }
func (m *mockNodeGeneric[T]) RoutePolicy() g.RoutePolicy[T]                  { return nil }
func (m *mockNodeGeneric[T]) Role() g.NodeRole                               { return m.role }

func TestEdgeImplFactory_BasicCreation(t *testing.T) {
	fromNode := &mockNode{name: "from", role: g.IntermediateNode}
	toNode := &mockNode{name: "to", role: g.IntermediateNode}

	edge := graph.EdgeImplFactory[TestState](fromNode, toNode, g.IntermediateEdge)

	if edge == nil {
		t.Fatal("EdgeImplFactory returned nil")
	}

	if edge.From() != fromNode {
		t.Errorf("Expected From() to return fromNode, got different node")
	}

	if edge.To() != toNode {
		t.Errorf("Expected To() to return toNode, got different node")
	}

	if edge.Role() != g.IntermediateEdge {
		t.Errorf("Expected Role() to return IntermediateEdge, got %v", edge.Role())
	}
}

func TestEdgeImplFactory_StartEdge(t *testing.T) {
	startNode := &mockNode{name: "start", role: g.StartNode}
	firstNode := &mockNode{name: "first", role: g.IntermediateNode}

	edge := graph.EdgeImplFactory[TestState](startNode, firstNode, g.StartEdge)

	if edge.Role() != g.StartEdge {
		t.Errorf("Expected Role() to return StartEdge, got %v", edge.Role())
	}

	if edge.From().Role() != g.StartNode {
		t.Errorf("Expected From() to be StartNode, got %v", edge.From().Role())
	}

	if edge.To().Role() != g.IntermediateNode {
		t.Errorf("Expected To() to be IntermediateNode, got %v", edge.To().Role())
	}
}

func TestEdgeImplFactory_EndEdge(t *testing.T) {
	lastNode := &mockNode{name: "last", role: g.IntermediateNode}
	endNode := &mockNode{name: "end", role: g.EndNode}

	edge := graph.EdgeImplFactory[TestState](lastNode, endNode, g.EndEdge)

	if edge.Role() != g.EndEdge {
		t.Errorf("Expected Role() to return EndEdge, got %v", edge.Role())
	}

	if edge.From().Role() != g.IntermediateNode {
		t.Errorf("Expected From() to be IntermediateNode, got %v", edge.From().Role())
	}

	if edge.To().Role() != g.EndNode {
		t.Errorf("Expected To() to be EndNode, got %v", edge.To().Role())
	}
}

func TestEdgeImplFactory_WithSingleLabels(t *testing.T) {
	fromNode := &mockNode{name: "from", role: g.IntermediateNode}
	toNode := &mockNode{name: "to", role: g.IntermediateNode}
	labels := map[string]string{
		"type":     "conditional",
		"priority": "high",
	}

	edge := graph.EdgeImplFactory[TestState](fromNode, toNode, g.IntermediateEdge, labels)

	// Check that labels are correctly stored
	if val, ok := edge.LabelByKey("type"); !ok || val != "conditional" {
		t.Errorf("Expected label 'type' to be 'conditional', got '%v' (ok=%v)", val, ok)
	}

	if val, ok := edge.LabelByKey("priority"); !ok || val != "high" {
		t.Errorf("Expected label 'priority' to be 'high', got '%v' (ok=%v)", val, ok)
	}
}

func TestEdgeImplFactory_WithMultipleLabels(t *testing.T) {
	fromNode := &mockNode{name: "from", role: g.IntermediateNode}
	toNode := &mockNode{name: "to", role: g.IntermediateNode}
	labels1 := map[string]string{
		"type": "conditional",
		"env":  "dev",
	}
	labels2 := map[string]string{
		"priority": "high",
		"team":     "backend",
	}

	edge := graph.EdgeImplFactory[TestState](fromNode, toNode, g.IntermediateEdge, labels1, labels2)

	// Check that all labels from both maps are present
	if val, ok := edge.LabelByKey("type"); !ok || val != "conditional" {
		t.Errorf("Expected label 'type' to be 'conditional', got '%v' (ok=%v)", val, ok)
	}

	if val, ok := edge.LabelByKey("env"); !ok || val != "dev" {
		t.Errorf("Expected label 'env' to be 'dev', got '%v' (ok=%v)", val, ok)
	}

	if val, ok := edge.LabelByKey("priority"); !ok || val != "high" {
		t.Errorf("Expected label 'priority' to be 'high', got '%v' (ok=%v)", val, ok)
	}

	if val, ok := edge.LabelByKey("team"); !ok || val != "backend" {
		t.Errorf("Expected label 'team' to be 'backend', got '%v' (ok=%v)", val, ok)
	}
}

func TestEdgeImplFactory_OverlappingLabels(t *testing.T) {
	fromNode := &mockNode{name: "from", role: g.IntermediateNode}
	toNode := &mockNode{name: "to", role: g.IntermediateNode}
	labels1 := map[string]string{
		"type": "conditional",
		"env":  "dev",
	}
	labels2 := map[string]string{
		"type": "sequential", // This should override the first "type"
		"team": "backend",
	}

	edge := graph.EdgeImplFactory[TestState](fromNode, toNode, g.IntermediateEdge, labels1, labels2)

	// The second label map should override the first for "type"
	if val, ok := edge.LabelByKey("type"); !ok || val != "sequential" {
		t.Errorf("Expected label 'type' to be 'sequential' (overridden), got '%v' (ok=%v)", val, ok)
	}

	if val, ok := edge.LabelByKey("env"); !ok || val != "dev" {
		t.Errorf("Expected label 'env' to be 'dev', got '%v' (ok=%v)", val, ok)
	}

	if val, ok := edge.LabelByKey("team"); !ok || val != "backend" {
		t.Errorf("Expected label 'team' to be 'backend', got '%v' (ok=%v)", val, ok)
	}
}

func TestEdgeImplFactory_NoLabels(t *testing.T) {
	fromNode := &mockNode{name: "from", role: g.IntermediateNode}
	toNode := &mockNode{name: "to", role: g.IntermediateNode}

	edge := graph.EdgeImplFactory[TestState](fromNode, toNode, g.IntermediateEdge)

	// Should return false for any key when no labels provided
	if val, ok := edge.LabelByKey("nonexistent"); ok {
		t.Errorf("Expected LabelByKey('nonexistent') to return false, got true with value '%v'", val)
	}

	if val, ok := edge.LabelByKey("type"); ok {
		t.Errorf("Expected LabelByKey('type') to return false, got true with value '%v'", val)
	}
}

func TestEdgeImplFactory_EmptyLabelMap(t *testing.T) {
	fromNode := &mockNode{name: "from", role: g.IntermediateNode}
	toNode := &mockNode{name: "to", role: g.IntermediateNode}
	emptyLabels := map[string]string{}

	edge := graph.EdgeImplFactory[TestState](fromNode, toNode, g.IntermediateEdge, emptyLabels)

	// Should return false for any key when empty label map provided
	if val, ok := edge.LabelByKey("nonexistent"); ok {
		t.Errorf("Expected LabelByKey('nonexistent') to return false, got true with value '%v'", val)
	}
}

func TestEdgeImplFactory_LabelByKeyNonExistent(t *testing.T) {
	fromNode := &mockNode{name: "from", role: g.IntermediateNode}
	toNode := &mockNode{name: "to", role: g.IntermediateNode}
	labels := map[string]string{
		"type": "conditional",
	}

	edge := graph.EdgeImplFactory[TestState](fromNode, toNode, g.IntermediateEdge, labels)

	// Should return false for non-existent key
	if val, ok := edge.LabelByKey("nonexistent"); ok {
		t.Errorf("Expected LabelByKey('nonexistent') to return false, got true with value '%v'", val)
	}

	// Should return true for existing key
	if val, ok := edge.LabelByKey("type"); !ok || val != "conditional" {
		t.Errorf("Expected LabelByKey('type') to return 'conditional', got '%v' (ok=%v)", val, ok)
	}
}

func TestEdgeImplFactory_DifferentStateTypes(t *testing.T) {
	type AnotherState struct {
		Counter int
	}

	fromNode := &mockNode{name: "from", role: g.IntermediateNode}
	toNode := &mockNode{name: "to", role: g.IntermediateNode}

	// Test that the factory works with TestState
	edge1 := graph.EdgeImplFactory[TestState](fromNode, toNode, g.IntermediateEdge)
	if edge1 == nil {
		t.Error("EdgeImplFactory failed to create edge with TestState")
	}

	// Test with AnotherState - demonstrates generic type safety
	fromNodeAnother := &mockNodeGeneric[AnotherState]{name: "from2", role: g.IntermediateNode}
	toNodeAnother := &mockNodeGeneric[AnotherState]{name: "to2", role: g.IntermediateNode}

	edge2 := graph.EdgeImplFactory[AnotherState](fromNodeAnother, toNodeAnother, g.IntermediateEdge)
	if edge2 == nil {
		t.Error("EdgeImplFactory failed to create edge with AnotherState")
	}
}

func TestEdgeImplFactory_AllRoles(t *testing.T) {
	fromNode := &mockNode{name: "from", role: g.IntermediateNode}
	toNode := &mockNode{name: "to", role: g.IntermediateNode}

	testCases := []struct {
		name string
		role g.EdgeRole
	}{
		{"StartEdge", g.StartEdge},
		{"IntermediateEdge", g.IntermediateEdge},
		{"EndEdge", g.EndEdge},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			edge := graph.EdgeImplFactory[TestState](fromNode, toNode, tc.role)

			if edge.Role() != tc.role {
				t.Errorf("Expected Role() to return %v, got %v", tc.role, edge.Role())
			}
		})
	}
}

func TestEdgeImplFactory_NodeReferences(t *testing.T) {
	fromNode := &mockNode{name: "source", role: g.IntermediateNode}
	toNode := &mockNode{name: "destination", role: g.IntermediateNode}

	edge := graph.EdgeImplFactory[TestState](fromNode, toNode, g.IntermediateEdge)

	// Verify that the edge maintains correct references
	if edge.From().Name() != "source" {
		t.Errorf("Expected From().Name() to be 'source', got '%v'", edge.From().Name())
	}

	if edge.To().Name() != "destination" {
		t.Errorf("Expected To().Name() to be 'destination', got '%v'", edge.To().Name())
	}

	// Verify references are to the same objects
	if edge.From() != fromNode {
		t.Error("From() does not reference the same node object")
	}

	if edge.To() != toNode {
		t.Error("To() does not reference the same node object")
	}
}
