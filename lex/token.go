package lex

import "fmt"

type Token struct {
	Kind  TokenType
	Value string
}

type TokenType byte

const (
	TokenUnknown TokenType = iota
	TokenLiteral
	TokenPostfixOperator
	TokenInfixOperator
	TokenOpenParen
	TokenCloseParen
	TokenError
	TokenEndOfStream
)

func (t TokenType) String() string {
	var s string
	switch t {
	case TokenUnknown:
		s = "ukn"
	case TokenLiteral:
		s = "lit"
	case TokenPostfixOperator:
		s = "pfo"
	case TokenInfixOperator:
		s = "ifo"
	case TokenOpenParen:
		s = "op"
	case TokenCloseParen:
		s = "cp"
	case TokenError:
		s = "err"
	case TokenEndOfStream:
		s = "eos"
	}
	return s
}

func (t Token) String() string {
	return fmt.Sprintf("%v:%v", t.Kind, t.Value)
}
