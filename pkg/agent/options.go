package agent

// ModelOptions defines the parameters for generating text completions.
type ModelOptions struct {
	// The model to use for generating completions.
	Model string
	// The prompt to generate completions for.
	Prompt string
	// The messages that make up the conversation history.
	Messages []Message
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
	// The maximum number of [tokens](/tokenizer) that can be used for the
	// completion portion of the conversation.
	MaxCompletionTokens int64
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

// ModelOption defines an interface for applying options to completion requests.
type ModelOption interface {
	// ApplyToCompletion applies the option to the given CompletionNewParams.
	//
	// Parameters:
	//   - r: A pointer to CompletionNewParams to modify.
	//
	// Returns:
	//   - An error if the application of the option fails, otherwise nil.
	ApplyToCompletion(r *ModelOptions) error
	// ApplyToConversation applies the option to the given ConversationParams.
	//
	// Parameters:
	//   - r: A pointer to ConversationParams to modify.
	//
	// Returns:
	//   - An error if the application of the option fails, otherwise nil.
	ApplyToConversation(r *ModelOptions) error
}

// ModelOptionFunc is a function type that implements the CompletionOption interface.
type ModelOptionFunc func(*ModelOptions) error

// ApplyToCompletion applies the ModelOptionFunc to the given CompletionNewParams.
//
// Parameters:
//   - r: A pointer to CompletionNewParams to modify.
//
// Returns:
//   - An error if the application of the option fails, otherwise nil.
func (s ModelOptionFunc) ApplyToCompletion(r *ModelOptions) error { return s(r) }

// ApplyToConversation applies the ModelOptionFunc to the given ConversationParams.
//
// Parameters:
//   - r: A pointer to ConversationParams to modify.
//
// Returns:
//   - An error if the application of the option fails, otherwise nil.
func (s ModelOptionFunc) ApplyToConversation(r *ModelOptions) error { return s(r) }

// WithBestOf sets the BestOf option for completion requests.
//
// Parameters:
//   - bestOf: The number of best completions to generate.
//
// Returns:
//   - A CompletionOption that sets the BestOf parameter.
//
// Example usage:
//
//	option := WithBestOf(3)
func WithBestOf(bestOf int64) ModelOption {
	return ModelOptionFunc(func(r *ModelOptions) error {
		if bestOf < 1 {
			return ErrorInvalidBestOf
		}
		r.BestOf = bestOf
		return nil
	})
}

// WithEcho sets the Echo option for completion requests.
//
// Parameters:
//   - echo: A boolean indicating whether to echo the prompt.
//
// Returns:
//   - A CompletionOption that sets the Echo parameter.
//
// Example usage:
//
//	option := WithEcho(true)
func WithEcho(echo bool) ModelOption {
	return ModelOptionFunc(func(r *ModelOptions) error {
		r.Echo = echo
		return nil
	})
}

// WithFrequencyPenalty sets the FrequencyPenalty option for completion requests.
//
// Parameters:
//   - frequencyPenalty: The frequency penalty value.
//
// Returns:
//   - A CompletionOption that sets the FrequencyPenalty parameter.
//
// Example usage:
//
//	option := WithFrequencyPenalty(0.5)
func WithFrequencyPenalty(frequencyPenalty float64) ModelOption {
	return ModelOptionFunc(func(r *ModelOptions) error {
		if frequencyPenalty < -2.0 || frequencyPenalty > 2.0 {
			return ErrorInvalidFrequencyPenalty
		}
		r.FrequencyPenalty = frequencyPenalty
		return nil
	})
}

// WithLogprobs sets the Logprobs option for completion requests.
//
// Parameters:
//   - logprobs: The number of log probabilities to include.
//
// Returns:
//   - A CompletionOption that sets the Logprobs parameter.
//
// Example usage:
//
//	option := WithLogprobs(5)
func WithLogprobs(logprobs int64) ModelOption {
	return ModelOptionFunc(func(r *ModelOptions) error {
		if logprobs < 0 || logprobs > 5 {
			return ErrorInvalidLogprobs
		}
		r.Logprobs = logprobs
		return nil
	})
}

// WithMaxTokens sets the MaxTokens option for completion requests.
//
// Parameters:
//   - maxTokens: The maximum number of tokens to generate.
//
// Returns:
//   - A CompletionOption that sets the MaxTokens parameter.
//
// Example usage:
//
//	option := WithMaxTokens(150)
func WithMaxTokens(maxTokens int64) ModelOption {
	return ModelOptionFunc(func(r *ModelOptions) error {
		if maxTokens < 1 {
			return ErrorInvalidMaxTokens
		}
		r.MaxTokens = maxTokens
		return nil
	})
}

// WithN sets the N option for completion requests.
//
// Parameters:
//   - n: The number of completions to generate.
//
// Returns:
//   - A CompletionOption that sets the N parameter.
//
// Example usage:
//
//	option := WithN(2)
func WithN(n int64) ModelOption {
	return ModelOptionFunc(func(r *ModelOptions) error {
		if n < 1 {
			return ErrorInvalidN
		}
		r.N = n
		return nil
	})
}

// WithPresencePenalty sets the PresencePenalty option for completion requests.
//
// Parameters:
//   - presencePenalty: The presence penalty value.
//
// Returns:
//   - A CompletionOption that sets the PresencePenalty parameter.
//
// Example usage:
//
//	option := WithPresencePenalty(0.3)
func WithPresencePenalty(presencePenalty float64) ModelOption {
	return ModelOptionFunc(func(r *ModelOptions) error {
		if presencePenalty < -2.0 || presencePenalty > 2.0 {
			return ErrorInvalidPresencePenalty
		}
		r.PresencePenalty = presencePenalty
		return nil
	})
}

// WithSeed sets the Seed option for completion requests.
//
// Parameters:
//   - seed: The seed value for deterministic sampling.
//
// Returns:
//   - A CompletionOption that sets the Seed parameter.
//
// Example usage:
//
//	option := WithSeed(42)
func WithSeed(seed int64) ModelOption {
	return ModelOptionFunc(func(r *ModelOptions) error {
		r.Seed = seed
		return nil
	})
}

// WithTemperature sets the Temperature option for completion requests.
//
// Parameters:
//   - temperature: The temperature value.
//
// Returns:
//   - A CompletionOption that sets the Temperature parameter.
//
// Example usage:
//
//	option := WithTemperature(0.7)
func WithTemperature(temperature float64) ModelOption {
	return ModelOptionFunc(func(r *ModelOptions) error {
		if temperature < 0.0 || temperature > 2.0 {
			return ErrorInvalidTemperature
		}
		r.Temperature = temperature
		return nil
	})
}

// WithTopP sets the TopP option for completion requests.
//
// Parameters:
//   - topP: The top_p value.
//
// Returns:
//   - A CompletionOption that sets the TopP parameter.
//
// Example usage:
//
//	option := WithTopP(0.9)
func WithTopP(topP float64) ModelOption {
	return ModelOptionFunc(func(r *ModelOptions) error {
		if topP < 0.0 || topP > 1.0 {
			return ErrorInvalidTopP
		}
		r.TopP = topP
		return nil
	})
}

// WithUser sets the User option for completion requests.
//
// Parameters:
//   - user: A unique identifier representing the end-user.
//
// Returns:
//   - A CompletionOption that sets the User parameter.
//
// Example usage:
//
//	option := WithUser("user-1234")
func WithUser(user string) ModelOption {
	return ModelOptionFunc(func(r *ModelOptions) error {
		r.User = user
		return nil
	})
}

// WithMaxCompletionTokens sets the MaxCompletionTokens option for conversation requests.
//
// Parameters:
//   - maxCompletionTokens: The maximum number of tokens for the completion part of the conversation.
//
// Returns:
//   - A ConversationOption that sets the MaxCompletionTokens parameter.
//
// Example usage:
//
//	option := WithMaxCompletionTokens(200)
func WithMaxCompletionTokens(maxCompletionTokens int64) ModelOption {
	return ModelOptionFunc(func(r *ModelOptions) error {
		if maxCompletionTokens < 1 {
			return ErrorInvalidMaxTokens
		}
		r.MaxCompletionTokens = maxCompletionTokens
		return nil
	})
}
