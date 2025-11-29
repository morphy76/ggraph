package builders

import (
	i "github.com/morphy76/ggraph/internal/graph"
	g "github.com/morphy76/ggraph/pkg/graph"
)

// CreateRuntime creates a new graph runtime with a zero-value initial state.
//
// The runtime is the execution engine that processes the graph workflow. It manages
// node execution, state transitions, edge traversal, and provides monitoring capabilities
// through a state monitoring channel. This function initializes the runtime with a
// zero-value state of type T, which is suitable when the initial state will be provided
// via the userInput parameter in the Invoke method.
//
// The state monitor channel receives entries after each node execution, providing
// visibility into the graph execution flow, state changes, errors, and completion status.
//
// Type Parameters:
//   - T: The SharedState type that will be passed through the graph execution.
//
// Parameters:
//   - startEdge: The edge that connects to the first operational node. This defines
//     the entry point of the graph workflow.
//   - stateMonitorCh: A buffered channel that receives state monitoring entries during
//     execution. Use a buffer size appropriate for your graph complexity (e.g., 10-100).
//   - opts: Optional configuration options for the runtime.
//
// Returns:
//   - A new Runtime instance ready to execute the graph workflow.
//   - An error if the runtime cannot be created.
//
// Example:
//
//	startEdge := CreateStartEdge(firstNode)
//	stateMonitorCh := make(chan StateMonitorEntry[MyState], 10)
//	runtime, err := CreateRuntime(startEdge, stateMonitorCh)
//	if err != nil {
//	    log.Fatal(err)
//	}
//	defer runtime.Shutdown()
//
//	runtime.AddEdge(edge1, edge2, edge3)
//	runtime.Validate()
//	runtime.Invoke(MyState{Value: 42})
func CreateRuntime[T g.SharedState](
	startEdge g.Edge[T],
	stateMonitorCh chan g.StateMonitorEntry[T],
	opts ...g.RuntimeOption[T],
) (g.Runtime[T], error) {

	var zeroState T
	useOpts := &g.RuntimeOptions[T]{
		InitialState: zeroState,
		Settings:     g.RuntimeSettings{},
	}
	for _, opt := range opts {
		opt.Apply(useOpts)
	}

	return i.RuntimeFactory(startEdge, stateMonitorCh, useOpts)
}
