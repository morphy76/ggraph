package main

import (
	"bufio"
	"fmt"
	"log"
	"os"
	"strings"

	"github.com/morphy76/ggraph/pkg/graph"
	"github.com/morphy76/ggraph/pkg/llm/ollama"
	"github.com/ollama/ollama/api"
)

func main() {
	// Set Ollama environment to localhost
	os.Setenv("OLLAMA_HOST", "http://localhost:11434")

	// Create Ollama chat node
	chatNode, err := ollama.CreateOLLamaChatNodeFromEnvironment("ChatNode", "Almawave/Velvet:2B")
	if err != nil {
		log.Fatalf("Failed to create chat node: %v", err)
	}

	// Build simple graph: chat -> end
	startEdge := graph.CreateStartEdge(chatNode)
	stateMonitorCh := make(chan graph.StateMonitorEntry[ollama.ChatModel], 10)

	// Initialize with system message
	initialState := ollama.NewChatModel(api.Message{
		Role:    "system",
		Content: "You are a helpful cooking assistant. You provide advice, recipes, and tips about cooking. Keep responses concise and friendly.",
	})

	g, err := graph.CreateRuntimeWithMergerAndInitialState(startEdge, stateMonitorCh, ollama.MergeChatModels, initialState)
	if err != nil {
		log.Fatalf("Runtime creation failed: %v", err)
	}
	defer g.Shutdown()

	g.AddEdge(graph.CreateEndEdge(chatNode))

	if err := g.Validate(); err != nil {
		log.Fatalf("Validation failed: %v", err)
	}

	fmt.Println("🍳 Cooking Chat Assistant")
	fmt.Println("💡 Ask me anything about cooking! Type 'exit' to quit.")
	fmt.Println()

	reader := bufio.NewReader(os.Stdin)
	currentState := initialState

	// Main conversation loop
	for {
		fmt.Print("💬 You: ")
		text, _ := reader.ReadString('\n')
		text = strings.TrimSpace(text)

		if text == "exit" || text == "" {
			fmt.Println("\n👋 Goodbye! Happy cooking!")
			break
		}

		// Add user message
		currentState.AddMessage(api.Message{
			Role:    "user",
			Content: text,
		})

		// Invoke graph with current state
		g.Invoke(currentState)

		// Wait for completion and get response
		for entry := range stateMonitorCh {
			if entry.Error != nil {
				fmt.Printf("⚠️  Error: %v\n", entry.Error)
				break
			}

			if entry.Partial {
				currentState = entry.CurrentState
				messages := currentState.Messages()
				if len(messages) > 0 {
					lastMsg := messages[len(messages)-1]
					if lastMsg.Role == "assistant" {
						fmt.Print(lastMsg.Content)
					}
				}
			}

			if !entry.Running {
				currentState = entry.CurrentState
				messages := currentState.Messages()
				if len(messages) > 0 {
					lastMsg := messages[len(messages)-1]
					if lastMsg.Role == "assistant" {
						fmt.Printf("🤖 Assistant: %s\n\n", lastMsg.Content)
					}
				}
				break
			}
		}
	}
}
