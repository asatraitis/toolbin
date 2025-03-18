package main

import (
	"encoding/json"
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
