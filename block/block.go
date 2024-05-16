package block

import (
	"errors"

	"github.com/ricochhet/gomake/object"
)

var errBlockNotFound = errors.New("function block was not found")

func GetBlock(blocks []object.FunctionBlock, blockName string) (object.FunctionBlock, error) {
	for _, block := range blocks {
		if block.Name == blockName {
			return block, nil
		}
	}

	return object.FunctionBlock{}, errBlockNotFound
}
