package main

import (
	"fmt"
	"log"

	b "github.com/morphy76/ggraph/pkg/builders"
	g "github.com/morphy76/ggraph/pkg/graph"
)

var _ g.SharedState = (*MyState)(nil)

type MyState struct {
	Message string
}

func merge(originalState, newState MyState) MyState {
	if newState.Message != "" {
		originalState.Message = newState.Message
	}
	return originalState
}

func main() {

	helloNode, err := b.CreateNode("HelloNode", func(userInput MyState, currentState MyState, notify func(MyState)) (MyState, error) {
		currentState.Message = fmt.Sprintf("Hello %s!!!", userInput.Message)
		return currentState, nil
	})
	if err != nil {
		log.Fatalf("Node creation failed: %v", err)
	}

	goodbyeNode, err := b.CreateNode("GoodbyeNode", func(userInput MyState, currentState MyState, notify func(MyState)) (MyState, error) {
		currentState.Message = fmt.Sprintf("Goodbye %s!!!", userInput.Message)
		return currentState, nil
	})
	if err != nil {
		log.Fatalf("Node creation failed: %v", err)
	}

	startEdge := b.CreateStartEdge(helloNode)
	midEdge := b.CreateEdge(helloNode, goodbyeNode)
	endEdge := b.CreateEndEdge(goodbyeNode)

	initialState := MyState{Message: ""}
	stateMonitorCh := make(chan g.StateMonitorEntry[MyState], 10)
	graph, err := b.CreateRuntimeWithMergerAndInitialState(startEdge, stateMonitorCh, merge, initialState)
	if err != nil {
		log.Fatalf("Runtime creation failed: %v", err)
	}
	defer graph.Shutdown()
	graph.AddEdge(midEdge, endEdge)

	err = graph.Validate()
	if err != nil {
		log.Fatalf("Graph validation failed: %v", err)
	}

	newState := MyState{Message: "Bob"}
	graph.Invoke(newState)

	for {
		entry := <-stateMonitorCh
		fmt.Printf("State Monitor Entry: %+v\n", entry)
		if !entry.Running {
			break
		}
	}
}
