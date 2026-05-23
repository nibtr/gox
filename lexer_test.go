package main

import (
	"testing"
)

func TestLexer(t *testing.T) {
	tests := []struct {
		input    string
		expected []Token
	}{
		{
			input: `!=;; // this is a comment`,
			expected: []Token{
				NewToken(BANG_EQUAL, "!=", nil, 1),
				NewToken(SEMICOLON, ";", nil, 1),
				NewToken(SEMICOLON, ";", nil, 1),
				NewToken(EOF, "", nil, 1),
			},
		},
		{
			input: `!=;; // this is a comment

// this is a new line
+<=>=<>

(( )) {}`,
			expected: []Token{
				NewToken(BANG_EQUAL, "!=", nil, 1),
				NewToken(SEMICOLON, ";", nil, 1),
				NewToken(SEMICOLON, ";", nil, 1),

				NewToken(PLUS, "+", nil, 4),
				NewToken(LESS_EQUAL, "<=", nil, 4),
				NewToken(GREATER_EQUAL, ">=", nil, 4),
				NewToken(LESS, "<", nil, 4),
				NewToken(GREATER, ">", nil, 4),

				NewToken(LEFT_PAREN, "(", nil, 6),
				NewToken(LEFT_PAREN, "(", nil, 6),
				NewToken(RIGHT_PAREN, ")", nil, 6),
				NewToken(RIGHT_PAREN, ")", nil, 6),
				NewToken(LEFT_BRACE, "{", nil, 6),
				NewToken(RIGHT_BRACE, "}", nil, 6),

				NewToken(EOF, "", nil, 6),
			},
		},
		{
			input: `"This is a string"`,
			expected: []Token{
				NewToken(STRING, `"This is a string"`, "This is a string", 1),
				NewToken(EOF, "", nil, 1),
			},
		},
		{
			input: `12.5 + 6.9;`,
			expected: []Token{
				NewToken(NUMBER, "12.5", 12.5, 1),
				NewToken(PLUS, "+", nil, 1),
				NewToken(NUMBER, "6.9", 6.9, 1),
				NewToken(SEMICOLON, ";", nil, 1),
				NewToken(EOF, "", nil, 1),
			},
		},
		{
			input: `/*
func disabled() {}
*/

"hello world";`,
			expected: []Token{
				NewToken(STRING, `"hello world"`, "hello world", 5),
				NewToken(SEMICOLON, ";", nil, 5),
				NewToken(EOF, "", nil, 5),
			},
		},
	}

	for _, test := range tests {
		t.Run("Lexing", func(t *testing.T) {
			l := NewLexer(test.input)
			tokens, err := l.scanTokens()
			if err != nil {
				t.Errorf("%v", err)
			}

			// compare tokens
			for i, token := range tokens {
				// check type
				if token.tokenType != test.expected[i].tokenType {
					t.Errorf("Token type mismatch at index %d: expected %v, got %v", i, test.expected[i].tokenType, token.tokenType)
				}

				// check lexeme
				if token.lexeme != test.expected[i].lexeme {
					t.Errorf("Token lexeme mismatch at index %d: expected %v, got %v", i, test.expected[i].lexeme, token.lexeme)
				}

				// check literal
				if token.literal != test.expected[i].literal {
					t.Errorf("Token literal mismatch at index %d: expected %v, got %v", i, test.expected[i].literal, token.literal)
				}

				// check line
				if token.line != test.expected[i].line {
					t.Errorf("Token line mismatch at index %d: expected %v, got %v", i, test.expected[i].line, token.line)
				}
			}
		})
	}
}
