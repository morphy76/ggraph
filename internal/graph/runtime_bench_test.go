package graph

import (
	"context"
	"fmt"
	"sync"
	"testing"

	g "github.com/morphy76/ggraph/pkg/graph"
)

// BenchmarkRuntimeFactory tests the performance of creating runtime instances
func BenchmarkRuntimeFactory(b *testing.B) {
	policy, _ := RouterPolicyImplFactory(AnyRoute[RuntimeTestState])
	startNode := newMockRuntimeNode("StartNode", g.StartNode, nil, policy)
	node1 := newMockRuntimeNode("Node1", g.IntermediateNode, nil, policy)
	startEdge := &mockRuntimeEdge{from: startNode, to: node1, role: g.StartEdge}
	initialState := RuntimeTestState{Value: "initial", Counter: 0}

	b.ResetTimer()
	b.ReportAllocs()

	for i := 0; i < b.N; i++ {
		stateMonitorCh := make(chan g.StateMonitorEntry[RuntimeTestState], 10)
		// Drain channel in background
		go func(ch chan g.StateMonitorEntry[RuntimeTestState]) {
			for range ch {
			}
		}(stateMonitorCh)
		runtime, _ := RuntimeFactory(startEdge, stateMonitorCh, &g.RuntimeOptions[RuntimeTestState]{InitialState: initialState})
		runtime.Shutdown()
		close(stateMonitorCh)
	}
}

// BenchmarkRuntime_AddEdge tests the performance of adding edges
func BenchmarkRuntime_AddEdge(b *testing.B) {
	policy, _ := RouterPolicyImplFactory(AnyRoute[RuntimeTestState])
	startNode := newMockRuntimeNode("StartNode", g.StartNode, nil, policy)
	node1 := newMockRuntimeNode("Node1", g.IntermediateNode, nil, policy)
	startEdge := &mockRuntimeEdge{from: startNode, to: node1, role: g.StartEdge}

	stateMonitorCh := make(chan g.StateMonitorEntry[RuntimeTestState], 10)
	// Drain channel in background
	go func() {
		for range stateMonitorCh {
		}
	}()
	runtime, _ := RuntimeFactory(startEdge, stateMonitorCh, &g.RuntimeOptions[RuntimeTestState]{})
	defer func() {
		runtime.Shutdown()
		close(stateMonitorCh)
	}()

	// Create edges to add
	edges := make([]g.Edge[RuntimeTestState], 100)
	for i := 0; i < 100; i++ {
		node := newMockRuntimeNode("Node", g.IntermediateNode, nil, policy)
		edges[i] = &mockRuntimeEdge{from: node1, to: node, role: g.IntermediateEdge}
	}

	b.ResetTimer()
	b.ReportAllocs()

	for i := 0; i < b.N; i++ {
		runtime.AddEdge(edges[i%100])
	}
}

// BenchmarkRuntime_Validate tests the performance of graph validation
func BenchmarkRuntime_Validate(b *testing.B) {
	policy, _ := RouterPolicyImplFactory(AnyRoute[RuntimeTestState])
	startNode := newMockRuntimeNode("StartNode", g.StartNode, nil, policy)
	node1 := newMockRuntimeNode("Node1", g.IntermediateNode, nil, policy)
	node2 := newMockRuntimeNode("Node2", g.IntermediateNode, nil, policy)
	node3 := newMockRuntimeNode("Node3", g.IntermediateNode, nil, policy)
	endNode := newMockRuntimeNode("EndNode", g.EndNode, nil, nil)

	startEdge := &mockRuntimeEdge{from: startNode, to: node1, role: g.StartEdge}
	edge1 := &mockRuntimeEdge{from: node1, to: node2, role: g.IntermediateEdge}
	edge2 := &mockRuntimeEdge{from: node2, to: node3, role: g.IntermediateEdge}
	edge3 := &mockRuntimeEdge{from: node3, to: endNode, role: g.EndEdge}

	stateMonitorCh := make(chan g.StateMonitorEntry[RuntimeTestState], 10)
	// Drain channel in background
	go func() {
		for range stateMonitorCh {
		}
	}()
	runtime, _ := RuntimeFactory(startEdge, stateMonitorCh, &g.RuntimeOptions[RuntimeTestState]{})
	defer func() {
		runtime.Shutdown()
		close(stateMonitorCh)
	}()

	runtime.AddEdge(edge1, edge2, edge3)

	b.ResetTimer()
	b.ReportAllocs()

	for i := 0; i < b.N; i++ {
		_ = runtime.Validate()
	}
}

// BenchmarkRuntime_CurrentState tests the performance of retrieving current state
func BenchmarkRuntime_CurrentState(b *testing.B) {
	policy, _ := RouterPolicyImplFactory(AnyRoute[RuntimeTestState])
	startNode := newMockRuntimeNode("StartNode", g.StartNode, nil, policy)
	node1 := newMockRuntimeNode("Node1", g.IntermediateNode, nil, policy)
	startEdge := &mockRuntimeEdge{from: startNode, to: node1, role: g.StartEdge}

	stateMonitorCh := make(chan g.StateMonitorEntry[RuntimeTestState], 10)
	// Drain channel in background
	go func() {
		for range stateMonitorCh {
		}
	}()
	runtime, _ := RuntimeFactory(startEdge, stateMonitorCh, &g.RuntimeOptions[RuntimeTestState]{
		InitialState: RuntimeTestState{Value: "test", Counter: 42},
	})
	defer func() {
		runtime.Shutdown()
		close(stateMonitorCh)
	}()

	threadID := "test-thread"
	runtime.Invoke(RuntimeTestState{}, g.InvokeConfigThreadID(threadID))

	b.ResetTimer()
	b.ReportAllocs()

	for i := 0; i < b.N; i++ {
		_ = runtime.CurrentState(threadID)
	}
}

// BenchmarkRuntime_SimpleInvoke tests the performance of simple graph execution
func BenchmarkRuntime_SimpleInvoke(b *testing.B) {
	policy, _ := RouterPolicyImplFactory(AnyRoute[RuntimeTestState])

	startNode := newMockRuntimeNode("StartNode", g.StartNode, nil, policy)
	node1 := newMockRuntimeNode("Node1", g.IntermediateNode, func(userInput, currentState RuntimeTestState, notify g.NotifyPartialFn[RuntimeTestState]) (RuntimeTestState, error) {
		currentState.Counter++
		return currentState, nil
	}, policy)
	endNode := newMockRuntimeNode("EndNode", g.EndNode, nil, nil)

	startEdge := &mockRuntimeEdge{from: startNode, to: node1, role: g.StartEdge}
	endEdge := &mockRuntimeEdge{from: node1, to: endNode, role: g.EndEdge}

	stateMonitorCh := make(chan g.StateMonitorEntry[RuntimeTestState], 1000)
	// Drain state monitor channel
	go func() {
		for range stateMonitorCh {
		}
	}()
	runtime, _ := RuntimeFactory(startEdge, stateMonitorCh, &g.RuntimeOptions[RuntimeTestState]{
		InitialState: RuntimeTestState{Counter: 0},
	})
	defer func() {
		runtime.Shutdown()
		close(stateMonitorCh)
	}()

	runtime.AddEdge(endEdge)

	b.ResetTimer()
	b.ReportAllocs()

	for i := 0; i < b.N; i++ {
		threadID := runtime.Invoke(RuntimeTestState{}, g.InvokeConfigThreadID("thread"))
		// Wait for completion by checking state
		for {
			state := runtime.CurrentState(threadID)
			if state.Counter > 0 {
				break
			}
		}
	}
}

// BenchmarkRuntime_MultiNodeInvoke tests the performance of execution through multiple nodes
func BenchmarkRuntime_MultiNodeInvoke(b *testing.B) {
	policy, _ := RouterPolicyImplFactory(AnyRoute[RuntimeTestState])

	startNode := newMockRuntimeNode("StartNode", g.StartNode, nil, policy)
	nodes := make([]*mockRuntimeNode, 5)
	for i := 0; i < 5; i++ {
		nodes[i] = newMockRuntimeNode("Node", g.IntermediateNode, func(userInput, currentState RuntimeTestState, notify g.NotifyPartialFn[RuntimeTestState]) (RuntimeTestState, error) {
			currentState.Counter++
			return currentState, nil
		}, policy)
	}
	endNode := newMockRuntimeNode("EndNode", g.EndNode, nil, nil)

	stateMonitorCh := make(chan g.StateMonitorEntry[RuntimeTestState], 1000)
	// Drain state monitor channel
	go func() {
		for range stateMonitorCh {
		}
	}()
	runtime, _ := RuntimeFactory(&mockRuntimeEdge{from: startNode, to: nodes[0], role: g.StartEdge}, stateMonitorCh, &g.RuntimeOptions[RuntimeTestState]{
		InitialState: RuntimeTestState{Counter: 0},
	})
	defer func() {
		runtime.Shutdown()
		close(stateMonitorCh)
	}()

	// Connect nodes in chain
	for i := 0; i < 4; i++ {
		runtime.AddEdge(&mockRuntimeEdge{from: nodes[i], to: nodes[i+1], role: g.IntermediateEdge})
	}
	runtime.AddEdge(&mockRuntimeEdge{from: nodes[4], to: endNode, role: g.EndEdge})

	b.ResetTimer()
	b.ReportAllocs()

	for i := 0; i < b.N; i++ {
		threadID := runtime.Invoke(RuntimeTestState{}, g.InvokeConfigThreadID("thread"))
		// Wait for completion
		for {
			state := runtime.CurrentState(threadID)
			if state.Counter >= 5 {
				break
			}
		}
	}
}

// BenchmarkRuntime_StateReplace tests the performance of state updates
func BenchmarkRuntime_StateReplace(b *testing.B) {
	policy, _ := RouterPolicyImplFactory(AnyRoute[RuntimeTestState])
	startNode := newMockRuntimeNode("StartNode", g.StartNode, nil, policy)
	node1 := newMockRuntimeNode("Node1", g.IntermediateNode, nil, policy)
	startEdge := &mockRuntimeEdge{from: startNode, to: node1, role: g.StartEdge}

	stateMonitorCh := make(chan g.StateMonitorEntry[RuntimeTestState], 10)
	// Drain channel in background
	go func() {
		for range stateMonitorCh {
		}
	}()
	runtime, _ := RuntimeFactory(startEdge, stateMonitorCh, &g.RuntimeOptions[RuntimeTestState]{
		InitialState: RuntimeTestState{Counter: 0},
	})
	defer func() {
		runtime.Shutdown()
		close(stateMonitorCh)
	}()

	threadID := "test-thread"
	runtime.Invoke(RuntimeTestState{}, g.InvokeConfigThreadID(threadID))

	runtimeImpl := runtime.(*runtimeImpl[RuntimeTestState])

	b.ResetTimer()
	b.ReportAllocs()

	for i := 0; i < b.N; i++ {
		newState := RuntimeTestState{Counter: i, Value: "updated"}
		runtimeImpl.replace(threadID, newState, Replacer[RuntimeTestState])
	}
}

// BenchmarkRuntime_WithPersistence tests the performance with persistence enabled
func BenchmarkRuntime_WithPersistence(b *testing.B) {
	policy, _ := RouterPolicyImplFactory(AnyRoute[RuntimeTestState])

	startNode := newMockRuntimeNode("StartNode", g.StartNode, nil, policy)
	node1 := newMockRuntimeNode("Node1", g.IntermediateNode, func(userInput, currentState RuntimeTestState, notify g.NotifyPartialFn[RuntimeTestState]) (RuntimeTestState, error) {
		currentState.Counter++
		return currentState, nil
	}, policy)
	endNode := newMockRuntimeNode("EndNode", g.EndNode, nil, nil)

	startEdge := &mockRuntimeEdge{from: startNode, to: node1, role: g.StartEdge}
	endEdge := &mockRuntimeEdge{from: node1, to: endNode, role: g.EndEdge}

	memory := &benchMemory{
		states: make(map[string]RuntimeTestState),
		mu:     sync.RWMutex{},
	}

	stateMonitorCh := make(chan g.StateMonitorEntry[RuntimeTestState], 10000)
	completions := make(map[string]chan struct{})
	var completionsMu sync.Mutex

	// Single goroutine to handle all state monitor events
	go func() {
		for entry := range stateMonitorCh {
			if !entry.Running && entry.Error == nil {
				completionsMu.Lock()
				if ch, exists := completions[entry.ThreadID]; exists {
					close(ch)
					delete(completions, entry.ThreadID)
				}
				completionsMu.Unlock()
			}
		}
	}()

	runtime, _ := RuntimeFactory(startEdge, stateMonitorCh, &g.RuntimeOptions[RuntimeTestState]{
		InitialState: RuntimeTestState{Counter: 0},
		Memory:       memory,
	})
	defer func() {
		runtime.Shutdown()
		close(stateMonitorCh)
	}()

	runtime.AddEdge(endEdge)

	b.ResetTimer()
	b.ReportAllocs()

	for i := 0; i < b.N; i++ {
		completed := make(chan struct{})
		threadID := fmt.Sprintf("thread-%d", i)

		completionsMu.Lock()
		completions[threadID] = completed
		completionsMu.Unlock()

		runtime.Invoke(RuntimeTestState{}, g.InvokeConfigThreadID(threadID))
		<-completed
	}
}

// benchMemory is a simple in-memory persistence implementation for benchmarks
type benchMemory struct {
	states map[string]RuntimeTestState
	mu     sync.RWMutex
}

func (m *benchMemory) PersistFn() g.PersistFn[RuntimeTestState] {
	return func(ctx context.Context, threadID string, state RuntimeTestState) error {
		m.mu.Lock()
		m.states[threadID] = state
		m.mu.Unlock()
		return nil
	}
}

func (m *benchMemory) RestoreFn() g.RestoreFn[RuntimeTestState] {
	return func(ctx context.Context, threadID string) (RuntimeTestState, error) {
		m.mu.RLock()
		state, exists := m.states[threadID]
		m.mu.RUnlock()
		if !exists {
			return RuntimeTestState{}, nil
		}
		return state, nil
	}
}

// BenchmarkRuntime_ListThreads tests the performance of listing threads
func BenchmarkRuntime_ListThreads(b *testing.B) {
	policy, _ := RouterPolicyImplFactory(AnyRoute[RuntimeTestState])
	startNode := newMockRuntimeNode("StartNode", g.StartNode, nil, policy)
	node1 := newMockRuntimeNode("Node1", g.IntermediateNode, nil, policy)
	startEdge := &mockRuntimeEdge{from: startNode, to: node1, role: g.StartEdge}

	stateMonitorCh := make(chan g.StateMonitorEntry[RuntimeTestState], 10)
	// Drain channel in background
	go func() {
		for range stateMonitorCh {
		}
	}()
	runtime, _ := RuntimeFactory(startEdge, stateMonitorCh, &g.RuntimeOptions[RuntimeTestState]{})
	defer close(stateMonitorCh)
	defer runtime.Shutdown()

	// Create multiple threads
	for i := 0; i < 100; i++ {
		runtime.Invoke(RuntimeTestState{}, g.InvokeConfigThreadID(string(rune(i))))
	}

	b.ResetTimer()
	b.ReportAllocs()

	for i := 0; i < b.N; i++ {
		_ = runtime.ListThreads()
	}
}

// BenchmarkRuntime_ConditionalRouting tests the performance of conditional routing
func BenchmarkRuntime_ConditionalRouting(b *testing.B) {
	anyPolicy, _ := RouterPolicyImplFactory(AnyRoute[RuntimeTestState])

	conditionalPolicy, _ := RouterPolicyImplFactory(func(userInput, currentState RuntimeTestState, edges []g.Edge[RuntimeTestState]) g.Edge[RuntimeTestState] {
		for _, edge := range edges {
			if label, ok := edge.LabelByKey("route"); ok && label == currentState.Value {
				return edge
			}
		}
		return edges[0] // fallback
	})

	startNode := newMockRuntimeNode("StartNode", g.StartNode, nil, anyPolicy)
	routerNode := newMockRuntimeNode("RouterNode", g.IntermediateNode, func(userInput, currentState RuntimeTestState, notify g.NotifyPartialFn[RuntimeTestState]) (RuntimeTestState, error) {
		currentState.Value = "left"
		return currentState, nil
	}, conditionalPolicy)
	leftNode := newMockRuntimeNode("LeftNode", g.IntermediateNode, func(userInput, currentState RuntimeTestState, notify g.NotifyPartialFn[RuntimeTestState]) (RuntimeTestState, error) {
		currentState.Counter = 100
		return currentState, nil
	}, anyPolicy)
	rightNode := newMockRuntimeNode("RightNode", g.IntermediateNode, func(userInput, currentState RuntimeTestState, notify g.NotifyPartialFn[RuntimeTestState]) (RuntimeTestState, error) {
		currentState.Counter = 200
		return currentState, nil
	}, anyPolicy)
	endNode := newMockRuntimeNode("EndNode", g.EndNode, nil, nil)

	startEdge := &mockRuntimeEdge{from: startNode, to: routerNode, role: g.StartEdge}
	leftEdge := &mockRuntimeEdge{from: routerNode, to: leftNode, role: g.IntermediateEdge, labels: map[string]string{"route": "left"}}
	rightEdge := &mockRuntimeEdge{from: routerNode, to: rightNode, role: g.IntermediateEdge, labels: map[string]string{"route": "right"}}
	endEdgeLeft := &mockRuntimeEdge{from: leftNode, to: endNode, role: g.EndEdge}
	endEdgeRight := &mockRuntimeEdge{from: rightNode, to: endNode, role: g.EndEdge}

	stateMonitorCh := make(chan g.StateMonitorEntry[RuntimeTestState], 1000)
	// Drain state monitor channel
	go func() {
		for range stateMonitorCh {
		}
	}()
	runtime, _ := RuntimeFactory(startEdge, stateMonitorCh, &g.RuntimeOptions[RuntimeTestState]{
		InitialState: RuntimeTestState{Counter: 0},
	})
	defer func() {
		runtime.Shutdown()
		close(stateMonitorCh)
	}()

	runtime.AddEdge(leftEdge, rightEdge, endEdgeLeft, endEdgeRight)

	b.ResetTimer()
	b.ReportAllocs()

	for i := 0; i < b.N; i++ {
		threadID := runtime.Invoke(RuntimeTestState{}, g.InvokeConfigThreadID("thread"))
		// Wait for completion
		for {
			state := runtime.CurrentState(threadID)
			if state.Counter > 0 {
				break
			}
		}
	}
}
