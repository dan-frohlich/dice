package lex

import (
	"fmt"
	"io"
)

type stateFn func(l *lexer) stateFn

func detector(l *lexer) stateFn {
	switch l.byte() {
	case '0', '1', '2', '3', '4', '5', '6', '7', '8', '9':
		l.token = nil
		return readingNumber
	case '(':
		l.token = &Token{Kind: TokenOpenParen, Value: string(l.buf)}
		return advanceOneByte
	case ')':
		l.token = &Token{Kind: TokenCloseParen, Value: string(l.buf)}
		return advanceOneByte
	case '+', '-', '*', '/':
		l.token = &Token{Kind: TokenInfixOperator, Value: string(l.buf)}
		return advanceOneByte
	case '!':
		l.token = &Token{Kind: TokenPostfixOperator, Value: string(l.buf)}
		return advanceOneByte
	case 'd':
		l.token = nil
		return readingDiceRoller
	default:
		return l.handleError(fmt.Errorf("unhandled char: %c @ offset %d", l.byte(), l.pos))
	}
}

func advanceOneByte(l *lexer) stateFn {
	return l.readOne()
}

func (l *lexer) readOne() stateFn {
	l.token = nil
	if err := l.next(); err != nil {
		return l.handleReadError(err)
	}
	return detector
}

func readingNumber(l *lexer) stateFn {
	bytes := make([]byte, 0)
	var err error

	for l.byte() >= '0' && l.byte() <= '9' && err == nil {
		bytes = append(bytes, l.byte())
		_, err = l.read()
	}
	l.token = &Token{Kind: TokenLiteral, Value: string(bytes)}
	if err != nil {
		return l.handleReadError(err)
	}
	return detector
}

func (l *lexer) read() (int, error) {
	n, err := l.in.Read(l.buf)
	l.pos += n
	return n, err
}

func readingDiceRoller(l *lexer) stateFn {
	bytes := make([]byte, 0)
	var err error

	acceptable := func(b byte) bool {
		return b == 'd' || b == 'F' || b == '%'
	}
	for acceptable(l.byte()) && err == nil {
		bytes = append(bytes, l.byte())
		_, err = l.read()
	}
	tt := TokenInfixOperator
	if len(bytes) > 1 {
		tt = TokenPostfixOperator
	}
	l.token = &Token{Kind: tt, Value: string(bytes)}
	if err != nil {
		return l.handleReadError(err)
	}
	return detector
}

func (l *lexer) handleReadError(err error) stateFn {
	if err == io.EOF {
		return endOfStream
	}
	return l.handleError(err)
}

func endOfStream(l *lexer) stateFn {
	l.token = &Token{Kind: TokenEndOfStream}
	return terminal
}

func (l *lexer) handleError(err error) stateFn {
	l.token = &Token{Kind: TokenError, Value: err.Error()}
	return terminal
}

func terminal(l *lexer) stateFn {
	l.token = nil
	return nil
}
