package graph

import (
	"fmt"
)

// EdgeSelectionFn defines a function type for conditional routing based on the current state and available edges.
type EdgeSelectionFn[T SharedState] func(state T, edges []Edge[T]) Edge[T]

// CreateAnyRoutePolicy creates a new instance of AnyRoutePolicy.
func CreateAnyRoutePolicy[T SharedState]() (RoutePolicy[T], error) {
	return CreateConditionalRoutePolicy(func(state T, edges []Edge[T]) Edge[T] {
		if len(edges) == 0 {
			return nil
		}
		return edges[0]
	})
}

// CreateConditionalRoutePolicy creates a new instance of ConditionalRoutePolicy with the specified conditional function.
func CreateConditionalRoutePolicy[T SharedState](selectionFn EdgeSelectionFn[T]) (RoutePolicy[T], error) {
	if selectionFn == nil {
		return nil, fmt.Errorf("conditional route policy creation failed: selection function cannot be nil")
	}
	return &routePolicyImpl[T]{
		selectionFunc: selectionFn,
	}, nil
}

// RoutePolicy defines a policy for routing between nodes in the graph.
type RoutePolicy[T SharedState] interface {
	// SelectEdge selects an edge from the available edges based on the current state.
	SelectEdge(state T, edges []Edge[T]) Edge[T]
}

var _ RoutePolicy[SharedState] = (*routePolicyImpl[SharedState])(nil)

type routePolicyImpl[T SharedState] struct {
	selectionFunc EdgeSelectionFn[T]
}

func (p *routePolicyImpl[T]) SelectEdge(state T, edges []Edge[T]) Edge[T] {
	return p.selectionFunc(state, edges)
}
