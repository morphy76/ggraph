package builders

import (
	i "github.com/morphy76/ggraph/internal/graph"
	g "github.com/morphy76/ggraph/pkg/graph"
)

// NewMemMemory creates a new in-memory Memory implementation.
//
// Parameters:
//   - opts ...g.MemoryOption: Optional memory configuration options.
//
// Returns:
//   - g.Memory[T]: In-memory Memory implementation.
func NewMemMemory[T g.SharedState](opts ...g.MemoryOption) g.Memory[T] {
	useOpts := &g.MemoryOptions{}
	for _, opt := range opts {
		if err := opt.Apply(useOpts); err != nil {
			panic(err)
		}
	}
	return i.MemMemoryFactory[T](useOpts)
}
