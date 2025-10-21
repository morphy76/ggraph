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
) (g.Runtime[T], error) {
	var zero T
	return CreateRuntimeWithInitialState(startEdge, stateMonitorCh, zero)
}

// CreateRuntimeWithInitialState creates a new graph runtime with a custom initial state.
//
// This function creates a runtime with a pre-populated initial state, which is useful
// when you want to establish default values, configuration, or context that should be
// available before the first node executes. The initial state is merged with or used
// alongside the userInput provided to the Invoke method.
//
// The runtime manages the complete lifecycle of graph execution:
//   - Maintains the current state throughout execution
//   - Routes execution through nodes based on edges and routing policies
//   - Monitors and reports state changes via the monitoring channel
//   - Handles errors and ensures graceful shutdown
//
// Type Parameters:
//   - T: The SharedState type that will be passed through the graph execution.
//
// Parameters:
//   - startEdge: The edge that connects to the first operational node. This defines
//     the entry point of the graph workflow.
//   - stateMonitorCh: A buffered channel that receives state monitoring entries during
//     execution. Each entry contains the node name, previous state, current state,
//     error (if any), and execution status flags.
//   - initialState: The starting state value for the graph. This state is available
//     when the first node begins execution.
//
// Returns:
//   - A new Runtime instance configured with the initial state.
//   - An error if the runtime cannot be created.
//
// Example:
//
//	initialState := MyState{
//	    Config: "production",
//	    Counter: 0,
//	    Data: make(map[string]string),
//	}
//	stateMonitorCh := make(chan StateMonitorEntry[MyState], 10)
//	runtime, err := CreateRuntimeWithInitialState(startEdge, stateMonitorCh, initialState)
//	if err != nil {
//	    log.Fatal(err)
//	}
//	defer runtime.Shutdown()
//
//	runtime.AddEdge(edge1, edge2)
//	if err := runtime.Validate(); err != nil {
//	    log.Fatal(err)
//	}
//
//	runtime.Invoke(MyState{UserInput: "request data"})
//
//	// Monitor execution
//	for entry := range stateMonitorCh {
//	    fmt.Printf("Node: %s, Running: %v\n", entry.Node, entry.Running)
//	    if !entry.Running {
//	        break
//	    }
//	}
func CreateRuntimeWithInitialState[T g.SharedState](
	startEdge g.Edge[T],
	stateMonitorCh chan g.StateMonitorEntry[T],
	initialState T,
) (g.Runtime[T], error) {
	return i.RuntimeFactory(startEdge, stateMonitorCh, initialState)
}
