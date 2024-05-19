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

package parser

import (
	"errors"

	"github.com/ricochhet/gomake/object"
	"github.com/ricochhet/gomake/scanner"
)

var errTooFewArgumentsInExpression = errors.New("too few arguments in expression")

func ParseExpressionResult(expr object.Expression) object.Expression {
	parsedExpr := expr

	switch expr.Operation {
	case 0:
		parsedExpr.Result = expr.OperandA == expr.OperandB
	case 1:
		parsedExpr.Result = expr.OperandA != expr.OperandB
	default:
		parsedExpr.Result = true
	}

	return parsedExpr
}

func ParseExpression(scanner *scanner.Scanner, currentBlock *object.StatefulFunctionBlock, operation int) error {
	scanner.ReadNext()
	operands := scanner.ScanParams()
	scanner.ScanToEndOfLine()

	if len(operands) != 2 { //nolint:mnd // wontfix
		return errTooFewArgumentsInExpression
	}

	currentBlock.Expression = object.Expression{
		OperandA:  object.SetEnvironmentVariables(operands[0]),
		OperandB:  object.SetEnvironmentVariables(operands[1]),
		Operation: operation,
		Result:    true,
	}

	return nil
}
