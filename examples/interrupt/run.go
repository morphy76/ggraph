package main

import (
	"context"
	"fmt"
	"log"
	"math/rand"

	b "github.com/morphy76/ggraph/pkg/builders"
	g "github.com/morphy76/ggraph/pkg/graph"
)

var _ g.SharedState = (*gameState)(nil)

type gameState struct {
	Target  int
	Guess   int
	Tries   int
	Success bool
	Hint    string
	Low     int
	High    int
}

func main() {
	// Node 1: Determine target number (only if not already set from restored state)
	initNode, _ := b.NewNodeBuilder("InitNode", func(userInput gameState, currentState gameState, notifyPartial g.NotifyPartialFn[gameState]) (gameState, error) {
		// If target is already set (from restored state), don't reinitialize
		if currentState.Target != 0 {
			fmt.Printf("🔄 Resuming game with existing target\n")
			return currentState, nil
		}
		currentState.Target = rand.Intn(100) + 1
		currentState.Tries = 0
		currentState.Low = 1
		currentState.High = 100
		fmt.Printf("🎯 Target set (hidden)\n")
		return currentState, nil
	}).Build()

	// Node 2: Make a guess using binary search
	guessNode, _ := b.NewNodeBuilder("GuessNode", func(userInput gameState, currentState gameState, notifyPartial g.NotifyPartialFn[gameState]) (gameState, error) {
		currentState.Tries++
		currentState.Guess = (currentState.Low + currentState.High) / 2
		currentState.Success = (currentState.Guess == currentState.Target)
		fmt.Printf("🤔 Try #%d: Guessed %d (range: %d-%d)\n", currentState.Tries, currentState.Guess, currentState.Low, currentState.High)
		return currentState, nil
	}).Build()

	// Node 3: Provide hint and adjust range
	hintNode, _ := b.NewNodeBuilder("HintNode", func(userInput gameState, currentState gameState, notifyPartial g.NotifyPartialFn[gameState]) (gameState, error) {
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
	}).Build()

	// Router: Check success
	routingPolicy, _ := b.CreateConditionalRoutePolicy(func(userInput, currentState gameState, edges []g.Edge[gameState]) g.Edge[gameState] {
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
	stateMonitorCh := make(chan g.StateMonitorEntry[gameState], 100)
	graph, _ := b.CreateRuntime(startEdge, stateMonitorCh)
	defer graph.Shutdown()

	graph.AddEdge(
		b.CreateEdge(initNode, guessNode),
		b.CreateEdge(guessNode, router),
		b.CreateEdge(router, hintNode, map[string]string{"path": "fail"}),
		b.CreateEdge(hintNode, guessNode), // Loop back
		b.CreateEndEdge(router, map[string]string{"path": "success"}),
	)

	if err := graph.Validate(); err != nil {
		log.Fatalf("Validation failed: %v", err)
	}

	memory := b.NewMemMemory[gameState]()
	graph.SetPersistentState(
		memory.PersistFn(),
		memory.RestoreFn(),
	)

	ctx, cancel := context.WithCancel(context.Background())
	monitor := func() {
		for entry := range stateMonitorCh {
			fmt.Printf("📊 Monitor: Node=%s, Running=%v, Error=%v, State=%+v\n", entry.Node, entry.Running, entry.Error, entry.NewState)
			if !entry.Running {
				if entry.Error != nil {
					fmt.Printf("❌ Error: %v\n", entry.Error)
				} else {
					fmt.Printf("✅ Success! Target was %d, found in %d tries\n", entry.NewState.Target, entry.NewState.Tries)
				}
				break
			} else if entry.Node == "HintNode" {
				cancel()
				fmt.Printf("🛑 Cancellation requested!\n")
				break
			}
		}
	}

	go monitor()

	threadID := graph.Invoke(
		gameState{},
		g.InvokeConfigContext(ctx),
	)

	fmt.Printf("🚀 Game stopped! Memory state: %+v\n", memory)
	<-ctx.Done()

	// Create a new context for resuming execution
	graph.Invoke(
		gameState{},
		g.InvokeConfigThreadID(threadID),
		g.InvokeConfigContext(context.Background()),
	)

	for entry := range stateMonitorCh {
		fmt.Printf("📊 Monitor: Node=%s, Running=%v, Error=%v, State=%+v\n", entry.Node, entry.Running, entry.Error, entry.NewState)
		if !entry.Running {
			if entry.Error != nil {
				fmt.Printf("❌ Error: %v\n", entry.Error)
			} else {
				fmt.Printf("✅ Success! Target was %d, found in %d tries\n", entry.NewState.Target, entry.NewState.Tries)
			}
			break
		}
	}
}
