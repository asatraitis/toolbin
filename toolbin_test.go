package toolbin

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func Test_NewBin(t *testing.T) {
	bin := NewBin("testBin")
	assert.Equal(t, "testBin", bin.Name)
	assert.Len(t, bin.tools, 0)
}

func Test_NewBin_Add(t *testing.T) {
	bin := NewBin("testBin")
	testTool, _ := NewTool("test", "test", func(s struct{}) (string, error) {
		return "", nil
	})
	bin.Add(testTool)
	assert.Len(t, bin.tools, 1)
}

func Test_Tools(t *testing.T) {
	bin := NewBin("testBin")
	testTool, _ := NewTool("test", "testDescription", func(s struct{}) (string, error) {
		return "result from test", nil
	})
	testTool2, _ := NewTool("test2", "test2Description", func(s struct{}) (string, error) {
		return "result from test2", nil
	})
	bin.SetTools([]Tool{testTool, testTool2})
	assert.Len(t, bin.tools, 2)

	defs := bin.GetToolDefs()
	assert.Len(t, defs, 2)

	result, err := bin.UseTool(`{"name": "test", "args":{}}`)
	assert.NoError(t, err)
	assert.Equal(t, "result from test", result.Content)

	result, err = bin.UseTool(`{"name": "test2", "args":{}}`)
	assert.NoError(t, err)
	assert.Equal(t, "result from test2", result.Content)
}
