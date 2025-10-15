package openai

import (
	"time"

	"github.com/morphy76/ggraph/pkg/llm"
	"github.com/openai/openai-go/v3"
)

func ToOpenAIModel(models ...llm.AgentModel) []openai.ChatCompletionMessageParamUnion {
	var messages []openai.ChatCompletionMessageParamUnion
	for _, m := range models {
		for _, msg := range m.Messages {
			messages = append(messages, ToOpenAIMessage(msg))
		}
	}
	return messages
}

func ToOpenAIMessage(msg llm.Message) openai.ChatCompletionMessageParamUnion {
	switch msg.Role {
	case llm.System:
		return openai.SystemMessage(msg.Content)
	case llm.User:
		return openai.UserMessage(msg.Content)
	case llm.Assistant:
		return openai.AssistantMessage(msg.Content)
	default:
		return openai.UserMessage(msg.Content)
	}
}

func FromOpenAIMessage(msg openai.ChatCompletionMessage) llm.Message {
	return llm.Message{
		Ts:      time.Now(),
		Role:    llm.Assistant,
		Content: msg.Content,
	}
}
