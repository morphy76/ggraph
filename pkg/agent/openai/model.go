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
func ConvertCompletionOptions(opts *a.CompletionOptions) openai.CompletionNewParams {
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
