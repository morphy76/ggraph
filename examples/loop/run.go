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

func merge(original, new GameState) GameState {
	if new.Target != 0 {
		original.Target = new.Target
	}
	if new.Guess != 0 {
		original.Guess = new.Guess
	}
	if new.Tries > 0 {
		original.Tries = new.Tries
	}
	original.Success = new.Success
	if new.Hint != "" {
		original.Hint = new.Hint
	}
	if new.Low > 0 {
		original.Low = new.Low
	}
	if new.High > 0 {
		original.High = new.High
	}
	return original
}

func main() {
	// Node 1: Determine target number
	initNode, _ := b.CreateNode("InitNode", func(state GameState, notify func(GameState)) (GameState, error) {
		state.Target = rand.Intn(100) + 1
		state.Tries = 0
		state.Low = 1
		state.High = 100
		fmt.Printf("🎯 Target set (hidden)\n")
		return state, nil
	})

	// Node 2: Make a guess using binary search
	guessNode, _ := b.CreateNode("GuessNode", func(state GameState, notify func(GameState)) (GameState, error) {
		state.Tries++
		state.Guess = (state.Low + state.High) / 2
		state.Success = (state.Guess == state.Target)
		fmt.Printf("🤔 Try #%d: Guessed %d (range: %d-%d)\n", state.Tries, state.Guess, state.Low, state.High)
		return state, nil
	})

	// Node 3: Provide hint and adjust range
	hintNode, _ := b.CreateNode("HintNode", func(state GameState, notify func(GameState)) (GameState, error) {
		if state.Guess < state.Target {
			state.Low = state.Guess + 1
			state.Hint = "higher"
			fmt.Printf("💡 Hint: Try higher!\n")
		} else {
			state.High = state.Guess - 1
			state.Hint = "lower"
			fmt.Printf("💡 Hint: Try lower!\n")
		}
		return state, nil
	})

	// Router: Check success
	routingPolicy, _ := b.CreateConditionalRoutePolicy(func(state GameState, edges []g.Edge[GameState]) g.Edge[GameState] {
		for _, edge := range edges {
			if state.Success {
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
	router, _ := b.CreateRouter("CheckRouter", routingPolicy)

	// End node for success
	endNode, _ := b.CreateNode("EndNode", func(state GameState, notify func(GameState)) (GameState, error) {
		fmt.Printf("🎉 Correct! The answer was %d\n", state.Target)
		return state, nil
	})

	// Build graph
	startEdge := b.CreateStartEdge(initNode)
	stateMonitorCh := make(chan g.StateMonitorEntry[GameState], 10)
	g, _ := b.CreateRuntimeWithMergerAndInitialState(startEdge, stateMonitorCh, merge, GameState{})
	defer g.Shutdown()

	g.AddEdge(
		b.CreateEdge(initNode, guessNode),
		b.CreateEdge(guessNode, router),
		b.CreateEdge(router, hintNode, map[string]string{"path": "fail"}),
		b.CreateEdge(hintNode, guessNode), // Loop back
		b.CreateEdge(router, endNode, map[string]string{"path": "success"}),
		b.CreateEndEdge(endNode),
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
