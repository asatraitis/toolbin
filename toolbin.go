package toolbin

import (
	"encoding/json"
	"errors"
	"reflect"
	"sync"
)

type ToolBin struct {
	Name string

	mu    sync.RWMutex
	tools map[string]Tool
}

func NewBin(name string) *ToolBin {
	return &ToolBin{
		Name:  name,
		tools: make(map[string]Tool),
	}
}

func (tb *ToolBin) Add(tool Tool) {
	tb.mu.Lock()
	defer tb.mu.Unlock()
	tb.tools[tool.Name] = tool
}
func (tb *ToolBin) SetTools(tools []Tool) {
	if len(tools) < 1 {
		return
	}
	tb.mu.Lock()
	defer tb.mu.Unlock()
	for _, tool := range tools {
		tb.tools[tool.Name] = tool
	}
}
func (tb *ToolBin) GetToolDefs() []ToolDef {
	toolDefs := []ToolDef{}
	if len(tb.tools) < 1 {
		return toolDefs
	}
	tb.mu.RLock()
	defer tb.mu.RUnlock()
	for _, tool := range tb.tools {
		toolDefs = append(toolDefs, tool.ToolDef)
	}
	return toolDefs
}
func (tb *ToolBin) UseTool(jsonString string) (*ToolResponse, error) {
	if jsonString == "" {
		return nil, errors.New("missing input")
	}
	var req ToolRequest[json.RawMessage]
	err := json.Unmarshal([]byte(jsonString), &req)
	if err != nil {
		return nil, err
	}
	// add concurrently safety with mutex
	tb.mu.RLock()
	defer tb.mu.RUnlock()

	// get the tool
	tool, ok := tb.tools[req.Name]
	if !ok {
		return nil, errors.New("tool not found")
	}
	if tool.Exec == nil {
		return nil, errors.New("tool does not have a defined function")
	}

	// Obtain the function's reflect.Value and its expected argument type.
	fnVal := reflect.ValueOf(tool.Exec)
	if fnVal.Type().Kind() != reflect.Func {
		return nil, errors.New("tool exec is not a function")
	}
	argType := fnVal.Type().In(0)

	// Create a new instance (pointer) of the argument type.
	argPtr := reflect.New(argType)

	// Unmarshal the raw JSON arguments into the newly created argument instance.
	if err := json.Unmarshal(req.Args, argPtr.Interface()); err != nil {
		return nil, err
	}

	results := fnVal.Call([]reflect.Value{argPtr.Elem()})
	resultStr := results[0].Interface().(string)
	errInterface := results[1].Interface()
	var errResult error
	if errInterface != nil {
		errResult = errInterface.(error)
	}

	return &ToolResponse{Content: resultStr}, errResult
}
