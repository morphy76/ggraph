package ollama

import (
	"github.com/ollama/ollama/api"
)

// ChatModel represents a chat-based language model using the Ollama API.
type ChatModel struct {
	messages []api.Message
}

// NewChatModel creates a new ChatModel with the given initial messages.
func NewChatModel(messages ...api.Message) ChatModel {
	return ChatModel{messages: messages}
}

// AddMessage appends a message to the chat model.
func (c *ChatModel) AddMessage(msg api.Message) {
	c.messages = append(c.messages, msg)
}

// Messages returns a copy of all messages in the chat model.
func (c ChatModel) Messages() []api.Message {
	return append([]api.Message{}, c.messages...)
}

// MergeChatModels merges two ChatModel instances by appending messages from the new model to the original.
func MergeChatModels(original, new ChatModel) ChatModel {
	if len(new.messages) > 0 {
		original.messages = new.messages
	}
	return original
}
