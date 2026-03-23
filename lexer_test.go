package main

import (
	"testing"
)

func TestLexer(t *testing.T) {
	tests := []struct {
		input    string
		expected []token
	}{
		{
			input: `!=;; // this is a comment`,
			expected: []token{
				newToken(BANG_EQUAL, "!=", nil, 1),
				newToken(SEMICOLON, ";", nil, 1),
				newToken(SEMICOLON, ";", nil, 1),
				newToken(EOF, "", nil, 1),
			},
		},
		{
			input: `!=;; // this is a comment

// this is a new line
+<=>=<>

(( )) {}`,
			expected: []token{
				newToken(BANG_EQUAL, "!=", nil, 1),
				newToken(SEMICOLON, ";", nil, 1),
				newToken(SEMICOLON, ";", nil, 1),

				newToken(PLUS, "+", nil, 4),
				newToken(LESS_EQUAL, "<=", nil, 4),
				newToken(GREATER_EQUAL, ">=", nil, 4),
				newToken(LESS, "<", nil, 4),
				newToken(GREATER, ">", nil, 4),

				newToken(LEFT_PAREN, "(", nil, 6),
				newToken(LEFT_PAREN, "(", nil, 6),
				newToken(RIGHT_PAREN, ")", nil, 6),
				newToken(RIGHT_PAREN, ")", nil, 6),
				newToken(LEFT_BRACE, "{", nil, 6),
				newToken(RIGHT_BRACE, "}", nil, 6),

				newToken(EOF, "", nil, 6),
			},
		},
		{
			input: `"This is a string"`,
			expected: []token{
				newToken(STRING, `"This is a string"`, "This is a string", 1),
				newToken(EOF, "", nil, 1),
			},
		},
		{
			input: `12.5 + 6.9;`,
			expected: []token{
				newToken(NUMBER, "12.5", 12.5, 1),
				newToken(PLUS, "+", nil, 1),
				newToken(NUMBER, "6.9", 6.9, 1),
				newToken(SEMICOLON, ";", nil, 1),
				newToken(EOF, "", nil, 1),
			},
		},
		{
			input: `/*
func disabled() {}
*/

"hello world";`,
			expected: []token{
				newToken(STRING, `"hello world"`, "hello world", 5),
				newToken(SEMICOLON, ";", nil, 5),
				newToken(EOF, "", nil, 5),
			},
		},
	}

	for _, test := range tests {
		t.Run("Lexing", func(t *testing.T) {
			l := newLexer(test.input)
			tokens := l.scanTokens()

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
