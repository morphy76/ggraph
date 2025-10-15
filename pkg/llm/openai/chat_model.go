package openai

import (
	"time"

	"github.com/morphy76/ggraph/pkg/llm"
	"github.com/openai/openai-go/v3"
)

func ToOpenAIModel(model llm.AgentModel) []openai.ChatCompletionMessageParamUnion {
	rv := make([]openai.ChatCompletionMessageParamUnion, len(model.Messages))
	for i, msg := range model.Messages {
		rv[i] = ToOpenAIMessage(msg)
	}
	return rv
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

func FromOpenAIMessage(msg openai.ChatCompletionResponse) llm.Message {
	return llm.Message{
		Ts:      time.Now(),
		Role:    encodeRole(msg.Role),
		Content: msg.Content,
	}
}

func encodeRole(role string) openai {
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
