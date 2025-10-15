package graph

import (
	"fmt"

	g "github.com/morphy76/ggraph/pkg/graph"
)

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

func (p *routePolicyImpl[T]) SelectEdge(state T, edges []g.Edge[T]) g.Edge[T] {
	return p.selectionFunc(state, edges)
}
