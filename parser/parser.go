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

func ParseBlock(block object.StatefulFunctionBlock) object.FunctionBlock {
	return object.FunctionBlock{
		Name:     block.Name,
		Params:   block.Params,
		Commands: block.Commands,
	}
}

func ParseStatefulBlock(block object.StatefulFunctionBlock, args []string) (object.StatefulFunctionBlock, error) {
	parsedBlock := object.StatefulFunctionBlock{
		Name:        block.Name,
		Params:      block.Params,
		Commands:    make([]object.Command, 0),
		OS:          block.OS,
		Directory:   block.Directory,
		Expression:  block.Expression,
		Environment: block.Environment,
	}

	if len(block.Params) != len(args) {
		return object.StatefulFunctionBlock{}, ErrTooFewArgumentsInBlock
	}

	for _, cmd := range block.Commands {
		envParsedCmd, err := object.SetKeyValueVariables(object.SetFunctionParams(cmd.Command, block.Params, args), cmd.Environment)
		if err != nil {
			return object.StatefulFunctionBlock{}, err
		}

		parsedBlock.Commands = append(parsedBlock.Commands, object.Command{
			OS:          cmd.OS,
			Command:     envParsedCmd,
			Directory:   cmd.Directory,
			Expression:  ParseExpressionResult(cmd.Expression, block.Params, args),
			Environment: cmd.Environment,
		})
	}

	parsedBlock.Params = []string{}

	return parsedBlock, nil
}

//nolint:gocognit,gocyclo,cyclop,funlen // wontfix
func ParseText(text string) ([]object.StatefulFunctionBlock, error) {
	blocks := []object.StatefulFunctionBlock{}

	var currentBlock *object.StatefulFunctionBlock

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

			currentBlock = &object.StatefulFunctionBlock{
				Name:        blockName,
				Params:      blockParams,
				Commands:    make([]object.Command, 0),
				OS:          "all",
				Directory:   cwd,
				Expression:  object.Expression{}, //nolint:exhaustruct // wontfix
				Environment: make([]string, 0),
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

						if err := ParseDirectory(scanner, currentBlock, cwd); err != nil {
							return nil, err
						}
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

						if err := ParseExpression(scanner, currentBlock, 0); err != nil {
							return nil, err
						}
					}

					if scanner.PeekAhead(4) == "env:" {
						scanner.ReadAhead(4)
						scanner.SkipWhitespace()

						ParseEnvironment(scanner, currentBlock)
					}
				case 'n':
					if scanner.PeekAhead(4) == "neq:" {
						scanner.ReadAhead(4)
						scanner.SkipWhitespace()

						if err := ParseExpression(scanner, currentBlock, 1); err != nil {
							return nil, err
						}
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
