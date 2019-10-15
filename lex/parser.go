package lex

import (
	"fmt"
	"io"
	"strconv"
)

type Parser interface {
	Parse() (AST, error)
}

type parser struct {
	l        Lexer
	stack    []*node
	registry map[TokenType]tokenProcessor
	err      error
}

func NewParser(in io.Reader) Parser {
	p := &parser{l: NewLexer(in)}
	p.registry = map[TokenType]tokenProcessor{
		TokenLiteral:         handleLiteral,
		TokenEndOfStream:     handleEOS,
		TokenInfixOperator:   handleIFO,
		TokenPostfixOperator: handlePFO,
		TokenOpenParen:       p.handleOP,
		TokenCloseParen:      p.handleCP,
	}
	return p
}

func (p *parser) pop() *node {
	if len(p.stack) == 0 {
		return nil
	}
	n := p.stack[len(p.stack)-1]
	p.stack = p.stack[:len(p.stack)-1]
	return n
}

func (p *parser) push(n *node) {
	p.stack = append(p.stack, n)
}

func (p *parser) Parse() (AST, error) {
	p.l.Lex(p.accumulator)
	for i, n := range p.stack {
		if i > 0 {
			p.stack[i-1].operand2 = n
		}
	}
	var n *node
	if len(p.stack) > 0 {
		n = p.stack[0]
	}
	return n, p.err
}

func (p *parser) accumulator(t Token) {
	if p.err == nil {
		fn, ok := p.registry[t.Kind]
		if ok {
			ast, err := fn(p.pop(), t)
			p.err = err
			p.push(ast)
		} else {
			p.err = fmt.Errorf("unregistered Token Type: %v", t)
		}
	}
}

type tokenProcessor func(ast *node, t Token) (*node, error)

func handleLiteral(prev *node, t Token) (*node, error) {
	i, err := strconv.ParseInt(t.Value, 10, 64)
	n := &node{
		kind: NodeTypeLeaf,
		v:    int(i),
	}

	if prev == nil {
		return n, err
	}
	switch (prev.kind) {
	case NodeTypeLeaf:
		return nil, fmt.Errorf("parse error: %v %v", prev, n)
	case NodeTypeInfixOperator:
		if prev.operand2 == nil {
			prev.operand2 = n
			return prev, nil
		}
		return nil, fmt.Errorf("parse error: %v %v", prev, n)
	default:
		return nil, fmt.Errorf("parse error: %v %v", prev, n)
	}
}

func handlePFO(ast *node, t Token) (i *node, e error) {
	return operator(NodeTypePostfixOperator, ast, t), nil
}

func handleIFO(ast *node, t Token) (i *node, e error) {
	return operator(NodeTypeInfixOperator, ast, t), nil
}

func operator(kind NodeType, ast *node, t Token) *node {
	o1 := ast
	if o1 == nil {
		o1 = &node{kind: NodeTypeLeaf, v: 1}
	}
	n := &node{
		kind:     kind,
		operator: t.Value,
		operand1: o1,
	}
	return n
}

func handleEOS(ast *node, t Token) (i *node, e error) {

	return ast, nil
}

func (p *parser) handleOP(ast *node, t Token) (i *node, e error) {
	p.push(ast)
	return nil, nil
}

func (p *parser) handleCP(ast *node, t Token) (i *node, e error) {
	return ast, nil
}
