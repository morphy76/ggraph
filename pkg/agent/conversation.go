package agent

import (
	"time"

	t "github.com/morphy76/ggraph/pkg/agent/tool"
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
	// Timestamp of the message.
	Ts time.Time
	// Role of the message (System, User, Assistant, Tool).
	Role MessageRole
	// Content of the message.
	Content string
	// Tool calls made in the message.
	ToolCalls []t.FnCall
}

// Conversation represents a chat-based language model for an agent.
type Conversation struct {
	// Messages holds the sequence of messages in the conversation.
	Messages []Message
	// CurrentToolCalls holds the current tool calls to be executed.
	CurrentToolCalls []t.FnCall
}
