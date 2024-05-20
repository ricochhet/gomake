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
	"flag"
	"strings"

	aflag "github.com/ricochhet/gomake/flag"
	"github.com/ricochhet/gomake/util"
)

var ErrTooFewArguments = errors.New("too few arguments for execution")

var (
	flags    *aflag.Flags = Newflag()    //nolint:gochecknoglobals // wontfix
	defaults              = aflag.Flags{ //nolint:gochecknoglobals // wontfix
		Path:      "",
		Function:  "",
		Arguments: []string{},
		Extension: ".gomake",
		Dump:      false,
	}
)

func Newflag() *aflag.Flags {
	return &defaults
}

//nolint:gochecknoinits // wontfix
func init() {
	flag.BoolVar(&flags.Dump, "dump", false, "dump parsed function block to console")
	flag.StringVar(&flags.Function, "run", "", "specify the task run")
	flag.StringVar(&flags.Path, "path", "", "specify the gomake file to use")
	args := flag.String("args", "", "specify the arguments to pass into the function block")
	flag.Parse()

	if *args != "" {
		flags.Arguments = append(flags.Arguments, strings.Split(*args, util.Separator())...)
	}
}
