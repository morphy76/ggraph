package builders

import (
	i "github.com/morphy76/ggraph/internal/graph"
	g "github.com/morphy76/ggraph/pkg/graph"
)

// CreateAnyRoutePolicy creates a new instance of AnyRoutePolicy.
func CreateAnyRoutePolicy[T g.SharedState]() (g.RoutePolicy[T], error) {
	return CreateConditionalRoutePolicy(func(state T, edges []g.Edge[T]) g.Edge[T] {
		if len(edges) == 0 {
			return nil
		}
		return edges[0]
	})
}

// CreateConditionalRoutePolicy creates a new instance of ConditionalRoutePolicy with the specified conditional function.
func CreateConditionalRoutePolicy[T g.SharedState](selectionFn g.EdgeSelectionFn[T]) (g.RoutePolicy[T], error) {
	return i.RouterPolicyImplFactory(selectionFn)
}
