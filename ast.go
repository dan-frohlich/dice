package dice

import (
	"fmt"
	"math/rand"
	"sort"
	"strings"
)

type number int

type node interface {
	Eval(r *rand.Rand) (number, bool)
	String() string
	Equals(other node) bool
}

type sided interface {
	Sides() []int
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
		if s, ok := n.right.(sided); ok {
			n.value, n.values = rollSides(r, left, s.Sides())
			return n.value, true
		}
		n.value, n.values = roll(r, left, right)
		return n.value, true
	case 'b':
		if right == 0 || left == 0 {
			n.value = 0
			return n.value, true
		}
		b, ok := n.right.(*binary)
		if !ok {
			return 0, false
		}
		lv, _ := n.left.Eval(r)
		if len(b.values) < int(lv) {
			return 0, false
		}
		n.values = sortNumbers(b.values)[len(b.values)-int(left):]
		n.value = 0
		for _, v := range n.values {
			n.value += v
		}
		return n.value, true
	case 'w':
		if right == 0 || left == 0 {
			n.value = 0
			return n.value, true
		}
		b, ok := n.right.(*binary)
		if !ok {
			return 0, false
		}
		lv, _ := n.left.Eval(r)
		if len(b.values) < int(lv) {
			return 0, false
		}
		n.values = sortNumbers(b.values)[:left]
		n.value = 0
		for _, v := range n.values {
			n.value += v
		}
		return n.value, true
	}
	n.value = 0
	return n.value, false
}

func (n *binary) Equals(other node) bool {
	o, ok := other.(*binary)
	return ok && o.op == n.op && n.left.Equals(o.left) && n.right.Equals(o.right)
}

func sortNumbers(z []number) []number {
	ints := make([]int, len(z))
	for i, v := range z {
		ints[i] = int(v)
	}
	sort.Ints(ints)
	for i, v := range ints {
		z[i] = number(v)
	}
	return z
}

func rollSides(r *rand.Rand, diceCount number, sides []int) (result number, parts []number) {
	acc := 0
	parts = make([]number, diceCount)
	for i := 0; i < int(diceCount); i++ {
		roll := r.Intn(len(sides))
		rawValue := sides[roll]
		parts[i] = number(rawValue)
		acc += rawValue
	}
	return number(acc), parts
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
	sides []int
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

func (n *leaf) Equals(other node) bool {
	l, ok := other.(*leaf)
	return ok && l.value == n.value
}

// fudge-sided value AST node
type fudge struct{}

func (n *fudge) Sides() []int {
	return []int{-1, 0, 1}
}

func (n *fudge) Eval(r *rand.Rand) (number, bool) {
	return 3, true
}

func (n *fudge) String() string {
	return "F"
}

func (n *fudge) Equals(other node) bool {
	_, ok := other.(*fudge)
	return ok
}
