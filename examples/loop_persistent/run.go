package main

import (
	"context"
	"encoding/json"
	"fmt"
	"log"
	"math/rand"
	"os"
	"path/filepath"

	"github.com/google/uuid"
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

type filePersistence struct {
	baseDir string
}

func newFilePersistence(baseDir string) *filePersistence {
	if err := os.MkdirAll(baseDir, 0755); err != nil {
		log.Printf("Failed to create base directory: %v", err)
	}
	return &filePersistence{
		baseDir: baseDir,
	}
}

func (fp *filePersistence) getFilePath(threadID string) string {
	return filepath.Join(fp.baseDir, fmt.Sprintf("game_state_%s.json", threadID))
}

func (fp *filePersistence) Persist(ctx context.Context, threadID string, state gameState) error {
	filePath := fp.getFilePath(threadID)

	data, err := json.MarshalIndent(state, "", "  ")
	if err != nil {
		return fmt.Errorf("failed to marshal state: %w", err)
	}

	err = os.WriteFile(filePath, data, 0644)
	if err != nil {
		return fmt.Errorf("failed to write state file: %w", err)
	}

	fmt.Printf("💾 State persisted to %s\n", filePath)
	return nil
}

func (fp *filePersistence) Restore(ctx context.Context, threadID string) (gameState, error) {
	filePath := fp.getFilePath(threadID)

	data, err := os.ReadFile(filePath)
	if err != nil {
		if os.IsNotExist(err) {
			return gameState{}, fmt.Errorf("no saved state found: %w", err)
		}
		return gameState{}, fmt.Errorf("failed to read state file: %w", err)
	}

	var state gameState
	err = json.Unmarshal(data, &state)
	if err != nil {
		return gameState{}, fmt.Errorf("failed to unmarshal state: %w", err)
	}

	fmt.Printf("📂 State restored from %s\n", filePath)
	return state, nil
}

var _ g.Memory[gameState] = (*fileMemory)(nil)

type fileMemory struct {
	persistence *filePersistence
}

func (m *fileMemory) PersistFn() g.PersistFn[gameState] {
	return m.persistence.Persist
}

func (m *fileMemory) RestoreFn() g.RestoreFn[gameState] {
	return m.persistence.Restore
}

func main() {
	stateDir := "game_states"
	persistence := newFilePersistence(stateDir)
	memory := &fileMemory{persistence: persistence}

	// Node 1: Determine target number
	initNode, _ := b.NewNode("InitNode", func(userInput gameState, currentState gameState, notifyPartial g.NotifyPartialFn[gameState]) (gameState, error) {
		currentState.Target = rand.Intn(100) + 1
		currentState.Tries = 0
		currentState.Low = 1
		currentState.High = 100
		fmt.Printf("🎯 Target set (hidden)\n")
		return currentState, nil
	})

	// Node 2: Make a guess using binary search
	guessNode, _ := b.NewNode("GuessNode", func(userInput gameState, currentState gameState, notifyPartial g.NotifyPartialFn[gameState]) (gameState, error) {
		currentState.Tries++
		currentState.Guess = (currentState.Low + currentState.High) / 2
		currentState.Success = (currentState.Guess == currentState.Target)
		fmt.Printf("🤔 Try #%d: Guessed %d (range: %d-%d)\n", currentState.Tries, currentState.Guess, currentState.Low, currentState.High)
		return currentState, nil
	})

	// Node 3: Provide hint and adjust range
	hintNode, _ := b.NewNode("HintNode", func(userInput gameState, currentState gameState, notifyPartial g.NotifyPartialFn[gameState]) (gameState, error) {
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
	stateMonitorCh := make(chan g.StateMonitorEntry[gameState], 10)
	runtime, _ := b.CreateRuntime(startEdge, stateMonitorCh, g.WithMemory(memory))
	defer runtime.Shutdown()

	runtime.AddEdge(
		b.CreateEdge(initNode, guessNode),
		b.CreateEdge(guessNode, router),
		b.CreateEdge(router, hintNode, map[string]string{"path": "fail"}),
		b.CreateEdge(hintNode, guessNode), // Loop back
		b.CreateEndEdge(router, map[string]string{"path": "success"}),
	)

	if err := runtime.Validate(); err != nil {
		log.Fatalf("Validation failed: %v", err)
	}

	// Configure persistence with runtime ID
	invokeConfig := g.InvokeConfigThreadID(uuid.NewString())

	// The runtime will automatically try to restore state on first invoke
	// If no state exists, it will start with initial state
	fmt.Println("\n=== Starting game (will auto-restore if previous state exists) ===")

	runtime.Invoke(gameState{}, invokeConfig)

	for entry := range stateMonitorCh {
		if !entry.Running {
			if entry.Error == nil {
				fmt.Printf("✅ Success! Target was %d, found in %d tries\n", entry.NewState.Target, entry.NewState.Tries)
				fmt.Printf("\n💡 Try running this example again to see persistence in action!\n")
				fmt.Printf("   The state is saved asynchronously in: %s/\n", stateDir)
			} else {
				fmt.Printf("❌ Error: %v\n", entry.Error)
			}
			break
		}
	}
}
