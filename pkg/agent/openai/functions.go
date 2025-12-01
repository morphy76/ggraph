package openai

import (
	"context"
	"fmt"
	"strings"

	"github.com/openai/openai-go/v3"

	a "github.com/morphy76/ggraph/pkg/agent"
	t "github.com/morphy76/ggraph/pkg/agent/tool"
	g "github.com/morphy76/ggraph/pkg/graph"
)

var completionNodeFn CompletionNodeFn = func(completionService openai.CompletionService, completionOptions a.ModelOptions) g.NodeFn[a.Completion] {
	return func(userInput, currentState a.Completion, notify g.NotifyPartialFn[a.Completion]) (a.Completion, error) {
		var usePrompt string
		if completionOptions.PromptFormat != "" && strings.Contains(completionOptions.PromptFormat, "%s") {
			usePrompt = fmt.Sprintf(completionOptions.PromptFormat, userInput.Text)
		} else {
			usePrompt = userInput.Text
		}

		completionOptions.Prompt = usePrompt

		openAIOpts := ConvertCompletionOptions(completionOptions)
		resp, err := completionService.New(context.Background(), openAIOpts)
		if err != nil {
			return currentState, fmt.Errorf("failed to generate completion: %w", err)
		}

		if len(resp.Choices) > 0 {
			currentState = a.CreateCompletion(resp.Choices[0].Text)
		} else {
			return currentState, fmt.Errorf("no completion choices returned")
		}

		return currentState, nil
	}
}

var conversationNodeFn ConversationNodeFn = func(chatService openai.ChatService, conversationOptions a.ModelOptions) g.NodeFn[a.Conversation] {
	return func(userInput, currentState a.Conversation, notify g.NotifyPartialFn[a.Conversation]) (a.Conversation, error) {
		useMessages := []a.Message{}
		if len(currentState.Messages) > 0 {
			useMessages = currentState.Messages
		} else {
			useMessages = append(useMessages, userInput.Messages...)
		}

		filteredMessages := []a.Message{}
		for _, msg := range useMessages {
			if msg.Role != a.System {
				filteredMessages = append(filteredMessages, msg)
			}
		}

		systemMessages := []a.Message{}
		for _, msg := range conversationOptions.Messages {
			if msg.Role == a.System {
				systemMessages = append(systemMessages, msg)
			}
		}
		useMessages = filteredMessages
		if len(systemMessages) > 0 {
			useMessages = append(systemMessages, useMessages...)
		}

		currentState.Messages = useMessages
		conversationOptions.Messages = currentState.Messages

		openAIOpts := ConvertConversationOptions(conversationOptions)

		resp, err := chatService.Completions.New(context.Background(), openAIOpts)
		if err != nil {
			return currentState, fmt.Errorf("failed to generate tool calls: %w", err)
		}

		answer := resp.Choices[0].Message
		useAnswer := a.CreateMessage(a.Assistant, answer.Content)
		requestedToolCalls := resp.Choices[0].Message.ToolCalls
		if len(requestedToolCalls) > 0 {
			toolCalls := make([]t.FnCall, 0, len(requestedToolCalls))
			for _, openAIToolCall := range requestedToolCalls {
				toolCall, err := ConvertToolCall(openAIToolCall)
				if err != nil {
					return currentState, fmt.Errorf("failed to convert tool call: %w", err)
				}
				toolCalls = append(toolCalls, *toolCall)
			}
			useAnswer.ToolCalls = toolCalls
			currentState.CurrentToolCalls = toolCalls
		}
		currentState.Messages = append(currentState.Messages, useAnswer)

		return currentState, nil
	}
}
