package main

import (
	"context"
	"fmt"
	"log"

	"github.com/openai/openai-go/v3"

	a "github.com/morphy76/ggraph/pkg/agent"
	o "github.com/morphy76/ggraph/pkg/agent/openai"
	t "github.com/morphy76/ggraph/pkg/agent/tool"
	b "github.com/morphy76/ggraph/pkg/builders"
	g "github.com/morphy76/ggraph/pkg/graph"
)

var _ g.SharedState = (*myState)(nil)

type myState struct {
	Message string
}

func additionTool(a int, b int) (int, error) {
	return a + b, nil
}

func divisionTool(a int, b int) (int, error) {
	if b == 0 {
		return 0, fmt.Errorf("division by zero")
	}
	return a / b, nil
}

func main() {
	apiKey := o.APIKeyFromEnv()
	if apiKey == "" {
		log.Fatal("API key environment variable not set.")
	}

	client := o.NewOpenAIClient(apiKey)

	llmFn := func(chatService openai.ChatService, model string, conversationOptions ...a.ModelOption) g.NodeFn[a.Conversation] {
		return func(userInput, currentState a.Conversation, notify g.NotifyPartialFn[a.Conversation]) (a.Conversation, error) {

			useOpts, err := a.CreateConversationOptions(
				model,
				userInput.Messages,
				conversationOptions...,
			)
			if err != nil {
				return currentState, fmt.Errorf("failed to create conversation options: %w", err)
			}

			openAIOpts := o.ConvertConversationOptions(useOpts)

			resp, err := chatService.Completions.New(context.Background(), openAIOpts)
			if err != nil {
				return currentState, fmt.Errorf("failed to generate tool calls: %w", err)
			}

			answer := resp.Choices[0].Message.Content
			currentState.Messages = append(userInput.Messages,
				a.CreateMessage(a.Assistant, answer))

			return currentState, nil
		}
	}

	tool1, err := t.CreateTool[int](additionTool, "Prompt: Adds two integers.", "Input: (int, int)", "Output: int")
	if err != nil {
		log.Fatalf("Failed to create addition tool: %v", err)
	}

	tool2, err := t.CreateTool[int](divisionTool, "Prompt: Divides the first integer by the second.", "Input: (int, int)", "Output: int")
	if err != nil {
		log.Fatalf("Failed to create division tool: %v", err)
	}

	llmWithTools, err := o.CreateConversationNode(
		"AgentWithTools",
		openai.ChatModelGPT4_1Nano,
		client,
		llmFn,
		a.WithUser("ggraph"),
		a.WithTools(tool1, tool2),
	)
	if err != nil {
		log.Fatalf("Failed to create agent with tools node: %v", err)
	}

	startEdge := b.CreateStartEdge(llmWithTools)
	endEdge := b.CreateEndEdge(llmWithTools)

	stateMonitorCh := make(chan g.StateMonitorEntry[a.Conversation], 10)
	graph, err := b.CreateRuntime(startEdge, stateMonitorCh)
	if err != nil {
		log.Fatalf("Runtime creation failed: %v", err)
	}
	graph.AddEdge(endEdge)
	defer graph.Shutdown()

	err = graph.Validate()
	if err != nil {
		log.Fatalf("Graph validation failed: %v", err)
	}

	graph.Invoke(a.CreateConversation(
		a.CreateMessage(a.User, "Can you please add 4 and 5 and divide the result by 2?"),
	))

	for {
		entry := <-stateMonitorCh
		fmt.Printf("[%s - Running: %t], Graph state message: %v, Error: %v\n", entry.Node, entry.Running, entry.NewState, entry.Error)
		if !entry.Running {
			break
		}
	}
}
