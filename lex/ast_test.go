package lex

import (
	"errors"
	"fmt"
	"math/rand"
	"sort"
	"strings"
	"testing"
)

type anyASTResult interface {
	Equal(other anyASTResult) bool
}

type astTestCases map[string]anyASTResult

type simpleASTResult struct {
	v int
	z []int
	e error
}

func (r simpleASTResult) String() string {
	return fmt.Sprintf("{n:%v, z:%v, e:%v}", r.v, r.z, r.e)
}

func (r simpleASTResult) Equal(other anyASTResult) bool {
	o, ok := other.(simpleASTResult)
	if !ok {
		return false
	}

	if r.e == nil || o.e == nil {
		return r.v == o.v && r.e == nil && o.e == nil
	}

	return r.v == o.v &&
		eq(r.z, o.z) &&
		r.e != nil &&
		o.e != nil &&
		strings.EqualFold(r.e.Error(), o.e.Error())
}

func eq(a, b []int) bool {
	if len(a) != len(b) {
		return false
	}
	for i, v := range b {
		if a[i] != v {
			return false
		}
	}
	return true
}

type diceASTExpectedResult struct {
	min, max int
	e        error
}

func (r diceASTExpectedResult) String() string {
	return fmt.Sprintf("{[%v,%v], e:%v}", r.min, r.max, r.e)
}
func (r diceASTExpectedResult) Equal(other anyASTResult) bool {
	o, ok := other.(simpleASTResult)
	if !ok {
		return false
	}

	if r.e == nil || o.e == nil {
		return r.min <= o.v && r.max >= o.v && r.e == nil && o.e == nil
	}
	return r.min <= o.v &&
		r.max >= o.v &&
		r.e != nil &&
		o.e != nil &&
		strings.EqualFold(r.e.Error(), o.e.Error())
}

func Test_ast(t *testing.T) {
	tests := astTestCases{
		"":       simpleASTResult{z: []int{}, e: errors.New("nill node")},
		"1wd%":   diceASTExpectedResult{min: 1, max: 100},
		"1":      simpleASTResult{v: 1},
		"1+3":    simpleASTResult{v: 4},
		"1*3":    simpleASTResult{v: 3},
		"1w2d20": diceASTExpectedResult{min: 1, max: 20},
		"1w3d6":  diceASTExpectedResult{min: 1, max: 6},
		"2b3d6":  diceASTExpectedResult{min: 2, max: 12},
		"2d6+12": diceASTExpectedResult{min: 14, max: 24},
		"3b4d6":  simpleASTResult{v: 17, z: []int{5, 6, 6}},
		"3d6":    diceASTExpectedResult{min: 3, max: 18},
		"4/2":    simpleASTResult{v: 2},
		"4b2d6":  diceASTExpectedResult{e: errors.New("(4b(2d6)) can't gather 4 best items from a slice of 2 items")},
		"4b3d10": diceASTExpectedResult{e: errors.New("(4b(3d10)) can't gather 4 best items from a slice of 3 items")},
		"5d6":    diceASTExpectedResult{min: 5, max: 30},
		"d%":     diceASTExpectedResult{min: 1, max: 100},
	}
	runASTTestCases(tests, t)
}

func runASTTestCases(tests astTestCases, t *testing.T) {

	size := len(tests)
	orderedTestKeys := make([]string, size, size)
	i := 0
	for test := range tests {
		orderedTestKeys[i] = test
		i++
	}
	sort.Strings(orderedTestKeys)
	for _, test := range orderedTestKeys {
		runASTTestCase(test, tests[test], t)
	}
}

func runASTTestCase(test string, expected anyASTResult, t *testing.T) {

	p := NewParser(strings.NewReader(test))
	ast, err := p.Parse()
	if err != nil {
		t.Error(err)
		return
	}
	n, ok := ast.(*node)
	if !ok {
		t.Errorf("ERROR %s : result %v (%T) is not type %T", test, ast, ast, &node{})
	}
	v, z, e := n.Evaluate(rand.New(rand.NewSource(11)))
	actual := simpleASTResult{v: v, z: z, e: e}

	if !expected.Equal(actual) {
		t.Errorf("ERROR %v\texpected\t%v\tgot\t%v", test, expected, actual)
	} else {
		t.Logf("OK %12v evaluated as %v", test, actual)
	}
}
