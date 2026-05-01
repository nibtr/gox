package main

import (
	"testing"
)

func TestAstPrinter(t *testing.T) {
	ex := &binary{
		left: &unary{
			operator: newToken(MINUS, "-", nil, 1),
			right:    &literal{value: 123},
		},
		operator: newToken(STAR, "*", nil, 1),
		right:    &grouping{expression: &literal{45.67}},
	}

	printer := astPrinter{}
	output, err := printer.Print(ex)
	if err != nil {
		t.Errorf("%w", err)
	}
	expected := "(* (- 123) (group 45.67))"

	if output != expected {
		t.Errorf("AST printer output incorrect.\nGot:      %s\nExpected: %s", output, expected)
	}
}
