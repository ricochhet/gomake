package block

import (
	"errors"

	"github.com/ricochhet/gomake/object"
)

var (
	errBlockNotFound   = errors.New("function block was not found")
	errTooFewArguments = errors.New("too few params for function block")
)

func GetBlock(blocks []object.FunctionBlock, blockName string) (object.FunctionBlock, error) {
	for _, block := range blocks {
		if block.Name == blockName {
			return block, nil
		}
	}

	return object.FunctionBlock{}, errBlockNotFound
}

func BlockParamsToCommand(block object.FunctionBlock, params []string) (object.FunctionBlock, error) {
	parsedBlock := block

	if len(parsedBlock.Params) == 0 {
		return block, nil
	}

	if len(parsedBlock.Params) > len(params) {
		return object.FunctionBlock{}, errTooFewArguments
	}

	parsedCommands := []object.Command{}

	for _, command := range parsedBlock.Commands {
		replacedCommand := object.SetFunctionParams(command.Command, block.Params, params)
		parsedCommands = append(parsedCommands, object.Command{Directory: command.Directory, Command: replacedCommand})
	}

	parsedBlock.Commands = parsedCommands

	return parsedBlock, nil
}
