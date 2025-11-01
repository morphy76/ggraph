package agent

// Completion represents a completion response from a language model.
type Completion struct {
	// Text is the generated text from the language model.
	Text string
}

// CreateCompletion is a helper function to create a Completion instance.
//
// Parameters:
//   - text: The text of the completion.
//
// Returns:
//   - An instance of Completion containing the provided text.
//
// Example usage:
//
//	comp := CreateCompletion("This is the generated response.")
func CreateCompletion(text string) Completion {
	return Completion{Text: text}
}
