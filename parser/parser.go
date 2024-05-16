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
	"os"

	"github.com/ricochhet/gomake/object"
	"github.com/ricochhet/gomake/scanner"
	"github.com/ricochhet/gomake/token"
)

var ErrTooFewArgumentsInBlock = errors.New("too few arguments in block")

func ParseBlock(block object.FunctionBlock, args []string) (object.FunctionBlock, error) {
	parsedBlock := object.FunctionBlock{
		Name:      block.Name,
		Params:    block.Params,
		Commands:  make([]object.Command, 0),
		Directory: block.Directory,
	}

	if len(block.Params) != len(args) {
		return object.FunctionBlock{}, ErrTooFewArgumentsInBlock
	}

	for _, cmd := range block.Commands {
		parsedBlock.Commands = append(parsedBlock.Commands, object.Command{
			Command:   object.SetFunctionParams(cmd.Command, block.Params, args),
			Directory: cmd.Directory,
		})
	}

	parsedBlock.Params = []string{}

	return parsedBlock, nil
}

//nolint:gocognit,gocyclo,cyclop // wontfix
func ParseText(text string) ([]object.FunctionBlock, error) {
	blocks := []object.FunctionBlock{}

	var currentBlock *object.FunctionBlock

	cwd, err := os.Getwd()
	if err != nil {
		return nil, err
	}

	scanner := scanner.NewScanner(text)
	function := false

	for scanner.CurrentRune != 0 {
		scanner.SkipWhitespace()

		if scanner.CurrentRune == token.TokenComment {
			for scanner.CurrentRune != token.TokenNewLine && scanner.CurrentRune != 0 {
				scanner.ReadNext()
			}

			continue
		}

		if scanner.CurrentRune == token.TokenLeftBracket {
			scanner.ReadNext()

			function = true

			continue
		}

		if scanner.CurrentRune == token.TokenRightBracket {
			scanner.ReadNext()

			if currentBlock != nil {
				blocks = append(blocks, *currentBlock)
				currentBlock = nil
			}

			function = false

			continue
		}

		if scanner.IsIndentifiable(scanner.CurrentRune) && !function {
			blockName, blockParams := scanner.ScanBlockWithParams()

			currentBlock = &object.FunctionBlock{
				Name:      blockName,
				Params:    blockParams,
				Commands:  make([]object.Command, 0),
				Directory: cwd,
			}

			continue
		}

		if currentBlock == nil || !function {
			scanner.ReadNext()
			continue
		}

		switch scanner.CurrentRune {
		case token.TokenDirectory:
			scanner.ReadNext()

			identifier := scanner.ScanToEndOfLine()

			if identifier == "" {
				currentBlock.Directory = cwd
			} else {
				currentBlock.Directory = identifier
			}

			continue
		case token.TokenCaller:
			scanner.ReadNext()

			callerName, callerParams := scanner.ScanBlockWithParams()
			if err := currentBlock.SetCallerBlock(blocks, callerName, callerParams); err != nil {
				return nil, err
			}

			continue
		default:
			command := scanner.ScanToEndOfLine()

			if directory, err := object.SetBlockDirectory(*currentBlock); err == nil {
				currentBlock.Commands = append(currentBlock.Commands, object.Command{Command: command, Directory: directory})
			} else {
				return nil, err
			}

			continue
		}
	}

	return blocks, nil
}
