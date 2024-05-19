package parser

import (
	"errors"
	"slices"

	"github.com/ricochhet/gomake/object"
	"github.com/ricochhet/gomake/process"
	"github.com/ricochhet/gomake/scanner"
	"github.com/ricochhet/gomake/token"
)

var errUnknownPlatformIdentifier = errors.New("unknown platform identifier")

func ParseCommand(scanner *scanner.Scanner, currentBlock *object.FunctionBlock) error {
	command := scanner.ScanToEndOfLine()

	directory, err := object.SetBlockDirectory(*currentBlock)
	if err != nil {
		return err
	}

	os := object.SetBlockOperatingSystem(*currentBlock)

	currentBlock.Commands = append(currentBlock.Commands, object.Command{Command: command, OS: os, Directory: directory})

	return nil
}

func ParseCaller(scanner *scanner.Scanner, currentBlock *object.FunctionBlock, blocks []object.FunctionBlock) error {
	callerName, callerParams := scanner.ScanBlockWithParams()
	if err := currentBlock.SetCallerBlock(blocks, callerName, callerParams); err != nil {
		return err
	}

	return nil
}

func ParseDirectory(scanner *scanner.Scanner, currentBlock *object.FunctionBlock, cwd string) {
	identifier := scanner.ScanToLastOccurrence(token.TokenRightParen)
	scanner.ReadNext()

	if identifier == "" {
		currentBlock.Directory = cwd
	} else {
		currentBlock.Directory = identifier
	}
}

func ParseOperatingSystem(scanner *scanner.Scanner, currentBlock *object.FunctionBlock) error {
	identifier := scanner.ScanToLastOccurrence(token.TokenRightParen)
	scanner.ReadNext()

	if !slices.Contains(process.KnownOS, identifier) && identifier != "all" {
		return errUnknownPlatformIdentifier
	}

	if identifier == "" {
		currentBlock.OS = "all"
	} else {
		currentBlock.OS = identifier
	}

	return nil
}
