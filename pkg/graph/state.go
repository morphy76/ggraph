package graph

// SharedState is the base interface for all state types used in graph processing.
//
// Any struct can implement SharedState by simply embedding it or using it as a type
// constraint. The state flows through the graph as nodes execute, with each node
// potentially transforming the state before passing it to the next node.
//
// Example:
//
//	type MyState struct {
//	    Counter int
//	    Data    string
//	    Results []string
//	}
//	var _ SharedState = (*MyState)(nil) // Verify MyState implements SharedState
type SharedState interface {
}

// StateMonitorEntry represents a single state transition event in the graph execution.
//
// Each time a node executes (or attempts to execute), a StateMonitorEntry is sent
// through the monitoring channel, providing visibility into the graph's execution flow.
// This enables debugging, logging, and tracking of how state evolves through the workflow.
//
// Fields:
//   - Node: The name of the node that just executed or attempted to execute.
//   - PreviousState: The state before the node's execution function ran.
//   - CurrentState: The state after the node's execution function completed.
//   - Error: Any error that occurred during node execution. nil if successful.
//   - Running: true while the graph is still executing, false when execution completes.
//   - Partial: true if this is a partial state update (from NotifyPartialFn), false
//     if this is the final state after node completion.
//
// Example usage:
//
//	stateMonitorCh := make(chan StateMonitorEntry[MyState], 10)
//	// ... create and invoke runtime ...
//	for entry := range stateMonitorCh {
//	    if entry.Error != nil {
//	        log.Printf("Error in node %s: %v", entry.Node, entry.Error)
//	    }
//	    if entry.Partial {
//	        log.Printf("Partial update from %s", entry.Node)
//	    }
//	    if !entry.Running {
//	        log.Printf("Graph execution completed at %s", entry.Node)
//	        break
//	    }
//	}
type StateMonitorEntry[T SharedState] struct {
	// Node is the name of the node that just executed or attempted to execute.
	Node string
	// NewState is the state after the node's execution function completed.
	NewState T
	// Error is any error that occurred during node execution. nil if successful.
	Error error
	// Running is true while the graph is still executing, false when execution completes.
	Running bool
	// Partial is true if this is a partial state update (from NotifyPartialFn), false
	// if this is the final state after node completion.
	Partial bool
	// ReducerFn is the function used to combine state updates.
	ReducerFn ReducerFn[T]
}
