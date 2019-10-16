package lex

import (
	"errors"
	"fmt"
	"math/rand"
	"sort"
)

//AST Abstract Syntax Tree
type AST interface {
	Evaluate(*rand.Rand) (int, []int, error)
	Plan() string
	String() string
}

//NodeType identifies the node type
type NodeType byte

const (
	//NodeTypeLeaf termin node of the tree. A Literal.
	NodeTypeLeaf NodeType = iota
	//NodeTypeInfixOperator is a binary operator.
	NodeTypeInfixOperator
	//NodeTypePrefixOperator is a unary operator.
	NodeTypePrefixOperator
	//NodeTypePostfixOperator is a unary operator.
	NodeTypePostfixOperator
)

type node struct {
	kind     NodeType
	v        int
	vs       []int
	operand1 *node
	operand2 *node
	operator string
}

//Evaluate evaluates the AST
func (n *node) Evaluate(r *rand.Rand) (int, []int, error) {
	switch n.kind {
	case NodeTypeLeaf:
		return n.v, []int{n.v}, nil
	case NodeTypeInfixOperator:
		return n.evalInfix(r)
	case NodeTypePostfixOperator:
		return n.evalPostfix(r)
	default:
		return 0, []int{0}, fmt.Errorf("unknown node type: %v", n)
	}
}

func (n *node) evalPostfix(r *rand.Rand) (int, []int, error) {
	result := 0
	results := []int{result}
	//right, _, err := n.operand2.Evaluate(r)
	if n.operand1 == nil {
		n.operand1 = &node{kind: NodeTypeLeaf, v: 1, vs: []int{1}}
	}
	left, _, err := n.operand1.Evaluate(r)
	if err != nil {
		return result, results, err
	}
	switch n.operator {
	case "d%":
		result, results, err = n.evalDice(r, left, 100)
	case "dF":
		_, results, err = n.evalDice(r, left, 3)
		result =0
		for i, v := range results {
			results[i] = v-2
			result += results[i]
		}
	default:
		err = fmt.Errorf("operator not implemented: %s", n.operator)
	}
	return result, results, err
}

func (n *node) evalInfix(r *rand.Rand) (int, []int, error) {
	result := 0
	results := []int{result}
	left, _, err := n.operand1.Evaluate(r)

	if err != nil {
		return result, results, err
	}
	right, rights, err := n.operand2.Evaluate(r)
	if err != nil {
		return result, results, err
	}
	switch n.operator {
	case "+", "-", "*", "/":
		result, err = n.evalMathOperators(r, left, right)
		results = []int{result}
	case "d":
		result, results, err = n.evalDice(r, left, right)
	case "b":
		result, results, err = n.evalBest(r, left, rights)
	case "w":
		result, results, err = n.evalWorst(r, left, rights)
	default:
		err = fmt.Errorf("unhandled operator: %v", n.operator)
	}
	return result, results, err
}

func (n *node) evalBest(r *rand.Rand, left int, rights []int) (int, []int, error) {
	if left > len(rights) {
		return 0, []int{0}, fmt.Errorf("%v can't gather %d best items from a slice of %d items", n, left, len(rights))
	}
	s := make([]int, len(rights))
	for i, v := range rights {
		s[i] = v
	}
	sort.Ints(s)
	n.v = 0
	n.vs = s[len(s)-left:]
	if len(n.vs) != left {
		return 0, nil, errors.New("FOOBAR")
	}
	for _, v := range n.vs {
		n.v += v
	}
	return n.v, n.vs, nil
}

func (n *node) evalWorst(r *rand.Rand, left int, rights []int) (int, []int, error) {
	if left > len(rights) {
		return 0, []int{0}, fmt.Errorf("%v can't gather %d worst items from a slice of %d items", n, left, len(rights))
	}
	s := make([]int, len(rights))
	n.v = 0
	for i, v := range rights {
		s[i] = v
	}
	sort.Ints(s)
	n.vs = s[:left]
	// fmt.Println(rights)
	// fmt.Println(n.vs)
	for i, v := range n.vs {
		s[i] = v
		n.v += v
	}
	return n.v, n.vs, nil
}

func (n *node) evalDice(r *rand.Rand, left int, right int) (int, []int, error) {
	acc := 0
	results := make([]int, left)
	for i := 0; i < left; i++ {
		results[i] = r.Intn(right) + 1
		acc += results[i]
	}
	n.v = acc
	n.vs = results
	// fmt.Println(n.vs)
	return acc, results, nil
}

func (n *node) evalMathOperators(r *rand.Rand, left int, right int) (int, error) {
	var result int
	var err error
	switch n.operator {
	case "+":
		result = left + right
	case "-":
		result = left - right
	case "*":
		result = left * right
	case "/":
		if right == 0 {
			err = fmt.Errorf("divide by zero in %v", n)
		} else {
			result = left / right
		}
	default:
		err = fmt.Errorf("unhandled operator: %v", n.operator)
	}
	n.v = result
	return n.v, err
}
func (n *node) String() string {
	if n == nil {
		return "<nil>"
	}
	switch n.kind {
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

func (n *node) Plan() string {
	if n == nil {
		return "<nil>"
	}
	switch n.kind {
	case NodeTypeLeaf:
		return fmt.Sprintf("%d", n.v)
	case NodeTypePostfixOperator:
		return fmt.Sprintf("(%v%s %v)", n.operand1, n.operator, n.vs)
	case NodeTypeInfixOperator:
		return fmt.Sprintf("(%v%s%v %v)", n.operand1, n.operator, n.operand2, n.vs)
	default:
		return fmt.Sprintf("[unhandled node type: %v]", n.kind)
	}
}
