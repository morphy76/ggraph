package graph

// StateMergeFn defines a function type for merging two SharedState instances.
type StateMergeFn[T SharedState] func(current, other T) T

// SharedState is an interface for shared state management in graph processing.
type SharedState interface {
}

// GraphRunning creates a StateMonitorEntry indicating progress in the graph processing.
func GraphRunning[T SharedState](node string, previousState, currentState T) StateMonitorEntry[T] {
	return StateMonitorEntry[T]{
		Node:          node,
		PreviousState: previousState,
		CurrentState:  currentState,
		Running:       true,
		Partial:       false,
	}
}

// GraphError creates a StateMonitorEntry indicating an error in the graph processing.
func GraphError[T SharedState](node string, currentState T, err error) StateMonitorEntry[T] {
	return StateMonitorEntry[T]{
		Node:          node,
		PreviousState: currentState,
		CurrentState:  currentState,
		Error:         err,
		Running:       false,
		Partial:       false,
	}
}

func GraphPartial[T SharedState](node string, currentState T) StateMonitorEntry[T] {
	return StateMonitorEntry[T]{
		Node:          node,
		PreviousState: currentState,
		CurrentState:  currentState,
		Running:       true,
		Partial:       true,
	}
}

// GraphCompleted creates a StateMonitorEntry indicating the completion of the graph processing.
func GraphCompleted[T SharedState](finalState T) StateMonitorEntry[T] {
	return StateMonitorEntry[T]{
		CurrentState: finalState,
		Running:      false,
		Partial:      false,
	}
}

// StateMonitorEntry represents an entry in the state monitoring log.
type StateMonitorEntry[T SharedState] struct {
	Node          string
	PreviousState T
	CurrentState  T
	Error         error
	Running       bool
	Partial       bool
}
