package openai

import (
	"encoding/json"
	"strings"

	"github.com/openai/openai-go/v3"

	a "github.com/morphy76/ggraph/pkg/agent"
	t "github.com/morphy76/ggraph/pkg/agent/tool"
)

// ConvertCompletionOptions converts internal CompletionOptions to OpenAI CompletionNewParams.
//
// Parameters:
//   - opts: The internal CompletionOptions to be converted.
//
// Returns:
//   - An openai.CompletionNewParams struct populated with the values from the internal CompletionOptions.
//
// Example usage:
//
//	internalOpts := a.CreateCompletionOptions(...)
//
//	openAIParams := ConvertCompletionOptions(internalOpts)
func ConvertCompletionOptions(opts *a.ModelOptions) openai.CompletionNewParams {
	rv := openai.CompletionNewParams{
		Prompt: openai.CompletionNewParamsPromptUnion{OfString: openai.String(opts.Prompt)},
		Model:  openai.CompletionNewParamsModel(opts.Model),
	}

	if opts.BestOf != nil {
		rv.BestOf = openai.Int(*opts.BestOf)
	}
	if opts.FrequencyPenalty != nil {
		rv.FrequencyPenalty = openai.Float(*opts.FrequencyPenalty)
	}
	if opts.Logprobs != nil {
		rv.Logprobs = openai.Int(*opts.Logprobs)
	}
	if opts.MaxTokens != nil {
		rv.MaxTokens = openai.Int(*opts.MaxTokens)
	}
	if opts.N != nil {
		rv.N = openai.Int(*opts.N)
	}
	if opts.PresencePenalty != nil {
		rv.PresencePenalty = openai.Float(*opts.PresencePenalty)
	}
	if opts.Temperature != nil {
		rv.Temperature = openai.Float(*opts.Temperature)
	}
	if opts.TopP != nil {
		rv.TopP = openai.Float(*opts.TopP)
	}
	if opts.User != nil {
		rv.User = openai.String(*opts.User)
	}
	if opts.Seed != nil {
		rv.Seed = openai.Int(*opts.Seed)
	}

	return rv
}

// ConvertConversationOptions converts internal ModelOptions for conversations to OpenAI ChatCompletionNewParams.
//
// Parameters:
//   - modelOptions: The internal ModelOptions to be converted.
//
// Returns:
//   - An openai.ChatCompletionNewParams struct populated with the values from the internal ModelOptions.
//
// Example usage:
//
//	internalOpts := a.CreateCompletionOptions(...)
//
//	openAIChatParams := ConvertConversationOptions(internalOpts)
func ConvertConversationOptions(modelOptions *a.ModelOptions) openai.ChatCompletionNewParams {
	messages := make([]openai.ChatCompletionMessageParamUnion, len(modelOptions.Messages))
	for i, msg := range modelOptions.Messages {
		var union openai.ChatCompletionMessageParamUnion
		switch msg.Role {
		case a.System:
			union = openai.SystemMessage(msg.Content)
		case a.User:
			union = openai.UserMessage(msg.Content)
		case a.Assistant:
			union = openai.AssistantMessage(msg.Content)
			useUnion := union.OfAssistant
			if len(msg.ToolCalls) > 0 {
				useUnion.ToolCalls = make([]openai.ChatCompletionMessageToolCallUnionParam, len(msg.ToolCalls))
				for i, tc := range msg.ToolCalls {
					argsAsString, _ := json.Marshal(tc.Arguments)
					useUnion.ToolCalls[i] = openai.ChatCompletionMessageToolCallUnionParam{
						OfFunction: &openai.ChatCompletionMessageFunctionToolCallParam{
							ID: tc.ID,
							Function: openai.ChatCompletionMessageFunctionToolCallFunctionParam{
								Name:      tc.ToolName,
								Arguments: string(argsAsString),
							},
						},
					}
				}

				union.OfAssistant = useUnion
			}
		case a.Tool:
			toolAnswer := strings.SplitN(msg.Content, ":", 2)
			union = openai.ToolMessage(toolAnswer[1], toolAnswer[0])
		}

		messages[i] = union
	}

	tools := make([]openai.ChatCompletionToolUnionParam, len(modelOptions.Tools))
	for i, tool := range modelOptions.Tools {
		fn := tool2Fn(tool)
		tools[i] = openai.ChatCompletionToolUnionParam{
			OfFunction: fn,
		}
	}

	rv := openai.ChatCompletionNewParams{
		Model:    openai.ChatModel(modelOptions.Model),
		Messages: messages,
	}

	if len(tools) > 0 {
		rv.ParallelToolCalls = openai.Bool(true)
		rv.Tools = tools
	}

	if modelOptions.FrequencyPenalty != nil {
		rv.FrequencyPenalty = openai.Float(*modelOptions.FrequencyPenalty)
	}
	if modelOptions.MaxTokens != nil {
		rv.MaxTokens = openai.Int(*modelOptions.MaxTokens)
	}
	if modelOptions.N != nil {
		rv.N = openai.Int(*modelOptions.N)
	}
	if modelOptions.PresencePenalty != nil {
		rv.PresencePenalty = openai.Float(*modelOptions.PresencePenalty)
	}
	if modelOptions.Temperature != nil {
		rv.Temperature = openai.Float(*modelOptions.Temperature)
	}
	if modelOptions.TopP != nil {
		rv.TopP = openai.Float(*modelOptions.TopP)
	}
	if modelOptions.User != nil {
		rv.User = openai.String(*modelOptions.User)
	}
	if modelOptions.Logprobs != nil {
		rv.Logprobs = openai.Bool(*modelOptions.Logprobs > 0)
	}
	if modelOptions.MaxCompletionTokens != nil {
		rv.MaxCompletionTokens = openai.Int(*modelOptions.MaxCompletionTokens)
	}
	if modelOptions.Seed != nil {
		rv.Seed = openai.Int(*modelOptions.Seed)
	}

	return rv
}

func tool2Fn(tool *t.Tool) *openai.ChatCompletionFunctionToolParam {
	toolProps := make(map[string]interface{})
	for _, arg := range tool.Args {
		useType := convertToSupportedJSONType(arg.Type)
		toolProps[arg.Name] = map[string]interface{}{
			"type": useType,
		}
	}

	return &openai.ChatCompletionFunctionToolParam{
		Function: openai.FunctionDefinitionParam{
			Name:        tool.Name,
			Description: openai.String(tool.BuildToolPrompt()),
			Parameters: openai.FunctionParameters{
				"type":       "object",
				"properties": toolProps,
				"required":   tool.RequiredArgs(),
			},
		},
	}
}

func convertToSupportedJSONType(argType string) string {
	switch argType {
	case "string":
		return "string"
	case "int", "int8", "int16", "int32", "int64",
		"uint", "uint8", "uint16", "uint32", "uint64":
		return "integer"
	case "float32", "float64":
		return "number"
	case "bool":
		return "boolean"
	case "[]string", "[]int", "[]float64", "[]bool":
		return "array"
	case "map[string]interface{}":
		return "object"
	default:
		return "string"
	}
}
