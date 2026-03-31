package main

type parser struct {
	// `tokens` is the list of tokens from the lexer
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

func (p *parser) expression() expr {
	return p.equality()
}

func (p *parser) equality() expr {
	expr := p.comparison()

	for p.match(BANG_EQUAL, EQUAL_EQUAL) {
		operator := p.previous()
		right := p.comparison()
		expr = &binary{
			left:     expr,
			operator: operator,
			right:    right,
		}
	}

	return expr
}

// ============== Helpers ====================

func (p *parser) match(types ...tokenType) bool {
	for _, tokenType := range types {
		if p.check(tokenType) {
			p.advance()
			return true
		}
	}

	return false
}

func (p *parser) check(t tokenType) bool {
	if p.isAtEnd() {
		return false
	}

	return p.peek().tokenType == t
}

func (p *parser) advance() token {
	if !p.isAtEnd() {
		p.current++
	}
	return p.previous()
}

func (p *parser) isAtEnd() bool {
	return p.peek().tokenType == EOF
}

func (p *parser) peek() token {
	return p.tokens[p.current]
}

func (p *parser) previous() token {
	return p.tokens[p.current-1]
}
