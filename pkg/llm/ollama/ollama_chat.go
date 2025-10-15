package ollama

import (
	"context"
	"time"

	"github.com/ollama/ollama/api"

	b "github.com/morphy76/ggraph/pkg/builders"
	g "github.com/morphy76/ggraph/pkg/graph"
	"github.com/morphy76/ggraph/pkg/llm"
)

// CreateOLLamaChatNode creates a graph node that interacts with the Ollama chat model.
func CreateOLLamaChatNodeFromEnvironment(name string, model string) (g.Node[llm.AgentModel], error) {
	client, err := api.ClientFromEnvironment()
	if err != nil {
		return nil, err
	}

	chatFunction := func(state llm.AgentModel, notify func(llm.AgentModel)) (llm.AgentModel, error) {
		ctx, cancel := context.WithTimeout(context.Background(), 60*time.Second)
		defer cancel()

		chatRequest := &api.ChatRequest{
			Model:    model,
			Messages: ToLLamaModel(state),
		}

		init := false
		respFunc := func(response api.ChatResponse) error {
			mex := FromLLamaMessage(response.Message)
			if !init {
				state.Messages = append(state.Messages)
				init = true
			} else {
				state.Messages[len(state.Messages)-1].Content += mex.Content
			}
			notify(llm.AgentModel{
				Messages: []llm.Message{mex},
			})
			return nil
		}

		err = client.Chat(ctx, chatRequest, respFunc)
		if err != nil {
			return state, err
		}
		return state, nil
	}

	return b.CreateNode(name, chatFunction)
}
