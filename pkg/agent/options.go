package agent

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
func WithBestOf(bestOf int64) CompletionOption {
	return CompletionOptionFunc(func(r *CompletionOptions) error {
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
func WithEcho(echo bool) CompletionOption {
	return CompletionOptionFunc(func(r *CompletionOptions) error {
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
func WithFrequencyPenalty(frequencyPenalty float64) CompletionOption {
	return CompletionOptionFunc(func(r *CompletionOptions) error {
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
func WithLogprobs(logprobs int64) CompletionOption {
	return CompletionOptionFunc(func(r *CompletionOptions) error {
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
func WithMaxTokens(maxTokens int64) CompletionOption {
	return CompletionOptionFunc(func(r *CompletionOptions) error {
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
func WithN(n int64) CompletionOption {
	return CompletionOptionFunc(func(r *CompletionOptions) error {
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
func WithPresencePenalty(presencePenalty float64) CompletionOption {
	return CompletionOptionFunc(func(r *CompletionOptions) error {
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
func WithSeed(seed int64) CompletionOption {
	return CompletionOptionFunc(func(r *CompletionOptions) error {
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
func WithTemperature(temperature float64) CompletionOption {
	return CompletionOptionFunc(func(r *CompletionOptions) error {
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
func WithTopP(topP float64) CompletionOption {
	return CompletionOptionFunc(func(r *CompletionOptions) error {
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
func WithUser(user string) CompletionOption {
	return CompletionOptionFunc(func(r *CompletionOptions) error {
		r.User = user
		return nil
	})
}
