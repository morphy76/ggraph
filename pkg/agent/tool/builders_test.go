package tool_test

import (
	"fmt"
	"testing"

	"github.com/morphy76/ggraph/pkg/agent/tool"
)

func addition[T int | float64](a, b T) (T, error) {
	return a + b, nil
}

func concat(a, b string) (string, error) {
	return a + b, nil
}

func reverse(s string) (string, error) {
	runes := []rune(s)
	for i, j := 0, len(runes)-1; i < j; i, j = i+1, j-1 {
		runes[i], runes[j] = runes[j], runes[i]
	}
	return string(runes), nil
}

func alwaysFail(s ...any) (any, error) {
	return nil, fmt.Errorf("always fail on purpose")
}

func TestTools(t *testing.T) {

	firstTool, err := tool.CreateTool[int](addition[int], "Prompt: Add two numbers together.", "Input: (a int, b int)", "Output: int")
	if err != nil {
		t.Fatalf("Failed to create firstTool: %v", err)
	}

	secondTool, err := tool.CreateTool[float64](addition[float64], "Prompt: Add two float64 numbers together.", "Input: (a float64, b float64)", "Output: float64")
	if err != nil {
		t.Fatalf("Failed to create secondTool: %v", err)
	}

	thirdTool, err := tool.CreateTool[string](concat, "Prompt: Concatenate two strings.", "Input: (a string, b string)", "Output: string")
	if err != nil {
		t.Fatalf("Failed to create thirdTool: %v", err)
	}

	fourthTool, err := tool.CreateTool[string](reverse, "Prompt: Reverse a string.", "Input: (s string)", "Output: string")
	if err != nil {
		t.Fatalf("Failed to create fourthTool: %v", err)
	}

	fifthTool, err := tool.CreateTool[any](alwaysFail, "Prompt: Always fail.", "Input: (...any)", "Output: (any, error)")
	if err != nil {
		t.Fatalf("Failed to create fifthTool: %v", err)
	}

	t.Run("exec_first_tool", func(t *testing.T) {
		result, err := firstTool.Call(3, 5)
		if err != nil {
			t.Errorf("firstTool.Call returned an error: %v", err)
		}
		expected := 8
		if result != expected {
			t.Errorf("firstTool.Call = %v; want %v", result, expected)
		}
	})

	t.Run("exec_second_tool", func(t *testing.T) {
		result, err := secondTool.Call(3.5, 2.5)
		if err != nil {
			t.Errorf("secondTool.Call returned an error: %v", err)
		}
		expected := 6.0
		if result != expected {
			t.Errorf("secondTool.Call = %v; want %v", result, expected)
		}
	})

	t.Run("exec_third_tool", func(t *testing.T) {
		result, err := thirdTool.Call("Hello, ", "world!")
		if err != nil {
			t.Errorf("thirdTool.Call returned an error: %v", err)
		}
		expected := "Hello, world!"
		if result != expected {
			t.Errorf("thirdTool.Call = %v; want %v", result, expected)
		}
	})

	t.Run("exec_fourth_tool", func(t *testing.T) {
		result, err := fourthTool.Call("Hello, world!")
		if err != nil {
			t.Errorf("fourthTool.Call returned an error: %v", err)
		}
		expected := "!dlrow ,olleH"
		if result != expected {
			t.Errorf("fourthTool.Call = %v; want %v", result, expected)
		}
	})

	t.Run("exec_fifth_tool", func(t *testing.T) {
		_, err := fifthTool.Call("This will fail")
		if err == nil {
			t.Errorf("fifthTool.Call expected to return an error, but got nil")
		}
	})

	t.Run("call_with_invalid_args_count", func(t *testing.T) {
		_, err := firstTool.Call(1) // expects 2 args but only 1 provided
		if err == nil {
			t.Errorf("Expected error for invalid argument count, but got nil")
		}
		if err != tool.ErrCallingToolInvalidArgsCount {
			t.Errorf("Expected ErrCallingToolInvalidArgsCount, got: %v", err)
		}
	})

	t.Run("call_with_too_many_args", func(t *testing.T) {
		_, err := firstTool.Call(1, 2, 3) // expects 2 args but 3 provided
		if err == nil {
			t.Errorf("Expected error for invalid argument count, but got nil")
		}
		if err != tool.ErrCallingToolInvalidArgsCount {
			t.Errorf("Expected ErrCallingToolInvalidArgsCount, got: %v", err)
		}
	})
}

func TestCreateToolErrors(t *testing.T) {
	t.Run("not_a_function", func(t *testing.T) {
		notAFunc := "this is a string, not a function"
		_, err := tool.CreateTool[string](notAFunc, "Prompt: Test")
		if err == nil {
			t.Errorf("Expected error when creating tool with non-function, but got nil")
		}
		if err != tool.ErrToolFnNotFunction {
			t.Errorf("Expected ErrToolFnNotFunction, got: %v", err)
		}
	})

	t.Run("invalid_return_count_zero", func(t *testing.T) {
		noReturn := func() {}
		_, err := tool.CreateTool[any](noReturn, "Prompt: Test")
		if err == nil {
			t.Errorf("Expected error for function with zero return values, but got nil")
		}
		if err != tool.ErrToolFnInvalidReturnCount {
			t.Errorf("Expected ErrToolFnInvalidReturnCount, got: %v", err)
		}
	})

	t.Run("invalid_return_count_one", func(t *testing.T) {
		oneReturn := func() int { return 1 }
		_, err := tool.CreateTool[int](oneReturn, "Prompt: Test")
		if err == nil {
			t.Errorf("Expected error for function with one return value, but got nil")
		}
		if err != tool.ErrToolFnInvalidReturnCount {
			t.Errorf("Expected ErrToolFnInvalidReturnCount, got: %v", err)
		}
	})

	t.Run("invalid_return_count_three", func(t *testing.T) {
		threeReturn := func() (int, int, error) { return 1, 2, nil }
		_, err := tool.CreateTool[int](threeReturn, "Prompt: Test")
		if err == nil {
			t.Errorf("Expected error for function with three return values, but got nil")
		}
		if err != tool.ErrToolFnInvalidReturnCount {
			t.Errorf("Expected ErrToolFnInvalidReturnCount, got: %v", err)
		}
	})

	t.Run("second_return_not_error", func(t *testing.T) {
		wrongSecondReturn := func() (int, int) { return 1, 2 }
		_, err := tool.CreateTool[int](wrongSecondReturn, "Prompt: Test")
		if err == nil {
			t.Errorf("Expected error for function where second return is not error, but got nil")
		}
		if err != tool.ErrToolFnInvalidReturnCount {
			t.Errorf("Expected ErrToolFnInvalidReturnCount, got: %v", err)
		}
	})

	t.Run("invalid_descriptor_format_no_colon", func(t *testing.T) {
		validFunc := func(a int) (int, error) { return a, nil }
		_, err := tool.CreateTool[int](validFunc, "Invalid descriptor without colon")
		if err == nil {
			t.Errorf("Expected error for invalid descriptor format, but got nil")
		}
		if err != tool.ErrInvalidDescriptorFormat {
			t.Errorf("Expected ErrInvalidDescriptorFormat, got: %v", err)
		}
	})

	t.Run("invalid_descriptor_format_multiple_colons", func(t *testing.T) {
		validFunc := func(a int) (int, error) { return a, nil }
		_, err := tool.CreateTool[int](validFunc, "Role:Description:Extra")
		if err == nil {
			t.Errorf("Expected error for invalid descriptor format, but got nil")
		}
		if err != tool.ErrInvalidDescriptorFormat {
			t.Errorf("Expected ErrInvalidDescriptorFormat, got: %v", err)
		}
	})

	t.Run("valid_tool_with_no_descriptors", func(t *testing.T) {
		validFunc := func(a int) (int, error) { return a * 2, nil }
		tool, err := tool.CreateTool[int](validFunc)
		if err != nil {
			t.Errorf("Expected no error for valid tool with no descriptors, got: %v", err)
		}
		if tool == nil {
			t.Errorf("Expected valid tool, got nil")
		}
	})

	t.Run("valid_tool_with_multiple_descriptors", func(t *testing.T) {
		validFunc := func(a int) (int, error) { return a * 2, nil }
		tool, err := tool.CreateTool[int](validFunc,
			"Prompt: Double the number",
			"Input: a int",
			"Output: int")
		if err != nil {
			t.Errorf("Expected no error for valid tool with multiple descriptors, got: %v", err)
		}
		if tool == nil {
			t.Errorf("Expected valid tool, got nil")
		}
		// Test the tool works
		result, err := tool.Call(5)
		if err != nil {
			t.Errorf("Unexpected error calling tool: %v", err)
		}
		if result != 10 {
			t.Errorf("Expected result 10, got: %v", result)
		}
	})
}
