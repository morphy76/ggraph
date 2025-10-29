package aiw

import (
	"os"

	"github.com/openai/openai-go/v3/option"

	a "github.com/morphy76/ggraph/pkg/agent"
	o "github.com/morphy76/ggraph/pkg/agent/openai"
	b "github.com/morphy76/ggraph/pkg/builders"
	g "github.com/morphy76/ggraph/pkg/graph"
)

const (
	// AIWBaseURL is the base URL for the Almawave AIW Platform.
	AIWBaseURL = "https://portal.aiwave.ai/llm/api"
	// EnvKeyPAT is the environment variable key for the AIW API key.
	EnvKeyPAT = "AIW_API_KEY"
)

// PATFromEnv retrieves the AIW API key from the environment variable "AIW_API_KEY
//
// Returns:
//   - The Personal Access Token (PAT) as a string.
func PATFromEnv() string {
	return os.Getenv(EnvKeyPAT)
}

// CreateCompletionNode creates a graph node for an AIW-based chat agent.
//
// Parameters:
//   - name: The unique name for the node.
//   - PAT: The API key for authentication.
//   - model: The OpenAI model to be used for the chat agent.
//   - completionNodeFn: A function that creates the node function for the AIW chat agent.
//   - opts: Additional request options for the OpenAI API calls.
//
// Returns:
//   - An instance of g.Node[a.Conversation] configured for the OpenAI chat agent.
//   - An error if the node creation fails.
//
// Example usage:
//
//	node, err := CreateCompletionNode("ChatNode",  "your-api-key", "velvet-2b", myOpenAINodeFn)
func CreateCompletionNode(
	name, PAT, model string,
	completionNodeFn o.CompletionNodeFn,
	opts ...option.RequestOption,
) (g.Node[a.Completion], error) {
	client := o.NewClient(AIWBaseURL, PAT, opts...)
	openAIFn := completionNodeFn(client.Completions, model, opts...)

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
