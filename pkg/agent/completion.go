package agent

import (
	"errors"
)

var (
	// ErrorInvalidBestOf is returned when the BestOf parameter is less than 1.
	ErrorInvalidBestOf = errors.New("bestOf must be at least 1")
	// ErrorInvalidFrequencyPenalty is returned when the FrequencyPenalty parameter is out of range.
	ErrorInvalidFrequencyPenalty = errors.New("frequencyPenalty must be between -2.0 and 2.0")
	// ErrorInvalidLogprobs is returned when the Logprobs parameter is out of range.
	ErrorInvalidLogprobs = errors.New("logprobs must be between 0 and 5")
	// ErrorInvalidMaxTokens is returned when the MaxTokens parameter is less than 1.
	ErrorInvalidMaxTokens = errors.New("maxTokens must be at least 1")
	// ErrorInvalidN is returned when the N parameter is less than 1.
	ErrorInvalidN = errors.New("n must be at least 1")
	// ErrorInvalidPresencePenalty is returned when the PresencePenalty parameter is out of range.
	ErrorInvalidPresencePenalty = errors.New("presencePenalty must be between -2.0 and 2.0")
	// ErrorInvalidTemperature is returned when the Temperature parameter is out of range.
	ErrorInvalidTemperature = errors.New("temperature must be between 0.0 and 2.0")
	// ErrorInvalidTopP is returned when the TopP parameter is out of range.
	ErrorInvalidTopP = errors.New("topP must be between 0.0 and 1.0")
)

// Completion represents a completion response from a language model.
type Completion struct {
	// Text is the generated text from the language model.
	Text string
}
