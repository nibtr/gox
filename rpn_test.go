package main

import (
	"testing"
)

func TestRPN(t *testing.T) {
	// (1 + -2) * (4 - 3) -> 1 2 + 4 3 - *
	ex := &binary{
		left: &grouping{
			expression: &binary{
				left:     &literal{value: 1},
				operator: newToken(PLUS, "+", nil, 1),
				right:    &unary{operator: newToken(MINUS, "-", nil, 1), right: &literal{value: 2}},
			},
		},
		operator: newToken(STAR, "*", nil, 1),
		right: &grouping{
			expression: &binary{
				left:     &literal{value: 4},
				operator: newToken(MINUS, "-", nil, 1),
				right:    &literal{value: 3},
			},
		}}

	r := rpn{}
	output := r.convert(ex)
	expected := "1 2 - + 4 3 - *"

	if output != expected {
		t.Errorf("Reverse Polish Notation (RPN) output incorrect.\nGot:      %s\nExpected: %s", output, expected)
	}
}
