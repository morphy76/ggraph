package main

import (
	"fmt"
	"log"

	b "github.com/morphy76/ggraph/pkg/builders"
	g "github.com/morphy76/ggraph/pkg/graph"
)

var _ g.SharedState = (*myState)(nil)

type myState struct {
	Message string
}

func main() {

	helloNode, err := b.NewNode("HelloNode", func(userInput myState, currentState myState, notifyPartial g.NotifyPartialFn[myState]) (myState, error) {
		currentState.Message = fmt.Sprintf("Hello %s!!!", userInput.Message)
		return currentState, nil
	})
	if err != nil {
		log.Fatalf("Node creation failed: %v", err)
	}

	goodbyeNode, err := b.NewNode("GoodbyeNode", func(userInput myState, currentState myState, notifyPartial g.NotifyPartialFn[myState]) (myState, error) {
		currentState.Message = fmt.Sprintf("Goodbye %s!!!", userInput.Message)
		return currentState, nil
	})
	if err != nil {
		log.Fatalf("Node creation failed: %v", err)
	}

	startEdge := b.CreateStartEdge(helloNode)
	midEdge := b.CreateEdge(helloNode, goodbyeNode)
	endEdge := b.CreateEndEdge(goodbyeNode)

	initialState := myState{Message: ""}
	stateMonitorCh := make(chan g.StateMonitorEntry[myState], 10)
	graph, err := b.CreateRuntime(startEdge, stateMonitorCh, g.WithInitialState(initialState))
	if err != nil {
		log.Fatalf("Runtime creation failed: %v", err)
	}
	defer graph.Shutdown()
	graph.AddEdge(midEdge, endEdge)

	err = graph.Validate()
	if err != nil {
		log.Fatalf("Graph validation failed: %v", err)
	}

	userInput := myState{Message: "Bob"}
	graph.Invoke(userInput)

	for {
		entry := <-stateMonitorCh
		fmt.Printf("[%s - Running: %t], Graph state message: %s, Error: %v\n", entry.Node, entry.Running, entry.NewState, entry.Error)
		if !entry.Running {
			break
		}
	}
}
