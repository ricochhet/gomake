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
	"regexp"
	"strings"
)

type FunctionBlock struct {
	Name     string
	Params   []string
	Commands []string
}

const (
	LeftBracket  = "{"
	RightBracket = "}"
	LeftParen    = "("
	RightParen   = ")"
	Comment      = "#"
	Caller       = "@"
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

func ParseText(s string) ([]FunctionBlock, error) {
	var blocks []FunctionBlock

	var currentBlock *FunctionBlock

	scanner := bufio.NewScanner(strings.NewReader(s))

	for scanner.Scan() {
		line := scanner.Text()
		line = strings.TrimSpace(line)

		if strings.HasPrefix(line, Comment) {
			continue
		}

		if strings.HasSuffix(line, LeftBracket) {
			blockNameWithParentheses := strings.TrimSpace(strings.TrimRight(line, LeftBracket))
			blockName, blockParams := parseBlockWithParams(blockNameWithParentheses)
			currentBlock = &FunctionBlock{Name: blockName, Params: blockParams, Commands: []string{}}

			continue
		}

		if line == RightBracket && currentBlock != nil {
			blocks = append(blocks, *currentBlock)
			currentBlock = nil

			continue
		}

		if currentBlock != nil {
			currentBlock.Commands = parseCallers(*currentBlock, blocks, line)
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

	parsedParams := []string{}
	for _, command := range parsedBlock.Commands {
		parsedParams = append(parsedParams, replaceArrayWithArray(command, block.Params, params))
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

func parseCallers(block FunctionBlock, blocks []FunctionBlock, line string) []string {
	var commands []string

	for _, command := range append(block.Commands, line) {
		if callerName := strings.TrimPrefix(command, Caller); strings.HasPrefix(command, Caller) {
			for _, b := range blocks {
				if b.Name == callerName {
					commands = append(commands, b.Commands...)
				}
			}
		} else {
			commands = append(commands, command)
		}
	}

	return commands
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
