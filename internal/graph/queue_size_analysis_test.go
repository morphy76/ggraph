package graph

import (
	"fmt"
	"testing"

	g "github.com/morphy76/ggraph/pkg/graph"
)

// BenchmarkQueueSizeRatio tests different queue size configurations
// to validate the 3/4 ratio approach
func BenchmarkQueueSizeRatio_RuntimeFactory(b *testing.B) {
	configs := []struct {
		name      string
		queueSize int
		workers   int // 0 means use 3/4 ratio
		note      string
	}{
		// Testing the 3/4 ratio with different queue sizes
		{"Queue_100_Workers_75", 100, 0, "Default: 3/4 of 100 = 75 workers"},
		{"Queue_50_Workers_37", 50, 0, "3/4 of 50 = 37 workers"},
		{"Queue_20_Workers_15", 20, 0, "3/4 of 20 = 15 workers"},
		{"Queue_10_Workers_7", 10, 0, "3/4 of 10 = 7 workers"},

		// Testing different ratios for comparison
		{"Queue_100_Workers_100", 100, 100, "1:1 ratio"},
		{"Queue_100_Workers_50", 100, 50, "1:2 ratio"},
		{"Queue_100_Workers_25", 100, 25, "1:4 ratio"},
		{"Queue_100_Workers_10", 100, 10, "1:10 ratio"},
		{"Queue_100_Workers_4", 100, 4, "1:25 ratio (previous optimal)"},

		// Testing queue size impact
		{"Queue_200_Workers_150", 200, 0, "Large queue: 3/4 of 200 = 150"},
		{"Queue_20_Workers_4", 20, 4, "Small queue with fixed 4 workers"},
		{"Queue_10_Workers_4", 10, 4, "Tiny queue with fixed 4 workers"},
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
					WorkerCount:     config.workers,
					WorkerQueueSize: config.queueSize,
				})
				runtime.Shutdown()
				close(stateMonitorCh)
			}
		})
	}
}

// BenchmarkQueueSizeRatio_Throughput tests execution throughput with different ratios
func BenchmarkQueueSizeRatio_Throughput(b *testing.B) {
	configs := []struct {
		name      string
		queueSize int
		workers   int
	}{
		{"Queue_100_Workers_75", 100, 0},  // 3/4 ratio default
		{"Queue_100_Workers_50", 100, 50}, // 1/2 ratio
		{"Queue_100_Workers_25", 100, 25}, // 1/4 ratio
		{"Queue_100_Workers_10", 100, 10}, // 1/10 ratio
		{"Queue_100_Workers_4", 100, 4},   // Previous optimal
		{"Queue_50_Workers_37", 50, 0},    // 3/4 of smaller queue
		{"Queue_20_Workers_15", 20, 0},    // 3/4 of tiny queue
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
				WorkerCount:     config.workers,
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

// BenchmarkQueueSizeRatio_HighConcurrency tests parallel execution
func BenchmarkQueueSizeRatio_HighConcurrency(b *testing.B) {
	configs := []struct {
		name      string
		queueSize int
		workers   int
	}{
		{"Queue_100_Workers_75", 100, 0},
		{"Queue_50_Workers_37", 50, 0},
		{"Queue_20_Workers_15", 20, 0},
		{"Queue_100_Workers_25", 100, 25},
		{"Queue_100_Workers_10", 100, 10},
		{"Queue_100_Workers_4", 100, 4},
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
				WorkerCount:     config.workers,
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
