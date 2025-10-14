package ollama

import (
	"context"
	"time"

	"github.com/ollama/ollama/api"

	"github.com/morphy76/ggraph/pkg/graph"
)

// CreateOLLamaChatNode creates a graph node that interacts with the Ollama chat model.
func CreateOLLamaChatNodeFromEnvironment(name string, model string) (graph.Node[ChatModel], error) {
	client, err := api.ClientFromEnvironment()
	if err != nil {
		return nil, err
	}

	chatFunction := func(state ChatModel, notify func(ChatModel)) (ChatModel, error) {
		ctx, cancel := context.WithTimeout(context.Background(), 60*time.Second)
		defer cancel()

		chatRequest := &api.ChatRequest{
			Model:    model,
			Messages: state.messages,
		}

		init := false
		respFunc := func(response api.ChatResponse) error {
			if !init {
				state.messages = append(state.messages, response.Message)
				init = true
			} else {
				state.messages[len(state.messages)-1].Content += response.Message.Content
			}
			notify(ChatModel{
				messages: []api.Message{response.Message},
			})
			return nil
		}

		err = client.Chat(ctx, chatRequest, respFunc)
		if err != nil {
			return state, err
		}
		return state, nil
	}

	return graph.CreateNode(name, chatFunction)
}
