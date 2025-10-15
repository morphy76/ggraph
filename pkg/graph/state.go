package graph

// SharedState is an interface for shared state management in graph processing.
type SharedState interface {
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
