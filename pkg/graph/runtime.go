package graph

import (
	"context"
	"errors"

	"github.com/google/uuid"
)

var (
	// ErrRuntimeExecuting indicates that the runtime is already executing and cannot accept another invocation.
	ErrRuntimeExecuting = errors.New("runtime is already executing")
	// ErrStartEdgeNil indicates that the provided start edge is nil.
	ErrStartEdgeNil = errors.New("start edge cannot be nil")
	// ErrNoPathToEnd indicates that there is no path from the start edge to any end edge.
	ErrNoPathToEnd = errors.New("no path from start edge to any end edge")
	// ErrRestoreNotSet indicates that the restore function is not set.
	ErrRestoreNotSet = errors.New("restore function is not set")
	// ErrPersistenceQueueFull indicates that the persistence queue is full.
	ErrPersistenceQueueFull = errors.New("persistence queue is full")
	// ErrEvictionByInactivity indicates that a thread was evicted due to inactivity.
	ErrEvictionByInactivity = errors.New("thread evicted due to inactivity")
	// ErrUnknownThreadID indicates that the provided thread ID is unknown.
	ErrUnknownThreadID = errors.New("unknown thread ID")
	// ErrRuntimeOptionsNil indicates that the provided runtime options are nil.
	ErrRuntimeOptionsNil = errors.New("runtime options cannot be nil")
)

// NodeExecutor defines an interface for submitting tasks to be executed.
type NodeExecutor interface {
	// Submit adds a task to be executed.
	//
	// Parameters:
	//   - task: A function representing the task to execute.
	//
	// Example:
	//
	//	executor.Submit(func() {
	//	    // Task logic here
	//	})
	Submit(task func())
}

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

// InvokeConfig holds configuration options for invoking the runtime.
//
// This struct allows you to specify parameters that influence the execution
// of the graph during an invocation. Currently, it supports setting a ThreadID
// to enable concurrent executions within the same runtime instance.
//
// Example:
//
//	config := InvokeConfig{
//	    ThreadID: "thread-123",
//	}
//	runtime.Invoke(userInput, config)
type InvokeConfig struct {
	// ThreadID is the identifier for the thread executing the invocation.
	ThreadID string
	// Context is the context for the invocation.
	Context context.Context
}

// MergeInvokeConfig merges multiple InvokeConfig instances into one.
//
// When merging, non-empty fields from later configurations override those
// from earlier ones. This allows for flexible configuration by combining
// multiple sources.
//
// Parameters:
//   - config: One or more InvokeConfig instances to merge.
//
// Returns:
//   - A single InvokeConfig instance containing the merged settings.
//
// Example:
//
//	baseConfig := InvokeConfig{ThreadID: "base-thread"}
//	overrideConfig := InvokeConfig{ThreadID: "override-thread"}
//
//	mergedConfig := MergeInvokeConfig(baseConfig, overrideConfig)
//	// mergedConfig.ThreadID will be "override-thread"
func MergeInvokeConfig(config ...InvokeConfig) InvokeConfig {
	merged := InvokeConfig{}
	for _, c := range config {
		if c.ThreadID != "" {
			merged.ThreadID = c.ThreadID
		}
		if c.Context != nil {
			merged.Context = c.Context
		}
	}
	return merged
}

// DefaultInvokeConfig creates an InvokeConfig with default settings.
//
// The default configuration generates a unique ThreadID using a UUID.
// This ensures that each invocation has its own distinct thread context.
//
// Returns:
//   - An InvokeConfig instance with default settings.
//
// Example:
//
//	defaultConfig := DefaultInvokeConfig()
//	runtime.Invoke(userInput, defaultConfig)
func DefaultInvokeConfig() InvokeConfig {
	return InvokeConfig{
		ThreadID: uuid.NewString(),
		Context:  context.TODO(),
	}
}

// InvokeConfigThreadID creates an InvokeConfig with the specified ThreadID.
//
// This helper function simplifies the creation of an InvokeConfig when only
// the ThreadID needs to be set.
// // Parameters:
//   - threadID: The identifier for the thread executing the invocation.
//
// Returns:
//   - An InvokeConfig instance with the specified ThreadID.
//
// Example:
//
//	threadConfig := InvokeConfigThreadID("custom-thread-1")
//	runtime.Invoke(userInput, threadConfig)
func InvokeConfigThreadID(threadID string) InvokeConfig {
	return InvokeConfig{ThreadID: threadID}
}

// InvokeConfigContext creates an InvokeConfig with the specified Context.
//
// This helper function simplifies the creation of an InvokeConfig when only
// the Context needs to be set.
//
// Parameters:
//   - ctx: The context for the invocation.
//
// Returns:
//   - An InvokeConfig instance with the specified Context.
//
// Example:
//
//	ctx := context.WithValue(context.Background(), "key", "value")
//	contextConfig := InvokeConfigContext(ctx)
//	runtime.Invoke(userInput, contextConfig)
func InvokeConfigContext(ctx context.Context) InvokeConfig {
	return InvokeConfig{Context: ctx}
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
	// Embeds StateObserver to provide state change notification capabilities.
	StateObserver[T]

	// Embeds Connected to provide graph construction and validation methods.
	Connected[T]

	// Embeds Persistent to provide state persistence capabilities.
	Persistent[T]

	// Embeds Threaded to provide active thread retrieval capabilities.
	Threaded

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
	//   - config: Optional configuration settings for this invocation.
	//
	// Returns:
	//   - The ThreadID used for this invocation, which can be specified in the
	//     InvokeConfig.
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
	Invoke(userInput T, config ...InvokeConfig) string

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
}
