package graph

// Memory interface defines methods for persisting and restoring shared state.
type Memory[T SharedState] interface {
	// PersistFn returns a function to persist the shared state.
	PersistFn() PersistFn[T]
	// RestoreFn returns a function to restore the shared state.
	RestoreFn() RestoreFn[T]
}
