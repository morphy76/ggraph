package main

import (
	"bufio"
	"fmt"
	"log"
	"os"

	a "github.com/morphy76/ggraph/pkg/agent"
	aiw "github.com/morphy76/ggraph/pkg/agent/aiw"
	o "github.com/morphy76/ggraph/pkg/agent/openai"
	b "github.com/morphy76/ggraph/pkg/builders"
	g "github.com/morphy76/ggraph/pkg/graph"
)

func main() {
	fmt.Println("=== Conversation Summarization Example ===")
	fmt.Println("This example demonstrates how an agent can automatically summarize")
	fmt.Println("its conversation history when it gets too long, saving tokens and context.")
	fmt.Println("The threshold is artificially set low (3 messages) to demonstrate.")

	// Get AIW API key from environment
	pat := aiw.PATFromEnv()
	if pat == "" {
		log.Fatal("AIW_API_KEY environment variable not set; visit https://portal.aiwave.ai to get your API key.")
	}

	aiwClient := aiw.NewAIWClient(pat)

	// Configure Summarization with very low thresholds for demonstration
	summaryConfig := &a.SummarizationConfig{
		Enabled:          true,
		MessageThreshold: 3, // Summarize after 3 messages
		KeepRecentCount:  1, // Keep only 1 recent message
	}

	// Create a single node that responds to the user and applies summarization
	chatNode, err := o.CreateConversationNode(
		"ChatNode",
		"velvet-14b",
		aiwClient,
		a.WithMessages(
			a.CreateMessage(a.System, "Sei un assistente utile che risponde in modo conciso in italiano."),
		),
		a.WithSummarization(summaryConfig),
	)
	if err != nil {
		log.Fatalf("Failed to create chat node: %v", err)
	}

	// Interactive terminal node
	terminalNode, err := b.NewNode("TerminalNode", func(userInput, currentState a.Conversation, notify g.NotifyPartialFn[a.Conversation]) (a.Conversation, error) {
		fmt.Println("\n[Current Conversation Messages Length]:", len(currentState.Messages))
		for i, m := range currentState.Messages {
			fmt.Printf(" MSG %d [%v]: %.80s...\n", i, m.Role, m.Content)
		}

		fmt.Print("\n[User]: ")
		reader := bufio.NewReader(os.Stdin)
		input, _ := reader.ReadString('\n')

		return a.CreateConversation(
			a.CreateMessage(a.User, input),
		), nil
	})

	// Create edges connecting the nodes in a loop
	startEdge := b.CreateStartEdge(terminalNode)
	terminalToChat := b.CreateEdge(terminalNode, chatNode)
	chatToTerminal := b.CreateEdge(chatNode, terminalNode)

	initialState := a.CreateConversation()
	stateMonitorCh := make(chan g.StateMonitorEntry[a.Conversation], 100)

	graph, err := b.CreateRuntime(startEdge, stateMonitorCh, g.WithInitialState(initialState))
	if err != nil {
		log.Fatalf("Runtime creation failed: %v", err)
	}
	defer graph.Shutdown()

	graph.AddEdge(terminalToChat, chatToTerminal)

	err = graph.Validate()
	if err != nil {
		log.Fatalf("Graph validation failed: %v", err)
	}

	// Consume and ignore state monitor events in background
	go func() {
		for entry := range stateMonitorCh {
			if entry.Error != nil {
				fmt.Printf("\n[Error in node %s]: %v\n", entry.Node, entry.Error)
			}
		}
	}()

	// Start graph (it will loop indefinitely between terminal and chat)
	userInput := a.CreateConversation(a.CreateMessage(a.User, "Start"))
	graph.Invoke(userInput)
}
