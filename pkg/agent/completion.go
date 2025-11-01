package agent

import "fmt"

var (
	// ErrorInvalidBestOf is returned when the BestOf parameter is less than 1.
	ErrorInvalidBestOf = fmt.Errorf("bestOf must be at least 1")
	// ErrorInvalidFrequencyPenalty is returned when the FrequencyPenalty parameter is out of range.
	ErrorInvalidFrequencyPenalty = fmt.Errorf("frequencyPenalty must be between -2.0 and 2.0")
	// ErrorInvalidLogprobs is returned when the Logprobs parameter is out of range.
	ErrorInvalidLogprobs = fmt.Errorf("logprobs must be between 0 and 5")
	// ErrorInvalidMaxTokens is returned when the MaxTokens parameter is less than 1.
	ErrorInvalidMaxTokens = fmt.Errorf("maxTokens must be at least 1")
	// ErrorInvalidN is returned when the N parameter is less than 1.
	ErrorInvalidN = fmt.Errorf("n must be at least 1")
	// ErrorInvalidPresencePenalty is returned when the PresencePenalty parameter is out of range.
	ErrorInvalidPresencePenalty = fmt.Errorf("presencePenalty must be between -2.0 and 2.0")
	// ErrorInvalidTemperature is returned when the Temperature parameter is out of range.
	ErrorInvalidTemperature = fmt.Errorf("temperature must be between 0.0 and 2.0")
	// ErrorInvalidTopP is returned when the TopP parameter is out of range.
	ErrorInvalidTopP = fmt.Errorf("topP must be between 0.0 and 1.0")
)

// Completion represents a completion response from a language model.
type Completion struct {
	// Text is the generated text from the language model.
	Text string
}

// CompletionOptions defines the parameters for generating text completions.
type CompletionOptions struct {
	// The prompt to generate completions for.
	Prompt string
	// The model to use for generating completions.
	Model string
	// Generates `best_of` completions server-side and returns the "best" (the one with
	// the highest log probability per token). Results cannot be streamed.
	//
	// When used with `n`, `best_of` controls the number of candidate completions and
	// `n` specifies how many to return – `best_of` must be greater than `n`.
	BestOf int64
	// Echo back the prompt in addition to the completion
	Echo bool
	// Number between -2.0 and 2.0. Positive values penalize new tokens based on their
	// existing frequency in the text so far, decreasing the model's likelihood to
	// repeat the same line verbatim.
	FrequencyPenalty float64
	// Include the log probabilities on the `logprobs` most likely output tokens, as
	// well the chosen tokens. For example, if `logprobs` is 5, the API will return a
	// list of the 5 most likely tokens. The API will always return the `logprob` of
	// the sampled token, so there may be up to `logprobs+1` elements in the response.
	//
	// The maximum value for `logprobs` is 5.
	Logprobs int64
	// The maximum number of [tokens](/tokenizer) that can be generated in the
	// completion.
	MaxTokens int64
	// How many completions to generate for each prompt.
	N int64
	// Number between -2.0 and 2.0. Positive values penalize new tokens based on
	// whether they appear in the text so far, increasing the model's likelihood to
	// talk about new topics.
	PresencePenalty float64
	// If specified, our system will make a best effort to sample deterministically,
	// such that repeated requests with the same `seed` and parameters should return
	// the same result.
	//
	// Determinism is not guaranteed, and you should refer to the `system_fingerprint`
	// response parameter to monitor changes in the backend.
	Seed int64
	// What sampling temperature to use, between 0 and 2. Higher values like 0.8 will
	// make the output more random, while lower values like 0.2 will make it more
	// focused and deterministic.
	//
	// We generally recommend altering this or `top_p` but not both.
	Temperature float64
	// An alternative to sampling with temperature, called nucleus sampling, where the
	// model considers the results of the tokens with top_p probability mass. So 0.1
	// means only the tokens comprising the top 10% probability mass are considered.
	//
	// We generally recommend altering this or `temperature` but not both.
	TopP float64
	// A unique identifier representing your end-user.
	User string
}

// CompletionOption defines an interface for applying options to completion requests.
type CompletionOption interface {
	// Apply applies the option to the given CompletionNewParams.
	//
	// Parameters:
	//   - r: A pointer to CompletionNewParams to modify.
	//
	// Returns:
	//   - An error if the application of the option fails, otherwise nil.
	Apply(*CompletionOptions) error
}

// CompletionOptionFunc is a function type that implements the CompletionOption interface.
type CompletionOptionFunc func(*CompletionOptions) error

// Apply applies the CompletionOptionFunc to the given CompletionParams.
//
// Parameters:
//   - r: A pointer to CompletionParams to modify.
//
// Returns:
//   - An error if the application of the option fails, otherwise nil.
func (s CompletionOptionFunc) Apply(r *CompletionOptions) error { return s(r) }
