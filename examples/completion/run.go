package main

import (
	"bufio"
	"context"
	"fmt"
	"log"
	"os"
	"strings"

	"github.com/openai/openai-go/v3"

	a "github.com/morphy76/ggraph/pkg/agent"
	"github.com/morphy76/ggraph/pkg/agent/aiw"
	o "github.com/morphy76/ggraph/pkg/agent/openai"
	b "github.com/morphy76/ggraph/pkg/builders"
	g "github.com/morphy76/ggraph/pkg/graph"
)

// CompletionNodeFn creates a completion node that generates text based on a prompt using chat completions
var CompletionNodeFn o.CompletionNodeFn = func(completionService openai.CompletionService, model string, completionOptions ...a.ModelOption) g.NodeFn[a.Completion] {
	return func(userInput, currentState a.Completion, notify g.NotifyPartialFn[a.Completion]) (a.Completion, error) {
		// Check if there's a user prompt in the conversation, or use a default
		prompt := "Use the family guy style: answer in the same way Peter Griffin would. The request is within boundaries defined by !!![ and ]!!!. Do not put boundaries in the final answer. Requested completion is:\n!!![\n%s\n]!!!"

		useOpts, err := a.CreateCompletionOptions(model, fmt.Sprintf(prompt, userInput.Text), completionOptions...)
		if err != nil {
			return currentState, fmt.Errorf("failed to create completion options: %w", err)
		}

		openAIOpts := o.ConvertCompletionOptions(useOpts)
		resp, err := completionService.New(context.Background(), openAIOpts)
		if err != nil {
			return currentState, fmt.Errorf("failed to generate completion: %w", err)
		}

		// Print troubleshooting info for the response
		fmt.Printf("=== Response Debug Info ===\n")
		fmt.Printf("Response ID: %s\n", resp.ID)
		fmt.Printf("Model: %s\n", resp.Model)
		fmt.Printf("Object: %s\n", resp.Object)
		fmt.Printf("Created: %d\n", resp.Created)
		fmt.Printf("Number of choices: %d\n", len(resp.Choices))
		if len(resp.Choices) > 0 {
			fmt.Printf("First choice text length: %d chars\n", len(resp.Choices[0].Text))
			fmt.Printf("First choice finish reason: %s\n", resp.Choices[0].FinishReason)
		}
		fmt.Printf("===========================\n")

		// Update the current state with the final completion
		if len(resp.Choices) > 0 {
			currentState = a.CreateCompletion(resp.Choices[0].Text)
		} else {
			return currentState, fmt.Errorf("no completion choices returned")
		}

		return currentState, nil
	}
}

func main() {
	fmt.Println("=== Completion Agent Example ===")
	fmt.Println()

	// Get the API key from environment variable
	apiKey := aiw.PATFromEnv()
	if apiKey == "" {
		log.Fatal("Environment variable not set to fetch the API key.")
	}

	client := aiw.NewAIWClient(apiKey)

	// Create the completion node
	completionNode, err := aiw.CreateCompletionNode(
		"CompletionNode",
		"velvet-2b",
		client,
		CompletionNodeFn,
		a.WithTemperature(0.7),
		a.WithFrequencyPenalty(0.5),
	)
	if err != nil {
		log.Fatalf("Failed to create completion node: %v", err)
	}

	// Create edges connecting the nodes
	startEdge := b.CreateStartEdge(completionNode)
	endEdge := b.CreateEndEdge(completionNode)

	// Initialize the conversation state
	stateMonitorCh := make(chan g.StateMonitorEntry[a.Completion], 10)

	// Create the runtime graph
	graph, err := b.CreateRuntime(startEdge, stateMonitorCh)
	if err != nil {
		log.Fatalf("Runtime creation failed: %v", err)
	}
	defer graph.Shutdown()

	// Add edges to the graph
	graph.AddEdge(endEdge)

	// Validate the graph
	err = graph.Validate()
	if err != nil {
		log.Fatalf("Graph validation failed: %v", err)
	}

	// Interactive prompt input
	scanner := bufio.NewScanner(os.Stdin)

	fmt.Println("Completion Agent is ready!")
	fmt.Println("Enter your text prompts for completion. Type 'quit', 'exit', or 'q' to stop.")
	fmt.Println("Press Enter without text to use the default prompt.")
	fmt.Println("━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━")

	promptCount := 0
	for {
		promptCount++
		fmt.Printf("\n[Prompt %d] Enter your text: ", promptCount)

		if !scanner.Scan() {
			break
		}

		input := strings.TrimSpace(scanner.Text())

		// Check for exit commands
		if input == "quit" || input == "exit" || input == "q" {
			fmt.Println("\nGoodbye!")
			break
		}

		// Display the input
		if input == "" {
			fmt.Println("Using default prompt...")
		} else {
			fmt.Printf("Processing: %s\n", input)
		}
		fmt.Println()

		// Invoke the graph
		graph.Invoke(a.CreateCompletion(input))

		// Monitor the execution
		for {
			entry := <-stateMonitorCh

			fmt.Printf("[%s - Running: %t]\n", entry.Node, entry.Running)

			if entry.Error != nil {
				fmt.Printf("❌ Error: %v\n", entry.Error)
				break
			}

			if !entry.Running {
				// Display the completion result
				fmt.Printf("✅ Generated completion: %s\n", entry.NewState.Text)
				fmt.Println("━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━")
				break
			}
		}
	}

	if scanner.Err() != nil {
		log.Printf("Error reading input: %v", scanner.Err())
	}

	fmt.Println("Session completed.")
}
