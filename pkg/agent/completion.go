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
