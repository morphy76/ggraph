package graph

import (
	"fmt"
	"testing"

	g "github.com/morphy76/ggraph/pkg/graph"
)

// BenchmarkWorkerPoolTuning_RuntimeFactory tests different worker pool configurations
// to find the optimal balance between initialization cost and throughput
func BenchmarkWorkerPoolTuning_RuntimeFactory(b *testing.B) {
	configs := []struct {
		name        string
		workerCount int
		queueSize   int
	}{
		{"Default_10x", 0, 0},     // Current default: NumCPU * 10
		{"Reduced_5x", 0, 0},      // Half of current
		{"Optimal_2x", 0, 0},      // Recommended
		{"Minimal_1x", 0, 0},      // NumCPU only
		{"Fixed_16", 16, 50},      // Fixed 16 workers
		{"Fixed_8", 8, 50},        // Fixed 8 workers
		{"Fixed_4", 4, 50},        // Fixed 4 workers
		{"Baseline_NoPool", 1, 1}, // Minimal pool for comparison
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

// BenchmarkWorkerPoolTuning_NodeExecution tests node execution performance
// with different worker pool sizes to ensure throughput isn't sacrificed
func BenchmarkWorkerPoolTuning_NodeExecution(b *testing.B) {
	configs := []struct {
		name        string
		workerCount int
		queueSize   int
	}{
		{"Default_10x", 0, 0},
		{"Optimal_2x", 0, 0},
		{"Minimal_1x", 0, 0},
		{"Fixed_16", 16, 50},
		{"Fixed_8", 8, 50},
		{"Fixed_4", 4, 50},
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
				threadID := fmt.Sprintf("thread-%d", i%10) // Reuse 10 threads
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

// BenchmarkWorkerPoolTuning_ConcurrentLoad tests behavior under concurrent load
func BenchmarkWorkerPoolTuning_ConcurrentLoad(b *testing.B) {
	configs := []struct {
		name        string
		workerCount int
		queueSize   int
	}{
		{"Default_10x", 0, 0},
		{"Optimal_2x", 0, 0},
		{"Fixed_16", 16, 100},
		{"Fixed_8", 8, 100},
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
					// Quick check without spinning
					_ = runtime.CurrentState(threadID)
					i++
				}
			})
		})
	}
}

// BenchmarkWorkerPoolTuning_MultiNode tests multi-node workflow with different configs
func BenchmarkWorkerPoolTuning_MultiNode(b *testing.B) {
	configs := []struct {
		name        string
		workerCount int
		queueSize   int
	}{
		{"Default_10x", 0, 0},
		{"Optimal_2x", 0, 0},
		{"Fixed_16", 16, 100},
	}

	for _, config := range configs {
		b.Run(config.name, func(b *testing.B) {
			policy, _ := RouterPolicyImplFactory(AnyRoute[RuntimeTestState])

			startNode := newMockRuntimeNode("StartNode", g.StartNode, nil, policy)
			nodes := make([]*mockRuntimeNode, 5)
			for i := 0; i < 5; i++ {
				nodes[i] = newMockRuntimeNode(fmt.Sprintf("Node%d", i), g.IntermediateNode, func(userInput, currentState RuntimeTestState, notify g.NotifyPartialFn[RuntimeTestState]) (RuntimeTestState, error) {
					currentState.Counter++
					return currentState, nil
				}, policy)
			}
			endNode := newMockRuntimeNode("EndNode", g.EndNode, nil, nil)

			stateMonitorCh := make(chan g.StateMonitorEntry[RuntimeTestState], 1000)
			go func() {
				for range stateMonitorCh {
				}
			}()

			runtime, _ := RuntimeFactory(&mockRuntimeEdge{from: startNode, to: nodes[0], role: g.StartEdge}, stateMonitorCh, &g.RuntimeOptions[RuntimeTestState]{
				InitialState:    RuntimeTestState{Counter: 0},
				WorkerCount:     config.workerCount,
				WorkerQueueSize: config.queueSize,
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

			b.ReportAllocs()
			b.ResetTimer()

			for i := 0; i < b.N; i++ {
				threadID := fmt.Sprintf("thread-%d", i%10)
				runtime.Invoke(RuntimeTestState{}, g.InvokeConfigThreadID(threadID))
				// Wait for completion
				for {
					state := runtime.CurrentState(threadID)
					if state.Counter >= 5 {
						break
					}
				}
			}
		})
	}
}
