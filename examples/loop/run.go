package main

import (
	"fmt"
	"log"
	"math/rand"

	b "github.com/morphy76/ggraph/pkg/builders"
	g "github.com/morphy76/ggraph/pkg/graph"
)

var _ g.SharedState = (*GameState)(nil)

type GameState struct {
	Target  int
	Guess   int
	Tries   int
	Success bool
	Hint    string
	Low     int
	High    int
}

func main() {
	// Node 1: Determine target number
	initNode, _ := b.CreateNode("InitNode", func(userInput GameState, currentState GameState, notifyPartial g.NotifyPartialFn[GameState]) (GameState, error) {
		currentState.Target = rand.Intn(100) + 1
		currentState.Tries = 0
		currentState.Low = 1
		currentState.High = 100
		fmt.Printf("🎯 Target set (hidden)\n")
		return currentState, nil
	})

	// Node 2: Make a guess using binary search
	guessNode, _ := b.CreateNode("GuessNode", func(userInput GameState, currentState GameState, notifyPartial g.NotifyPartialFn[GameState]) (GameState, error) {
		currentState.Tries++
		currentState.Guess = (currentState.Low + currentState.High) / 2
		currentState.Success = (currentState.Guess == currentState.Target)
		fmt.Printf("🤔 Try #%d: Guessed %d (range: %d-%d)\n", currentState.Tries, currentState.Guess, currentState.Low, currentState.High)
		return currentState, nil
	})

	// Node 3: Provide hint and adjust range
	hintNode, _ := b.CreateNode("HintNode", func(userInput GameState, currentState GameState, notifyPartial g.NotifyPartialFn[GameState]) (GameState, error) {
		if currentState.Guess < currentState.Target {
			currentState.Low = currentState.Guess + 1
			currentState.Hint = "higher"
			fmt.Printf("💡 Hint: Try higher!\n")
		} else {
			currentState.High = currentState.Guess - 1
			currentState.Hint = "lower"
			fmt.Printf("💡 Hint: Try lower!\n")
		}
		return currentState, nil
	})

	// Router: Check success
	routingPolicy, _ := b.CreateConditionalRoutePolicy(func(userInput, currentState GameState, edges []g.Edge[GameState]) g.Edge[GameState] {
		for _, edge := range edges {
			if currentState.Success {
				if label, ok := edge.LabelByKey("path"); ok && label == "success" {
					return edge
				}
			} else {
				if label, ok := edge.LabelByKey("path"); ok && label == "fail" {
					return edge
				}
			}
		}
		return nil
	})
	router, err := b.CreateRouter("CheckRouter", routingPolicy)
	if err != nil {
		log.Fatalf("Router creation failed: %v", err)
	}

	// Build graph
	startEdge := b.CreateStartEdge(initNode)
	stateMonitorCh := make(chan g.StateMonitorEntry[GameState], 10)
	g, _ := b.CreateRuntime(startEdge, stateMonitorCh)
	defer g.Shutdown()

	g.AddEdge(
		b.CreateEdge(initNode, guessNode),
		b.CreateEdge(guessNode, router),
		b.CreateEdge(router, hintNode, map[string]string{"path": "fail"}),
		b.CreateEdge(hintNode, guessNode), // Loop back
		b.CreateEndEdge(router, map[string]string{"path": "success"}),
	)

	if err := g.Validate(); err != nil {
		log.Fatalf("Validation failed: %v", err)
	}

	g.Invoke(GameState{})

	for entry := range stateMonitorCh {
		if !entry.Running {
			fmt.Printf("✅ Success! Target was %d, found in %d tries\n", entry.CurrentState.Target, entry.CurrentState.Tries)
			break
		}
	}
}
