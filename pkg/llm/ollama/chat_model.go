package ollama

import (
	"time"

	"github.com/morphy76/ggraph/pkg/llm"
	"github.com/ollama/ollama/api"
)

func ToLLamaModel(model ...llm.AgentModel) []api.Message {
	var messages []api.Message
	for _, m := range model {
		for _, msg := range m.Messages {
			messages = append(messages, ToLLamaMessage(msg))
		}
	}
	return messages
}

func ToLLamaMessage(msg llm.Message) api.Message {
	return api.Message{
		Role:    decodeRole(msg.Role),
		Content: msg.Content,
	}
}

func FromLLamaMessage(msg api.Message) llm.Message {
	return llm.Message{
		Ts:      time.Now(),
		Role:    encodeRole(msg.Role),
		Content: msg.Content,
	}
}

func encodeRole(role string) llm.MessageRole {
	switch role {
	case "system":
		return llm.System
	case "user":
		return llm.User
	case "assistant":
		return llm.Assistant
	default:
		return llm.User
	}
}

func decodeRole(role llm.MessageRole) string {
	switch role {
	case llm.System:
		return "system"
	case llm.User:
		return "user"
	case llm.Assistant:
		return "assistant"
	default:
		return "user"
	}
}
