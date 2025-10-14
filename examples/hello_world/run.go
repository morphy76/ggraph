package main

import (
	"fmt"
	"log"

	"github.com/morphy76/ggraph/pkg/graph"
)

var _ graph.SharedState = (*MyState)(nil)

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

	helloNode, err := graph.CreateNode("HelloNode", func(state MyState) (MyState, error) {
		state.Message = fmt.Sprintf("Hello %s!!!", state.Message)
		return state, nil
	})
	if err != nil {
		log.Fatalf("Node creation failed: %v", err)
	}

	goodbyeNode, err := graph.CreateNode("GoodbyeNode", func(state MyState) (MyState, error) {
		state.Message = fmt.Sprintf("Goodbye %s!!!", state.Message)
		return state, nil
	})
	if err != nil {
		log.Fatalf("Node creation failed: %v", err)
	}

	startEdge := graph.CreateStartEdge(helloNode)
	midEdge := graph.CreateEdge(helloNode, goodbyeNode)
	endEdge := graph.CreateEndEdge(goodbyeNode)

	initialState := MyState{Message: ""}
	stateMonitorCh := make(chan graph.StateMonitorEntry[MyState], 10)
	graph := graph.CreateRuntime(startEdge, initialState, merge, stateMonitorCh)
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
