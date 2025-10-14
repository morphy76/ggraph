package main

import (
	"fmt"
	"log"

	"github.com/morphy76/ggraph/pkg/graph"
)

var _ graph.SharedState = (*MyState)(nil)

type MyState struct {
	num1   int
	op     string
	num2   int
	Result int
}

func main() {

	routingPolicy, err := graph.CreateConditionalRoutePolicy(func(state MyState, edges []graph.Edge[MyState]) graph.Edge[MyState] {
		op := state.op
		for _, edge := range edges {
			if label, ok := edge.LabelByKey("operation"); ok && label == op {
				return edge
			}
		}
		return nil
	})
	routerNode, err := graph.CreateRouter("operation_routing", routingPolicy)

	adder, err := graph.CreateNode("Adder", func(state MyState) (MyState, error) {
		state.Result = state.num1 + state.num2
		return state, nil
	})
	if err != nil {
		log.Fatalf("Node creation failed: %v", err)
	}

	subtractor, err := graph.CreateNode("Subtractor", func(state MyState) (MyState, error) {
		state.Result = state.num1 - state.num2
		return state, nil
	})
	if err != nil {
		log.Fatalf("Node creation failed: %v", err)
	}

	startEdge := graph.CreateStartEdge(routerNode)
	additionEdge := graph.CreateEdge(routerNode, adder, map[string]string{"operation": "+"})
	subtractionEdge := graph.CreateEdge(routerNode, subtractor, map[string]string{"operation": "-"})
	additionEndEdge := graph.CreateEndEdge(adder)
	subtractionEndEdge := graph.CreateEndEdge(subtractor)

	stateMonitorCh := make(chan graph.StateMonitorEntry[MyState], 10)
	graph, err := graph.CreateRuntime(startEdge, stateMonitorCh)
	if err != nil {
		log.Fatalf("Runtime creation failed: %v", err)
	}
	defer graph.Shutdown()
	graph.AddEdge(additionEdge, subtractionEdge, additionEndEdge, subtractionEndEdge)

	err = graph.Validate()
	if err != nil {
		log.Fatalf("Graph validation failed: %v", err)
	}

	newState := MyState{num1: 10, op: "-", num2: 5}
	graph.Invoke(newState)

	for {
		entry := <-stateMonitorCh
		fmt.Printf("State Monitor Entry: %+v\n", entry)
		if !entry.Running {
			break
		}
	}

	newState = MyState{num1: 10, op: "+", num2: 5}
	graph.Invoke(newState)

	for {
		entry := <-stateMonitorCh
		fmt.Printf("State Monitor Entry: %+v\n", entry)
		if !entry.Running {
			break
		}
	}
}
