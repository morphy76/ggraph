package graph

import (
	"fmt"
	"testing"

	g "github.com/morphy76/ggraph/pkg/graph"
)

// BenchmarkAggressiveTuning_RuntimeFactory tests very small worker counts
// and lazy initialization to minimize overhead
func BenchmarkAggressiveTuning_RuntimeFactory(b *testing.B) {
	configs := []struct {
		name        string
		workerCount int
		queueSize   int
	}{
		// Very small fixed counts
		{"Fixed_1", 1, 10},
		{"Fixed_2", 2, 10},
		{"Fixed_3", 3, 20},
		{"Fixed_4", 4, 20},
		{"Fixed_6", 6, 30},
		{"Fixed_8", 8, 50},

		// Reference points
		{"Current_Default", 0, 0},
		{"Baseline_Minimal", 1, 1},
	}

	policy, _ := RouterPolicyImplFactory(AnyRoute[RuntimeTestState])
	startNode := newMockRuntimeNode("StartNode", g.StartNode, nil, policy)
	node1 := newMockRuntimeNode("Node1", g.IntermediateNode, nil, policy)
	startEdge := &mockRuntimeEdge{from: startNode, to: node1, role: g.StartEdge}
	initialState := RuntimeTestState{Value: "initial", Counter: 0}

	for _, config := range configs {
		b.Run(config.name, func(b *testing.B) {
			b.ReportAllocs()
			b.ResetTimer()

			for i := 0; i < b.N; i++ {
				stateMonitorCh := make(chan g.StateMonitorEntry[RuntimeTestState], 10)
				go func(ch chan g.StateMonitorEntry[RuntimeTestState]) {
					for range ch {
					}
				}(stateMonitorCh)

				runtime, _ := RuntimeFactory(startEdge, stateMonitorCh, &g.RuntimeOptions[RuntimeTestState]{
					InitialState:    initialState,
					WorkerCount:     config.workerCount,
					WorkerQueueSize: config.queueSize,
				})
				runtime.Shutdown()
				close(stateMonitorCh)
			}
		})
	}
}

// BenchmarkAggressiveTuning_FullWorkflow tests end-to-end performance
// with small worker counts to ensure no throughput degradation
func BenchmarkAggressiveTuning_FullWorkflow(b *testing.B) {
	configs := []struct {
		name        string
		workerCount int
		queueSize   int
	}{
		{"Fixed_1", 1, 10},
		{"Fixed_2", 2, 10},
		{"Fixed_4", 4, 20},
		{"Fixed_8", 8, 50},
		{"Current_Default", 0, 0},
	}

	for _, config := range configs {
		b.Run(config.name, func(b *testing.B) {
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
			go func() {
				for range stateMonitorCh {
				}
			}()

			runtime, _ := RuntimeFactory(startEdge, stateMonitorCh, &g.RuntimeOptions[RuntimeTestState]{
				InitialState:    RuntimeTestState{Counter: 0},
				WorkerCount:     config.workerCount,
				WorkerQueueSize: config.queueSize,
			})
			defer func() {
				runtime.Shutdown()
				close(stateMonitorCh)
			}()

			runtime.AddEdge(endEdge)

			b.ReportAllocs()
			b.ResetTimer()

			for i := 0; i < b.N; i++ {
				threadID := fmt.Sprintf("thread-%d", i%10)
				runtime.Invoke(RuntimeTestState{}, g.InvokeConfigThreadID(threadID))
				// Wait for completion
				for {
					state := runtime.CurrentState(threadID)
					if state.Counter > i/10 {
						break
					}
				}
			}
		})
	}
}

// BenchmarkAggressiveTuning_NodeLatency measures pure node execution latency
func BenchmarkAggressiveTuning_NodeLatency(b *testing.B) {
	configs := []struct {
		name        string
		workerCount int
		queueSize   int
	}{
		{"Fixed_1", 1, 10},
		{"Fixed_2", 2, 10},
		{"Fixed_4", 4, 20},
		{"Fixed_8", 8, 50},
		{"Current_Default", 0, 0},
	}

	for _, config := range configs {
		b.Run(config.name, func(b *testing.B) {
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
			go func() {
				for range stateMonitorCh {
				}
			}()

			runtime, _ := RuntimeFactory(startEdge, stateMonitorCh, &g.RuntimeOptions[RuntimeTestState]{
				InitialState:    RuntimeTestState{Counter: 0},
				WorkerCount:     config.workerCount,
				WorkerQueueSize: config.queueSize,
			})
			defer func() {
				runtime.Shutdown()
				close(stateMonitorCh)
			}()

			runtime.AddEdge(endEdge)

			b.ReportAllocs()
			b.ResetTimer()

			for i := 0; i < b.N; i++ {
				threadID := "bench-thread"
				runtime.Invoke(RuntimeTestState{}, g.InvokeConfigThreadID(threadID))
				// Quick wait
				for {
					state := runtime.CurrentState(threadID)
					if state.Counter > 0 {
						break
					}
				}
			}
		})
	}
}

// BenchmarkAggressiveTuning_HighConcurrency tests behavior with many parallel invocations
func BenchmarkAggressiveTuning_HighConcurrency(b *testing.B) {
	configs := []struct {
		name        string
		workerCount int
		queueSize   int
	}{
		{"Fixed_2", 2, 50},
		{"Fixed_4", 4, 50},
		{"Fixed_8", 8, 100},
		{"Current_Default", 0, 0},
	}

	for _, config := range configs {
		b.Run(config.name, func(b *testing.B) {
			policy, _ := RouterPolicyImplFactory(AnyRoute[RuntimeTestState])

			startNode := newMockRuntimeNode("StartNode", g.StartNode, nil, policy)
			node1 := newMockRuntimeNode("Node1", g.IntermediateNode, func(userInput, currentState RuntimeTestState, notify g.NotifyPartialFn[RuntimeTestState]) (RuntimeTestState, error) {
				currentState.Counter++
				return currentState, nil
			}, policy)
			endNode := newMockRuntimeNode("EndNode", g.EndNode, nil, nil)

			startEdge := &mockRuntimeEdge{from: startNode, to: node1, role: g.StartEdge}
			endEdge := &mockRuntimeEdge{from: node1, to: endNode, role: g.EndEdge}

			stateMonitorCh := make(chan g.StateMonitorEntry[RuntimeTestState], 10000)
			go func() {
				for range stateMonitorCh {
				}
			}()

			runtime, _ := RuntimeFactory(startEdge, stateMonitorCh, &g.RuntimeOptions[RuntimeTestState]{
				InitialState:    RuntimeTestState{Counter: 0},
				WorkerCount:     config.workerCount,
				WorkerQueueSize: config.queueSize,
			})
			defer func() {
				runtime.Shutdown()
				close(stateMonitorCh)
			}()

			runtime.AddEdge(endEdge)

			b.ReportAllocs()
			b.ResetTimer()

			b.RunParallel(func(pb *testing.PB) {
				i := 0
				for pb.Next() {
					threadID := fmt.Sprintf("thread-%d", i%100)
					runtime.Invoke(RuntimeTestState{}, g.InvokeConfigThreadID(threadID))
					i++
				}
			})
		})
	}
}
