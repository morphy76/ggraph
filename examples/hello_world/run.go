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

	graph := graph.CreateRuntime(startEdge)
	graph.AddEdge(midEdge)
	graph.AddEdge(endEdge)

	err = graph.Validate()
	if err != nil {
		log.Fatalf("Graph validation failed: %v", err)
	}

	// Invoke the runtime with an initial state
	initialState := MyState{Message: "Bob"}
	finalState, err := graph.Invoke(initialState)
	if err != nil {
		log.Fatalf("Runtime invocation failed: %v", err)
	}

	log.Printf("Final State: %+v", finalState)
}
