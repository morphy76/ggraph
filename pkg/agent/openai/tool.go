package openai

import (
	"encoding/json"
	"fmt"

	tool "github.com/morphy76/ggraph/pkg/agent/tool"
	"github.com/openai/openai-go/v3"
)

// ConvertToolCall converts an OpenAI ChatCompletionMessageToolCallUnion to our internal ToolCall structure.
//
// Parameters:
//   - openAIToolCall: The OpenAI tool call to convert
//   - tools: A map of tool names to Tool objects for lookup
//
// Returns:
//   - A ToolCall structure with the converted data
//   - An error if conversion fails
func ConvertToolCall(openAIToolCall openai.ChatCompletionMessageToolCallUnion, tools map[string]*tool.Tool) (*tool.ToolCall, error) {
	// Extract the function tool call (the most common type)
	functionToolCall := openAIToolCall.AsFunction()

	// Get the tool from the registry
	toolName := functionToolCall.Function.Name
	usingTool, exists := tools[toolName]
	if !exists {
		return nil, fmt.Errorf("tool '%s' not found in tool registry", toolName)
	}

	// Parse the arguments JSON string
	var arguments map[string]any
	if functionToolCall.Function.Arguments != "" {
		err := json.Unmarshal([]byte(functionToolCall.Function.Arguments), &arguments)
		if err != nil {
			return nil, fmt.Errorf("failed to parse tool arguments: %w", err)
		}
	} else {
		arguments = make(map[string]any)
	}

	// Create and return the ToolCall
	return &tool.ToolCall{
		Id:        openAIToolCall.ID,
		UsingTool: *usingTool,
		Arguments: arguments,
	}, nil
}
