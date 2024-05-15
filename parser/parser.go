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
	Directory string
	Command   string
}

type FunctionBlock struct {
	Name      string
	Params    []string
	Commands  []Command
	Directory string
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

//nolint:cyclop // wontfix
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
				currentBlock.Commands = append(currentBlock.Commands, Command{Directory: currentBlock.Directory, Command: cmd})
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

	parsedParams := []Command{}

	for _, command := range parsedBlock.Commands {
		replacedCommand := replaceArrayWithArray(command.Command, block.Params, params)
		parsedParams = append(parsedParams, Command{Directory: command.Directory, Command: replacedCommand})
	}

	parsedBlock.Commands = parsedParams

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
	for i, oldString := range oldArray {
		original = strings.ReplaceAll(original, oldString, newArray[i])
	}

	return original
}

func parseCallers(blocks []FunctionBlock, line string) ([]string, []string) {
	var commands []string

	var params []string

	if strings.HasPrefix(line, Caller) {
		callerName := strings.TrimPrefix(line, Caller)
		for _, block := range blocks {
			if block.Name == callerName {
				for _, cmd := range block.Commands {
					commands = append(commands, cmd.Command)
				}

				params = append(params, block.Params...)
			}
		}
	} else {
		commands = append(commands, line)
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
