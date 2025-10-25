package builders

import (
	i "github.com/morphy76/ggraph/internal/graph"
	g "github.com/morphy76/ggraph/pkg/graph"
)

// NewMemMemory creates a new in-memory Memory implementation.
//
// Returns:
//   - g.Memory[T]: In-memory Memory implementation.
func NewMemMemory[T g.SharedState]() g.Memory[T] {
	return i.MemMemoryFactory[T]()
}
