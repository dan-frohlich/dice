package lex

import (
	"fmt"
	"strings"
)

type AST interface {
	Evaluate() int
}

type NodeType byte

const (
	NodeTypeLeaf NodeType = iota
	NodeTypeInfixOperator
	NodeTypePostfixOperator
)

type node struct {
	kind     NodeType
	v        int
	operand1 *node
	operand2 *node
	operator string
}

func (n *node) Evaluate() int {
	return n.v
}

func (n *node) String() string {
	switch (n.kind) {
	case NodeTypeLeaf:
		return fmt.Sprintf("%d", n.v)
	case NodeTypePostfixOperator:
		return fmt.Sprintf("(%v%s)", n.operand1, n.operator)
	case NodeTypeInfixOperator:
		return fmt.Sprintf("(%v%s%v)", n.operand1, n.operator, n.operand2)
	default:
		return fmt.Sprintf("[unhandled node type: %v]", n.kind)
	}
}

func (n *node) Equal(m *node) bool {
	return equal(n, m)
}

func equal(m, n *node) bool {
	if m == nil || n == nil {
		return n == nil && m == nil
	}
	return n.kind == m.kind &&
		strings.EqualFold(n.operator, m.operator) &&
		equal(n.operand1, m.operand1) &&
		equal(n.operand2, m.operand2)

}
