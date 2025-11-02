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
	// Tool represents a tool message.
	Tool
)

// Message represents a single message in a chat conversation.
type Message struct {
	Ts      time.Time
	Role    MessageRole
	Content string
}

// Conversation represents a chat-based language model for an agent.
type Conversation struct {
	Messages []Message
}
