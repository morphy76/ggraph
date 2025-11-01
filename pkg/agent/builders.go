package agent

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

// CreateCompletionOptions creates a CompletionOptions instance by applying the provided options.
//
// Parameters:
//   - prompt: The prompt for the completion.
//   - model: The model to use for the completion.
//   - opts: A variadic list of CompletionOption to apply.
//
// Returns:
//   - An instance of CompletionOptions with the applied options.
//   - An error if any option application fails.
//
// Example usage:
//
//	options, err := CreateCompletionOptions("Hello, world!", "gpt-4", WithMaxTokens(100))
func CreateCompletionOptions(
	prompt, model string,
	opts ...CompletionOption,
) (*CompletionOptions, error) {
	useOptions := CompletionOptions{
		Prompt: prompt,
		Model:  model,
	}
	for _, opt := range opts {
		if err := opt.Apply(&useOptions); err != nil {
			return &CompletionOptions{}, err
		}
	}
	return &useOptions, nil
}
