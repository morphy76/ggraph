package tool

import (
	"reflect"
	"runtime"
	"strings"
)

// ExecFn represents a generic function type that can be used as a tool.
//
// An ExecFn has the signature: func(args...) (T, error)
//
// Type parameters:
//   - T: The return type of the function as specified in the tool creation.
//
// Parameters:
//   - args: A variable number of arguments of any type.
//
// Returns:
//   - T: The result of the function execution.
//   - error: An error if the function fails during execution.
type ExecFn any

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
//   - *Tool: A pointer to the created Tool instance.
//   - error: An error if the function is not valid or descriptors are incorrectly formatted.
//
// Example:
//
//	myTool, err := tool.CreateTool[string](myFunc, "Prompt: This tool does X", "Usage: Call with Y")
func CreateTool[T any](fn ExecFn, descriptors ...string) (*Tool, error) {
	// TODO fn is a builder which accepts builder args (pooled http clients, dbconnections,...) and returns a ExecFn

	toolDesc := make(map[string]string, len(descriptors))
	for _, desc := range descriptors {
		parts := strings.Split(desc, ":")

		if len(parts) < 2 {
			return nil, ErrInvalidDescriptorFormat
		}

		role := strings.ToLower(strings.TrimSpace(parts[0]))
		value := strings.TrimSpace(strings.Join(parts[1:], ":"))

		toolDesc[role] = value
	}

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

	var t T
	toolName := extractToolName(fn, reflect.TypeOf(t))

	rv := &Tool{
		Name:         toolName,
		descriptions: toolDesc,
		callable:     toolFn,
	}

	args := make([]Arg, fnType.NumIn())
	for i := 0; i < fnType.NumIn(); i++ {
		argType := fnType.In(i)
		args[i] = Arg{
			Name: rv.InputNameByIdx(i),
			Type: argType.String(),
		}
	}

	rv.Args = args

	return rv, nil
}

func extractToolName(fn any, genericType reflect.Type) string {
	fnValue := reflect.ValueOf(fn)

	fullName := runtime.FuncForPC(fnValue.Pointer()).Name()

	lastSlash := strings.LastIndex(fullName, "/")
	if lastSlash != -1 {
		fullName = fullName[lastSlash+1:]
	}

	baseName := fullName
	if dotIdx := strings.Index(fullName, "."); dotIdx != -1 {
		baseName = fullName[dotIdx+1:]
	}

	isGeneric := strings.HasSuffix(baseName, "[...]")

	if isGeneric {
		baseName = strings.TrimSuffix(baseName, "[...]")

		if genericType != nil {
			typeName := genericType.String()
			if idx := strings.LastIndex(typeName, "."); idx != -1 {
				typeName = typeName[idx+1:]
			}
			return baseName + "_" + typeName
		}
	}

	return baseName
}
