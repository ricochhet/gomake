package main

import (
	"errors"
	"os"

	aflag "github.com/ricochhet/gomake/flag"
)

var errTooFewArguments = errors.New("too few arguments for execution")

var (
	flags    *aflag.Flags = Newflag()    //nolint:gochecknoglobals // wontfix
	defaults              = aflag.Flags{ //nolint:gochecknoglobals // wontfix
		Path:      "",
		Function:  "",
		Arguments: []string{},
	}
)

func Newflag() *aflag.Flags {
	return &defaults
}

//nolint:gochecknoinits // wontfix
func init() {
	required := 2

	if len(os.Args) < required {
		panic(errTooFewArguments)
	}

	if len(os.Args) == required {
		flags.Function = os.Args[1]
		flags.Arguments = os.Args[2:]
	} else {
		flags.Path = os.Args[1]
		flags.Function = os.Args[2]
		flags.Arguments = os.Args[3:]
	}
}
