package main

import (
	"bufio"
	"fmt"
	"log"
	"os"
	"strings"

	b "github.com/morphy76/ggraph/pkg/builders"
	g "github.com/morphy76/ggraph/pkg/graph"
	"github.com/morphy76/ggraph/pkg/llm"
	"github.com/morphy76/ggraph/pkg/llm/openai"
)

func main() {
	// Create OpenAI chat node
	chatNode, err := openai.CreateOpenAIChatNodeFromEnvironment("ChatNode", "gpt-5-nano")
	if err != nil {
		log.Fatalf("Failed to create chat node: %v", err)
	}

	// Build simple graph: chat -> end
	startEdge := b.CreateStartEdge(chatNode)
	stateMonitorCh := make(chan g.StateMonitorEntry[llm.AgentModel], 10)

	// Initialize with system message
	initialState := llm.CreateModel(
		llm.CreateMessage(llm.System, "You are a helpful cooking assistant. You provide advice, recipes, and tips about cooking. Keep responses concise and friendly."),
	)

	g, err := b.CreateRuntimeWithInitialState(startEdge, stateMonitorCh, initialState)
	if err != nil {
		log.Fatalf("Runtime creation failed: %v", err)
	}
	defer g.Shutdown()

	g.AddEdge(b.CreateEndEdge(chatNode))

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
		currentState.AddUserMessage(text)

		// Invoke graph with current state
		g.Invoke(currentState)

		// Wait for completion and get response
		wasStreaming := false
		fmt.Print("🤖 Assistant: ")
		for entry := range stateMonitorCh {
			if entry.Error != nil {
				fmt.Printf("⚠️  Error: %v\n", entry.Error)
				break
			}

			if entry.Partial {
				wasStreaming = true
				currentState = entry.CurrentState
				messages := currentState.Messages
				if len(messages) > 0 {
					lastMsg := messages[len(messages)-1]
					if lastMsg.Role == llm.Assistant {
						fmt.Print(lastMsg.Content)
					}
				}
			}

			if !entry.Running {
				if wasStreaming {
					fmt.Printf("\n\n")
				} else {
					currentState = entry.CurrentState
					messages := currentState.Messages
					if len(messages) > 0 {
						lastMsg := messages[len(messages)-1]
						if lastMsg.Role == llm.Assistant {
							fmt.Printf("%s\n\n", lastMsg.Content)
						}
					}
				}
				break
			}
		}
	}
}
