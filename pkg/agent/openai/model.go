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
		}
		messages[i] = union
	}

	rv := openai.ChatCompletionNewParams{
		Model:    openai.ChatModel(modelOptions.Model),
		Messages: messages,
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
