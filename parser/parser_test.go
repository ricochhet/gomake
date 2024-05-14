package parser_test

import (
	"testing"

	"github.com/ricochhet/gomake/parser"
)

func TestParse(t *testing.T) {
	t.Parallel()

	testMakefile1 := `
test() {
    echo test
}`

	if blocks, err := parser.ParseText(testMakefile1); err != nil {
		t.Fatal(err)
	} else if len(blocks) == 0 {
		t.Fatal("no function blocks when there should be")
	}

	if _, err := parser.GetBlock(testMakefile1, "test", []string{}); err != nil {
		t.Fatal()
	}

	if _, err := parser.GetBlock(testMakefile1, "null", []string{}); err == nil {
		t.Fatal()
	}
}

func TestParseWithParams(t *testing.T) {
	t.Parallel()

	testMakefile2 := `
test({param}) {
	echo {param}
}`

	if blocks, err := parser.ParseText(testMakefile2); err != nil {
		t.Fatal(err)
	} else if len(blocks) == 0 {
		t.Fatal("no function blocks when there should be")
	}

	if _, err := parser.GetBlock(testMakefile2, "test", []string{"test"}); err != nil {
		t.Fatal(err)
	}

	if _, err := parser.GetBlock(testMakefile2, "test", []string{}); err == nil {
		t.Fatal("function block requires params but none exist")
	}
}

func TestBlockParamsToCommand(t *testing.T) {
	t.Parallel()

	block := parser.FunctionBlock{
		Name:     "test",
		Commands: []string{"echo {param}"},
		Params:   []string{"test"},
	}

	if _, err := parser.BlockParamsToCommand(block, []string{"test"}); err != nil {
		t.Fatal(err)
	}

	if _, err := parser.BlockParamsToCommand(block, []string{}); err == nil {
		t.Fatal("function block requires params but none exist")
	}
}
