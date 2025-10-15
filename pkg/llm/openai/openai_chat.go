package openai

import (
	"context"
	"time"

	"github.com/openai/openai-go/v3"

	b "github.com/morphy76/ggraph/pkg/builders"
	g "github.com/morphy76/ggraph/pkg/graph"
)

// CreateOpenAIChatNode creates a graph node that interacts with the OpenAI chat model.
func CreateOpenAIChatNodeFromEnvironment(name string, model string) (g.Node[ChatModel], error) {
	client := openai.NewClient()

	chatFunction := func(state ChatModel, notify func(ChatModel)) (ChatModel, error) {
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
