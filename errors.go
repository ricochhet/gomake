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
	"errors"
	"fmt"
)

var (
	ErrNoMakefilePath  = errors.New("no makefile path specified")
	ErrNoFunctionName  = errors.New("no function name specified")
	ErrInvalidFileType = errors.New("invalid file type, expects: .gomake")
)

func Errr(err error) {
	fmt.Println("gomake: " + err.Error())
}
