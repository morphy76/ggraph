package graph

import (
	g "github.com/morphy76/ggraph/pkg/graph"
)

func monitorRunning[T g.SharedState](node string, stateChange T) g.StateMonitorEntry[T] {
	return g.StateMonitorEntry[T]{
		Node:        node,
		StateChange: stateChange,
		Running:     true,
		Partial:     false,
	}
}

func monitorNonFatalError[T g.SharedState](node string, err error) g.StateMonitorEntry[T] {
	return g.StateMonitorEntry[T]{
		Node:    node,
		Error:   err,
		Running: true,
		Partial: false,
	}
}

func monitorError[T g.SharedState](node string, err error) g.StateMonitorEntry[T] {
	return g.StateMonitorEntry[T]{
		Node:    node,
		Error:   err,
		Running: false,
		Partial: false,
	}
}

func monitorPartial[T g.SharedState](node string, stateChange T) g.StateMonitorEntry[T] {
	return g.StateMonitorEntry[T]{
		Node:        node,
		StateChange: stateChange,
		Running:     true,
		Partial:     true,
	}
}

func monitorCompleted[T g.SharedState](node string) g.StateMonitorEntry[T] {
	return g.StateMonitorEntry[T]{
		Node:    node,
		Running: false,
		Partial: false,
	}
}
