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

var (
	ErrTooFewArgumentsInBlock  = errors.New("too few arguments in block")
	ErrOutdatedDirectorySetter = errors.New("outdated directory setter, use @(cd:./path/)")
)

func ParseBlock(block object.FunctionBlock, args []string) (object.FunctionBlock, error) {
	parsedBlock := object.FunctionBlock{
		Name:       block.Name,
		Params:     block.Params,
		Commands:   make([]object.Command, 0),
		OS:         block.OS,
		Directory:  block.Directory,
		Expression: block.Expression,
	}

	if len(block.Params) != len(args) {
		return object.FunctionBlock{}, ErrTooFewArgumentsInBlock
	}

	for _, cmd := range block.Commands {
		parsedBlock.Commands = append(parsedBlock.Commands, object.Command{
			OS:         cmd.OS,
			Command:    object.SetFunctionParams(cmd.Command, block.Params, args),
			Directory:  cmd.Directory,
			Expression: ParseExpressionResult(cmd.Expression),
		})
	}

	parsedBlock.Params = []string{}

	return parsedBlock, nil
}

//nolint:gocognit,gocyclo,cyclop,funlen // wontfix
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
				Name:       blockName,
				Params:     blockParams,
				Commands:   make([]object.Command, 0),
				OS:         "all",
				Directory:  cwd,
				Expression: object.Expression{}, //nolint:exhaustruct // wontfix
			}

			continue
		}

		if currentBlock == nil || !function {
			scanner.ReadNext()
			continue
		}

		switch scanner.CurrentRune {
		case token.TokenCaller:
			scanner.ReadNext()

			switch scanner.CurrentRune {
			case token.TokenDirectory:
				return nil, ErrOutdatedDirectorySetter
			case token.TokenLeftParen:
				//nolint:mnd // wontfix
				switch scanner.Peek(0) {
				case 'c':
					if scanner.PeekAhead(3) == "cd:" {
						scanner.ReadAhead(3)
						scanner.SkipWhitespace()
						ParseDirectory(scanner, currentBlock, cwd)
					}
				case 'o':
					if scanner.PeekAhead(3) == "os:" {
						scanner.ReadAhead(3)
						scanner.SkipWhitespace()

						if err := ParseOperatingSystem(scanner, currentBlock); err != nil {
							return nil, err
						}
					}
				case 'e':
					if scanner.PeekAhead(3) == "eq:" {
						scanner.ReadAhead(3)
						scanner.SkipWhitespace()

						ParseExpression(scanner, currentBlock, 0)
					}
				case 'n':
					if scanner.PeekAhead(4) == "neq:" {
						scanner.ReadAhead(4)
						scanner.SkipWhitespace()

						ParseExpression(scanner, currentBlock, 1)
					}
				default:
					scanner.ScanToEndOfLine()
				}

				continue
			default:
				if err := ParseCaller(scanner, currentBlock, blocks); err != nil {
					return nil, err
				}

				continue
			}
		default:
			if err := ParseCommand(scanner, currentBlock); err != nil {
				return nil, err
			}

			continue
		}
	}

	return blocks, nil
}
