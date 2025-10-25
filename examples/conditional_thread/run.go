package main

import (
	"fmt"
	"log"

	"github.com/google/uuid"
	b "github.com/morphy76/ggraph/pkg/builders"
	g "github.com/morphy76/ggraph/pkg/graph"
)

var _ g.SharedState = (*myState)(nil)

type myState struct {
	op     string
	num2   int
	Result int
}

func main() {

	routingPolicy, err := b.CreateConditionalRoutePolicy(func(userInput, currentState myState, edges []g.Edge[myState]) g.Edge[myState] {
		op := userInput.op
		for _, edge := range edges {
			if label, ok := edge.LabelByKey("operation"); ok && label == op {
				return edge
			}
		}
		return nil
	})
	routerNode, err := b.CreateRouter("operation_routing", routingPolicy)
	if err != nil {
		log.Fatalf("Router creation failed: %v", err)
	}

	adder, err := b.NewNodeBuilder("Adder", func(userInput myState, currentState myState, notifyPartial g.NotifyPartialFn[myState]) (myState, error) {
		currentState.Result = currentState.Result + userInput.num2
		return currentState, nil
	}).Build()
	if err != nil {
		log.Fatalf("Node creation failed: %v", err)
	}

	subtractor, err := b.NewNodeBuilder("Subtractor", func(userInput myState, currentState myState, notifyPartial g.NotifyPartialFn[myState]) (myState, error) {
		currentState.Result = currentState.Result - userInput.num2
		return currentState, nil
	}).Build()
	if err != nil {
		log.Fatalf("Node creation failed: %v", err)
	}

	startEdge := b.CreateStartEdge(routerNode)
	stateMonitorCh := make(chan g.StateMonitorEntry[myState], 10)
	myGraph, err := b.CreateRuntimeWithInitialState(startEdge, stateMonitorCh, myState{Result: 10})
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

	invocationConfig1 := g.InvokeConfigThreadID(uuid.NewString())
	invocationConfig2 := g.InvokeConfigThreadID(uuid.NewString())

	myGraph.Invoke(myState{op: "+", num2: 5}, invocationConfig1)
	myGraph.Invoke(myState{op: "+", num2: 5}, invocationConfig2)
	twice := 2
	for twice > 0 {
		entry := <-stateMonitorCh
		if !entry.Running {
			twice--
			fmt.Printf("State Monitor Node: %s Thread: %s Entry: %+v Error: %v\n", entry.Node, entry.ThreadID, entry.NewState.Result, entry.Error)
		}
	}

	twice = 2
	myGraph.Invoke(myState{op: "-", num2: 5}, invocationConfig1)
	myGraph.Invoke(myState{op: "-", num2: 5}, invocationConfig2)
	for twice > 0 {
		entry := <-stateMonitorCh
		if !entry.Running {
			twice--
			fmt.Printf("State Monitor Node: %s Thread: %s Entry: %+v Error: %v\n", entry.Node, entry.ThreadID, entry.NewState.Result, entry.Error)
		}
	}

	twice = 2
	myGraph.Invoke(myState{op: "-", num2: 10}, invocationConfig1)
	myGraph.Invoke(myState{op: "-", num2: 10}, invocationConfig2)
	for twice > 0 {
		entry := <-stateMonitorCh
		if !entry.Running {
			twice--
			fmt.Printf("State Monitor Node: %s Thread: %s Entry: %+v Error: %v\n", entry.Node, entry.ThreadID, entry.NewState.Result, entry.Error)
		}
	}
}
