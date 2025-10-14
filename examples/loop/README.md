# Loop Example - Number Guessing Game

This example demonstrates a **loop graph** with a number guessing game using binary search.

## Prompt

Create an example using the other examples as inspiration which demonstrates a loop graph where:

- the first node determines a number
- the second node tries to guess
- on success, exit
- on fail, an additional node provides a hint to improve the next try to the guess node
- be concise
- use a stateful graph

## Flow

1. **InitNode**: Sets a random target number (1-100)
2. **GuessNode**: Makes a guess using binary search
3. **CheckRouter**: Evaluates if the guess is correct
   - ✅ **Success**: Routes to EndNode
   - ❌ **Fail**: Routes to HintNode
4. **HintNode**: Provides hint (higher/lower) and adjusts search range
5. **Loop back** to GuessNode

## Key Features

- **Stateful graph** with merge function
- **Conditional routing** based on success/failure
- **Loop structure** for iterative guessing
- **Binary search** ensures convergence in ~7 tries

## Run

```bash
go run ./examples/loop/run.go
```
