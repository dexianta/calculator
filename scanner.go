package main

import (
	"fmt"
	"log"
	"strconv"
)

type TokenType int

const (
	MINUS TokenType = iota + 1
	PLUS
	SLASH
	STAR
	LEFT_PAREN
	RIGHT_PAREN

	NUMBER

	// end of the file
	EOF
)

type Token struct {
	ttype   TokenType
	lexeme  string
	literal interface{}
	line    int
}

type Scanner struct {
	source  string
	tokens  []Token
	start   int
	current int
	line    int
}

func NewScanner(s string) Scanner {
	scanner := Scanner{
		source: s,
		line:   1,
	}

	scanner.scanTokens()
	return scanner
}

func (s *Scanner) scanTokens() []Token {
	for s.current < len(s.source) {
		s.start = s.current
		s.scanToken()
	}

	s.tokens = append(s.tokens, Token{EOF, "", nil, s.line})
	return s.tokens
}

func (s *Scanner) scanToken() {
	c := s.advance()
	switch c {
	case '(':
		s.addToken(LEFT_PAREN, nil)
	case ')':
		s.addToken(RIGHT_PAREN, nil)
	case '+':
		s.addToken(PLUS, nil)
	case '-':
		s.addToken(MINUS, nil)
	case '*':
		s.addToken(STAR, nil)
	case '/':
		s.addToken(SLASH, nil)
	case ' ', '\r', '\t': // ignore all these
	case '\n':
		s.line++
	default:
		// getting the actual number
		if s.isDigit(c) {
			s.number()
		} else {
			log.Fatal("error parsing number: ", fmt.Sprintf("%c", c))
		}
	}
}

func (s *Scanner) addToken(tokenType TokenType, literal interface{}) {
	lexeme := s.getLexeme()
	s.tokens = append(s.tokens, Token{
		tokenType,
		lexeme,
		literal,
		s.line,
	})
}

func (s *Scanner) isDigit(c byte) bool {
	return c <= '9' && c >= '0'
}

func (s *Scanner) number() {
	for s.isDigit(s.peek()) {
		s.advance()
	}

	if s.match('.') && s.isDigit(s.peekNext()) {
		s.advance()
		for s.isDigit(s.peek()) {
			s.advance()
		}
	}

	n, err := strconv.ParseFloat(s.getLexeme(), 64)
	if err != nil {
		log.Fatal("get number error: ", err.Error())
	}

	s.addToken(NUMBER, n)
}

func (s *Scanner) getLexeme() string {
	return s.source[s.start:s.current]
}

func (s *Scanner) peekNext() byte {
	if s.isAtEnd() {
		panic("parse error, expect symbol but at end of file")
	}
	return s.source[s.current+1]
}

func (s *Scanner) isAtEnd() bool {
	return s.current >= len(s.source)
}

func (s *Scanner) peek() byte {
	return s.source[s.current]
}

func (s *Scanner) match(c byte) bool {
	return s.peek() == c
}

func (s *Scanner) advance() byte {
	s.current++
	return s.source[s.current-1]
}
