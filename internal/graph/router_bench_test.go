package graph_test

import (
	"testing"

	"github.com/morphy76/ggraph/internal/graph"
	g "github.com/morphy76/ggraph/pkg/graph"
)

// BenchmarkAnyRoute tests the performance of the AnyRoute function
func BenchmarkAnyRoute(b *testing.B) {
	edges := []g.Edge[RouterTestState]{
		&mockEdge{from: "node1", to: "node2", role: g.IntermediateEdge},
		&mockEdge{from: "node1", to: "node3", role: g.IntermediateEdge},
		&mockEdge{from: "node1", to: "node4", role: g.IntermediateEdge},
	}

	userInput := RouterTestState{Value: "input", Counter: 5}
	currentState := RouterTestState{Value: "current", Counter: 10}

	b.ResetTimer()
	b.ReportAllocs()

	for i := 0; i < b.N; i++ {
		_ = graph.AnyRoute(userInput, currentState, edges)
	}
}

// BenchmarkAnyRoute_SingleEdge tests performance with a single edge
func BenchmarkAnyRoute_SingleEdge(b *testing.B) {
	edges := []g.Edge[RouterTestState]{
		&mockEdge{from: "node1", to: "node2", role: g.IntermediateEdge},
	}

	userInput := RouterTestState{Value: "input"}
	currentState := RouterTestState{Value: "current"}

	b.ResetTimer()
	b.ReportAllocs()

	for i := 0; i < b.N; i++ {
		_ = graph.AnyRoute(userInput, currentState, edges)
	}
}

// BenchmarkAnyRoute_ManyEdges tests performance with many edges
func BenchmarkAnyRoute_ManyEdges(b *testing.B) {
	edges := make([]g.Edge[RouterTestState], 100)
	for i := 0; i < 100; i++ {
		edges[i] = &mockEdge{from: "node1", to: "node2", role: g.IntermediateEdge}
	}

	userInput := RouterTestState{Value: "input"}
	currentState := RouterTestState{Value: "current"}

	b.ResetTimer()
	b.ReportAllocs()

	for i := 0; i < b.N; i++ {
		_ = graph.AnyRoute(userInput, currentState, edges)
	}
}

// BenchmarkAnyRoute_EmptyEdges tests performance with empty edges
func BenchmarkAnyRoute_EmptyEdges(b *testing.B) {
	edges := []g.Edge[RouterTestState]{}

	userInput := RouterTestState{Value: "input"}
	currentState := RouterTestState{Value: "current"}

	b.ResetTimer()
	b.ReportAllocs()

	for i := 0; i < b.N; i++ {
		_ = graph.AnyRoute(userInput, currentState, edges)
	}
}

// BenchmarkRouterPolicyImplFactory tests the performance of creating router policies
func BenchmarkRouterPolicyImplFactory(b *testing.B) {
	selectionFn := func(userInput, currentState RouterTestState, edges []g.Edge[RouterTestState]) g.Edge[RouterTestState] {
		if len(edges) > 0 {
			return edges[0]
		}
		return nil
	}

	b.ResetTimer()
	b.ReportAllocs()

	for i := 0; i < b.N; i++ {
		_, _ = graph.RouterPolicyImplFactory[RouterTestState](selectionFn)
	}
}

// BenchmarkRouterPolicyImplFactory_WithAnyRoute tests creating policy with AnyRoute
func BenchmarkRouterPolicyImplFactory_WithAnyRoute(b *testing.B) {
	b.ResetTimer()
	b.ReportAllocs()

	for i := 0; i < b.N; i++ {
		_, _ = graph.RouterPolicyImplFactory[RouterTestState](graph.AnyRoute[RouterTestState])
	}
}

// BenchmarkRoutePolicy_SelectEdge_Simple tests simple edge selection
func BenchmarkRoutePolicy_SelectEdge_Simple(b *testing.B) {
	edges := []g.Edge[RouterTestState]{
		&mockEdge{from: "node1", to: "node2"},
		&mockEdge{from: "node1", to: "node3"},
	}

	selectionFn := func(userInput, currentState RouterTestState, edges []g.Edge[RouterTestState]) g.Edge[RouterTestState] {
		if len(edges) > 0 {
			return edges[0]
		}
		return nil
	}

	policy, _ := graph.RouterPolicyImplFactory[RouterTestState](selectionFn)
	userInput := RouterTestState{Value: "input"}
	currentState := RouterTestState{Value: "current"}

	b.ResetTimer()
	b.ReportAllocs()

	for i := 0; i < b.N; i++ {
		_ = policy.SelectEdge(userInput, currentState, edges)
	}
}

// BenchmarkRoutePolicy_SelectEdge_Conditional tests conditional edge selection
func BenchmarkRoutePolicy_SelectEdge_Conditional(b *testing.B) {
	edges := []g.Edge[RouterTestState]{
		&mockEdge{from: "node1", to: "node2"},
		&mockEdge{from: "node1", to: "node3"},
	}

	selectionFn := func(userInput, currentState RouterTestState, edges []g.Edge[RouterTestState]) g.Edge[RouterTestState] {
		if currentState.Counter > 10 {
			return edges[0]
		}
		return edges[1]
	}

	policy, _ := graph.RouterPolicyImplFactory[RouterTestState](selectionFn)
	userInput := RouterTestState{Value: "input"}
	currentState := RouterTestState{Counter: 15}

	b.ResetTimer()
	b.ReportAllocs()

	for i := 0; i < b.N; i++ {
		_ = policy.SelectEdge(userInput, currentState, edges)
	}
}

// BenchmarkRoutePolicy_SelectEdge_LabelBased tests label-based edge selection
func BenchmarkRoutePolicy_SelectEdge_LabelBased(b *testing.B) {
	edges := []g.Edge[RouterTestState]{
		&mockEdge{from: "node1", to: "node2", labels: map[string]string{"priority": "high"}},
		&mockEdge{from: "node1", to: "node3", labels: map[string]string{"priority": "low"}},
		&mockEdge{from: "node1", to: "node4", labels: map[string]string{"priority": "medium"}},
	}

	selectionFn := func(userInput, currentState RouterTestState, edges []g.Edge[RouterTestState]) g.Edge[RouterTestState] {
		targetPriority := "medium"
		for _, edge := range edges {
			if label, ok := edge.LabelByKey("priority"); ok && label == targetPriority {
				return edge
			}
		}
		return edges[0]
	}

	policy, _ := graph.RouterPolicyImplFactory[RouterTestState](selectionFn)
	userInput := RouterTestState{}
	currentState := RouterTestState{}

	b.ResetTimer()
	b.ReportAllocs()

	for i := 0; i < b.N; i++ {
		_ = policy.SelectEdge(userInput, currentState, edges)
	}
}

// BenchmarkRoutePolicy_SelectEdge_ComplexLogic tests complex routing logic
func BenchmarkRoutePolicy_SelectEdge_ComplexLogic(b *testing.B) {
	edges := []g.Edge[RouterTestState]{
		&mockEdge{from: "node1", to: "retry-node", labels: map[string]string{"type": "retry"}},
		&mockEdge{from: "node1", to: "success-node", labels: map[string]string{"type": "success"}},
		&mockEdge{from: "node1", to: "error-node", labels: map[string]string{"type": "error"}},
	}

	selectionFn := func(userInput, currentState RouterTestState, edges []g.Edge[RouterTestState]) g.Edge[RouterTestState] {
		if currentState.Counter < 3 {
			for _, edge := range edges {
				if label, ok := edge.LabelByKey("type"); ok && label == "retry" {
					return edge
				}
			}
		} else if currentState.Counter >= 10 {
			for _, edge := range edges {
				if label, ok := edge.LabelByKey("type"); ok && label == "success" {
					return edge
				}
			}
		}
		for _, edge := range edges {
			if label, ok := edge.LabelByKey("type"); ok && label == "error" {
				return edge
			}
		}
		return edges[0]
	}

	policy, _ := graph.RouterPolicyImplFactory[RouterTestState](selectionFn)
	userInput := RouterTestState{}
	currentState := RouterTestState{Counter: 5}

	b.ResetTimer()
	b.ReportAllocs()

	for i := 0; i < b.N; i++ {
		_ = policy.SelectEdge(userInput, currentState, edges)
	}
}

// BenchmarkRoutePolicy_SelectEdge_ManyEdges tests performance with many edges
func BenchmarkRoutePolicy_SelectEdge_ManyEdges(b *testing.B) {
	edges := make([]g.Edge[RouterTestState], 100)
	for i := 0; i < 100; i++ {
		edges[i] = &mockEdge{
			from:   "node1",
			to:     "node2",
			labels: map[string]string{"index": string(rune(i))},
		}
	}

	selectionFn := func(userInput, currentState RouterTestState, edges []g.Edge[RouterTestState]) g.Edge[RouterTestState] {
		// Select the middle edge
		if len(edges) > 50 {
			return edges[50]
		}
		return edges[0]
	}

	policy, _ := graph.RouterPolicyImplFactory[RouterTestState](selectionFn)
	userInput := RouterTestState{}
	currentState := RouterTestState{}

	b.ResetTimer()
	b.ReportAllocs()

	for i := 0; i < b.N; i++ {
		_ = policy.SelectEdge(userInput, currentState, edges)
	}
}

// BenchmarkRoutePolicy_SelectEdge_WithAnyRoute tests AnyRoute through policy
func BenchmarkRoutePolicy_SelectEdge_WithAnyRoute(b *testing.B) {
	edges := []g.Edge[RouterTestState]{
		&mockEdge{from: "node1", to: "node2"},
		&mockEdge{from: "node1", to: "node3"},
		&mockEdge{from: "node1", to: "node4"},
	}

	policy, _ := graph.RouterPolicyImplFactory[RouterTestState](graph.AnyRoute[RouterTestState])
	userInput := RouterTestState{Value: "input"}
	currentState := RouterTestState{Value: "current"}

	b.ResetTimer()
	b.ReportAllocs()

	for i := 0; i < b.N; i++ {
		_ = policy.SelectEdge(userInput, currentState, edges)
	}
}

// BenchmarkRoutePolicy_SelectEdge_FlagBased tests flag-based routing
func BenchmarkRoutePolicy_SelectEdge_FlagBased(b *testing.B) {
	edges := []g.Edge[RouterTestState]{
		&mockEdge{from: "node1", to: "success-node"},
		&mockEdge{from: "node1", to: "failure-node"},
	}

	selectionFn := func(userInput, currentState RouterTestState, edges []g.Edge[RouterTestState]) g.Edge[RouterTestState] {
		if currentState.Flag {
			return edges[0]
		}
		return edges[1]
	}

	policy, _ := graph.RouterPolicyImplFactory[RouterTestState](selectionFn)
	userInput := RouterTestState{}
	currentState := RouterTestState{Flag: true}

	b.ResetTimer()
	b.ReportAllocs()

	for i := 0; i < b.N; i++ {
		_ = policy.SelectEdge(userInput, currentState, edges)
	}
}

// BenchmarkRoutePolicy_SelectEdge_UserInputBased tests user input-based routing
func BenchmarkRoutePolicy_SelectEdge_UserInputBased(b *testing.B) {
	edges := []g.Edge[RouterTestState]{
		&mockEdge{from: "node1", to: "node2"},
		&mockEdge{from: "node1", to: "node3"},
	}

	selectionFn := func(userInput, currentState RouterTestState, edges []g.Edge[RouterTestState]) g.Edge[RouterTestState] {
		if userInput.Value == "special" {
			return edges[1]
		}
		return edges[0]
	}

	policy, _ := graph.RouterPolicyImplFactory[RouterTestState](selectionFn)
	userInput := RouterTestState{Value: "special"}
	currentState := RouterTestState{}

	b.ResetTimer()
	b.ReportAllocs()

	for i := 0; i < b.N; i++ {
		_ = policy.SelectEdge(userInput, currentState, edges)
	}
}
