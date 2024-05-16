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
	"bufio"
	"errors"
	"os"
	"regexp"
	"strings"
)

type Command struct {
	Directory string `json:"directory"`
	Command   string `json:"command"`
}

type FunctionBlock struct {
	Name      string    `json:"name"`
	Params    []string  `json:"params"`
	Commands  []Command `json:"commands"`
	Directory string    `json:"directory"`
}

const (
	LeftBracket  = "{"
	RightBracket = "}"
	LeftParen    = "("
	RightParen   = ")"
	Comment      = "#"
	Caller       = "@"
	Directory    = "~"
)

var (
	errTooFewArguments    = errors.New("too few params for function block")
	errEmptyFunctionBlock = errors.New("empty function block")
	errBlockNotFound      = errors.New("function block was not found")
)

func GetBlock(s string, fname string, params []string) (FunctionBlock, error) {
	blocks, err := ParseText(s)
	if err != nil {
		return FunctionBlock{}, err
	}

	block, err := getBlock(blocks, fname, params)
	if err != nil {
		return FunctionBlock{}, err
	}

	if err := checkForEmptyBlock(block); err != nil {
		return FunctionBlock{}, err
	}

	return block, nil
}

//nolint:cyclop,nestif // wontfix
func ParseText(text string) ([]FunctionBlock, error) {
	var blocks []FunctionBlock

	var currentBlock *FunctionBlock

	defaultDir, err := os.Getwd()
	if err != nil {
		return nil, err
	}

	scanner := bufio.NewScanner(strings.NewReader(text))

	for scanner.Scan() {
		line := scanner.Text()
		line = strings.TrimSpace(line)

		if strings.HasPrefix(line, Comment) {
			continue
		}

		if strings.HasSuffix(line, LeftBracket) {
			blockNameWithParentheses := strings.TrimSpace(strings.TrimRight(line, LeftBracket))
			blockName, blockParams := parseBlockWithParams(blockNameWithParentheses)
			currentBlock = &FunctionBlock{
				Name:      blockName,
				Params:    blockParams,
				Commands:  []Command{},
				Directory: defaultDir,
			}

			continue
		}

		if line == RightBracket && currentBlock != nil {
			blocks = append(blocks, *currentBlock)
			currentBlock = nil

			continue
		}

		if currentBlock != nil {
			if strings.HasPrefix(line, Directory) {
				newDir := strings.TrimSpace(strings.TrimPrefix(line, Directory))
				if newDir == "" {
					currentBlock.Directory = defaultDir
				} else {
					currentBlock.Directory = newDir
				}

				continue
			}

			commands, params := parseCallers(blocks, line)
			for _, cmd := range commands {
				if cmd.Directory == "" {
					cmd.Directory = currentBlock.Directory
				}

				currentBlock.Commands = append(currentBlock.Commands, cmd)
			}

			currentBlock.Params = append(currentBlock.Params, params...)
		}
	}

	if err := scanner.Err(); err != nil {
		return nil, err
	}

	return blocks, nil
}

func BlockParamsToCommand(block FunctionBlock, params []string) (FunctionBlock, error) {
	parsedBlock := block

	if len(parsedBlock.Params) == 0 {
		return block, nil
	}

	if len(parsedBlock.Params) > len(params) {
		return FunctionBlock{}, errTooFewArguments
	}

	parsedCommands := []Command{}

	for _, command := range parsedBlock.Commands {
		replacedCommand := replaceArrayWithArray(command.Command, block.Params, params)
		parsedCommands = append(parsedCommands, Command{Directory: command.Directory, Command: replacedCommand})
	}

	parsedBlock.Commands = parsedCommands

	return parsedBlock, nil
}

func getBlock(blocks []FunctionBlock, fn string, params []string) (FunctionBlock, error) {
	for _, block := range blocks {
		if block.Name == fn {
			parsedBlock, err := BlockParamsToCommand(block, params)
			if err != nil {
				return FunctionBlock{}, err
			}

			if err := checkForEmptyBlock(parsedBlock); err != nil {
				return FunctionBlock{}, err
			}

			return parsedBlock, nil
		}
	}

	return FunctionBlock{}, errBlockNotFound
}

func checkForEmptyBlock(block FunctionBlock) error {
	if block.Name == "" || len(block.Commands) == 0 {
		return errEmptyFunctionBlock
	}

	return nil
}

func replaceArrayWithArray(original string, oldArray []string, newArray []string) string {
	replacements := make(map[string]string)
	for i := range oldArray {
		replacements[oldArray[i]] = newArray[i]
	}

	for old, new := range replacements {
		original = strings.ReplaceAll(original, old, new)
	}

	return original
}

func parseCallers(blocks []FunctionBlock, line string) ([]Command, []string) {
	var commands []Command

	var params []string

	if strings.HasPrefix(line, Caller) {
		callerNameWithArgs, callerArgs := parseBlockWithParams(line)
		callerName := strings.TrimPrefix(callerNameWithArgs, Caller)

		for _, block := range blocks {
			if block.Name == callerName {
				var parsedBlock FunctionBlock

				parsedBlock.Name = block.Name
				parsedBlock.Directory = block.Directory

				for _, command := range block.Commands {
					if len(block.Params) != 0 && len(callerArgs) >= len(block.Params) {
						parsedCommand := replaceArrayWithArray(command.Command, block.Params, callerArgs)
						parsedBlock.Commands = append(parsedBlock.Commands, Command{
							Command:   parsedCommand,
							Directory: command.Directory,
						})
					} else {
						parsedBlock.Commands = append(parsedBlock.Commands, command)
						parsedBlock.Params = block.Params
					}

				}

				commands = append(commands, parsedBlock.Commands...)

				params = append(params, parsedBlock.Params...)
			}
		}
	} else {
		commands = append(commands, Command{Command: line}) //nolint:exhaustruct // wontfix
	}

	return commands, params
}

func parseBlockWithParams(s string) (string, []string) {
	regex := regexp.MustCompile(`^([^\(]+)(?:\(([^)]*)\))?$`)
	matches := regex.FindStringSubmatch(s)
	required := 2

	if len(matches) < required {
		return "", nil
	}

	blockName := strings.TrimSpace(matches[1])

	var params []string

	if len(matches) == 3 && matches[2] != "" {
		params = strings.Split(matches[2], ",")
		for i, param := range params {
			params[i] = strings.TrimSpace(param)
		}
	}

	return blockName, params
}
