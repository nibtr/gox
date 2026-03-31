package main

type parser struct {
	tokens []token
	// `current` is the pointer to the next token to be parsed
	current uint32
}

func newParser(tokens []token) *parser {
	return &parser{
		tokens:  tokens,
		current: 0,
	}
}
