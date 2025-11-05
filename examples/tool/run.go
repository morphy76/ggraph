package main

import (
	"context"
	"encoding/json"
	"fmt"
	"log"

	"github.com/openai/openai-go/v3"

	a "github.com/morphy76/ggraph/pkg/agent"
	o "github.com/morphy76/ggraph/pkg/agent/openai"
	t "github.com/morphy76/ggraph/pkg/agent/tool"
	b "github.com/morphy76/ggraph/pkg/builders"
	g "github.com/morphy76/ggraph/pkg/graph"
)

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

			systemMex := `
			Use all tools: feel free to use the available tools to answer the user's question through multiple tool calls.
			You must never perform arithmetic or reasoning operations yourself.
			You must always use the provided tools for every operation.
			If the output of one tool should be used by another tool:
			- reference it as a JSON object of the form {"from_call": "<previous_call_id>", "field": "result"};
			- Never substitute numeric results directly; always reference previous tool outputs using the {"from_call": "<previous_call_id>", "field": "result"} object format.
			Example of dependant tool calls where a divisionTool uses the output of an additionTool:
			[
				{
					"function": { "name": "additionTool", "arguments": "{"addend1": 4, "addend2": 5}" },
					"id": "call_1"
				},
				{
					"function": { "name": "divisionTool", "arguments": "{"dividend": {"from_call": "call_1", "field": "result"}, "divisor": 2}" }
				}
			]`

			useMessages := append(
				[]a.Message{
					a.CreateMessage(a.System, systemMex),
				},
				userInput.Messages...,
			)

			useOpts, err := a.CreateConversationOptions(
				model,
				useMessages,
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

			dump, _ := json.Marshal(resp.Choices[0].Message.ToolCalls)
			fmt.Printf("Rest:\n\n%s\n\n", dump)

			answer := resp.Choices[0].Message.Content
			currentState.Messages = append(userInput.Messages,
				a.CreateMessage(a.Assistant, answer))

			return currentState, nil
		}
	}

	tool1, err := t.CreateTool[int](additionTool, "Prompt: this tool is used to sum two integers.", "Input: addend1, addend2", "Required: addend1, addend2")
	if err != nil {
		log.Fatalf("Failed to create addition tool: %v", err)
	}

	tool2, err := t.CreateTool[int](divisionTool, "Prompt: this tool is used to divide a dividend by a divisor.", "Input: dividend, divisor", "Required: dividend, divisor")
	if err != nil {
		log.Fatalf("Failed to create division tool: %v", err)
	}

	llmWithTools, err := o.CreateConversationNode(
		"AgentWithTools",
		openai.ChatModelGPT5Nano,
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
		a.CreateMessage(a.User, "Can you please add 4 and 5 and then divide the result by 2?"),
	))

	for {
		entry := <-stateMonitorCh
		fmt.Printf("[%s - Running: %t], Graph state message: %v, Error: %v\n", entry.Node, entry.Running, entry.NewState, entry.Error)
		if !entry.Running {
			break
		}
	}
}
