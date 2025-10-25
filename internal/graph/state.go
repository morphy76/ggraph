package graph

import (
	g "github.com/morphy76/ggraph/pkg/graph"
)

func monitorRunning[T g.SharedState](node string, threadID string, newState T) g.StateMonitorEntry[T] {
	return g.StateMonitorEntry[T]{
		Node:     node,
		ThreadID: threadID,
		Running:  true,
		Partial:  false,
		NewState: newState,
	}
}

func monitorNonFatalError[T g.SharedState](node string, threadID string, err error) g.StateMonitorEntry[T] {
	return g.StateMonitorEntry[T]{
		Node:     node,
		ThreadID: threadID,
		Error:    err,
		Running:  true,
		Partial:  false,
	}
}

func monitorError[T g.SharedState](node string, threadID string, err error) g.StateMonitorEntry[T] {
	return g.StateMonitorEntry[T]{
		Node:     node,
		ThreadID: threadID,
		Error:    err,
		Running:  false,
		Partial:  false,
	}
}

func monitorPartial[T g.SharedState](node string, threadID string, stateChange T) g.StateMonitorEntry[T] {
	return g.StateMonitorEntry[T]{
		Node:     node,
		ThreadID: threadID,
		NewState: stateChange,
		Running:  true,
		Partial:  true,
	}
}

func monitorCompleted[T g.SharedState](node string, threadID string, newState T) g.StateMonitorEntry[T] {
	return g.StateMonitorEntry[T]{
		Node:     node,
		ThreadID: threadID,
		Running:  false,
		Partial:  false,
		NewState: newState,
	}
}
