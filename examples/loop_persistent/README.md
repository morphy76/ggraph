# Loop with Persistence Example - Number Guessing Game

This example demonstrates a **loop graph with persistence** using a number guessing game with binary search. It showcases how to implement and use the `Persistent` interface to save and restore graph runtime state.

## Prompt

```text
create a loop example which demonstrates the Persistence interface of a graph runtime:

- implement a persistent state based on a simple text file
- demonstrate both the persist and the restore capability
- do not touch business code
- add an example readme with some implementation details
- use the other examples as a pattern, do not create a go module
```

## Overview

This example extends the basic loop example by adding state persistence capabilities. The game state is automatically saved to a file after each node execution, allowing the game to be interrupted and resumed later from where it left off.

## Persistence Interface

The `Persistent` interface in ggraph provides two main capabilities:

1. **SetPersistentState**: Configure persistence with custom save/restore functions
2. **Restore**: Restore previously saved state before starting graph execution

## Implementation Details

### FilePersistence

A simple file-based persistence implementation that:

- **Stores state as JSON** - Human-readable format for easy debugging
- **Automatic persistence** - State is saved after each node execution
- **Error handling** - Gracefully handles missing or corrupted state files
- **File location** - Saves to `game_state.json` in the current directory

```go
type FilePersistence struct {
    filepath   string
    runtimeID  uuid.UUID
}

// Persist writes the state to a JSON file
func (fp *FilePersistence) Persist(state GameState) error {
    // Marshal state to JSON and write to file
}

// Restore reads the state from a JSON file
func (fp *FilePersistence) Restore() (GameState, error) {
    // Read file and unmarshal JSON to state
}
```

### Integration with Runtime

The runtime is configured with persistence functions using the `SetPersistentState` method:

```go
persistence := NewFilePersistence("game_state.json")
runtimeID := uuid.New()

runtime.SetPersistentState(
    persistence.Persist,  // Called after each node execution
    persistence.Restore,  // Called manually before Invoke
    runtimeID,           // Unique identifier for this runtime instance
)

// Attempt to restore previous state
if err := runtime.Restore(); err != nil {
    fmt.Println("No previous state, starting fresh")
}

runtime.Invoke(GameState{})
```

## Flow

1. **InitNode**: Sets a random target number (1-100) *(persisted)*
2. **GuessNode**: Makes a guess using binary search *(persisted)*
3. **CheckRouter**: Evaluates if the guess is correct *(persisted)*
   - ✅ **Success**: Routes to EndNode
   - ❌ **Fail**: Routes to HintNode
4. **HintNode**: Provides hint (higher/lower) and adjusts search range *(persisted)*
5. **Loop back** to GuessNode

Each step automatically persists the current state to disk.

## Key Features

- **Stateful graph** with automatic persistence
- **File-based storage** using JSON format
- **Resumable execution** - Interrupt and continue later
- **Conditional routing** based on success/failure
- **Loop structure** for iterative guessing
- **Binary search** ensures convergence in ~7 tries

## Run

```bash
# First run - completes the game
go run ./examples/loop_persistent/run.go

# The state is saved at each step in game_state.json

# Second run - starts fresh (game completes)
go run ./examples/loop_persistent/run.go

# To test restoration, you can manually interrupt the program mid-execution
# or modify the code to limit iterations
```

## Testing Persistence

To see persistence in action:

1. **Run the example once** - Let it complete and observe the state file created
2. **Inspect the state file** - `cat game_state.json` to see the saved state
3. **Run again** - The previous completed state will be detected and reported
4. **Manual testing** - Modify the code to exit early and run again to see mid-game restoration

## State File Structure

The `game_state.json` file contains:

```json
{
  "Target": 42,
  "Guess": 50,
  "Tries": 3,
  "Success": false,
  "Hint": "lower",
  "Low": 25,
  "High": 49
}
```

## Important Notes

- **No business logic changes** - Persistence is completely separate from game logic
- **Automatic saves** - The runtime handles persistence after each node
- **Thread-safe** - State merging is protected by mutex locks
- **Non-fatal errors** - Persistence errors don't stop graph execution
- **Runtime ID** - Used to identify specific runtime instances (useful for multi-instance scenarios)
