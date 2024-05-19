package parser

import (
	"github.com/ricochhet/gomake/object"
	"github.com/ricochhet/gomake/scanner"
)

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

func ParseExpression(scanner *scanner.Scanner, currentBlock *object.StatefulFunctionBlock, operation int) {
	scanner.ReadNext()
	operands := scanner.ScanParams()
	scanner.ScanToEndOfLine()

	if len(operands) != 2 { //nolint:mnd // wontfix
		return
	}

	currentBlock.Expression = object.Expression{
		OperandA:  operands[0],
		OperandB:  operands[1],
		Operation: operation,
		Result:    true,
	}
}
