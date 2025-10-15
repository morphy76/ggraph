package graph

// Connected provides access to the connected graph components.
type Connected[T SharedState] interface {
	// AddEdge adds an edge to the runtime's graph.
	AddEdge(edge ...Edge[T])
	// Validate checks the integrity of the graph structure.
	Validate() error
}

// Runtime represents the runtime environment for graph processing.
type Runtime[T SharedState] interface {
	Connected[T]
	// Invoke executes the graph processing with the given user input.
	Invoke(userInput T)
	// TODO invoke with context
	// InvokeWithContext(ctx context.Context, entryState T)
	// Shutdown gracefully stops the runtime processing.
	Shutdown()
	// StartEdge returns the starting edge of the graph.
	StartEdge() Edge[T]
}
