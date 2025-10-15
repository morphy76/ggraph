package llm

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

// AgentModel represents a chat-based language model for an agent.
type AgentModel struct {
	Messages []Message
}

// AddUserMessage adds a new user message to the AgentModel.
func (m *AgentModel) AddUserMessage(content string) {
	m.Messages = append(m.Messages, CreateMessage(User, content))
}

// CreateMessage is a helper function to create a Message instance.
func CreateMessage(role MessageRole, content string) Message {
	return Message{
		Ts:      time.Now(),
		Role:    role,
		Content: content,
	}
}

// CreateModel is a helper function to create an AgentModel instance.
func CreateModel(messages ...Message) AgentModel {
	return AgentModel{Messages: messages}
}

// MergeAgentModel merges two AgentModel instances by concatenating their messages.
func MergeAgentModel(a, b AgentModel) AgentModel {
	return AgentModel{
		Messages: append(a.Messages, b.Messages...),
	}
}
