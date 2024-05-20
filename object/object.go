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

package object

import (
	"errors"
	"os"
	"strings"

	"github.com/ricochhet/gomake/scanner"
	"github.com/ricochhet/gomake/token"
)

type Command struct {
	OS          string     `json:"os"`
	Directory   string     `json:"directory"`
	Command     string     `json:"command"`
	Expression  Expression `json:"expression"`
	Environment []string   `json:"environment"`
}

type StatefulFunctionBlock struct {
	Name        string     `json:"name"`
	Params      []string   `json:"params"`
	Commands    []Command  `json:"commands"`
	OS          string     `json:"os"`
	Directory   string     `json:"directory"`
	Expression  Expression `json:"expression"`
	Environment []string   `json:"environment"`
}

type FunctionBlock struct {
	Name     string    `json:"name"`
	Params   []string  `json:"params"`
	Commands []Command `json:"commands"`
}

type Expression struct {
	OperandA  string `json:"operandA"`
	OperandB  string `json:"operandB"`
	Operation int    `json:"operation"`
	Result    bool   `json:"result"`
}

var (
	ErrBlockNotFound       = errors.New("function block was not found or is malformed")
	ErrInvalidKeyValuePair = errors.New("invalid key=value pair")
)

func GetBlock(blocks []StatefulFunctionBlock, blockName string) (StatefulFunctionBlock, error) {
	for _, block := range blocks {
		if block.Name == blockName {
			return block, nil
		}
	}

	return StatefulFunctionBlock{}, ErrBlockNotFound
}

//nolint:cyclop // wontfix
func (currentBlock *StatefulFunctionBlock) SetCallerBlock(blocks []StatefulFunctionBlock, callerName string, callerParams []string) error {
	for _, block := range blocks {
		//nolint:nestif // wontfix
		if block.Name == callerName {
			directory, err := SetBlockDirectory(block)
			if err != nil {
				return err
			}

			for _, cmd := range block.Commands {
				commandText := cmd.Command
				commandExpr := cmd.Expression

				if len(callerParams) != 0 {
					commandText = SetFunctionParams(cmd.Command, block.Params, callerParams)
					commandExpr = SetOperandFunctionParams(cmd.Expression, block.Params, callerParams)
				}

				commandDirectory := cmd.Directory
				if commandDirectory == "" {
					commandDirectory = directory
				}

				commandOS := cmd.OS
				if commandOS == "" {
					commandOS = block.Directory
				}

				if commandExpr.OperandA == "" && commandExpr.OperandB == "" {
					commandExpr = block.Expression
				}

				commandEnv := cmd.Environment
				if len(commandEnv) == 0 {
					commandEnv = block.Environment
				}

				currentBlock.Commands = append(currentBlock.Commands, Command{
					Command:     commandText,
					Directory:   commandDirectory,
					OS:          commandOS,
					Expression:  commandExpr,
					Environment: commandEnv,
				})
			}

			continue
		}
	}

	return nil
}

func SetBlockDirectory(block StatefulFunctionBlock) (string, error) {
	if block.Directory == "" {
		cwd, err := os.Getwd()
		if err != nil {
			return "", err
		}

		return cwd, nil
	}

	return block.Directory, nil
}

func SetBlockOperatingSystem(block StatefulFunctionBlock) string {
	if block.OS == "" {
		return "all"
	}

	return block.OS
}

func SetFunctionParams(original string, oldArray []string, newArray []string) string {
	replacements := make(map[string]string)
	for i := range oldArray {
		replacements[oldArray[i]] = newArray[i]
	}

	for old, new := range replacements {
		original = strings.ReplaceAll(original, string(token.TokenLeftBracket)+old+string(token.TokenRightBracket), new)
	}

	return SetEnvironmentVariables(original)
}

func SetOperandFunctionParams(expr Expression, oldArray []string, newArray []string) Expression {
	expr.OperandA = SetFunctionParams(expr.OperandA, oldArray, newArray)
	expr.OperandB = SetFunctionParams(expr.OperandA, oldArray, newArray)

	return expr
}

func SetEnvironmentVariables(original string) string {
	variables := scanner.ScanVariables(original)

	for _, variable := range variables {
		replacement := string(token.TokenString) + string(token.TokenLeftBracket) + variable + string(token.TokenRightBracket)
		env := os.Getenv(variable)

		if env == "" {
			env = replacement
		}

		original = strings.ReplaceAll(original, replacement, env)
	}

	return original
}

func SetArrayFunctionParams(original, oldArray, newArray []string) []string {
	var replacements []string

	for _, str := range original {
		replacements = append(replacements, SetFunctionParams(str, oldArray, newArray))
	}

	return replacements
}

func SetKeyValueVariables(original string, pairs []string) (string, error) {
	variables := scanner.ScanVariables(original)

	for _, variable := range variables {
		replacement := string(token.TokenString) + string(token.TokenLeftBracket) + variable + string(token.TokenRightBracket)

		//nolint:mnd // wontfix
		for _, pair := range pairs {
			kvp := strings.SplitN(pair, "=", 2)

			if len(kvp) != 2 {
				return "", ErrInvalidKeyValuePair
			}

			if variable == kvp[0] {
				original = strings.ReplaceAll(original, replacement, kvp[1])
			}
		}
	}

	return original, nil
}
