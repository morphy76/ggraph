package aiw

import (
	"context"
	"os"
	"time"

	"github.com/openai/openai-go/v3"

	b "github.com/morphy76/ggraph/pkg/builders"
	g "github.com/morphy76/ggraph/pkg/graph"
	"github.com/morphy76/ggraph/pkg/llm"
)

// CreateAIWChatNodeFromEnvironment creates a graph node that interacts with the AIW chat model.
func CreateAIWChatNodeFromEnvironment(name string, model string) (g.Node[llm.AgentModel], error) {

	// TODO, this is not the right way to set the base URL
	os.Setenv("OPENAI_BASE_URL", "https://portal.aiwave.ai/llm/api")

	client := openai.NewClient()

	chatFunction := func(userInput, currentState llm.AgentModel, notifyPartial g.NotifyPartialFn[llm.AgentModel]) (llm.AgentModel, error) {
		ctx, cancel := context.WithTimeout(context.Background(), 60*time.Second)
		defer cancel()

		chatCompletion, err := client.Chat.Completions.New(ctx, openai.ChatCompletionNewParams{
			Messages:    ToOpenAIModel(currentState, userInput),
			Model:       model,
			Temperature: openai.Float(0.5),
		})
		if err != nil {
			return currentState, err
		}

		currentState.Messages = append(currentState.Messages, FromOpenAIMessage(chatCompletion.Choices[0].Message))

		return currentState, nil
	}

	return b.CreateNode(name, chatFunction)
}
