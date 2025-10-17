package main

import (
	"encoding/json"
	"fmt"
	"log"
	"math/rand"
	"os"

	"github.com/google/uuid"
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

// FilePersistence implements simple file-based persistence for the game state
type FilePersistence struct {
	filepath  string
	runtimeID uuid.UUID
}

func NewFilePersistence(filepath string) *FilePersistence {
	return &FilePersistence{
		filepath: filepath,
	}
}

// Persist writes the state to a JSON file
func (fp *FilePersistence) Persist(state GameState) error {
	data, err := json.MarshalIndent(state, "", "  ")
	if err != nil {
		return fmt.Errorf("failed to marshal state: %w", err)
	}

	err = os.WriteFile(fp.filepath, data, 0644)
	if err != nil {
		return fmt.Errorf("failed to write state file: %w", err)
	}

	fmt.Printf("💾 State persisted to %s\n", fp.filepath)
	return nil
}

// Restore reads the state from a JSON file
func (fp *FilePersistence) Restore() (GameState, error) {
	data, err := os.ReadFile(fp.filepath)
	if err != nil {
		if os.IsNotExist(err) {
			return GameState{}, fmt.Errorf("no saved state found: %w", err)
		}
		return GameState{}, fmt.Errorf("failed to read state file: %w", err)
	}

	var state GameState
	err = json.Unmarshal(data, &state)
	if err != nil {
		return GameState{}, fmt.Errorf("failed to unmarshal state: %w", err)
	}

	fmt.Printf("📂 State restored from %s\n", fp.filepath)
	return state, nil
}

func main() {
	stateFile := "game_state.json"
	persistence := NewFilePersistence(stateFile)

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

	// End node for success
	endNode := b.CreateEndNode[GameState]()

	// Build graph
	startEdge := b.CreateStartEdge(initNode)
	stateMonitorCh := make(chan g.StateMonitorEntry[GameState], 10)
	runtime, _ := b.CreateRuntime(startEdge, stateMonitorCh)
	defer runtime.Shutdown()

	runtime.AddEdge(
		b.CreateEdge(initNode, guessNode),
		b.CreateEdge(guessNode, router),
		b.CreateEdge(router, hintNode, map[string]string{"path": "fail"}),
		b.CreateEdge(hintNode, guessNode), // Loop back
		b.CreateEdge(router, endNode, map[string]string{"path": "success"}),
		b.CreateEndEdge(endNode),
	)

	if err := runtime.Validate(); err != nil {
		log.Fatalf("Validation failed: %v", err)
	}

	// Configure persistence with runtime ID
	runtimeID := uuid.New()
	runtime.SetPersistentState(
		persistence.Persist,
		persistence.Restore,
		runtimeID,
	)

	// Try to restore previous state
	fmt.Println("\n=== Attempting to restore previous state ===")
	if err := runtime.Restore(); err != nil {
		fmt.Printf("⚠️  No previous state to restore: %v\n", err)
		fmt.Println("Starting fresh game...")
	} else {
		fmt.Println("✅ Previous state restored! Continuing from where we left off...")
	}

	runtime.Invoke(GameState{})

	for entry := range stateMonitorCh {
		if !entry.Running {
			if entry.Error == nil {
				fmt.Printf("✅ Success! Target was %d, found in %d tries\n", entry.CurrentState.Target, entry.CurrentState.Tries)
				fmt.Printf("\n💡 Try running this example again to see persistence in action!\n")
				fmt.Printf("   The state is saved at each step in: %s\n", stateFile)
			} else {
				fmt.Printf("❌ Error: %v\n", entry.Error)
			}
			break
		}
	}
}
