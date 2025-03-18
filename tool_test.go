package toolbin

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func testingFn(args struct{}) (string, error) {
	return "", nil
}

func Test_NewTool_OK(t *testing.T) {
	type testStruct struct {
		TestFieldA string `description:"test field A"`
		TestFieldB string `description:"test field B"`
	}
	testTool, err := NewTool("test_tool", "tool for testing", func(args testStruct) (string, error) {
		return args.TestFieldA + args.TestFieldB, nil
	})
	assert.NoError(t, err)
	assert.Equal(t, testTool.Name, "test_tool")
	assert.Equal(t, testTool.Description, "tool for testing")
}

func Test_NewTool_FAIL_BadArgs(t *testing.T) {
	// no tool name
	_, err := NewTool("", "test", testingFn)
	assert.ErrorContains(t, err, "missing name")

	// no tool description
	_, err = NewTool("test", "", testingFn)
	assert.ErrorContains(t, err, "missing description")

	// no tool func
	_, err = NewTool("test", "description", nil)
	assert.ErrorContains(t, err, "missing toolFunc argument")

	// no tool func
	_, err = NewTool("test", "description", "whatever")
	assert.ErrorContains(t, err, "toolFunc must be of type Func")

	// tool func w/o args
	_, err = NewTool("test", "description", func() (string, error) { return "", nil })
	assert.ErrorContains(t, err, "function expects 1 argument")

	// tool func w/ wrong arg type (wants a struct)
	_, err = NewTool("test", "description", func(s []struct{}) (string, error) { return "", nil })
	assert.ErrorContains(t, err, "expects argument to be of type struct")

	// tool func w/ wrong return number (wants 2 values: (string, error))
	_, err = NewTool("test", "description", func(s struct{}) string { return "" })
	assert.ErrorContains(t, err, "expects to return 2 values")

	// tool func w/ wrong returns type (wants 2 values: (string, error))
	_, err = NewTool("test", "description", func(s struct{}) (int, error) { return 0, nil })
	assert.ErrorContains(t, err, "expects to return a string as a first value")

	// tool func w/ wrong returns type (wants 2 values: (string, error))
	_, err = NewTool("test", "description", func(s struct{}) (string, string) { return "", "" })
	assert.ErrorContains(t, err, "expects to return an error as a second value")
}
