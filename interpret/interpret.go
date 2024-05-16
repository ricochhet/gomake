package interpret

import (
	"github.com/ricochhet/gomake/block"
	"github.com/ricochhet/gomake/object"
	"github.com/ricochhet/gomake/parser"
)

func Interpret(text string, fname string, args []string) (object.FunctionBlock, error) {
	blocks, err := parser.ParseText(text)
	if err != nil {
		return object.FunctionBlock{}, err
	}

	block, err := block.GetBlock(blocks, fname)
	if err != nil {
		return object.FunctionBlock{}, err
	}

	parsedBlock, err := parser.ParseBlock(block, args)
	if err != nil {
		return object.FunctionBlock{}, err
	}

	return parsedBlock, nil
}
