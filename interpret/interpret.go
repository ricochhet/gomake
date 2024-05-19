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

package interpret

import (
	"github.com/ricochhet/gomake/block"
	"github.com/ricochhet/gomake/object"
	"github.com/ricochhet/gomake/parser"
)

func Interpret(text string, fname string, args []string) (object.FunctionBlock, error) {
	blocks, err := parser.ParseText(text)
	if err != nil {
		return object.FunctionBlock{}, err
	}

	block, err := block.GetBlock(blocks, fname)
	if err != nil {
		return object.FunctionBlock{}, err
	}

	parsedStatefulBlock, err := parser.ParseStatefulBlock(block, args)
	if err != nil {
		return object.FunctionBlock{}, err
	}

	return parser.ParseBlock(parsedStatefulBlock), nil
}
