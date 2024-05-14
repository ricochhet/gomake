package main

import "errors"

var (
	errNoMakefilePath = errors.New("no makefile path specified")
	errNoFunctionName = errors.New("no function name specified")
)
