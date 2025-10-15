package main

import (
	"fmt"
	"log"

	b "github.com/morphy76/ggraph/pkg/builders"
	g "github.com/morphy76/ggraph/pkg/graph"
)

var _ g.SharedState = (*MyState)(nil)

type MyState struct {
	op     string
	num2   int
	Result int
}

func merger(current, other MyState) MyState {
	if current.op != other.op {
		current.op = other.op
	}
	if current.num2 != other.num2 {
		current.num2 = other.num2
	}
	if current.Result != other.Result {
		current.Result = other.Result
	}
	return current
}

func main() {

	routingPolicy, err := b.CreateConditionalRoutePolicy(func(state MyState, edges []g.Edge[MyState]) g.Edge[MyState] {
		op := state.op
		for _, edge := range edges {
			if label, ok := edge.LabelByKey("operation"); ok && label == op {
				return edge
			}
		}
		return nil
	})
	routerNode, err := b.CreateRouter("operation_routing", routingPolicy)

	adder, err := b.CreateNode("Adder", func(state MyState, notify func(MyState)) (MyState, error) {
		state.Result = state.Result + state.num2
		return state, nil
	})
	if err != nil {
		log.Fatalf("Node creation failed: %v", err)
	}

	subtractor, err := b.CreateNode("Subtractor", func(state MyState, notify func(MyState)) (MyState, error) {
		state.Result = state.Result - state.num2
		return state, nil
	})
	if err != nil {
		log.Fatalf("Node creation failed: %v", err)
	}

	startEdge := b.CreateStartEdge(routerNode)
	stateMonitorCh := make(chan g.StateMonitorEntry[MyState], 10)
	myGraph, err := b.CreateRuntimeWithMergerAndInitialState(startEdge, stateMonitorCh, merger, MyState{Result: 10})
	if err != nil {
		log.Fatalf("Runtime creation failed: %v", err)
	}
	defer myGraph.Shutdown()

	additionEdge := b.CreateEdge(routerNode, adder, map[string]string{"operation": "+"})
	subtractionEdge := b.CreateEdge(routerNode, subtractor, map[string]string{"operation": "-"})
	additionEndEdge := b.CreateEndEdge(adder)
	subtractionEndEdge := b.CreateEndEdge(subtractor)
	myGraph.AddEdge(additionEdge, subtractionEdge, additionEndEdge, subtractionEndEdge)

	err = myGraph.Validate()
	if err != nil {
		log.Fatalf("Graph validation failed: %v", err)
	}

	myGraph.Invoke(MyState{op: "+", num2: 5})
	myGraph.Invoke(MyState{op: "-", num2: 5})

	breakLoop := 2
	for {
		entry := <-stateMonitorCh
		if !entry.Running {
			fmt.Printf("State Monitor Node: %s Entry: %+v Error: %v\n", entry.Node, entry.CurrentState, entry.Error)
			breakLoop--
			if breakLoop == 0 {
				break
			}
		}
	}
}
