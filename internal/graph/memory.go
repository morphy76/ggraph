package graph

import (
	"context"
	"sync"

	g "github.com/morphy76/ggraph/pkg/graph"
)

// MemMemoryFactory creates an in-memory Memory implementation.
func MemMemoryFactory[T g.SharedState]() g.Memory[T] {
	return &memMemory[T]{
		store: make(map[string]T),
		mu:    &sync.RWMutex{},
	}
}

// ------------------------------------------------------------------------------
// In-Memory Memory Implementation
// ------------------------------------------------------------------------------

var _ g.Memory[g.SharedState] = (*memMemory[g.SharedState])(nil)

type memMemory[T g.SharedState] struct {
	store map[string]T
	mu    *sync.RWMutex
}

func (m *memMemory[T]) PersistFn() g.PersistFn[T] {
	return func(ctx context.Context, key string, state T) error {
		m.mu.Lock()
		defer m.mu.Unlock()
		m.store[key] = state
		return nil
	}
}

func (m *memMemory[T]) RestoreFn() g.RestoreFn[T] {
	return func(ctx context.Context, key string) (T, error) {
		m.mu.RLock()
		defer m.mu.RUnlock()
		var zero T
		state, exists := m.store[key]
		if !exists {
			return zero, nil
		}
		return state, nil
	}
}
