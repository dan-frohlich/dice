package lex

import (
	"io"
)

type TokenReceiver func(t Token)

type Lexer interface {
	Lex(receiver TokenReceiver)
}

func NewLexer(in io.Reader) Lexer {
	return &lexer{
		buf: make([]byte, 1, 1),
		in:  in,
	}
}

type lexer struct {
	in    io.Reader
	model AST
	token *Token
	buf   []byte
	pos   int
}

func (l *lexer) byte() byte {
	return l.buf[0]
}

func (l *lexer) next() error {
	var n int
	var err error
	for n == 0 && err == nil {
		n, err = l.read()
		if err != nil {
			return err
		}
	}
	return err
}

func (l *lexer) Lex(receiver TokenReceiver) {
	emitter := func(t *Token) {
		if t != nil {
			receiver(*t)
		}
	}
	for state := advanceOneByte; state != nil; state = state(l) {
		t := l.token
		emitter(t)
	}
}
