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
		Extension: ".gomake",
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
