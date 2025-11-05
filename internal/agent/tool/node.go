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
	rv, err := b.NewNodeBuilder(name, runToolsFunc(tools...)).
		WithReducer(toolExecutionReducer).
		Build()
	if err != nil {
		return nil, fmt.Errorf("failed to create the tool executor node: %w", err)
	}
	return rv, nil
}

// ------------------------------------------------------------------------------
// Node Implementation
// ------------------------------------------------------------------------------

func runToolsFunc(tools ...*t.Tool) g.NodeFn[a.Conversation] {
	return func(userInput, currentState a.Conversation, notifyPartial g.NotifyPartialFn[a.Conversation]) (a.Conversation, error) {
		toolCalls := currentState.ToolCalls
		if len(toolCalls) == 0 || len(tools) == 0 {
			return a.CreateConversation(), nil
		}

		// TODO assuming so far that there are no dependencies among tool calls, then I run all tool calls in parallel
		wg := sync.WaitGroup{}
		callStateMutex := sync.Mutex{}

		callState := a.CreateConversation()
		for _, call := range toolCalls {
			wg.Add(1)
			go func(tc t.ToolCall) {
				defer wg.Done()

				var useTool *t.Tool
				for _, tool := range tools {
					if tool.Name == tc.ToolName {
						useTool = tool
						break
					}
				}
				useArgs := call.ArgsAsSortedSlice(useTool)
				rv, err := useTool.Call(useArgs...)
				if err != nil {
					errorToolMessage := fmt.Sprintf("{\"id\": \"%s\", \"name\": \"%s\", \"error\": \"%s\"}", call.Id, useTool.Name, err)
					callStateMutex.Lock()
					callState.Messages = append(callState.Messages, a.CreateMessage(a.Tool, errorToolMessage))
					callStateMutex.Unlock()
				}

				resultToolMessage := fmt.Sprintf("{\"id\": \"%s\", \"name\": \"%s\", \"result\": %v}", call.Id, useTool.Name, rv)
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
	currentState.ToolCalls = change.ToolCalls
	return currentState
}
