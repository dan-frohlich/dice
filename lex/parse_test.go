package lex

import (
	"errors"
	"fmt"
	"sort"
	"strings"
	"testing"
)

type parserResult struct {
	node *node
	err  error
}

func (r parserResult) String() string {
	return fmt.Sprintf("{n:%v, e:%v}", r.node, r.err)
}

type parserTestCases map[string]parserResult

func Test_parser(t *testing.T) {
	testCases := parserTestCases{
		"3":  parserResult{node: &node{kind: NodeTypeLeaf, v: 3}},
		"34": parserResult{node: &node{kind: NodeTypeLeaf, v: 33}},
		"z7": parserResult{err: errors.New("unregistered Token Type: err:unhandled char: z @ offset 1")},
		"3d6": parserResult{
			node: &node{
				kind:     NodeTypeInfixOperator,
				operator: "d",
				operand1: &node{kind: NodeTypeLeaf, v: 3},
				operand2: &node{kind: NodeTypeLeaf, v: 6},
			},
		},
		"d20": parserResult{
			node: &node{
				kind:     NodeTypeInfixOperator,
				operator: "d",
				operand1: &node{kind: NodeTypeLeaf, v: 1},
				operand2: &node{kind: NodeTypeLeaf, v: 20},
			},
		},
		"1d4d%": parserResult{
			node: &node{
				kind:     NodeTypePostfixOperator,
				operator: "d%",
				operand1: &node{
					kind:     NodeTypeInfixOperator,
					operator: "d",
					operand1: &node{kind: NodeTypeLeaf, v: 1},
					operand2: &node{kind: NodeTypeLeaf, v: 4},
				},
			},
		},
		"2d6!": parserResult{
			node: &node{
				kind:     NodeTypePostfixOperator,
				operator: "!",
				operand1: &node{
					kind:     NodeTypeInfixOperator,
					operator: "d",
					operand1: &node{kind: NodeTypeLeaf, v: 2},
					operand2: &node{kind: NodeTypeLeaf, v: 6},
				},
			},
		},
		"d%/2": parserResult{
			node: &node{
				kind:     NodeTypeInfixOperator,
				operator: "/",
				operand1: &node{
					kind:     NodeTypePostfixOperator,
					operator: "d%",
					operand1: &node{kind: NodeTypeLeaf, v: 1},
				},
				operand2: &node{kind: NodeTypeLeaf, v: 2},
			},
		},
		"3d(2+(1d4))": parserResult{
			node: &node{
				kind:     NodeTypeInfixOperator,
				operator: "d",
				operand1: &node{kind: NodeTypeLeaf, v: 3},
				operand2: &node{
					kind:     NodeTypeInfixOperator,
					operator: "+",
					operand1: &node{kind: NodeTypeLeaf, v: 2},
					operand2: &node{
						kind:     NodeTypeInfixOperator,
						operator: "d",
						operand1: &node{kind: NodeTypeLeaf, v: 1},
						operand2: &node{kind: NodeTypeLeaf, v: 4},
					},
				},
			},
		},
	}
	runParserTestCases(testCases, t)
}

func runParserTestCases(tests parserTestCases, t *testing.T) {

	size := len(tests)
	orderedTestKeys := make([]string, size, size)
	i := 0
	for test := range tests {
		orderedTestKeys[i] = test
		i++
	}
	sort.Strings(orderedTestKeys)
	for _, test := range orderedTestKeys {
		runParserTestCase(test, tests[test], t)
	}
}

func runParserTestCase(test string, expected parserResult, t *testing.T) {

	p := NewParser(strings.NewReader(test))
	ast, err := p.Parse()
	n, ok := ast.(*node)
	if ! ok {
		t.Errorf("ERROR %s : result %v (%T) is not type %T", test, ast, ast, &node{})
	}
	actual := parserResult{node: n, err: err}

	if ! expected.Equal(actual) {
		t.Errorf("ERROR %v\texpected\t%v\tgot\t%v", test, expected, actual)
	} else {
		t.Logf("OK %v \tParsed as %v", test, actual)
	}
}

func (r parserResult) Equal(other parserResult) bool {
	if r.err == nil && other.err == nil {
		return r.node.Equal(other.node)
	}
	return r.err != nil &&
		other.err != nil &&
		r.node.Equal(other.node) &&
		strings.EqualFold(r.err.Error(), other.err.Error())
}
