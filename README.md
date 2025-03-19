# Toolbin

[![License](https://img.shields.io/github/license/asatraitis/toolbin)](LICENSE) [![Go Report Card](https://goreportcard.com/badge/github.com/asatraitis/toolbin)](https://goreportcard.com/report/github.com/asatraitis/toolbin)

Utility that allows easy LLM tool setup and dynamic calling by abstracting gross `reflect` type validations. Uses struct2prop to create JSON schema for tools.

## Table of Contents

- [Installation](#installation)
- [Examples](#examples)
- [Usage](#usage)

## Installation

```sh
go get github.com/asatraitis/toolbin
```

## Examples

```go
package main

import (
	"fmt"

	"github.com/asatraitis/toolbin"
)

type AddInput struct {
	X float64 `description:"numeric value for addition"`
	Y float64 `description:"numeric value for addition"`
}
type SubInput struct {
	X float64 `description:"numeric value for subtraction"`
	Y float64 `description:"numeric value for subtraction"`
}

func add(args AddInput) (string, error) {
	return fmt.Sprintf("%v", args.X+args.Y), nil
}
func subtract(args SubInput) (string, error) {
	return fmt.Sprintf("%v", args.X-args.Y), nil
}

func main() {
	myTools := toolbin.NewBin("math")

	addTool, err := toolbin.NewTool("add", "sums X and Y", add)
	if err != nil {
		panic(err)
	}
	subTool, err := toolbin.NewTool("subtract", "subtracts Y from X", subtract)
	if err != nil {
		panic(err)
	}
	myTools.SetTools([]toolbin.Tool{addTool, subTool})

	result, err := myTools.UseTool(`{
		"name": "subtract",
		"args": {
			"x":5,
			"y":3
		}
	}`)
	if err != nil {
		panic(err)
	}
    b, _ := json.MarshalIndent(result, "", "  ")
	fmt.Println(string(b))
}
```

```
{
  "content": "2"
}
```

## Usage

### Creating a Bin

Bin instance provides helpful utilities to manage tools.

Creating a bin instance:

```go
myTools  :=  toolbin.NewBin("mathTools")
```

### Creating a tool.

- Define a struct that will be used as input for your tool. It **MUST** be a **struct**! It will be inspected to create a structure that can be marshaled into a valid JSON schema. Highly encouraged to add `description` tags to your struct fields. They are used for JSON schema `description` fields.

```go
type AddArgs struct {
	X float64 `description:"numeric value for addition"`
	Y float64 `description:"numeric value for addition"`
}
```

Above struct will be used to create a struct recursively with [github.com/asatraitis/struct2prop](https://github.com/asatraitis/struct2prop):

```go
type  Prop  struct {
	Type         PropType        `json:"type"`
	Description  string          `json:"description,omitempty"`
	Items        *Prop           `json:"items,omitempty"`
	Enum         []any           `json:"enum,omitempty"`
	Properties   map[string]Prop `json:"properties,omitempty"`
	Required     []string        `json:"required,omitempty"`
}
```

- Define a function that will use your struct as input and return (string, error). It **MUST** return **(string, error)**.

```go
func add(args AddArgs) (string, error) {
	return fmt.Sprintf("%s", args.X+args.Y), nil
}
```

- Use `NewTool()` to create a tool. 3rd argument is of type `any` which **MUST** be a `func(myStruct struct{}) (string, error)`. To reduce friction when creating and calling tools in a dynamic fashion, we lose some of conveniences of typed arguments. This method will validate that provided function is correct using `reflect` standard library. It will use the functions struct argument to create `Prop` struct.

```go
mathAddTool, err := toolbin.NewTool("add", "Add X and Y and returns the sum", add)
if err != nil {
	panic(err)
}
```

- Add your tool to your bin

```go
myTools.Add(mathAddTool)
```

- Call the bin using JSON string from LLM:

```go
result, err := myTools.UseTool(`{"name": "add", "args": {"x": 5, "y": 3}}`)
if err != nil {
	print(err)
}
b, _  :=  json.MarshalIndent(result, "", " ")
fmt.Println(string(b))
```

Print:

```
{
  "content": "8"
}
```

### Tool Definitions

When utilizing function calling (tools) with LLM’s, to “register” the tools an API expects a list of functions in a JSON schema format. These definitions are created when calling `NewTool`. Under the hood it will inspect the input argument of the function and create the parameters schema.

`GetToolDefs` will return definition structs that can be marshalled or typeconv'ed for use.

```go
package main

import (
	"encoding/json"
	"fmt"

	"github.com/asatraitis/toolbin"
)

type ComplexArg struct {
	Name    string   `description:"Name of the thing"`
	Values  []string `description:"Values of the thing"`
	IsThing bool     `description:"flag for the thing"`
	Fields  struct {
		A string `description:"field a"`
		B int    `description:"field b"`
	} `description:"fields of the thing"`
	ComplexFields []struct {
		FieldA string  `description:"fieldA of the complex fields in the thing"`
		FieldB float64 `description:"fieldB of the complex fields in the thing"`
	} `description:"complex field of the thing"`
}

type SimpleArg struct {
	A string `description:"just a simple string field"`
	B string `description:"just b simple string field"`
}

func main() {
	bin := toolbin.NewBin("tools")
	complexTool, err := toolbin.NewTool("thing", "awesome complex thing", func(c ComplexArg) (string, error) {
		return "", nil
	})
	if err != nil {
		panic(err)
	}
	simpleTool, err := toolbin.NewTool("thing2", "awesome simple thing", func(s SimpleArg) (string, error) {
		return "", nil
	})
	if err != nil {
		panic(err)
	}

	bin.Add(complexTool)
	bin.Add(simpleTool)

	defs := bin.GetToolDefs()
	b, _ := json.MarshalIndent(defs, "", "  ")

	fmt.Println(string(b))
}
```

prints:

```json
[
  {
    "name": "thing",
    "description": "awesome complex thing",
    "parameters": {
      "type": "object",
      "properties": {
        "complexFields": {
          "type": "array",
          "description": "complex field of the thing",
          "items": {
            "type": "object",
            "properties": {
              "fieldA": {
                "type": "string",
                "description": "fieldA of the complex fields in the thing"
              },
              "fieldB": {
                "type": "number",
                "description": "fieldB of the complex fields in the thing"
              }
            }
          }
        },
        "fields": {
          "type": "object",
          "description": "fields of the thing",
          "properties": {
            "a": {
              "type": "string",
              "description": "field a"
            },
            "b": {
              "type": "integer",
              "description": "field b"
            }
          }
        },
        "isThing": {
          "type": "boolean",
          "description": "flag for the thing"
        },
        "name": {
          "type": "string",
          "description": "Name of the thing"
        },
        "values": {
          "type": "array",
          "description": "Values of the thing",
          "items": {
            "type": "string"
          }
        }
      }
    }
          "type": "array",
          "description": "Values of the thing",
          "items": {
            "type": "string"
          }
        }
      }
    }
          "description": "Values of the thing",
          "items": {
            "type": "string"
          }
        }
      }
    }
            "type": "string"
          }
        }
      }
    }
        }
      }
    }
  },
  {
    "name": "thing2",
    "description": "awesome simple thing",
    "parameters": {
  },
  {
    "name": "thing2",
    "description": "awesome simple thing",
    "parameters": {
    "description": "awesome simple thing",
    "parameters": {
      "type": "object",
      "type": "object",
      "properties": {
        "a": {
        "a": {
          "type": "string",
          "description": "just a simple string field"
        },
        "b": {
          "type": "string",
          "description": "just b simple string field"
        }
      }
    }
  }
]
```
