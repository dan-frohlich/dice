package dice

import (
	"fmt"
	"math/rand"
	"strings"
)

type number int

type node interface {
	Eval(r *rand.Rand) (number, bool)
	String() string
}

// binary operator AST node
type binary struct {
	op     byte
	left   node
	right  node
	value  number
	values []number
}

func (n *binary) init(op byte, left, right node) node {
	n.op = op
	n.left = left
	n.right = right
	return n
}

func (n *binary) Eval(r *rand.Rand) (number, bool) {
	left, ok := n.left.Eval(r)
	if !ok {
		return 0, false
	}
	right, ok := n.right.Eval(r)
	if !ok {
		return 0, false
	}
	switch n.op {
	case '+':
		n.value = left + right
		return n.value, true
	case '-':
		n.value = left - right
		return n.value, true
	case '*':
		n.value = left * right
		return n.value, true
	case '/':
		if right == 0 {
			n.value = 0
			return n.value, false
		}
		n.value = left / right
		return n.value, true
	case 'd':
		if right == 0 || left == 0 {
			n.value = 0
			return n.value, true
		}
		n.value, n.values = roll(r, left, right)
		return n.value, true
	}
	n.value = 0
	return n.value, false
}

func roll(r *rand.Rand, diceCount, sides number) (result number, parts []number) {
	acc := 0
	parts = make([]number, diceCount)
	for i := 0; i < int(diceCount); i++ {
		roll := r.Intn(int(sides)) + 1
		parts[i] = number(roll)
		acc += roll
	}
	return number(acc), parts
}

func (n *binary) String() string {
	if n.values != nil && len(n.values) > 0 {
		s := strings.Replace(fmt.Sprint(n.values), " ", ",", -1)
		return fmt.Sprintf("(%s %c %s -> %d %s)", n.left, n.op, n.right, n.value, s)
	}
	return fmt.Sprintf("(%s %c %s -> %d)", n.left, n.op, n.right, n.value)
}

// leaf values AST node
type leaf struct {
	value number
}

func (n *leaf) init(value number) node {
	n.value = value
	return n
}

func (n *leaf) Eval(r *rand.Rand) (number, bool) {
	return n.value, true
}

func (n *leaf) String() string {
	return fmt.Sprintf("%v", n.value) // %v = default format
}
