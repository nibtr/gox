package main

import (
	"testing"
)

func TestAstPrinter(t *testing.T) {
	ex := &Binary{
		left: &Unary{
			operator: NewToken(MINUS, "-", nil, 1),
			right:    &Literal{value: 123},
		},
		operator: NewToken(STAR, "*", nil, 1),
		right:    &Grouping{expression: &Literal{45.67}},
	}

	printer := astPrinter{}
	output, err := printer.Print(ex)
	if err != nil {
		t.Errorf("%v", err)
	}
	expected := "(* (- 123) (group 45.67))"

	if output != expected {
		t.Errorf("AST printer output incorrect.\nGot:      %s\nExpected: %s", output, expected)
	}
}
