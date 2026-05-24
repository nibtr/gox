package runtime

import (
	"testing"

	"github.com/nibtr/gox/ast"
	"github.com/nibtr/gox/lexer"
)

func TestAstPrinter(t *testing.T) {
	ex := &ast.Binary{
		Left: &ast.Unary{
			Operator: lexer.NewToken(lexer.MINUS, "-", nil, 1),
			Right:    &ast.Literal{Value: 123},
		},
		Operator: lexer.NewToken(lexer.STAR, "*", nil, 1),
		Right:    &ast.Grouping{Expression: &ast.Literal{Value: 45.67}},
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
