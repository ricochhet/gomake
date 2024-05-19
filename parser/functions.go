/*
 * gomake
 * Copyright (C) 2024 gomake contributors
 *
 * This program is free software: you can redistribute it and/or modify
 * it under the terms of the GNU Affero General Public License as published
 * by the Free Software Foundation, either version 3 of the License, or
 * (at your option) any later version.
 *
 * This program is distributed in the hope that it will be useful,
 * but WITHOUT ANY WARRANTY; without even the implied warranty of
 * MERCHANTABILITY or FITNESS FOR A PARTICULAR PURPOSE.  See the
 * GNU Affero General Public License for more details.

 * You should have received a copy of the GNU Affero General Public License
 * along with this program.  If not, see <https://www.gnu.org/licenses/>.
 */

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
		Command:     command,
		OS:          os,
		Directory:   directory,
		Expression:  currentBlock.Expression,
		Environment: currentBlock.Environment,
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

	currentBlock.Directory = object.SetEnvironmentVariables(currentBlock.Directory)

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

func ParseEnvironment(scanner *scanner.Scanner, currentBlock *object.StatefulFunctionBlock) {
	scanner.ReadNext()
	variables := scanner.ScanParams()
	scanner.ScanToEndOfLine()

	currentBlock.Environment = append(currentBlock.Environment, variables...)
}
