package agent

import "time"

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
//   - completionOptions: A variadic list of ModelOption to apply.
//
// Returns:
//   - An instance of ModelOptions with the applied options.
//   - An error if any option application fails.
//
// Example usage:
//
//	options, err := CreateCompletionOptions("Hello, world!", "gpt-4", WithMaxTokens(100))
func CreateCompletionOptions(
	model, prompt string,
	completionOptions ...ModelOption,
) (*ModelOptions, error) {
	useOptions := ModelOptions{
		Prompt: prompt,
		Model:  model,
	}
	for _, opt := range completionOptions {
		if err := opt.ApplyToCompletion(&useOptions); err != nil {
			return &ModelOptions{}, err
		}
	}
	return &useOptions, nil
}

// CreateMessage is a helper function to create a Message instance.
//
// Parameters:
//   - role: The role of the message (System, User, Assistant).
//   - content: The content of the message.
//
// Returns:
//   - An instance of Message with the current timestamp.
//
// Example usage:
//
//	msg := CreateMessage(User, "Hello, how can I assist you?")
func CreateMessage(role MessageRole, content string) Message {
	return Message{
		Ts:      time.Now(),
		Role:    role,
		Content: content,
	}
}

// CreateConversation is a helper function to create an AgentModel instance.
//
// Parameters:
//   - messages: A variadic list of Message instances to initialize the model.
//
// Returns:
//   - An instance of AgentModel containing the provided messages.
//
// Example usage:
//
//	model := CreateConversation(
//	    CreateMessage(System, "You are a helpful assistant."),
//	    CreateMessage(User, "What's the weather like today?"),
//	)
func CreateConversation(messages ...Message) Conversation {
	return Conversation{Messages: messages}
}

// CreateConversationOptions creates a ConversationOptions instance by applying the provided options.
//
// Parameters:
//   - promptModel: The model to use for the conversation.
//   - messages: A slice of Message instances that make up the conversation history.
//   - modelOptions: A variadic list of ModelOption to apply.
//
// Returns:
//   - An instance of ModelOptions with the applied options.
//   - An error if any option application fails.
//
// Example usage:
//
//	options, err := CreateConversationOptions(
//	    "gpt-4-chat",
//	    []Message{
//	        CreateMessage(System, "You are a helpful assistant."),
//	        CreateMessage(User, "Tell me a joke."),
//	    },
//	    WithTemperature(0.7),
//	)
func CreateConversationOptions(
	promptModel string,
	messages []Message,
	modelOptions ...ModelOption,
) (*ModelOptions, error) {
	useOptions := ModelOptions{
		Model: promptModel,
	}
	for _, opt := range modelOptions {
		if err := opt.ApplyToConversation(&useOptions); err != nil {
			return &ModelOptions{}, err
		}
	}
	return &useOptions, nil
}
