package graph_test

import (
	"context"
	"testing"

	"github.com/morphy76/ggraph/internal/graph"
	g "github.com/morphy76/ggraph/pkg/graph"
)

// MemoryTestState is a state type for memory testing
type MemoryTestState struct {
	Value   string
	Counter int
	Data    map[string]interface{}
}

func TestMemMemoryFactory(t *testing.T) {
	t.Run("creates memory instance", func(t *testing.T) {
		memory := graph.MemMemoryFactory[MemoryTestState]()
		if memory == nil {
			t.Fatal("MemMemoryFactory returned nil")
		}
	})

	t.Run("implements Memory interface", func(t *testing.T) {
		memory := graph.MemMemoryFactory[MemoryTestState]()
		var _ g.Memory[MemoryTestState] = memory
	})

	t.Run("returns valid PersistFn", func(t *testing.T) {
		memory := graph.MemMemoryFactory[MemoryTestState]()
		persistFn := memory.PersistFn()
		if persistFn == nil {
			t.Fatal("PersistFn returned nil")
		}
	})

	t.Run("returns valid RestoreFn", func(t *testing.T) {
		memory := graph.MemMemoryFactory[MemoryTestState]()
		restoreFn := memory.RestoreFn()
		if restoreFn == nil {
			t.Fatal("RestoreFn returned nil")
		}
	})
}

func TestMemMemory_PersistAndRestore(t *testing.T) {
	t.Run("persists and restores simple state", func(t *testing.T) {
		memory := graph.MemMemoryFactory[MemoryTestState]()
		ctx := context.Background()
		key := "test-key"

		state := MemoryTestState{
			Value:   "hello",
			Counter: 42,
		}

		// Persist state
		persistFn := memory.PersistFn()
		err := persistFn(ctx, key, state)
		if err != nil {
			t.Fatalf("PersistFn failed: %v", err)
		}

		// Restore state
		restoreFn := memory.RestoreFn()
		restored, err := restoreFn(ctx, key)
		if err != nil {
			t.Fatalf("RestoreFn failed: %v", err)
		}

		// Verify restored state
		if restored.Value != state.Value {
			t.Errorf("Value mismatch: got %q, want %q", restored.Value, state.Value)
		}
		if restored.Counter != state.Counter {
			t.Errorf("Counter mismatch: got %d, want %d", restored.Counter, state.Counter)
		}
	})

	t.Run("persists and restores state with nested data", func(t *testing.T) {
		memory := graph.MemMemoryFactory[MemoryTestState]()
		ctx := context.Background()
		key := "nested-key"

		state := MemoryTestState{
			Value:   "complex",
			Counter: 100,
			Data: map[string]interface{}{
				"name": "test",
				"age":  30,
			},
		}

		// Persist state
		persistFn := memory.PersistFn()
		err := persistFn(ctx, key, state)
		if err != nil {
			t.Fatalf("PersistFn failed: %v", err)
		}

		// Restore state
		restoreFn := memory.RestoreFn()
		restored, err := restoreFn(ctx, key)
		if err != nil {
			t.Fatalf("RestoreFn failed: %v", err)
		}

		// Verify restored state
		if restored.Value != state.Value {
			t.Errorf("Value mismatch: got %q, want %q", restored.Value, state.Value)
		}
		if restored.Counter != state.Counter {
			t.Errorf("Counter mismatch: got %d, want %d", restored.Counter, state.Counter)
		}
		if restored.Data["name"] != state.Data["name"] {
			t.Errorf("Data name mismatch: got %v, want %v", restored.Data["name"], state.Data["name"])
		}
		if restored.Data["age"] != state.Data["age"] {
			t.Errorf("Data age mismatch: got %v, want %v", restored.Data["age"], state.Data["age"])
		}
	})

	t.Run("restores zero value for non-existent key", func(t *testing.T) {
		memory := graph.MemMemoryFactory[MemoryTestState]()
		ctx := context.Background()
		key := "non-existent-key"

		restoreFn := memory.RestoreFn()
		restored, err := restoreFn(ctx, key)
		if err != nil {
			t.Fatalf("RestoreFn failed: %v", err)
		}

		// Verify zero value
		var zero MemoryTestState
		if restored.Value != zero.Value {
			t.Errorf("Value should be zero value, got %q", restored.Value)
		}
		if restored.Counter != zero.Counter {
			t.Errorf("Counter should be zero value, got %d", restored.Counter)
		}
	})

	t.Run("restores zero value before any persist", func(t *testing.T) {
		memory := graph.MemMemoryFactory[MemoryTestState]()
		ctx := context.Background()
		key := "empty-key"

		restoreFn := memory.RestoreFn()
		restored, err := restoreFn(ctx, key)
		if err != nil {
			t.Fatalf("RestoreFn failed: %v", err)
		}

		// Verify zero value
		var zero MemoryTestState
		if restored.Value != zero.Value {
			t.Errorf("Value should be zero value, got %q", restored.Value)
		}
		if restored.Counter != zero.Counter {
			t.Errorf("Counter should be zero value, got %d", restored.Counter)
		}
		if restored.Data != nil {
			t.Errorf("Data should be nil, got %+v", restored.Data)
		}
	})

	t.Run("overwrites existing state with same key", func(t *testing.T) {
		memory := graph.MemMemoryFactory[MemoryTestState]()
		ctx := context.Background()
		key := "overwrite-key"

		// Persist first state
		state1 := MemoryTestState{Value: "first", Counter: 1}
		persistFn := memory.PersistFn()
		err := persistFn(ctx, key, state1)
		if err != nil {
			t.Fatalf("First PersistFn failed: %v", err)
		}

		// Persist second state with same key
		state2 := MemoryTestState{Value: "second", Counter: 2}
		err = persistFn(ctx, key, state2)
		if err != nil {
			t.Fatalf("Second PersistFn failed: %v", err)
		}

		// Restore and verify second state
		restoreFn := memory.RestoreFn()
		restored, err := restoreFn(ctx, key)
		if err != nil {
			t.Fatalf("RestoreFn failed: %v", err)
		}

		if restored.Value != state2.Value {
			t.Errorf("Value mismatch: got %q, want %q", restored.Value, state2.Value)
		}
		if restored.Counter != state2.Counter {
			t.Errorf("Counter mismatch: got %d, want %d", restored.Counter, state2.Counter)
		}
	})
}

func TestMemMemory_MultipleKeys(t *testing.T) {
	t.Run("handles multiple keys independently", func(t *testing.T) {
		memory := graph.MemMemoryFactory[MemoryTestState]()
		ctx := context.Background()

		// Persist multiple states with different keys
		states := map[string]MemoryTestState{
			"key1": {Value: "state1", Counter: 1},
			"key2": {Value: "state2", Counter: 2},
			"key3": {Value: "state3", Counter: 3},
		}

		persistFn := memory.PersistFn()
		for key, state := range states {
			err := persistFn(ctx, key, state)
			if err != nil {
				t.Fatalf("PersistFn failed for key %q: %v", key, err)
			}
		}

		// Restore and verify each state
		restoreFn := memory.RestoreFn()
		for key, expectedState := range states {
			restored, err := restoreFn(ctx, key)
			if err != nil {
				t.Fatalf("RestoreFn failed for key %q: %v", key, err)
			}

			if restored.Value != expectedState.Value {
				t.Errorf("Key %q: Value mismatch: got %q, want %q", key, restored.Value, expectedState.Value)
			}
			if restored.Counter != expectedState.Counter {
				t.Errorf("Key %q: Counter mismatch: got %d, want %d", key, restored.Counter, expectedState.Counter)
			}
		}
	})
}

func TestMemMemory_ContextPropagation(t *testing.T) {
	t.Run("persist respects context", func(t *testing.T) {
		memory := graph.MemMemoryFactory[MemoryTestState]()
		ctx := context.Background()
		key := "context-key"

		state := MemoryTestState{Value: "test", Counter: 1}

		persistFn := memory.PersistFn()
		err := persistFn(ctx, key, state)
		if err != nil {
			t.Fatalf("PersistFn failed: %v", err)
		}

		// Verify state was persisted
		restoreFn := memory.RestoreFn()
		restored, err := restoreFn(ctx, key)
		if err != nil {
			t.Fatalf("RestoreFn failed: %v", err)
		}

		if restored.Value != state.Value {
			t.Errorf("Value mismatch: got %q, want %q", restored.Value, state.Value)
		}
	})

	t.Run("restore respects context", func(t *testing.T) {
		memory := graph.MemMemoryFactory[MemoryTestState]()
		ctx := context.Background()
		key := "context-restore-key"

		state := MemoryTestState{Value: "test", Counter: 1}

		persistFn := memory.PersistFn()
		err := persistFn(ctx, key, state)
		if err != nil {
			t.Fatalf("PersistFn failed: %v", err)
		}

		// Restore with context
		restoreFn := memory.RestoreFn()
		restored, err := restoreFn(ctx, key)
		if err != nil {
			t.Fatalf("RestoreFn failed: %v", err)
		}

		if restored.Value != state.Value {
			t.Errorf("Value mismatch: got %q, want %q", restored.Value, state.Value)
		}
	})
}

func TestMemMemory_EmptyKey(t *testing.T) {
	t.Run("handles empty string key", func(t *testing.T) {
		memory := graph.MemMemoryFactory[MemoryTestState]()
		ctx := context.Background()
		key := ""

		state := MemoryTestState{Value: "empty-key-state", Counter: 99}

		// Persist with empty key
		persistFn := memory.PersistFn()
		err := persistFn(ctx, key, state)
		if err != nil {
			t.Fatalf("PersistFn failed: %v", err)
		}

		// Restore with empty key
		restoreFn := memory.RestoreFn()
		restored, err := restoreFn(ctx, key)
		if err != nil {
			t.Fatalf("RestoreFn failed: %v", err)
		}

		if restored.Value != state.Value {
			t.Errorf("Value mismatch: got %q, want %q", restored.Value, state.Value)
		}
		if restored.Counter != state.Counter {
			t.Errorf("Counter mismatch: got %d, want %d", restored.Counter, state.Counter)
		}
	})
}

func TestMemMemory_ZeroValueState(t *testing.T) {
	t.Run("persists and restores zero value state", func(t *testing.T) {
		memory := graph.MemMemoryFactory[MemoryTestState]()
		ctx := context.Background()
		key := "zero-value-key"

		var zeroState MemoryTestState

		// Persist zero value state
		persistFn := memory.PersistFn()
		err := persistFn(ctx, key, zeroState)
		if err != nil {
			t.Fatalf("PersistFn failed: %v", err)
		}

		// Restore zero value state
		restoreFn := memory.RestoreFn()
		restored, err := restoreFn(ctx, key)
		if err != nil {
			t.Fatalf("RestoreFn failed: %v", err)
		}

		if restored.Value != zeroState.Value {
			t.Errorf("Value should be zero, got %q", restored.Value)
		}
		if restored.Counter != zeroState.Counter {
			t.Errorf("Counter should be zero, got %d", restored.Counter)
		}
		if restored.Data != nil {
			t.Errorf("Data should be nil, got %+v", restored.Data)
		}
	})
}

func TestMemMemory_Concurrent(t *testing.T) {
	t.Run("handles concurrent persist and restore operations", func(t *testing.T) {
		memory := graph.MemMemoryFactory[MemoryTestState]()
		ctx := context.Background()

		const goroutines = 10
		done := make(chan bool, goroutines)

		// Launch concurrent goroutines
		for i := 0; i < goroutines; i++ {
			go func(id int) {
				defer func() { done <- true }()

				key := "concurrent-key"
				state := MemoryTestState{
					Value:   "concurrent",
					Counter: id,
				}

				// Persist
				persistFn := memory.PersistFn()
				err := persistFn(ctx, key, state)
				if err != nil {
					t.Errorf("Goroutine %d: PersistFn failed: %v", id, err)
					return
				}

				// Restore
				restoreFn := memory.RestoreFn()
				_, err = restoreFn(ctx, key)
				if err != nil {
					t.Errorf("Goroutine %d: RestoreFn failed: %v", id, err)
				}
			}(i)
		}

		// Wait for all goroutines to complete
		for i := 0; i < goroutines; i++ {
			<-done
		}
	})

	t.Run("handles concurrent operations on different keys", func(t *testing.T) {
		memory := graph.MemMemoryFactory[MemoryTestState]()
		ctx := context.Background()

		const goroutines = 10
		done := make(chan bool, goroutines)

		// Launch concurrent goroutines with different keys
		for i := 0; i < goroutines; i++ {
			go func(id int) {
				defer func() { done <- true }()

				key := string(rune('a' + id))
				state := MemoryTestState{
					Value:   key,
					Counter: id,
				}

				// Persist
				persistFn := memory.PersistFn()
				err := persistFn(ctx, key, state)
				if err != nil {
					t.Errorf("Goroutine %d: PersistFn failed: %v", id, err)
					return
				}

				// Restore
				restoreFn := memory.RestoreFn()
				restored, err := restoreFn(ctx, key)
				if err != nil {
					t.Errorf("Goroutine %d: RestoreFn failed: %v", id, err)
					return
				}

				// Verify
				if restored.Value != state.Value {
					t.Errorf("Goroutine %d: Value mismatch: got %q, want %q", id, restored.Value, state.Value)
				}
				if restored.Counter != state.Counter {
					t.Errorf("Goroutine %d: Counter mismatch: got %d, want %d", id, restored.Counter, state.Counter)
				}
			}(i)
		}

		// Wait for all goroutines to complete
		for i := 0; i < goroutines; i++ {
			<-done
		}
	})
}

// Test with a different state type to verify generic implementation
type AnotherSimpleState struct {
	ID   int
	Name string
}

func TestMemMemory_DifferentStateType(t *testing.T) {
	t.Run("works with different state type", func(t *testing.T) {
		memory := graph.MemMemoryFactory[AnotherSimpleState]()
		ctx := context.Background()
		key := "simple-key"

		state := AnotherSimpleState{
			ID:   123,
			Name: "test-state",
		}

		// Persist
		persistFn := memory.PersistFn()
		err := persistFn(ctx, key, state)
		if err != nil {
			t.Fatalf("PersistFn failed: %v", err)
		}

		// Restore
		restoreFn := memory.RestoreFn()
		restored, err := restoreFn(ctx, key)
		if err != nil {
			t.Fatalf("RestoreFn failed: %v", err)
		}

		if restored.ID != state.ID {
			t.Errorf("ID mismatch: got %d, want %d", restored.ID, state.ID)
		}
		if restored.Name != state.Name {
			t.Errorf("Name mismatch: got %q, want %q", restored.Name, state.Name)
		}
	})
}
