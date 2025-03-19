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
