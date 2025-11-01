package openai

import (
	"github.com/openai/openai-go/v3"

	a "github.com/morphy76/ggraph/pkg/agent"
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
	params := openai.CompletionNewParams{
		Prompt:           openai.CompletionNewParamsPromptUnion{OfString: openai.String(opts.Prompt)},
		Model:            openai.CompletionNewParamsModel(opts.Model),
		BestOf:           openai.Int(opts.BestOf),
		FrequencyPenalty: openai.Float(opts.FrequencyPenalty),
		Logprobs:         openai.Int(opts.Logprobs),
		MaxTokens:        openai.Int(opts.MaxTokens),
		N:                openai.Int(opts.N),
		PresencePenalty:  openai.Float(opts.PresencePenalty),
		Seed:             openai.Int(opts.Seed),
		Temperature:      openai.Float(opts.Temperature),
		TopP:             openai.Float(opts.TopP),
		User:             openai.String(opts.User),
	}

	return params
}

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
		}
		messages[i] = union
	}
	rv := openai.ChatCompletionNewParams{
		Model:               openai.ChatModel(modelOptions.Model),
		Messages:            messages,
		FrequencyPenalty:    openai.Float(modelOptions.FrequencyPenalty),
		MaxTokens:           openai.Int(modelOptions.MaxTokens),
		N:                   openai.Int(modelOptions.N),
		PresencePenalty:     openai.Float(modelOptions.PresencePenalty),
		Temperature:         openai.Float(modelOptions.Temperature),
		TopP:                openai.Float(modelOptions.TopP),
		User:                openai.String(modelOptions.User),
		Logprobs:            openai.Bool(modelOptions.Logprobs > 0),
		MaxCompletionTokens: openai.Int(modelOptions.MaxCompletionTokens),
		Seed:                openai.Int(modelOptions.Seed),
	}
	return rv
}
