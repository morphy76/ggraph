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
	return CreateRuntimeWithMerger(startEdge, stateMonitorCh, nil)
}

// CreateRuntimeWithMerger creates a new instance of Runtime with the specified SharedState type and state merger function.
func CreateRuntimeWithMerger[T g.SharedState](
	startEdge g.Edge[T],
	stateMonitorCh chan g.StateMonitorEntry[T],
	merger g.StateMergeFn[T],
) (g.Runtime[T], error) {
	var zero T
	return CreateRuntimeWithMergerAndInitialState(startEdge, stateMonitorCh, merger, zero)
}

// CreateRuntimeWithMergerAndInitialState creates a new instance of Runtime with the specified SharedState type, state merger function, and initial state.
func CreateRuntimeWithMergerAndInitialState[T g.SharedState](
	startEdge g.Edge[T],
	stateMonitorCh chan g.StateMonitorEntry[T],
	merger g.StateMergeFn[T],
	initialState T,
) (g.Runtime[T], error) {
	return i.RuntimeFactory(startEdge, stateMonitorCh, merger, initialState)
}
