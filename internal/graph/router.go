package graph

import (
	"fmt"

	g "github.com/morphy76/ggraph/pkg/graph"
)

// AnyRoute is a simple EdgeSelectionFn that selects the first available edge.
func AnyRoute[T g.SharedState](userInput T, currentState T, edges []g.Edge[T]) g.Edge[T] {
	if len(edges) > 0 {
		return edges[0]
	}
	return nil
}

// RouterPolicyImplFactory creates a new instance of RoutePolicy with the specified SharedState type and selection function.
func RouterPolicyImplFactory[T g.SharedState](selectionFn g.EdgeSelectionFn[T]) (g.RoutePolicy[T], error) {
	if selectionFn == nil {
		return nil, fmt.Errorf("conditional route policy creation failed: selection function cannot be nil")
	}
	return &routePolicyImpl[T]{
		selectionFunc: selectionFn,
	}, nil
}

// ------------------------------------------------------------------------------
// RoutePolicy Implementation
// ------------------------------------------------------------------------------

var _ g.RoutePolicy[g.SharedState] = (*routePolicyImpl[g.SharedState])(nil)

type routePolicyImpl[T g.SharedState] struct {
	selectionFunc g.EdgeSelectionFn[T]
}

func (p *routePolicyImpl[T]) SelectEdge(userInput T, currentState T, edges []g.Edge[T]) g.Edge[T] {
	return p.selectionFunc(userInput, currentState, edges)
}
