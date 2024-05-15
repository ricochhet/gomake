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

package main

import (
	"fmt"
	"os"
	"path/filepath"

	"github.com/ricochhet/gomake/parser"
	"github.com/ricochhet/gomake/process"
)

func main() {
	if flags.Function == "" {
		Errr(errNoFunctionName)
	}

	if flags.Path == "" {
		flags.Path = "./make.gomake"
	}

	fmt.Println()

	if filepath.Ext(flags.Path) != flags.Extension {
		Errr(errInvalidFileType)
		return
	}

	file, err := os.ReadFile(flags.Path)
	if err != nil && flags.Path == "" {
		Errr(errNoMakefilePath)
		return
	}

	if err != nil {
		Errr(err)
		return
	}

	block, err := parser.GetBlock(string(file), flags.Function, flags.Arguments)
	if err != nil {
		Errr(err)
		return
	}

	if err := process.Exec(block.Commands); err != nil {
		Errr(err)
		return
	}
}
