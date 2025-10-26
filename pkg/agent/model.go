package agent

import (
	"time"
)

// MessageRole defines the role of a message in a chat conversation.
type MessageRole int

const (
	// System represents a system message.
	System MessageRole = iota
	// User represents a user message.
	User
	// Assistant represents an assistant message.
	Assistant
)

// Message represents a single message in a chat conversation.
type Message struct {
	Ts      time.Time
	Role    MessageRole
	Content string
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

// Completion represents a completion response from a language model.
type Completion struct {
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

// Conversation represents a chat-based language model for an agent.
type Conversation struct {
	Messages []Message
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
