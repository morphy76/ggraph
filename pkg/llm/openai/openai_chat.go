package openai

import (
	"context"
	"time"

	"github.com/openai/openai-go/v3"

	b "github.com/morphy76/ggraph/pkg/builders"
	g "github.com/morphy76/ggraph/pkg/graph"
	"github.com/morphy76/ggraph/pkg/llm"
)

// CreateOpenAIChatNodeFromEnvironment creates a graph node that interacts with the OpenAI chat model.
func CreateOpenAIChatNodeFromEnvironment(name string, model string) (g.Node[llm.AgentModel], error) {
	client := openai.NewClient()

	chatFunction := func(userInput, currentState llm.AgentModel, notifyPartial g.NotifyPartialFn[llm.AgentModel]) (llm.AgentModel, error) {
		ctx, cancel := context.WithTimeout(context.Background(), 60*time.Second)
		defer cancel()

		acc := openai.ChatCompletionAccumulator{}
		stream := client.Chat.Completions.NewStreaming(ctx, openai.ChatCompletionNewParams{
			Model:    model,
			Messages: ToOpenAIModel(currentState, userInput),
			Seed:     openai.Int(0),
		})

		for stream.Next() {
			chunk := stream.Current()
			acc.AddChunk(chunk)
			if len(chunk.Choices) > 0 {
				notifyPartial(llm.AgentModel{
					Messages: []llm.Message{{Role: llm.Assistant, Content: chunk.Choices[0].Delta.Content}},
				})
			}
		}

		if stream.Err() != nil {
			return currentState, stream.Err()
		}

		currentState.Messages = append(currentState.Messages, FromOpenAIMessage(acc.Choices[0].Message))

		return currentState, nil
	}

	rv, err := b.NewNodeBuilder(name, chatFunction).Build()
	return rv, err
}
