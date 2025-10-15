package openai

import (
	"github.com/openai/openai-go/v3"
)

// ChatModel represents a chat-based language model using the OpenAI API.
type ChatModel struct {
	messages []openai.Message
}

// NewChatModel creates a new ChatModel with the given initial messages.
func NewChatModel(messages ...openai.Message) ChatModel {
	return ChatModel{messages: messages}
}

// AddMessage appends a message to the chat model.
func (c *ChatModel) AddMessage(msg openai.Message) {
	c.messages = append(c.messages, msg)
}

// Messages returns a copy of all messages in the chat model.
func (c ChatModel) Messages() []openai.Message {
	return append([]openai.Message{}, c.messages...)
}

// MergeChatModels merges two ChatModel instances by appending messages from the new model to the original.
func MergeChatModels(original, new ChatModel) ChatModel {
	if len(new.messages) > 0 {
		original.messages = new.messages
	}
	return original
}
