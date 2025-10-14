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

	adder, err := graph.CreateNode("Adder", func(state MyState, notify func(MyState)) (MyState, error) {
		state.Result = state.Result + state.num2
		return state, nil
	})
	if err != nil {
		log.Fatalf("Node creation failed: %v", err)
	}

	subtractor, err := graph.CreateNode("Subtractor", func(state MyState, notify func(MyState)) (MyState, error) {
		state.Result = state.Result - state.num2
		return state, nil
	})
	if err != nil {
		log.Fatalf("Node creation failed: %v", err)
	}

	startEdge := graph.CreateStartEdge(routerNode)
	stateMonitorCh := make(chan graph.StateMonitorEntry[MyState], 10)
	myGraph, err := graph.CreateRuntime(startEdge, stateMonitorCh)
	if err != nil {
		log.Fatalf("Runtime creation failed: %v", err)
	}
	defer myGraph.Shutdown()

	additionEdge := graph.CreateEdge(routerNode, adder, map[string]string{"operation": "+"})
	subtractionEdge := graph.CreateEdge(routerNode, subtractor, map[string]string{"operation": "-"})
	additionEndEdge := graph.CreateEndEdge(adder)
	subtractionEndEdge := graph.CreateEndEdge(subtractor)
	myGraph.AddEdge(additionEdge, subtractionEdge, additionEndEdge, subtractionEndEdge)

	err = myGraph.Validate()
	if err != nil {
		log.Fatalf("Graph validation failed: %v", err)
	}

	newState := MyState{op: "-", num2: 5}
	myGraph.Invoke(newState)

	newState = MyState{op: "+", num2: 5}
	myGraph.Invoke(newState)

	breakLoop := 2
	for {
		entry := <-stateMonitorCh
		fmt.Printf("State Monitor Node: %s Entry: %+v Error: %v\n", entry.Node, entry.CurrentState, entry.Error)
		if !entry.Running {
			breakLoop--
			if breakLoop == 0 {
				break
			}
		}
	}
}
