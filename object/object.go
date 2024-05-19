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

//nolint:cyclop // wontfix
func (currentBlock *StatefulFunctionBlock) SetCallerBlock(blocks []StatefulFunctionBlock, callerName string, callerParams []string) error {
	for _, block := range blocks {
		//nolint:nestif // wontfix
		if block.Name == callerName {
			directory, err := SetBlockDirectory(block)
			os := block.Directory //nolint:varnamelen // wontfix
			expr := block.Expression
			env := block.Environment

			if err != nil {
				return err
			}

			for _, cmd := range block.Commands {
				commandText := cmd.Command
				commandExpr := cmd.Expression

				if len(callerParams) != 0 {
					commandText = SetFunctionParams(cmd.Command, block.Params, callerParams)
					commandExpr.OperandA = SetFunctionParams(cmd.Expression.OperandA, block.Params, callerParams)
					commandExpr.OperandB = SetFunctionParams(cmd.Expression.OperandB, block.Params, callerParams)
				}

				commandDirectory := cmd.Directory
				if commandDirectory == "" {
					commandDirectory = directory
				}

				commandOS := cmd.OS
				if commandOS == "" {
					commandOS = os
				}

				if commandExpr.OperandA == "" && commandExpr.OperandB == "" {
					commandExpr = expr
				}

				commandEnv := cmd.Environment
				if len(commandEnv) == 0 {
					commandEnv = env
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

func SetEnvironmentVariables(original string) string {
	variables := scanner.ScanVariables(original)

	for _, variable := range variables {
		//nolint:lll // wontfix
		original = strings.ReplaceAll(original, string(token.TokenString)+string(token.TokenLeftBracket)+variable+string(token.TokenRightBracket), os.Getenv(variable))
	}

	return original
}
