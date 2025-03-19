package toolbin

import (
	"errors"
	"reflect"

	"github.com/asatraitis/struct2prop"
)

type ToolRequest[T any] struct {
	Name string `json:"name"`
	Args T
}
type ToolResponse struct {
	Content string `json:"content"`
}
type ToolDef struct {
	// Tool name; no spaces
	Name string `json:"name"`

	// Tool description
	Description string `json:"description"`

	// Schema params created from a struct
	Parameters *struct2prop.Prop `json:"parameters,omitempty"`
}
type Tool struct {
	ToolDef
	Exec any
}

func NewTool(name, description string, toolFunc any) (Tool, error) {
	var tool = Tool{
		ToolDef: ToolDef{},
	}
	if name == "" {
		return tool, errors.New("missing name")
	}
	if description == "" {
		return tool, errors.New("missing description")
	}

	if toolFunc == nil {
		return tool, errors.New("missing toolFunc argument")
	}

	fnType := reflect.TypeOf(toolFunc)
	if fnType.Kind() != reflect.Func {
		return tool, errors.New("toolFunc must be of type Func")
	}

	// Ensure the function has at least one parameter
	if fnType.NumIn() == 0 {
		return tool, errors.New("toolFunc function expects 1 argument of type struct")
	}
	argType := fnType.In(0)

	// Check if the argument is a struct
	if argType.Kind() != reflect.Struct {
		return tool, errors.New("toolFunc function expects argument to be of type struct")
	}

	// Validate number of return values
	if fnType.NumOut() != 2 {
		return tool, errors.New("toolFunc function expects to return 2 values: (string, error)")
	}
	// Validate first return type (should be string)
	if fnType.Out(0).Kind() != reflect.String {
		return tool, errors.New("toolFunc function expects to return a string as a first value: (string, error)")
	}
	// Validate second return type (should be error)
	errorType := reflect.TypeOf((*error)(nil)).Elem()
	if !fnType.Out(1).Implements(errorType) {
		return tool, errors.New("toolFunc function expects to return an error as a second value: (string, error)")
	}

	// Create a new instance of the struct
	argStructInstance := reflect.New(argType).Elem()

	// create param schema for input struct
	params, err := struct2prop.GetProperties(argStructInstance.Interface())
	if err != nil {
		return tool, err
	}

	tool.Name = name
	tool.Description = description
	tool.Parameters = params
	tool.Exec = toolFunc

	return tool, nil
}
