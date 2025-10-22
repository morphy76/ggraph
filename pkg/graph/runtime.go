package graph

// Connected provides methods for building and validating the graph structure.
//
// This interface allows you to add edges to construct the graph topology and
// validate that the resulting structure is correct and executable. It is embedded
// in the Runtime interface to provide graph construction capabilities.
type Connected[T SharedState] interface {
	// AddEdge adds one or more edges to the graph structure.
	//
	// Edges define the connections between nodes and determine how execution flows
	// through the graph. Multiple edges can be added in a single call for convenience.
	// Edges should be added after runtime creation but before calling Invoke().
	//
	// All edges except the StartEdge (which is provided during runtime creation)
	// must be added using this method before the graph can execute.
	//
	// Parameters:
	//   - edge: One or more Edge instances to add to the graph.
	//
	// Example:
	//
	//	runtime.AddEdge(edge1)
	//	runtime.AddEdge(edge2, edge3, edge4) // Multiple edges at once
	AddEdge(edge ...Edge[T])

	// Validate checks the integrity and correctness of the graph structure.
	//
	// This method performs structural validation to ensure the graph is properly
	// formed and executable. It checks for:
	//   - Exactly one StartEdge exists
	//   - At least one EndEdge exists
	//   - All nodes (except EndNode) have at least one outgoing edge
	//   - No unreachable nodes or edges exist
	//   - Graph topology is valid
	//
	// It is recommended to call Validate() after adding all edges and before
	// invoking the graph to catch configuration errors early.
	//
	// Returns:
	//   - nil if the graph structure is valid and executable.
	//   - An error describing the validation failure if the graph is invalid.
	//
	// Example:
	//
	//	runtime.AddEdge(edge1, edge2, edge3)
	//	if err := runtime.Validate(); err != nil {
	//	    log.Fatalf("Invalid graph: %v", err)
	//	}
	//	runtime.Invoke(userInput)
	Validate() error
}

// Runtime represents the execution engine for graph-based workflows.
//
// The Runtime is the central component that manages graph execution. It:
//   - Maintains the graph structure (nodes and edges)
//   - Executes nodes in the correct order based on routing policies
//   - Manages state transitions throughout execution
//   - Provides monitoring and observability through state channels
//   - Supports persistence for stateful, resumable workflows
//
// A typical workflow:
//  1. Create runtime with builders.CreateRuntime() or builders.CreateRuntimeWithInitialState()
//  2. Add edges using AddEdge()
//  3. Validate the graph structure with Validate()
//  4. (Optional) Configure persistence with SetPersistentState() and Restore()
//  5. Execute the graph with Invoke()
//  6. Monitor execution through the state monitoring channel
//  7. Shutdown gracefully when done
//
// Example:
//
//	stateMonitorCh := make(chan StateMonitorEntry[MyState], 10)
//	runtime, _ := builders.CreateRuntime(startEdge, stateMonitorCh)
//	defer runtime.Shutdown()
//
//	runtime.AddEdge(edge1, edge2, edge3)
//	runtime.Validate()
//	runtime.Invoke(MyState{Request: "process"})
//
//	for entry := range stateMonitorCh {
//	    fmt.Printf("Node %s executed\n", entry.Node)
//	    if !entry.Running {
//	        break
//	    }
//	}
type Runtime[T SharedState] interface {
	// Embeds Connected to provide graph construction and validation methods.
	Connected[T]

	// Embeds Persistent to provide state persistence capabilities.
	Persistent[T]

	// Invoke starts the graph execution with the provided user input.
	//
	// This method initiates the graph workflow by traversing the StartEdge to
	// reach the first operational node. Execution continues node-by-node,
	// following edges determined by routing policies, until an EndEdge is reached.
	//
	// The userInput is passed to every node's NodeFn and every routing policy's
	// SelectEdge method, remaining unchanged throughout execution. This allows
	// nodes to access the original request while working with the evolving
	// currentState.
	//
	// Invoke runs asynchronously. Monitor execution progress through the state
	// monitoring channel provided during runtime creation.
	//
	// Parameters:
	//   - userInput: The input state to process. This is passed to all nodes and
	//     routing policies but is never modified by the runtime.
	//
	// Example:
	//
	//	userInput := MyState{
	//	    Request: "process data",
	//	    Config:  myConfig,
	//	}
	//	runtime.Invoke(userInput)
	//
	//	// Monitor execution
	//	for entry := range stateMonitorCh {
	//	    if entry.Error != nil {
	//	        log.Printf("Error: %v", entry.Error)
	//	    }
	//	    if !entry.Running {
	//	        fmt.Printf("Final state: %+v\n", entry.CurrentState)
	//	        break
	//	    }
	//	}
	Invoke(userInput T)
	// TODO invoke with context
	// InvokeWithContext(ctx context.Context, entryState T)

	// Shutdown gracefully stops the runtime and cleans up resources.
	//
	// This method should be called when the runtime is no longer needed, typically
	// using defer immediately after runtime creation. It ensures that:
	//   - The state monitoring channel is properly closed
	//   - Internal goroutines are terminated
	//   - Resources are released
	//
	// After calling Shutdown(), the runtime cannot be used again.
	//
	// Example:
	//
	//	runtime, _ := builders.CreateRuntime(startEdge, stateMonitorCh)
	//	defer runtime.Shutdown() // Ensure cleanup
	//
	//	runtime.AddEdge(edges...)
	//	runtime.Invoke(input)
	Shutdown()

	// StartEdge returns the entry edge of the graph.
	//
	// This is the edge that was provided during runtime creation and defines
	// where graph execution begins. The StartEdge connects the implicit start
	// node to the first operational node.
	//
	// Returns:
	//   - The StartEdge of this runtime.
	StartEdge() Edge[T]

	// CurrentState returns the current state of the graph execution.
	//
	// This method provides access to the latest state after the most recent
	// node execution. It is safe for concurrent use and reflects the state
	// as it evolves through the graph.
	//
	// Returns:
	//   - The current state of type T.
	CurrentState() T
}
