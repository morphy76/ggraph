package graph

// NodeFunc defines a function type that processes a node with the given SharedState type.
type NodeFunc[T SharedState] func(state T, notify func(T)) (T, error)

// EdgeSelectionFn defines a function type for conditional routing based on the current state and available edges.
type EdgeSelectionFn[T SharedState] func(state T, edges []Edge[T]) Edge[T]

// StateMergeFn defines a function type for merging two SharedState instances.
type StateMergeFn[T SharedState] func(current, other T) T

// StateObserver is an interface for observing state changes in nodes during graph processing.
type StateObserver[T SharedState] interface {
	// NotifyStateChange is called when a node changes state during processing.
	NotifyStateChange(node Node[T], state T, err error, partial bool)
}
