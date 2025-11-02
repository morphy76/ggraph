package tool

import (
	"fmt"
	"reflect"
)

var (
	// ErrToolFnNotFunction indicates that the provided tool function is not a function.
	ErrToolFnNotFunction = fmt.Errorf("tool function must be a function")
	// ErrToolFnInvalidReturnCount indicates that the tool function does not have the correct return count.
	ErrToolFnInvalidReturnCount = fmt.Errorf("tool function must return exactly two values: (T, error)")
	// ErrInvalidDescriptorFormat indicates that the tool descriptor format is invalid.
	ErrInvalidDescriptorFormat = fmt.Errorf("invalid descriptor format (role:description expected)")
	// ErrCallingToolInvalidArgsCount indicates that the number of arguments provided to the tool function is incorrect.
	ErrCallingToolInvalidArgsCount = fmt.Errorf("invalid number of arguments provided to tool function")
)

type callable struct {
	fn reflect.Value
	in int
}

// Arg represents a single argument for a Tool.
type Arg struct {
	// Name is the name of the argument.
	Name string
	// Type is the type of the argument.
	Type string
}

// Tool represents a callable tool with metadata.
//
// The tool function must have the signature: func(args...) (T, error)
type Tool struct {
	// Name is the name of the tool.
	Name string
	// Args is the list of arguments for the tool.
	Args         []Arg
	descriptions map[string]string
	callable     callable
}

// Call invokes the tool function with the provided arguments.
//
// It returns the result of type T and an error if any occurred during the function call.
// An error is returned if the number of arguments does not match the expected count.
//
// Parameters:
//   - args: The arguments to pass to the tool function.
//
// Returns:
//   - T: The result of the tool function.
//   - error: An error if the call failed or if the argument count is incorrect.
//
// Example:
//
//	result, err := myTool.Call(arg1, arg2)
func (t Tool) Call(args ...any) (any, error) {
	if len(args) != t.callable.in {
		return nil, ErrCallingToolInvalidArgsCount
	}

	in := make([]reflect.Value, len(args))
	for i, arg := range args {
		in[i] = reflect.ValueOf(arg)
	}

	rvs := t.callable.fn.Call(in)

	if rvs[1].IsNil() {
		return rvs[0].Interface().(any), nil
	} else {
		return nil, rvs[1].Interface().(error)
	}
}

// Description returns the tool's description.
//
// It looks for common description roles in the following order:
//   - "prompt"
//   - "description"
//   - "usage"
//
// If none of these roles are found, it returns an empty string.
//
// Returns:
//   - string: The tool's description.
func (t Tool) Description() string {
	return t.descriptionForRoles("prompt", "description", "usage")
}

func (t Tool) descriptionForRoles(role ...string) string {
	for _, r := range role {
		if desc, ok := t.descriptions[r]; ok {
			return desc
		}
	}
	return ""
}
