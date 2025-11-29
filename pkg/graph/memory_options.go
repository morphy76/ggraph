package graph

// MemoryOptions defines configuration options for memory implementations.
type MemoryOptions struct {
}

// MemoryOption defines an interface for applying configuration options to MemoryOptions.
type MemoryOption interface {
	// Apply applies the option to the MemoryOptions.
	//
	// Parameters:
	//   - r: A pointer to MemoryOptions to modify.
	//
	// Returns:
	//   - An error if the application of the option fails, otherwise nil.
	Apply(r *MemoryOptions) error
}

// MemoryOptionFunc is a function type that implements the MemoryOption interface.
type MemoryOptionFunc func(*MemoryOptions) error

// TODO options to limit the memory size and eviction policies
