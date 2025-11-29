package tool

import (
	"fmt"
	"sync"

	a "github.com/morphy76/ggraph/pkg/agent"
	t "github.com/morphy76/ggraph/pkg/agent/tool"
	b "github.com/morphy76/ggraph/pkg/builders"
	g "github.com/morphy76/ggraph/pkg/graph"
)

// NodeToolFactory creates a new instance of a Node capable of processing tool calls within an agent conversation.
func NodeToolFactory(name string, tools ...*t.Tool) (g.Node[a.Conversation], error) {
	rv, err := b.NewNode(name, runToolsFunc(tools...),
		g.WithReducer(toolExecutionReducer))
	if err != nil {
		return nil, fmt.Errorf("failed to create the tool executor node: %w", err)
	}
	return rv, nil
}

// ------------------------------------------------------------------------------
// Node Implementation
// ------------------------------------------------------------------------------

func runToolsFunc(tools ...*t.Tool) g.NodeFn[a.Conversation] {

	mappedTools := make(map[string]*t.Tool)
	for _, tool := range tools {
		mappedTools[tool.Name] = tool
	}

	return func(userInput, currentState a.Conversation, notifyPartial g.NotifyPartialFn[a.Conversation]) (a.Conversation, error) {
		toolCalls := currentState.CurrentToolCalls
		if len(toolCalls) == 0 || len(tools) == 0 {
			return a.CreateConversation(), nil
		}

		// TODO assuming so far that there are no dependencies among tool calls, then I run all tool calls in parallel
		wg := sync.WaitGroup{}
		callStateMutex := sync.Mutex{}

		callState := a.CreateConversation()
		for _, call := range toolCalls {
			wg.Add(1)
			go func(tc t.FnCall) {
				defer wg.Done()

				useTool, found := mappedTools[tc.ToolName]
				if !found {
					errorToolMessage := fmt.Sprintf("%s:%s", call.ID, t.ErrToolNotFound)
					callStateMutex.Lock()
					callState.Messages = append(callState.Messages, a.CreateMessage(a.Tool, errorToolMessage))
					callStateMutex.Unlock()
					return
				}
				useArgs := call.ArgsAsSortedSlice(useTool)
				rv, err := useTool.Call(useArgs...)
				if err != nil {
					errorToolMessage := fmt.Sprintf("%s:%s", call.ID, err)
					callStateMutex.Lock()
					callState.Messages = append(callState.Messages, a.CreateMessage(a.Tool, errorToolMessage))
					callStateMutex.Unlock()
					return
				}

				resultToolMessage := fmt.Sprintf("%s:%v", call.ID, rv)
				callStateMutex.Lock()
				callState.Messages = append(callState.Messages, a.CreateMessage(a.Tool, resultToolMessage))
				callStateMutex.Unlock()
			}(call)
		}

		wg.Wait()

		return callState, nil
	}
}

func toolExecutionReducer(currentState, change a.Conversation) a.Conversation {
	currentState.Messages = append(currentState.Messages, change.Messages...)
	currentState.CurrentToolCalls = change.CurrentToolCalls

	return currentState
}
