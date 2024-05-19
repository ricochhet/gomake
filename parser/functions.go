package parser

import (
	"errors"
	"slices"

	"github.com/ricochhet/gomake/object"
	"github.com/ricochhet/gomake/process"
	"github.com/ricochhet/gomake/scanner"
)

var (
	errUnknownPlatformIdentifier = errors.New("unknown platform identifier")
	errUnknownParameterInCaller  = errors.New("unknown parameter in caller")
)

func ParseCommand(scanner *scanner.Scanner, currentBlock *object.StatefulFunctionBlock) error {
	command := scanner.ScanToEndOfLine()

	directory, err := object.SetBlockDirectory(*currentBlock)
	if err != nil {
		return err
	}

	os := object.SetBlockOperatingSystem(*currentBlock)

	currentBlock.Commands = append(currentBlock.Commands, object.Command{
		Command:    command,
		OS:         os,
		Directory:  directory,
		Expression: currentBlock.Expression,
	})

	return nil
}

func ParseCaller(scanner *scanner.Scanner, currentBlock *object.StatefulFunctionBlock, blocks []object.StatefulFunctionBlock) error {
	callerName, callerParams := scanner.ScanBlockWithParams()
	if err := currentBlock.SetCallerBlock(blocks, callerName, callerParams); err != nil {
		return err
	}

	return nil
}

func ParseDirectory(scanner *scanner.Scanner, currentBlock *object.StatefulFunctionBlock, cwd string) error {
	scanner.ReadNext()
	identifier := scanner.ScanParams()

	if len(identifier) != 1 {
		return errUnknownParameterInCaller
	}

	if identifier[0] == "" {
		currentBlock.Directory = cwd
	} else {
		currentBlock.Directory = identifier[0]
	}

	return nil
}

func ParseOperatingSystem(scanner *scanner.Scanner, currentBlock *object.StatefulFunctionBlock) error {
	scanner.ReadNext()
	identifier := scanner.ScanParams()

	if len(identifier) != 1 {
		return errUnknownParameterInCaller
	}

	if !slices.Contains(process.KnownOS, identifier[0]) && identifier[0] != "all" {
		return errUnknownPlatformIdentifier
	}

	if identifier[0] == "" {
		currentBlock.OS = "all"
	} else {
		currentBlock.OS = identifier[0]
	}

	return nil
}
