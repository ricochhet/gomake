package parser

import (
	"github.com/ricochhet/gomake/object"
	"github.com/ricochhet/gomake/scanner"
)

func CdCaller(s *scanner.Scanner, block *object.StatefulFunctionBlock, cwd string) error {
	if call(s, "cd:") {
		if err := ParseDirectory(s, block, cwd); err != nil {
			return err
		}
	}

	return nil
}

func OsCaller(s *scanner.Scanner, block *object.StatefulFunctionBlock) error {
	if call(s, "os:") {
		if err := ParseOperatingSystem(s, block); err != nil {
			return err
		}
	}

	return nil
}

func EqCaller(s *scanner.Scanner, block *object.StatefulFunctionBlock) error {
	if call(s, "eq:") {
		if err := ParseExpression(s, block, 0); err != nil {
			return err
		}
	}

	return nil
}

func NeqCaller(s *scanner.Scanner, block *object.StatefulFunctionBlock) error {
	if call(s, "neq:") {
		if err := ParseExpression(s, block, 1); err != nil {
			return err
		}
	}

	return nil
}

func EnvCaller(s *scanner.Scanner, block *object.StatefulFunctionBlock) error {
	if call(s, "env:") {
		ParseEnvironment(s, block)
	}

	return nil
}

func call(s *scanner.Scanner, caller string) bool {
	if s.PeekAhead(len(caller)) == caller {
		s.ReadAhead(len(caller))
		s.SkipWhitespace()

		return true
	}

	return false
}
