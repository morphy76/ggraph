package openai

import (
	"os"

	"github.com/openai/openai-go/v3"
	"github.com/openai/openai-go/v3/option"

	a "github.com/morphy76/ggraph/pkg/agent"
	b "github.com/morphy76/ggraph/pkg/builders"
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
//   - model: The OpenAI model to be used for the chat agent.
//   - opts: Additional request options for the OpenAI API calls.
//
// Returns:
//   - A g.CompletionNodeFn[a.Completion] function that handles the chat agent's completion logic.
//
// Example usage:
//
//	var chatNodeFn CompletionNodeFn = func(client openai.Client, model string, opts ...option.RequestOption) g.CompletionNodeFn[a.Completion] {
//	    return func(userInput, currentState a.Conversation, notify g.NotifyPartialFn[a.Conversation]) g.CompletionNodeFn[a.Completion] {
//	        // Implementation here...
//	    }
//	}
type CompletionNodeFn func(completionService openai.CompletionService, model string, opts ...option.RequestOption) g.NodeFn[a.Completion]

// ConversationNodeFn defines a function type that creates a node function for an OpenAI-based chat agent.
//
// Parameters:
//   - chatService: The OpenAI ChatService client.
//   - model: The OpenAI model to be used for the chat agent.
//   - opts: Additional request options for the OpenAI API calls.
//
// Returns:
//   - A g.NodeFn[a.Conversation] function that handles the chat agent's conversation logic.
//
// Example usage:
//
//	var chatNodeFn ConversationNodeFn = func(client openai.Client, model string, opts ...option.RequestOption) g.NodeFn[a.Conversation] {
//	    return func(userInput, currentState a.Conversation, notify g.NotifyPartialFn[a.Conversation]) g.NodeFn[a.Conversation] {
//	        // Implementation here...
//	    }
//	}
type ConversationNodeFn func(chatService openai.ChatService, model string, opts ...option.RequestOption) g.NodeFn[a.Conversation]

// NewClient creates a new OpenAI client with the specified base URL and API key.
//
// Additional request options can be provided as needed.
//
// Parameters:
//   - baseURL: The base URL for the OpenAI API.
//   - apiKey: The API key for authentication.
//   - opts: Additional request options.
//
// Returns:
//   - An instance of openai.Client configured with the provided parameters.
//
// Example usage:
//
//	client := NewClient("https://custom-openai-endpoint.com/v1", "your-api-key")
func NewClient(
	baseURL, apiKey string,
	opts ...option.RequestOption,
) openai.Client {
	useOpts := append(opts,
		option.WithBaseURL(""),
		option.WithAPIKey(""),
	)
	return openai.NewClient(useOpts...)
}

// CreateCompletionNode creates a graph node for an OpenAI-based chat agent.
//
// Parameters:
//   - name: The unique name for the node.
//   - apiKey: The API key for authentication.
//   - model: The OpenAI model to be used for the chat agent.
//   - completionNodeFn:
//   - opts: Additional request options for the OpenAI API calls.
//
// Returns:
//   - An instance of g.Node[a.Conversation] configured for the OpenAI chat agent.
//   - An error if the node creation fails.
//
// Example usage:
//
//	node, err := CreateCompletionNode("ChatNode",  "your-api-key", "gpt-4", myOpenAINodeFn)
func CreateCompletionNode(
	name, apiKey, model string,
	completionNodeFn CompletionNodeFn,
	opts ...option.RequestOption,
) (g.Node[a.Completion], error) {
	client := NewClient(OpenAIBaseURL, apiKey, opts...)
	openAIFn := completionNodeFn(client.Completions, model, opts...)

	rv, err := b.NewNodeBuilder(name, openAIFn).Build()
	return rv, err
}

// CreateConversationNode creates a graph node for an OpenAI-based chat agent.
//
// Parameters:
//   - name: The unique name for the node.
//   - apiKey: The API key for authentication.
//   - model: The OpenAI model to be used for the chat agent.
//   - conversationNodeFn: A function that creates the node function for the OpenAI chat agent.
//   - opts: Additional request options for the OpenAI API calls.
//
// Returns:
//   - An instance of g.Node[a.Conversation] configured for the OpenAI chat agent.
//   - An error if the node creation fails.
//
// Example usage:
//
//	node, err := CreateConversationNode("ChatNode",  "your-api-key", "gpt-4", myOpenAINodeFn)
func CreateConversationNode(
	name, apiKey, model string,
	conversationNodeFn ConversationNodeFn,
	opts ...option.RequestOption,
) (g.Node[a.Conversation], error) {
	client := NewClient(OpenAIBaseURL, apiKey, opts...)
	openAIFn := conversationNodeFn(client.Chat, model, opts...)

	rv, err := b.NewNodeBuilder(name, openAIFn).Build()
	return rv, err
}
