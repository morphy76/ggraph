package tool

import (
	"errors"
	"reflect"
	"strconv"
	"strings"
)

var (
	// ErrToolFnNotFunction indicates that the provided tool function is not a function.
	ErrToolFnNotFunction = errors.New("tool function must be a function")
	// ErrToolFnInvalidReturnCount indicates that the tool function does not have the correct return count.
	ErrToolFnInvalidReturnCount = errors.New("tool function must return exactly two values: (T, error)")
	// ErrInvalidDescriptorFormat indicates that the tool descriptor format is invalid.
	ErrInvalidDescriptorFormat = errors.New("invalid descriptor format (role:description expected)")
	// ErrCallingToolInvalidArgsCount indicates that the number of arguments provided to the tool function is incorrect.
	ErrCallingToolInvalidArgsCount = errors.New("invalid number of arguments provided to tool function")

	descriptions = []string{"prompt", "description", "usage"}
	requiredArgs = []string{"required", "required_args", "mandatory_args"}
	inputs       = []string{"input", "inputs", "parameters", "args"}
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

// ToolCall represents a single tool call in a conversation.
type ToolCall struct {
	// Id is the unique identifier for the tool call.
	Id string
	// ToolName is the name of the tool being called.
	ToolName string
	// Arguments are the arguments passed to the tool.
	Arguments map[string]any
}

func (t ToolCall) ArgsAsSortedSlice(tool *Tool) []any {
	args := make([]any, 0, len(tool.Args))
	for idx := range tool.Args {
		args = append(args, t.Arguments[tool.InputNameByIdx(idx)])
	}
	return args
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

	toolPrompt   string
	requiredArgs []string
	argNames     map[int]string
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
		return rvs[0].Interface(), nil
	}

	return nil, rvs[1].Interface().(error)
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
	return t.descriptionForRoles(descriptions...)
}

// BuildToolPrompt constructs a comprehensive tool prompt.
//
// It includes the tool's description, input parameters, and output parameters.
//
// Returns:
//   - string: The constructed tool prompt.
func (t *Tool) BuildToolPrompt() string {

	if t.toolPrompt != "" {
		return t.toolPrompt
	}

	desc := t.Description()
	if desc != "" {
		desc += "\n"
	}

	input := t.descriptionForRoles(inputs...)
	types := make([]string, len(t.Args))
	for i, arg := range t.Args {
		types[i] = arg.Type
	}
	typesStr := strings.Join(types, ", ")
	if input != "" {
		desc += "Input hint:[ " + input + "] of types [" + typesStr + "]\n"
	}

	t.toolPrompt = desc

	return t.toolPrompt
}

// RequiredArgs returns a list of required argument names for the tool.
//
// It looks for the "required", "required_args", or "mandatory_args" roles
// in the tool's descriptions.
//
// Returns:
//   - []string: A slice of required argument names.
func (t *Tool) RequiredArgs() []string {

	if t.requiredArgs != nil {
		return t.requiredArgs
	}

	required := t.descriptionForRoles(requiredArgs...)
	if required == "" {
		return []string{}
	}

	parts := strings.Split(required, ",")
	rv := make([]string, len(parts))
	for i, arg := range parts {
		rv[i] = strings.TrimSpace(arg)
	}

	t.requiredArgs = rv

	return t.requiredArgs
}

// InputNameByIdx returns the name of the input argument at the specified index.
//
// If no specific names are provided in the tool's descriptions, it defaults to "arg<index>".
//
// Parameters:
//   - idx: The index of the input argument.
//
// Returns:
func (t *Tool) InputNameByIdx(idx int) string {

	if t.argNames != nil && idx >= 0 && idx < len(t.argNames) {
		if val, ok := t.argNames[idx]; ok {
			return val
		}
	}

	if t.argNames == nil {
		t.argNames = make(map[int]string, len(t.Args))
	}

	inputNames := t.descriptionForRoles(inputs...)
	if inputNames == "" {
		t.argNames[idx] = "arg" + strconv.Itoa(idx)
		return t.argNames[idx]
	}

	parts := strings.Split(inputNames, ",")
	for i := range parts {
		t.argNames[i] = strings.TrimSpace(parts[i])
	}

	rv, ok := t.argNames[idx]
	if ok {
		return rv
	}
	return "arg" + strconv.Itoa(idx)
}

func (t Tool) descriptionForRoles(role ...string) string {
	for _, r := range role {
		if desc, ok := t.descriptions[r]; ok {
			return desc
		}
	}
	return ""
}
