package graph

// RuntimeOptions holds the configuration for a node.
type RuntimeOptions[T SharedState] struct {
	InitialState T
	Memory       Memory[T]

	WorkerCount     int
	WorkerQueueSize int
}

// RuntimeOption is a functional option for configuring a graph runtime.
type RuntimeOption[T SharedState] interface {
	// Apply applies the option to the RuntimeOptions.
	//
	// Parameters:
	//   - r: A pointer to RuntimeOptions to modify.
	//
	// Returns:
	//   - An error if the application of the option fails, otherwise nil.
	Apply(r *RuntimeOptions[T]) error
}

// RuntimeOptionFunc is a function type that implements the RuntimeOption interface.
type RuntimeOptionFunc[T SharedState] func(*RuntimeOptions[T]) error

// Apply applies the RuntimeOptionFunc to the given RuntimeOptions.
//
// Parameters:
//   - r: A pointer to RuntimeOptions to modify.
//
// Returns:
//   - An error if the application of the option fails, otherwise nil.
func (s RuntimeOptionFunc[T]) Apply(r *RuntimeOptions[T]) error { return s(r) }

// WithInitialState sets the initial shared state for the graph runtime.
//
// Parameters:
//   - initialState: The initial shared state of type T.
//
// Returns:
//   - A RuntimeOption that sets the initial state.
//
// Example:
//
//	initialState := MyState{Value: 100}
//	runtime, err := builders.CreateRuntimeWithInitialState(startEdge, stateMonitorCh, initialState)
func WithInitialState[T SharedState](initialState T) RuntimeOption[T] {
	return RuntimeOptionFunc[T](func(r *RuntimeOptions[T]) error {
		r.InitialState = initialState
		return nil
	})
}

// WithMemory sets the memory component for the graph runtime.
//
// Parameters:
//   - memory: An instance of Memory[T] to be used by the runtime.
//
// Returns:
//   - A RuntimeOption that sets the memory component.
//
// Example:
//
//	memory := NewInMemoryStorage[MyState]()
//	runtime, err := builders.CreateRuntimeWithMemory(startEdge, stateMonitorCh, memory)
func WithMemory[T SharedState](memory Memory[T]) RuntimeOption[T] {
	return RuntimeOptionFunc[T](func(r *RuntimeOptions[T]) error {
		r.Memory = memory
		return nil
	})
}

// WithWorkerPool configures the worker pool for the graph runtime.
//
// Parameters:
//   - workerCount: The number of workers in the pool.
//   - workerQueueSize: The size of the task queue for the worker pool.
//
// Returns:
//   - A RuntimeOption that sets the worker pool configuration.
//
// Example:
//
//	runtime, err := builders.CreateRuntime(startEdge, stateMonitorCh, WithWorkerPool(10, 200))
func WithWorkerPool[T SharedState](workerCount int, workerQueueSize int) RuntimeOption[T] {
	return RuntimeOptionFunc[T](func(r *RuntimeOptions[T]) error {
		r.WorkerCount = workerCount
		r.WorkerQueueSize = workerQueueSize
		return nil
	})
}

// TODO pluggable log
// TODO observability hooks
