package object

import (
	"os"
	"strings"

	"github.com/ricochhet/gomake/token"
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

func (currentBlock *FunctionBlock) SetCallerBlock(blocks []FunctionBlock, callerName string, callerParams []string) error {
	for _, block := range blocks {
		if block.Name == callerName {
			directory, err := SetBlockDirectory(block)
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

				currentBlock.Commands = append(currentBlock.Commands, Command{
					Command:   commandText,
					Directory: commandDirectory,
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
