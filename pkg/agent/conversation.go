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
	Ts      time.Time
	Role    MessageRole
	Content string
}

// Conversation represents a chat-based language model for an agent.
type Conversation struct {
	// Messages holds the sequence of messages in the conversation.
	Messages []Message
	// ToolCalls holds the sequence of tool calls made during the conversation.
	ToolCalls []t.ToolCall
}
