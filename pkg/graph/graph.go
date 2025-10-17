package graph

import (
	"context"

	"github.com/google/uuid"
)

// NotifyPartialFn defines a function type for notifying partial state updates.
type NotifyPartialFn[T SharedState] func(newState T)

// NodeFn defines a function type that processes a node with the given SharedState type.
type NodeFn[T SharedState] func(userInput T, currentState T, notify NotifyPartialFn[T]) (T, error)

// EdgeSelectionFn defines a function type for conditional routing based on the current state and available edges.
type EdgeSelectionFn[T SharedState] func(userInput T, currentState T, edges []Edge[T]) Edge[T]

// StateObserver is an interface for observing state changes in nodes during graph processing.
type StateObserver[T SharedState] interface {
	// NotifyStateChange is called when a node changes state during processing.
	NotifyStateChange(node Node[T], userInput T, newState T, err error, partial bool)
	// CurrentState returns the current state of the observer.
	CurrentState() T
}

type PersistFn[T SharedState] func(ctx context.Context, runtimeID uuid.UUID, state T) error

type RestoreFn[T SharedState] func(ctx context.Context, runtimeID uuid.UUID) (T, error)

// Persistent is an interface for managing persistent state in the graph runtime.
type Persistent[T SharedState] interface {
	// SetPersistentState sets up the persistent state management with the provided writer, reader, and runtime ID.
	SetPersistentState(
		persist PersistFn[T],
		restore RestoreFn[T],
		runtimeID uuid.UUID,
	)

	// Restore restores the state from persistent storage.
	Restore() error
}
