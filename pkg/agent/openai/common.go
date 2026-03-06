package openai

import (
	"os"

	"github.com/openai/openai-go/v3"

	a "github.com/morphy76/ggraph/pkg/agent"
	g "github.com/morphy76/ggraph/pkg/graph"
)

const (
	// OpenAIBaseURL is the base URL for the OpenAI API.
	OpenAIBaseURL = "https://api.openai.com/v1"
	// EnvKeyAPIKey is the environment variable key for the OpenAI API key.
	EnvKeyAPIKey = "OPENAI_API_KEY"
)

// APIKeyFromEnv retrieves the OpenAI API key from the environment variable "OPENAI_API_KEY".
//
// Returns:
//   - The OpenAI API key as a string.
func APIKeyFromEnv() string {
	return os.Getenv(EnvKeyAPIKey)
}

// CompletionNodeFn defines a function type that creates a node function for an OpenAI-based chat agent.
//
// Parameters:
//   - completionService: The OpenAI CompletionService client.
//   - modelOptions: Additional request options for the OpenAI API calls.
//
// Returns:
//   - A g.CompletionNodeFn[a.Completion] function that handles the chat agent's completion logic.
//
// Example usage:
//
//	var chatNodeFn CompletionNodeFn = func(client openai.Client, modelOptions a.ModelOptions) g.CompletionNodeFn[a.Completion] {
//	    return func(userInput, currentState a.Conversation, notify g.NotifyPartialFn[a.Conversation]) g.CompletionNodeFn[a.Completion] {
//	        // Implementation here...
//	    }
//	}
type CompletionNodeFn func(completionService openai.CompletionService, modelOptions a.ModelOptions) g.NodeFn[a.Completion]

// ConversationNodeFn defines a function type that creates a node function for an OpenAI-based chat agent.
//
// Parameters:
//   - chatService: The OpenAI ChatService client.
//   - modelOptions: Additional request options for the OpenAI API calls.
//
// Returns:
//   - A g.NodeFn[a.Conversation] function that handles the chat agent's conversation logic.
//
// Example usage:
//
//	var chatNodeFn ConversationNodeFn = func(chatService openai.ChatService, modelOptions a.ModelOptions) g.NodeFn[a.Conversation] {
//	    return func(userInput, currentState a.Conversation, notify g.NotifyPartialFn[a.Conversation]) g.NodeFn[a.Conversation] {
//	        // Implementation here...
//	    }
//	}
type ConversationNodeFn func(chatService openai.ChatService, modelOptions a.ModelOptions) g.NodeFn[a.Conversation]
