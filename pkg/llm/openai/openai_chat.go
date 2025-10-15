package openai

import (
	"context"
	"time"

	"github.com/openai/openai-go/v3"

	b "github.com/morphy76/ggraph/pkg/builders"
	g "github.com/morphy76/ggraph/pkg/graph"
	"github.com/morphy76/ggraph/pkg/llm"
)

// CreateOpenAIChatNode creates a graph node that interacts with the OpenAI chat model.
func CreateOpenAIChatNodeFromEnvironment(name string, model string) (g.Node[llm.AgentModel], error) {
	client := openai.NewClient()

	chatFunction := func(state llm.AgentModel, notify func(llm.AgentModel)) (llm.AgentModel, error) {
		ctx, cancel := context.WithTimeout(context.Background(), 60*time.Second)
		defer cancel()

		client.Chat.Completions.New(ctx, openai.ChatCompletionNewParams{
			Model:    model,
			Messages: []openai.ChatCompletionMessageParamUnion{},
		})
		return state, nil
	}

	return b.CreateNode(name, chatFunction)
}
