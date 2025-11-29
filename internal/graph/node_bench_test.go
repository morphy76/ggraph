package graph_test

import (
	"testing"

	"github.com/morphy76/ggraph/internal/graph"
	g "github.com/morphy76/ggraph/pkg/graph"
)

// BenchmarkNodeFactory tests the performance of creating node instances
func BenchmarkNodeFactory(b *testing.B) {
	nodeFn := func(userInput, currentState NodeTestState, notify g.NotifyPartialFn[NodeTestState]) (NodeTestState, error) {
		currentState.Counter++
		return currentState, nil
	}
	routePolicy, _ := graph.RouterPolicyImplFactory[NodeTestState](graph.AnyRoute[NodeTestState])
	reducer := graph.Replacer[NodeTestState]

	opts := &g.NodeOptions[NodeTestState]{
		RoutingPolicy: routePolicy,
		Reducer:       reducer,
	}

	b.ResetTimer()
	b.ReportAllocs()

	for i := 0; i < b.N; i++ {
		_ = graph.NodeImplFactory[NodeTestState](
			g.IntermediateNode,
			"test-node",
			nodeFn,
			opts,
		)
	}
}

// BenchmarkNode_Accept tests the performance of node execution
func BenchmarkNode_Accept(b *testing.B) {
	nodeFn := func(userInput, currentState NodeTestState, notify g.NotifyPartialFn[NodeTestState]) (NodeTestState, error) {
		currentState.Counter = userInput.Counter + 1
		return currentState, nil
	}
	routePolicy, _ := graph.RouterPolicyImplFactory[NodeTestState](graph.AnyRoute[NodeTestState])
	reducer := graph.Replacer[NodeTestState]

	opts := &g.NodeOptions[NodeTestState]{
		RoutingPolicy: routePolicy,
		Reducer:       reducer,
	}

	node := graph.NodeImplFactory[NodeTestState](
		g.IntermediateNode,
		"bench-node",
		nodeFn,
		opts,
	)

	observer := newMockStateObserver(NodeTestState{Value: "initial", Counter: 0})
	userInput := NodeTestState{Value: "input", Counter: 5}

	// Drain notifications
	go func() {
		for range observer.notificationsCh {
		}
	}()

	b.ResetTimer()
	b.ReportAllocs()

	for i := 0; i < b.N; i++ {
		node.Accept(userInput, observer, g.DefaultInvokeConfig())
	}
}

// BenchmarkNode_SimpleExecution tests performance of simple node function execution
func BenchmarkNode_SimpleExecution(b *testing.B) {
	nodeFn := func(userInput, currentState NodeTestState, notify g.NotifyPartialFn[NodeTestState]) (NodeTestState, error) {
		currentState.Counter++
		return currentState, nil
	}
	routePolicy, _ := graph.RouterPolicyImplFactory[NodeTestState](graph.AnyRoute[NodeTestState])
	reducer := graph.Replacer[NodeTestState]

	opts := &g.NodeOptions[NodeTestState]{
		RoutingPolicy: routePolicy,
		Reducer:       reducer,
	}

	node := graph.NodeImplFactory[NodeTestState](
		g.IntermediateNode,
		"simple-node",
		nodeFn,
		opts,
	)

	observer := newMockStateObserver(NodeTestState{Value: "initial", Counter: 0})
	userInput := NodeTestState{Value: "input", Counter: 0}

	// Drain notifications in background
	done := make(chan bool)
	go func() {
		notificationCount := 0
		for range observer.notificationsCh {
			notificationCount++
			if notificationCount >= b.N {
				close(done)
				return
			}
		}
	}()

	b.ResetTimer()
	b.ReportAllocs()

	for i := 0; i < b.N; i++ {
		node.Accept(userInput, observer, g.DefaultInvokeConfig())
	}

	<-done
}

// BenchmarkNode_ComplexStateTransformation tests performance with complex state changes
func BenchmarkNode_ComplexStateTransformation(b *testing.B) {
	nodeFn := func(userInput, currentState NodeTestState, notify g.NotifyPartialFn[NodeTestState]) (NodeTestState, error) {
		// Simulate complex transformation
		currentState.Counter = (userInput.Counter * 2) + (currentState.Counter * 3)
		currentState.Value = userInput.Value + "-" + currentState.Value
		return currentState, nil
	}
	routePolicy, _ := graph.RouterPolicyImplFactory[NodeTestState](graph.AnyRoute[NodeTestState])
	reducer := graph.Replacer[NodeTestState]

	opts := &g.NodeOptions[NodeTestState]{
		RoutingPolicy: routePolicy,
		Reducer:       reducer,
	}

	node := graph.NodeImplFactory[NodeTestState](
		g.IntermediateNode,
		"complex-node",
		nodeFn,
		opts,
	)

	observer := newMockStateObserver(NodeTestState{Value: "state", Counter: 10})
	userInput := NodeTestState{Value: "input", Counter: 5}

	// Drain notifications
	go func() {
		for range observer.notificationsCh {
		}
	}()

	b.ResetTimer()
	b.ReportAllocs()

	for i := 0; i < b.N; i++ {
		node.Accept(userInput, observer, g.DefaultInvokeConfig())
	}
}

// BenchmarkNode_WithCustomReducer tests performance with custom reducer
func BenchmarkNode_WithCustomReducer(b *testing.B) {
	customReducer := func(currentState, change NodeTestState) NodeTestState {
		currentState.Counter += change.Counter
		if change.Value != "" {
			currentState.Value = change.Value
		}
		return currentState
	}

	nodeFn := func(userInput, currentState NodeTestState, notify g.NotifyPartialFn[NodeTestState]) (NodeTestState, error) {
		return NodeTestState{Counter: userInput.Counter, Value: "updated"}, nil
	}
	routePolicy, _ := graph.RouterPolicyImplFactory[NodeTestState](graph.AnyRoute[NodeTestState])

	opts := &g.NodeOptions[NodeTestState]{
		RoutingPolicy: routePolicy,
		Reducer:       customReducer,
	}

	node := graph.NodeImplFactory[NodeTestState](
		g.IntermediateNode,
		"reducer-node",
		nodeFn,
		opts,
	)

	observer := newMockStateObserver(NodeTestState{Value: "initial", Counter: 10})
	userInput := NodeTestState{Value: "input", Counter: 5}

	// Drain notifications
	go func() {
		for range observer.notificationsCh {
		}
	}()

	b.ResetTimer()
	b.ReportAllocs()

	for i := 0; i < b.N; i++ {
		node.Accept(userInput, observer, g.DefaultInvokeConfig())
	}
}

// BenchmarkNode_PartialUpdates tests performance with partial state updates
func BenchmarkNode_PartialUpdates(b *testing.B) {
	nodeFn := func(userInput, currentState NodeTestState, notify g.NotifyPartialFn[NodeTestState]) (NodeTestState, error) {
		// Send 3 partial updates
		for i := 1; i <= 3; i++ {
			currentState.Counter = i
			notify(currentState)
		}
		currentState.Counter = 10
		return currentState, nil
	}
	routePolicy, _ := graph.RouterPolicyImplFactory[NodeTestState](graph.AnyRoute[NodeTestState])
	reducer := graph.Replacer[NodeTestState]

	opts := &g.NodeOptions[NodeTestState]{
		RoutingPolicy: routePolicy,
		Reducer:       reducer,
	}

	node := graph.NodeImplFactory[NodeTestState](
		g.IntermediateNode,
		"partial-node",
		nodeFn,
		opts,
	)

	observer := newMockStateObserver(NodeTestState{Value: "initial", Counter: 0})
	userInput := NodeTestState{Value: "input", Counter: 0}

	// Drain notifications (4 per execution: 3 partial + 1 final)
	go func() {
		for range observer.notificationsCh {
		}
	}()

	b.ResetTimer()
	b.ReportAllocs()

	for i := 0; i < b.N; i++ {
		node.Accept(userInput, observer, g.DefaultInvokeConfig())
	}
}

// BenchmarkNode_Name tests performance of Name() method
func BenchmarkNode_Name(b *testing.B) {
	nodeFn := func(userInput, currentState NodeTestState, notify g.NotifyPartialFn[NodeTestState]) (NodeTestState, error) {
		return currentState, nil
	}
	routePolicy, _ := graph.RouterPolicyImplFactory[NodeTestState](graph.AnyRoute[NodeTestState])
	reducer := graph.Replacer[NodeTestState]

	opts := &g.NodeOptions[NodeTestState]{
		RoutingPolicy: routePolicy,
		Reducer:       reducer,
	}

	node := graph.NodeImplFactory[NodeTestState](
		g.IntermediateNode,
		"bench-name-node",
		nodeFn,
		opts,
	)

	b.ResetTimer()
	b.ReportAllocs()

	for i := 0; i < b.N; i++ {
		_ = node.Name()
	}
}

// BenchmarkNode_Role tests performance of Role() method
func BenchmarkNode_Role(b *testing.B) {
	nodeFn := func(userInput, currentState NodeTestState, notify g.NotifyPartialFn[NodeTestState]) (NodeTestState, error) {
		return currentState, nil
	}
	routePolicy, _ := graph.RouterPolicyImplFactory[NodeTestState](graph.AnyRoute[NodeTestState])
	reducer := graph.Replacer[NodeTestState]

	opts := &g.NodeOptions[NodeTestState]{
		RoutingPolicy: routePolicy,
		Reducer:       reducer,
	}

	node := graph.NodeImplFactory[NodeTestState](
		g.IntermediateNode,
		"bench-role-node",
		nodeFn,
		opts,
	)

	b.ResetTimer()
	b.ReportAllocs()

	for i := 0; i < b.N; i++ {
		_ = node.Role()
	}
}

// BenchmarkNode_RoutePolicy tests performance of RoutePolicy() method
func BenchmarkNode_RoutePolicy(b *testing.B) {
	nodeFn := func(userInput, currentState NodeTestState, notify g.NotifyPartialFn[NodeTestState]) (NodeTestState, error) {
		return currentState, nil
	}
	routePolicy, _ := graph.RouterPolicyImplFactory[NodeTestState](graph.AnyRoute[NodeTestState])
	reducer := graph.Replacer[NodeTestState]

	opts := &g.NodeOptions[NodeTestState]{
		RoutingPolicy: routePolicy,
		Reducer:       reducer,
	}

	node := graph.NodeImplFactory[NodeTestState](
		g.IntermediateNode,
		"bench-policy-node",
		nodeFn,
		opts,
	)

	b.ResetTimer()
	b.ReportAllocs()

	for i := 0; i < b.N; i++ {
		_ = node.RoutePolicy()
	}
}

// BenchmarkNode_ConcurrentExecution tests parallel node execution performance
func BenchmarkNode_ConcurrentExecution(b *testing.B) {
	nodeFn := func(userInput, currentState NodeTestState, notify g.NotifyPartialFn[NodeTestState]) (NodeTestState, error) {
		currentState.Counter = userInput.Counter + 1
		return currentState, nil
	}
	routePolicy, _ := graph.RouterPolicyImplFactory[NodeTestState](graph.AnyRoute[NodeTestState])
	reducer := graph.Replacer[NodeTestState]

	opts := &g.NodeOptions[NodeTestState]{
		RoutingPolicy: routePolicy,
		Reducer:       reducer,
	}

	node := graph.NodeImplFactory[NodeTestState](
		g.IntermediateNode,
		"concurrent-node",
		nodeFn,
		opts,
	)

	observer := newMockStateObserver(NodeTestState{Value: "initial", Counter: 0})

	// Drain notifications
	go func() {
		for range observer.notificationsCh {
		}
	}()

	b.ResetTimer()
	b.ReportAllocs()

	b.RunParallel(func(pb *testing.PB) {
		i := 0
		for pb.Next() {
			userInput := NodeTestState{Value: "input", Counter: i}
			node.Accept(userInput, observer, g.DefaultInvokeConfig())
			i++
		}
	})
}

// BenchmarkNode_DifferentRoles tests performance across different node roles
func BenchmarkNode_DifferentRoles(b *testing.B) {
	roles := map[string]g.NodeRole{
		"StartNode":        g.StartNode,
		"IntermediateNode": g.IntermediateNode,
		"EndNode":          g.EndNode,
	}

	for name, role := range roles {
		b.Run(name, func(b *testing.B) {
			nodeFn := func(userInput, currentState NodeTestState, notify g.NotifyPartialFn[NodeTestState]) (NodeTestState, error) {
				currentState.Counter++
				return currentState, nil
			}
			routePolicy, _ := graph.RouterPolicyImplFactory[NodeTestState](graph.AnyRoute[NodeTestState])
			reducer := graph.Replacer[NodeTestState]

			opts := &g.NodeOptions[NodeTestState]{
				RoutingPolicy: routePolicy,
				Reducer:       reducer,
			}

			node := graph.NodeImplFactory[NodeTestState](
				role,
				"bench-node",
				nodeFn,
				opts,
			)

			observer := newMockStateObserver(NodeTestState{Value: "initial", Counter: 0})
			userInput := NodeTestState{Value: "input", Counter: 0}

			// Drain notifications
			go func() {
				for range observer.notificationsCh {
				}
			}()

			b.ResetTimer()
			b.ReportAllocs()

			for i := 0; i < b.N; i++ {
				node.Accept(userInput, observer, g.DefaultInvokeConfig())
			}
		})
	}
}
