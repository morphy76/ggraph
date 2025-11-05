package openai_test

import (
	"encoding/json"
	"strings"
	"testing"

	o "github.com/morphy76/ggraph/pkg/agent/openai"
	tool "github.com/morphy76/ggraph/pkg/agent/tool"
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

		// Create a mock tool for testing with the named function
		mockTool, err := tool.CreateTool[int](additionTool, "prompt: Addition tool for testing")
		if err != nil {
			t.Fatalf("Failed to create mock tool: %v", err)
		}

		// Create a tool registry
		tools := map[string]*tool.Tool{
			"additionTool": mockTool,
		}

		// Convert the OpenAI tool call to our internal ToolCall
		toolCall, err := o.ConvertToolCall(openAIToolCalls[0], tools)
		if err != nil {
			t.Fatalf("ConvertToolCall failed: %v", err)
		}

		// Verify the conversion
		expectedID := "call_DwJfIP6DcYFG3PiaVvxQoqAy"
		if toolCall.Id != expectedID {
			t.Errorf("Expected ID %s, got %s", expectedID, toolCall.Id)
		}

		if toolCall.UsingTool.Name != "additionTool" {
			t.Errorf("Expected tool name 'additionTool', got %s", toolCall.UsingTool.Name)
		}

		// Verify arguments were parsed correctly
		expectedAddend1 := float64(4) // JSON numbers are parsed as float64
		expectedAddend2 := float64(5)

		if toolCall.Arguments["addend1"] != expectedAddend1 {
			t.Errorf("Expected addend1 %v, got %v", expectedAddend1, toolCall.Arguments["addend1"])
		}

		if toolCall.Arguments["addend2"] != expectedAddend2 {
			t.Errorf("Expected addend2 %v, got %v", expectedAddend2, toolCall.Arguments["addend2"])
		}
	})

	t.Run("tool_not_found", func(t *testing.T) {
		// Test data with unknown tool name
		jsonData := `[{"id":"call_123","function":{"arguments":"{}","name":"unknownTool"},"type":"function"}]`

		var openAIToolCalls []openai.ChatCompletionMessageToolCallUnion
		err := json.Unmarshal([]byte(jsonData), &openAIToolCalls)
		if err != nil {
			t.Fatalf("Failed to parse JSON: %v", err)
		}

		// Empty tool registry
		tools := map[string]*tool.Tool{}

		// Convert should fail
		_, err = o.ConvertToolCall(openAIToolCalls[0], tools)
		if err == nil {
			t.Fatal("Expected error when tool not found, but got none")
		}

		expectedError := "tool 'unknownTool' not found in tool registry"
		if err.Error() != expectedError {
			t.Errorf("Expected error message '%s', got '%s'", expectedError, err.Error())
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

		mockTool, err := tool.CreateTool[int](additionTool, "prompt: Addition tool for testing")
		if err != nil {
			t.Fatalf("Failed to create mock tool: %v", err)
		}

		tools := map[string]*tool.Tool{
			"additionTool": mockTool,
		}

		toolCall, err := o.ConvertToolCall(openAIToolCalls[0], tools)
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

		mockTool, err := tool.CreateTool[int](additionTool, "prompt: Addition tool for testing")
		if err != nil {
			t.Fatalf("Failed to create mock tool: %v", err)
		}

		tools := map[string]*tool.Tool{
			"additionTool": mockTool,
		}

		// Convert should fail due to invalid arguments JSON
		_, err = o.ConvertToolCall(openAIToolCalls[0], tools)
		if err == nil {
			t.Fatal("Expected error when arguments JSON is invalid, but got none")
		}

		if !strings.Contains(err.Error(), "failed to parse tool arguments") {
			t.Errorf("Expected error message to contain 'failed to parse tool arguments', got '%s'", err.Error())
		}
	})
}
