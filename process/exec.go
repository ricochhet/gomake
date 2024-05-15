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

package process

import (
	"fmt"
	"os"
	"os/exec"
	"runtime"

	"github.com/ricochhet/gomake/parser"
)

func Exec(commands []parser.Command) error {
	var shell, flag string
	if runtime.GOOS == "windows" {
		shell = "cmd"
		flag = "/C"
	} else {
		shell = "bash"
		flag = "-c"
	}

	for _, cmd := range commands {
		fmt.Printf("gomake: executing command: %s in directory: %s\n", cmd.Command, cmd.Directory)
		command := exec.Command(shell, flag, cmd.Command)
		command.Stdout = os.Stdout
		command.Stderr = os.Stderr
		command.Dir = cmd.Directory

		if err := command.Run(); err != nil {
			return err
		}
	}

	return nil
}
