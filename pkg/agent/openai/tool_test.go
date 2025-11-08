package openai_test

import (
	"encoding/json"
	"strings"
	"testing"

	o "github.com/morphy76/ggraph/pkg/agent/openai"
	"github.com/morphy76/ggraph/pkg/agent/tool"
	"github.com/openai/openai-go/v3"
)

// Helper function for testing - named so the tool will have a proper name
func additionTool(addend1, addend2 int) (int, error) {
	return addend1 + addend2, nil
}

func TestConvertToolCall(t *testing.T) {
	t.Run("successful_conversion", func(t *testing.T) {
		// Test data from the user's request
		jsonData := `[{"id":"call_DwJfIP6DcYFG3PiaVvxQoqAy","function":{"arguments":"{\"addend1\":4,\"addend2\":5}","name":"additionTool"},"type":"function","custom":{"input":"","name":""}}]`

		// Parse the OpenAI tool calls
		var openAIToolCalls []openai.ChatCompletionMessageToolCallUnion
		err := json.Unmarshal([]byte(jsonData), &openAIToolCalls)
		if err != nil {
			t.Fatalf("Failed to parse JSON: %v", err)
		}

		// Convert the OpenAI tool call to our internal ToolCall
		toolCall, err := o.ConvertToolCall(openAIToolCalls[0])
		if err != nil {
			t.Fatalf("ConvertToolCall failed: %v", err)
		}

		// Verify the conversion
		expectedID := "call_DwJfIP6DcYFG3PiaVvxQoqAy"
		if toolCall.ID != expectedID {
			t.Errorf("Expected ID %s, got %s", expectedID, toolCall.ID)
		}

		if toolCall.ToolName != "additionTool" {
			t.Errorf("Expected tool name 'additionTool', got %s", toolCall.ToolName)
		}

		// additionTool expects two arguments: addend1 and addend2
		additionTool, err := tool.CreateTool[int](func(addend1, addend2 int) (int, error) {
			return addend1 + addend2, nil
		}, "input:addend1, addend2")
		if err != nil {
			t.Fatalf("Failed to create addition tool: %v", err)
		}

		// Verify arguments were parsed correctly
		expectedAddend1 := float64(4) // JSON numbers are parsed as float64
		expectedAddend2 := float64(5)

		useArgs := toolCall.ArgsAsSortedSlice(additionTool) // Using the predefined tool for argument order
		if useArgs[0] != expectedAddend1 {
			t.Errorf("Expected addend1 %v, got %v", expectedAddend1, useArgs[0])
		}

		if useArgs[1] != expectedAddend2 {
			t.Errorf("Expected addend2 %v, got %v", expectedAddend2, useArgs[1])
		}
	})

	t.Run("empty_arguments", func(t *testing.T) {
		// Test data with empty arguments
		jsonData := `[{"id":"call_456","function":{"arguments":"","name":"additionTool"},"type":"function"}]`

		var openAIToolCalls []openai.ChatCompletionMessageToolCallUnion
		err := json.Unmarshal([]byte(jsonData), &openAIToolCalls)
		if err != nil {
			t.Fatalf("Failed to parse JSON: %v", err)
		}

		toolCall, err := o.ConvertToolCall(openAIToolCalls[0])
		if err != nil {
			t.Fatalf("ConvertToolCall failed: %v", err)
		}

		// Should have empty arguments map
		if len(toolCall.Arguments) != 0 {
			t.Errorf("Expected empty arguments, got %v", toolCall.Arguments)
		}
	})

	t.Run("invalid_json_arguments", func(t *testing.T) {
		// Test data with invalid JSON arguments
		jsonData := `[{"id":"call_789","function":{"arguments":"{invalid_json","name":"additionTool"},"type":"function"}]`

		var openAIToolCalls []openai.ChatCompletionMessageToolCallUnion
		err := json.Unmarshal([]byte(jsonData), &openAIToolCalls)
		if err != nil {
			t.Fatalf("Failed to parse JSON: %v", err)
		}

		// Convert should fail due to invalid arguments JSON
		_, err = o.ConvertToolCall(openAIToolCalls[0])
		if err == nil {
			t.Fatal("Expected error when arguments JSON is invalid, but got none")
		}

		if !strings.Contains(err.Error(), "failed to parse tool arguments") {
			t.Errorf("Expected error message to contain 'failed to parse tool arguments', got '%s'", err.Error())
		}
	})
}
