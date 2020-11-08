package main

import (
	"errors"
	"fmt"
	"log"
)

type Parser struct {
	tokens  []Token
	current int
}

func NewParser(ts []Token) Parser {
	return Parser{tokens: ts}
}

// func (p *Parser) PrintASTbetter() {
// 	e, err := p.Parse()
// 	if err != nil {
// 		log.Fatal("parse failed: ", err)
// 	}
// }

func (p *Parser) PrintAST() {
	e, err := p.Parse()
	if err != nil {
		log.Fatal("parse failed: ", err)
	}

	fmt.Println(e.print(0))
}

func (p *Parser) Parse() (Expr, error) {
	return p.expression()
}

func (p *Parser) expression() (Expr, error) {
	return p.term()
}

func (p *Parser) term() (Expr, error) {
	expr, e1 := p.factor()
	if e1 != nil {
		return expr, e1
	}

	for p.match(MINUS, PLUS) {
		operator := p.previous()
		right, e2 := p.factor()
		if e2 != nil {
			return expr, e2
		}

		expr = Binary{expr, operator, right}
	}

	return expr, nil
}

func (p *Parser) factor() (Expr, error) {
	expr, e := p.unary()
	if e != nil {
		return expr, e
	}

	for p.match(SLASH, STAR) {
		operator := p.previous()
		right, e := p.unary()
		if e != nil {
			return expr, e
		}

		expr = Binary{expr, operator, right}
	}

	return expr, nil
}

func (p *Parser) unary() (Expr, error) {
	for p.match(MINUS) {
		operator := p.previous()
		right, e := p.unary()
		if e != nil {
			return right, e
		}

		return Unary{operator, right}, nil
	}

	return p.primary()
}

func (p *Parser) primary() (Expr, error) {
	if p.match(NUMBER) {
		return Literal{p.previous().literal}, nil
	}

	if p.match(LEFT_PAREN) {
		expr, e := p.expression()
		if e != nil {
			return nil, e
		}

		_, e = p.consume(RIGHT_PAREN, "Expect ')' after expression")
		if e != nil {
			return nil, e
		}

		return Grouping{expr}, nil
	}

	return nil, errors.New(fmt.Sprintf("%v, Expect expression", p.peek().lexeme))
}

/*
---------------------- helpers -----------------------
*/

func (p *Parser) consume(t TokenType, msg string) (Token, error) {
	if p.check(t) {
		return p.advance(), nil
	}

	return Token{}, errors.New(fmt.Sprintf("error: %s, %s", p.peek().lexeme, msg))
}

func (p *Parser) peek() Token {
	return p.tokens[p.current]
}

func (p *Parser) isAtEnd() bool {
	return p.peek().ttype == EOF
}

func (p *Parser) previous() Token {
	return p.tokens[p.current-1]
}

func (p *Parser) advance() Token {
	if !p.isAtEnd() {
		p.current++
	}
	return p.previous()
}

func (p *Parser) check(t TokenType) bool {
	if p.isAtEnd() {
		return false
	}
	return p.peek().ttype == t
}

func (p *Parser) match(types ...TokenType) bool {
	for _, t := range types {
		if p.check(t) {
			p.advance()
			return true
		}
	}

	return false
}
