package openai

import (
	"encoding/json"
	"errors"
	"fmt"

	"github.com/openai/openai-go/v3"

	t "github.com/morphy76/ggraph/pkg/agent/tool"
)

var (
	// ErrArgumentParse is returned when tool arguments cannot be parsed
	ErrArgumentParse = errors.New("failed to parse tool arguments")
)

// ConvertToolCall converts an OpenAI ChatCompletionMessageToolCallUnion to our internal ToolCall structure.
//
// Parameters:
//   - openAIToolCall: The OpenAI tool call to convert
//
// Returns:
//   - A ToolCall structure with the converted data
//   - An error if conversion fails
func ConvertToolCall(openAIToolCall openai.ChatCompletionMessageToolCallUnion) (*t.FnCall, error) {
	functionToolCall := openAIToolCall.AsFunction()
	toolName := functionToolCall.Function.Name

	var arguments map[string]any
	if functionToolCall.Function.Arguments != "" {
		err := json.Unmarshal([]byte(functionToolCall.Function.Arguments), &arguments)
		if err != nil {
			return nil, fmt.Errorf("failed to parse tool arguments: %w", err)
		}
	} else {
		arguments = make(map[string]any)
	}

	return &t.FnCall{
		ID:        openAIToolCall.ID,
		ToolName:  toolName,
		Arguments: arguments,
	}, nil
}
