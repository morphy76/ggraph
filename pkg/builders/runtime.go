package builders

import (
	i "github.com/morphy76/ggraph/internal/graph"
	g "github.com/morphy76/ggraph/pkg/graph"
)

// CreateRuntime creates a new instance of Runtime with the specified SharedState type.
func CreateRuntime[T g.SharedState](
	startEdge g.Edge[T],
	stateMonitorCh chan g.StateMonitorEntry[T],
) (g.Runtime[T], error) {
	var zero T
	return CreateRuntimeWithInitialState(startEdge, stateMonitorCh, zero)
}

// CreateRuntimeWithInitialState creates a new instance of Runtime with the specified SharedState type and an initial state.
func CreateRuntimeWithInitialState[T g.SharedState](
	startEdge g.Edge[T],
	stateMonitorCh chan g.StateMonitorEntry[T],
	initialState T,
) (g.Runtime[T], error) {
	return i.RuntimeFactory(startEdge, stateMonitorCh, initialState)
}
