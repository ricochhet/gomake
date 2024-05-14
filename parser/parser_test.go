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
