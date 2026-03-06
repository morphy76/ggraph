package main

import (
	"fmt"
	"log"

	"github.com/openai/openai-go/v3"

	a "github.com/morphy76/ggraph/pkg/agent"
	ag "github.com/morphy76/ggraph/pkg/agent/graph"
	o "github.com/morphy76/ggraph/pkg/agent/openai"
	t "github.com/morphy76/ggraph/pkg/agent/tool"
	b "github.com/morphy76/ggraph/pkg/builders"
	g "github.com/morphy76/ggraph/pkg/graph"
)

func additionTool(a int, b int) (int, error) {
	fmt.Printf("=======>>> additionTool with args %d, %d\n", a, b) // DEBUG
	return a + b, nil
}

func divisionTool(a int, b int) (int, error) {
	fmt.Printf("=======>>> divisionTool with args %d, %d\n", a, b) // DEBUG
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

	tool1, err := t.CreateTool[int](additionTool, "Prompt: this tool is used to sum two integers.", "Input: addend1, addend2", "Required: addend1, addend2")
	if err != nil {
		log.Fatalf("Failed to create addition tool: %v", err)
	}

	tool2, err := t.CreateTool[int](divisionTool, "Prompt: this tool is used to divide a dividend by a divisor.", "Input: dividend, divisor", "Required: dividend, divisor")
	if err != nil {
		log.Fatalf("Failed to create division tool: %v", err)
	}

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

	llmWithTools, err := o.CreateConversationNode(
		"AgentWithTools",
		openai.ChatModelGPT5Nano,
		client,
		a.WithMessages(
			a.CreateMessage(a.System, systemMex),
		),
		a.WithUser("ggraph"),
		a.WithTools(tool1, tool2),
	)
	if err != nil {
		log.Fatalf("Failed to create agent with tools node: %v", err)
	}

	toolProcessor, err := ag.CreateToolNode("ToolProcessor", tool1, tool2)
	if err != nil {
		log.Fatalf("Failed to create tool processor node: %v", err)
	}

	startEdge := b.CreateStartEdge(llmWithTools)
	toolRequestEdge := b.CreateEdge(llmWithTools, toolProcessor, map[string]string{a.RouteTagToolKey: a.RouteTagToolRequest})
	toolResponseEdge := b.CreateEdge(toolProcessor, llmWithTools, map[string]string{a.RouteTagToolKey: a.RouteTagToolResponse})
	endEdge := b.CreateEndEdge(llmWithTools)

	stateMonitorCh := make(chan g.StateMonitorEntry[a.Conversation], 10)
	graph, err := b.CreateRuntime(startEdge, stateMonitorCh)
	if err != nil {
		log.Fatalf("Runtime creation failed: %v", err)
	}
	graph.AddEdge(toolRequestEdge, toolResponseEdge, endEdge)
	defer graph.Shutdown()

	err = graph.Validate()
	if err != nil {
		log.Fatalf("Graph validation failed: %v", err)
	}

	userInput := "What is the result of adding 15 and 30, and then dividing the sum by 3?"
	graph.Invoke(a.CreateConversation(
		a.CreateMessage(a.User, userInput),
	))
	fmt.Println(userInput)

	for {
		entry := <-stateMonitorCh
		if !entry.Running {
			break
		}
		if entry.Error != nil {
			fmt.Printf("Error during graph execution: %v\n", entry.Error)
			break
		} else {
			if len(entry.NewState.Messages) > 0 {
				lastMessage := entry.NewState.Messages[len(entry.NewState.Messages)-1]
				if lastMessage.Role != a.Tool {
					fmt.Println(lastMessage.Content)
				}
			}
		}
	}
}
