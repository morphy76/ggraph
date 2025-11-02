package tool

import (
	"reflect"
	"strings"
)

// CreateTool creates a new Tool instance from the provided function and descriptors.
//
// T is the return type of the tool function along with an error.
// The tool function must have the signature: func(args...) (T, error)
//
// Common descriptors roles could include:
//   - "Prompt": A brief description of the tool's purpose.
//   - "Usage": Instructions on how to use the tool.
//   - "Input": Description of the expected input.
//   - "Output": Description of the expected output.
//   - "Example": An example of how to use the tool.
//
// Parameters:
//   - fn: The function to be wrapped as a tool.
//   - descriptors: A variable number of strings describing the tool's purpose and usage in the format "role:description".
//
// Returns:
//   - *Tool[T]: A pointer to the created Tool instance.
//   - error: An error if the function is not valid or descriptors are incorrectly formatted.
//
// Example:
//
//	myTool, err := tool.CreateTool[string](myFunc, "Prompt: This tool does X", "Usage: Call with Y")
func CreateTool[T any](fn any, descriptors ...string) (*Tool[T], error) {
	fnType := reflect.TypeOf(fn)
	fnValue := reflect.ValueOf(fn)

	if fnType.Kind() != reflect.Func {
		return nil, ErrToolFnNotFunction
	}
	if fnType.NumOut() != 2 {
		return nil, ErrToolFnInvalidReturnCount
	}
	if !fnType.Out(1).Implements(reflect.TypeOf((*error)(nil)).Elem()) {
		return nil, ErrToolFnInvalidReturnCount
	}

	toolFn := callable{fn: fnValue, in: fnType.NumIn()}

	toolDesc := make(map[string]string, len(descriptors))
	for _, desc := range descriptors {
		parts := strings.Split(desc, ":")

		if len(parts) != 2 {
			return nil, ErrInvalidDescriptorFormat
		}

		role := strings.TrimSpace(parts[0])
		description := strings.TrimSpace(parts[1])

		toolDesc[role] = description
	}

	return &Tool[T]{
		name:         fnType.Name(),
		descriptions: toolDesc,
		callable:     toolFn,
	}, nil
}
