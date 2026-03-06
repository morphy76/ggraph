package main

import (
	"bufio"
	"fmt"
	"log"
	"os"
	"strings"

	"github.com/openai/openai-go/v3"

	a "github.com/morphy76/ggraph/pkg/agent"
	o "github.com/morphy76/ggraph/pkg/agent/openai"
	b "github.com/morphy76/ggraph/pkg/builders"
	g "github.com/morphy76/ggraph/pkg/graph"
)

func main() {
	fmt.Println("=== Completion Agent Example ===")
	fmt.Println()

	// Get the API key from environment variable
	apiKey := o.APIKeyFromEnv()
	if apiKey == "" {
		log.Fatal("Environment variable not set to fetch the API key.")
	}

	client := o.NewOpenAIClient(apiKey)

	// Create the completion node
	completionNode, err := o.CreateCompletionNode(
		"CompletionNode",
		openai.ChatModelGPT4_1Nano,
		client,
		a.WithPromptFormat("Use the family guy style: answer in the same way Peter Griffin would. The request is within boundaries defined by !!![ and ]!!!. Do not put boundaries in the final answer. Requested completion is:\n!!![\n%s\n]!!!"),
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
