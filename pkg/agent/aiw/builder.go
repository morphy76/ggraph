package aiw

import (
	"github.com/openai/openai-go/v3"
	"github.com/openai/openai-go/v3/option"

	a "github.com/morphy76/ggraph/pkg/agent"
	o "github.com/morphy76/ggraph/pkg/agent/openai"
	b "github.com/morphy76/ggraph/pkg/builders"
	g "github.com/morphy76/ggraph/pkg/graph"
)

// NewAIWClient creates a new OpenAI client configured for the AIW platform.
//
// Parameters:
//   - PAT: The Personal Access Token (PAT) for authentication.
//   - opts: Additional request options for the OpenAI API calls.
//
// Returns:
//   - A pointer to an instance of openai.Client configured for AIW.
//
// Example usage:
//
//	client := NewAIWClient("your-api-key", option.WithTimeout(30*time.Second))
func NewAIWClient(
	PAT string,
	opts ...option.RequestOption,
) *openai.Client {
	return o.NewClient(AIWBaseURL, PAT, opts...)
}

// CreateCompletionNode creates a graph node for an AIW-based chat agent.
//
// Parameters:
//   - name: The unique name for the node.
//   - model: The OpenAI model to be used for the chat agent.
//   - client: The OpenAI client instance.
//   - completionNodeFn: A function that creates the node function for the AIW chat agent.
//   - completionOptions: Additional completion options for the OpenAI API calls.
//
// Returns:
//   - An instance of g.Node[a.Conversation] configured for the OpenAI chat agent.
//   - An error if the node creation fails.
//
// Example usage:
//
//	node, err := CreateCompletionNode("ChatNode",  "your-api-key", "velvet-2b", myOpenAINodeFn)
func CreateCompletionNode(
	name, model string,
	client *openai.Client,
	completionNodeFn o.CompletionNodeFn,
	completionOptions ...a.CompletionOption,
) (g.Node[a.Completion], error) {
	openAIFn := completionNodeFn(client.Completions, model, completionOptions...)

	rv, err := b.NewNodeBuilder(name, openAIFn).Build()
	return rv, err
}

// CreateConversationNode creates a graph node for an AIW-based chat agent.
//
// Parameters:
//   - name: The unique name for the node.
//   - PAT: The API key for authentication.
//   - model: The OpenAI model to be used for the chat agent.
//   - conversationNodeFn: A function that creates the node function for the AIW chat agent.
//   - opts: Additional request options for the OpenAI API calls.
//
// Returns:
//   - An instance of g.Node[a.Conversation] configured for the OpenAI chat agent.
//   - An error if the node creation fails.
//
// Example usage:
//
//	node, err := CreateConversationNode("ChatNode",  "your-api-key", "velvet-2b", myOpenAINodeFn)
func CreateConversationNode(
	name, PAT, model string,
	conversationNodeFn o.ConversationNodeFn,
	opts ...option.RequestOption,
) (g.Node[a.Conversation], error) {
	client := o.NewClient(AIWBaseURL, PAT, opts...)
	openAIFn := conversationNodeFn(client.Chat, model, opts...)

	rv, err := b.NewNodeBuilder(name, openAIFn).Build()
	return rv, err
}
