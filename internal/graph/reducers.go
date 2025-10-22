package graph

import (
	g "github.com/morphy76/ggraph/pkg/graph"
)

// Replacer is a simple state reducer function.
func Replacer[T g.SharedState](currentState, change T) T {
	return change
}
