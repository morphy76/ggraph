package openai

import (
	"fmt"

	"github.com/openai/openai-go/v3"
	"github.com/openai/openai-go/v3/option"

	a "github.com/morphy76/ggraph/pkg/agent"
	b "github.com/morphy76/ggraph/pkg/builders"
	g "github.com/morphy76/ggraph/pkg/graph"
)

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
) *openai.Client {
	useOpts := append(opts,
		option.WithBaseURL(baseURL),
		option.WithAPIKey(apiKey),
	)
	rv := openai.NewClient(useOpts...)
	return &rv
}

// NewOpenAIClient creates a new OpenAI client using the default OpenAI base URL.
//
// Parameters:
//   - apiKey: The API key for authentication.
//   - opts: Additional request options.
//
// Returns:
//   - An instance of openai.Client configured with the OpenAI base URL and provided parameters.
//
// Example usage:
//
//	client := NewOpenAIClient("your-api-key")
func NewOpenAIClient(
	apiKey string,
	opts ...option.RequestOption,
) *openai.Client {
	return NewClient(OpenAIBaseURL, apiKey, opts...)
}

// CreateCompletionNode creates a graph node for an OpenAI-based chat agent.
//
// Parameters:
//   - name: The unique name for the node.
//   - model: The OpenAI model to be used for the chat agent.
//   - client: The OpenAI client instance.
//   - completionNodeFn: A function that creates the node function for the OpenAI chat agent.
//   - completionOptions: Additional completion options for the OpenAI API calls.
//
// Returns:
//   - An instance of g.Node[a.Conversation] configured for the OpenAI chat agent.
//   - An error if the node creation fails.
//
// Example usage:
//
//	node, err := CreateCompletionNode("ChatNode",  "your-api-key", "gpt-4", myOpenAINodeFn)
func CreateCompletionNode(
	name, model string,
	client *openai.Client,
	completionNodeFn CompletionNodeFn,
	completionOptions ...a.ModelOption,
) (g.Node[a.Completion], error) {
	openAIFn := completionNodeFn(client.Completions, model, completionOptions...)

	rv, err := b.NewNodeBuilder(name, openAIFn).Build()
	return rv, err
}

// CreateConversationNode creates a graph node for an OpenAI-based chat agent.
//
// Parameters:
//   - name: The unique name for the node.
//   - model: The OpenAI model to be used for the chat agent.
//   - client: The OpenAI client instance.
//   - conversationNodeFn: A function that creates the node function for the OpenAI chat agent.
//   - conversationOptions: Additional conversation options for the OpenAI API calls.
//
// Returns:
//   - An instance of g.Node[a.Conversation] configured for the OpenAI chat agent.
//   - An error if the node creation fails.
//
// Example usage:
//
//	node, err := CreateConversationNode("ChatNode",  "your-api-key", "gpt-4", myOpenAINodeFn)
func CreateConversationNode(
	name, model string,
	client *openai.Client,
	conversationNodeFn ConversationNodeFn,
	conversationOptions ...a.ModelOption,
) (g.Node[a.Conversation], error) {
	openAIFn := conversationNodeFn(client.Chat, model, conversationOptions...)

	routingPolicy, err := b.CreateConditionalRoutePolicy(a.ToolProcessorRoutingFn)
	if err != nil {
		return nil, fmt.Errorf("cannot create a conversation node %w", err)
	}

	rv, err := b.NewNodeBuilder(name, openAIFn).
		WithRoutingPolicy(routingPolicy).
		Build()
	return rv, err
}
