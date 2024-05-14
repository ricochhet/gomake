package main

import (
	"os"

	"github.com/ricochhet/gomake/parser"
	"github.com/ricochhet/gomake/process"
)

func main() {
	if flags.Function == "" {
		process.Errr(errNoFunctionName)
	}

	if flags.Path == "" {
		flags.Path = "./Makefile"
	}

	file, err := os.ReadFile(flags.Path)
	if err != nil && flags.Path == "" {
		process.Errr(errNoMakefilePath)
		return
	}

	if err != nil {
		process.Errr(err)
		return
	}

	block, err := parser.GetBlock(string(file), flags.Function, flags.Arguments)
	if err != nil {
		process.Errr(err)
		return
	}

	if err := process.Exec(block.Commands); err != nil {
		process.Errr(err)
		return
	}
}
