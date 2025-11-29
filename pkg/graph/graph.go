package graph

import (
	"context"
)

// NotifyPartialFn is a callback function for sending partial state updates during node execution.
//
// When a node's processing takes significant time or produces incremental results,
// it can call this function to emit partial state updates before the final result.
// This enables streaming responses, progress tracking, and responsive user interfaces.
//
// The partial state updates are delivered through the state monitoring channel with
// the Partial flag set to true.
//
// Example usage within a node function:
//
//	func processData(userInput MyState, currentState MyState, notify NotifyPartialFn[MyState]) (MyState, error) {
//	    for i := 0; i < 10; i++ {
//	        currentState.Progress = i * 10
//	        notify(currentState) // Send partial update
//	        // ... do some work ...
//	    }
//	    currentState.Progress = 100
//	    return currentState, nil
//	}
type NotifyPartialFn[T SharedState] func(newState T)

// NodeFn is a function that implements the logic executed when a node is reached.
//
// This is the core processing function for operational nodes in the graph. It receives
// the original user input, the current state, and a notification function for partial
// updates. The function should process the state and return the updated state along
// with any error that occurred.
//
// Parameters:
//   - userInput: The original input provided to Runtime.Invoke(), unchanged throughout execution.
//   - currentState: The current state at the time this node executes, potentially modified
//     by previous nodes in the execution path.
//   - notify: A callback function to send partial state updates during processing.
//
// Returns:
//   - The updated state after processing.
//   - An error if processing failed, which will halt graph execution.
//
// Example:
//
//	func myNodeFunction(userInput MyState, currentState MyState, notify NotifyPartialFn[MyState]) (MyState, error) {
//	    currentState.Counter++
//	    currentState.Message = fmt.Sprintf("Processed: %s", userInput.Request)
//	    return currentState, nil
//	}
type NodeFn[T SharedState] func(userInput, currentState T, notify NotifyPartialFn[T]) (T, error)

// EdgeSelectionFn is a function that determines which edge to follow during graph execution.
//
// This function implements the routing logic for conditional branching, loops, and
// dynamic path selection. It examines the user input, current state, and available
// outgoing edges to decide which edge the execution should follow next.
//
// Parameters:
//   - userInput: The original input provided to Runtime.Invoke().
//   - currentState: The current state at the time of routing decision.
//   - edges: All available outgoing edges from the current node.
//
// Returns:
//   - The Edge to follow. Must be one of the edges from the provided slice.
//
// Example:
//
//	func routingLogic(userInput MyState, currentState MyState, edges []Edge[MyState]) Edge[MyState] {
//	    if currentState.Counter > 10 {
//	        return edges[0] // Exit edge
//	    }
//	    return edges[1] // Loop back edge
//	}
type EdgeSelectionFn[T SharedState] func(userInput, currentState T, edges []Edge[T]) Edge[T]

// PersistFn is a function that saves the current state of a thread to persistent storage.
//
// This function is called by the runtime to persist state to external storage
// (database, file system, etc.), enabling stateful workflows that can be resumed
// after interruption or across multiple execution sessions.
//
// Parameters:
//   - ctx: Context for cancellation and timeout control.
//   - threadID: Unique identifier for this thread instance.
//   - state: The current state to persist.
//
// Returns:
//   - An error if the persistence operation fails.
//
// Example:
//
//	func persistState(ctx context.Context, threadID string, state MyState) error {
//	    data, err := json.Marshal(state)
//	    if err != nil {
//	        return err
//	    }
//	    return db.Save(ctx, threadID, data)
//	}
type PersistFn[T SharedState] func(ctx context.Context, threadID string, state T) error

// RestoreFn is a function that retrieves previously persisted state from storage.
//
// This function is called by the runtime to restore state from external storage,
// allowing workflows to resume from where they left off. It should return the
// most recently persisted state for the given runtime ID.
//
// Parameters:
//   - ctx: Context for cancellation and timeout control.
//   - runtimeID: Unique identifier for the runtime instance to restore.
//
// Returns:
//   - The restored state.
//   - An error if the restoration operation fails.
//
// Example:
//
//	func restoreState(ctx context.Context, runtimeID uuid.UUID) (MyState, error) {
//	    data, err := db.Load(ctx, runtimeID.String())
//	    if err != nil {
//	        return MyState{}, err
//	    }
//	    var state MyState
//	    err = json.Unmarshal(data, &state)
//	    return state, err
//	}
type RestoreFn[T SharedState] func(ctx context.Context, threadID string) (T, error)

// ReducerFn is a function that combines multiple states into a single state.
//
// This function is used to merge state changes from different nodes or execution
// paths. It takes the current state and a new state change, and produces a combined
// state that reflects both.
//
// Parameters:
//   - currentState: The existing state before applying the change.
//   - change: The new state change to incorporate.
//
// Returns:
//   - The combined state after applying the change.
//
// Example:
//
//	func mergeStates(currentState MyState, change MyState) MyState {
//	    currentState.Counter += change.Counter
//	    if change.Message != "" {
//	        currentState.Message = change.Message
//	    }
//	    return currentState
//	}
type ReducerFn[T SharedState] func(currentState, change T) T

// StateObserver is an internal interface for tracking state changes during graph execution.
//
// This interface is primarily used by the runtime to monitor and record state transitions
// as nodes execute. It provides the mechanism for the state monitoring channel to receive
// updates about graph execution progress.
//
// Most users will not need to implement this interface directly; it is used internally
// by the runtime implementation.
type StateObserver[T SharedState] interface {
	// NotifyStateChange is called whenever a node modifies the state.
	//
	// Parameters:
	//   - node: The node that produced the state change.
	//   - config: The configuration settings for the invocation.
	//   - userInput: The original user input to the graph.
	//   - stateChange: The new state after the node's execution.
	//   - reducer: The reducer function used to combine states.
	//   - err: Any error that occurred during node execution.
	//   - partial: true if this is a partial update, false if final.
	NotifyStateChange(node Node[T], config InvokeConfig, userInput, stateChange T, reducer ReducerFn[T], err error, partial bool)

	// CurrentState returns the current state for the given thread ID.
	//
	// Parameters:
	//   - threadID: Unique identifier for the thread instance.
	//
	// Returns:
	//   - The current state associated with the specified thread ID.
	CurrentState(threadID string) T

	// InitialState returns the initial state used at the start of execution.
	//
	// Returns:
	//   - The initial state provided to the runtime.
	InitialState() T
}

// Persistent is an interface for managing state persistence in graph workflows.
//
// This interface enables stateful graph execution that can survive process restarts,
// handle long-running workflows, and implement checkpoint/resume patterns. It is
// embedded in the Runtime interface to provide persistence capabilities.
type Persistent[T SharedState] interface {
	// Restore loads and applies previously persisted state to the runtime.
	//
	// This method should be called before Invoke() to restore the state from
	// a previous execution session. It uses the RestoreFn provided in
	// SetPersistentState to retrieve the state from storage.
	//
	// Parameters:
	//   - threadID: Unique identifier for the thread instance to restore.
	//
	// Returns:
	//   - An error if the restoration fails.
	//
	// Example:
	//
	//	if err := runtime.Restore(); err != nil {
	//	    log.Printf("Failed to restore state: %v", err)
	//	}
	//	runtime.Invoke(userInput)
	Restore(threadID string) error
}

// Threaded is an interface for retrieving active thread identifiers in a runtime.
//
// This interface allows users to query the runtime for a list of currently
// active thread IDs, enabling management and monitoring of concurrent
// graph executions.
type Threaded interface {
	// ListThreads returns a slice of active thread IDs.
	//
	// Returns:
	//   - A slice of strings representing the active thread identifiers.
	//
	// Example:
	//
	//	threads := runtime.ListThreads()
	//	for _, threadID := range threads {
	//	    fmt.Println("Active thread:", threadID)
	//	}
	ListThreads() []string
}
