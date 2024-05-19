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

	"github.com/ricochhet/gomake/token"
)

type Command struct {
	OS        string `json:"os"`
	Directory string `json:"directory"`
	Command   string `json:"command"`
}

type FunctionBlock struct {
	Name      string    `json:"name"`
	Params    []string  `json:"params"`
	Commands  []Command `json:"commands"`
	OS        string    `json:"os"`
	Directory string    `json:"directory"`
}

func (currentBlock *FunctionBlock) SetCallerBlock(blocks []FunctionBlock, callerName string, callerParams []string) error {
	for _, block := range blocks {
		if block.Name == callerName {
			directory, err := SetBlockDirectory(block)
			os := block.Directory

			if err != nil {
				return err
			}

			for _, cmd := range block.Commands {
				commandText := cmd.Command
				if len(callerParams) != 0 {
					commandText = SetFunctionParams(cmd.Command, block.Params, callerParams)
				}

				commandDirectory := cmd.Directory
				if commandDirectory == "" {
					commandDirectory = directory
				}

				commandOS := cmd.OS
				if commandOS == "" {
					commandOS = os
				}

				currentBlock.Commands = append(currentBlock.Commands, Command{
					Command:   commandText,
					Directory: commandDirectory,
					OS:        commandOS,
				})
			}

			continue
		}
	}

	return nil
}

func SetBlockDirectory(block FunctionBlock) (string, error) {
	if block.Directory == "" {
		cwd, err := os.Getwd()
		if err != nil {
			return "", err
		}

		return cwd, nil
	}

	return block.Directory, nil
}

func SetBlockOperatingSystem(block FunctionBlock) string {
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

	return original
}
