package graph

import (
	g "github.com/morphy76/ggraph/pkg/graph"
)

// GraphRunning creates a StateMonitorEntry indicating progress in the graph processing.
func GraphRunning[T g.SharedState](node string, previousState, currentState T) g.StateMonitorEntry[T] {
	return g.StateMonitorEntry[T]{
		Node:          node,
		PreviousState: previousState,
		CurrentState:  currentState,
		Running:       true,
		Partial:       false,
	}
}

// GraphError creates a StateMonitorEntry indicating an error in the graph processing.
func GraphError[T g.SharedState](node string, currentState T, err error) g.StateMonitorEntry[T] {
	return g.StateMonitorEntry[T]{
		Node:          node,
		PreviousState: currentState,
		CurrentState:  currentState,
		Error:         err,
		Running:       false,
		Partial:       false,
	}
}

// GraphPartial creates a StateMonitorEntry indicating partial progress in the graph processing.
func GraphPartial[T g.SharedState](node string, currentState T) g.StateMonitorEntry[T] {
	return g.StateMonitorEntry[T]{
		Node:          node,
		PreviousState: currentState,
		CurrentState:  currentState,
		Running:       true,
		Partial:       true,
	}
}

// GraphCompleted creates a StateMonitorEntry indicating the completion of the graph processing.
func GraphCompleted[T g.SharedState](node string, finalState T) g.StateMonitorEntry[T] {
	return g.StateMonitorEntry[T]{
		Node:         node,
		CurrentState: finalState,
		Running:      false,
		Partial:      false,
	}
}
