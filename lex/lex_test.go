package lex

import (
	"fmt"
	"sort"
	"strings"
	"testing"
)

type tokenResult []Token
type lexerTestCases map[string]tokenResult

func (r tokenResult) String() string {
	s := make([]string, len(r))
	for i, v := range r {
		s[i] = v.String()
	}
	return fmt.Sprintf("{%s}", strings.Join(s, ", "))
}

func Test_lexer(t *testing.T) {
	tests := lexerTestCases{
		"3d6+12": {
			{Kind: TokenLiteral, Value: "3"},
			{Kind: TokenInfixOperator, Value: "d"},
			{Kind: TokenLiteral, Value: "6"},
			{Kind: TokenInfixOperator, Value: "+"},
			{Kind: TokenLiteral, Value: "12"},
			{Kind: TokenEndOfStream},
		},
		"4dF+2": {
			{Kind: TokenLiteral, Value: "4"},
			{Kind: TokenPostfixOperator, Value: "dF"},
			{Kind: TokenInfixOperator, Value: "+"},
			{Kind: TokenLiteral, Value: "2"},
			{Kind: TokenEndOfStream},
		},
		"3d(2+1d4)": {
			{Kind: TokenLiteral, Value: "3"},
			{Kind: TokenInfixOperator, Value: "d"},
			{Kind: TokenOpenParen, Value: "("},
			{Kind: TokenLiteral, Value: "2"},
			{Kind: TokenInfixOperator, Value: "+"},
			{Kind: TokenLiteral, Value: "1"},
			{Kind: TokenInfixOperator, Value: "d"},
			{Kind: TokenLiteral, Value: "4"},
			{Kind: TokenCloseParen, Value: ")"},
			{Kind: TokenEndOfStream},
		},
		"d%/2": {
			{Kind: TokenPostfixOperator, Value: "d%"},
			{Kind: TokenInfixOperator, Value: "/"},
			{Kind: TokenLiteral, Value: "2"},
			{Kind: TokenEndOfStream},
		},
		"2d6!": {
			{Kind: TokenLiteral, Value: "2"},
			{Kind: TokenInfixOperator, Value: "d"},
			{Kind: TokenLiteral, Value: "6"},
			{Kind: TokenPostfixOperator, Value: "!"},
			{Kind: TokenEndOfStream},
		},
		"fdx-2": {
			{Kind: TokenError, Value: "unhandled char: f @ offset 1"},
		},
		"3dx-2": {
			{Kind: TokenLiteral, Value: "3"},
			{Kind: TokenInfixOperator, Value: "d"},
			{Kind: TokenError, Value: "unhandled char: x @ offset 3"},
		},
		"1d6%2": {
			{Kind: TokenLiteral, Value: "1"},
			{Kind: TokenInfixOperator, Value: "d"},
			{Kind: TokenLiteral, Value: "6"},
			{Kind: TokenError, Value: "unhandled char: % @ offset 4"},
		},
	}
	runLexerTestCases(tests, t)
}

func runLexerTestCases(tests lexerTestCases, t *testing.T) {

	size := len(tests)
	orderedTestKeys := make([]string, size, size)
	i := 0
	for test := range tests {
		orderedTestKeys[i] = test
		i++
	}
	sort.Strings(orderedTestKeys)
	for _, test := range orderedTestKeys {
		runLexerTestCase(test, tests[test], t)
	}
}

func runLexerTestCase(test string, expected tokenResult, t *testing.T) {
	lexer := NewLexer(strings.NewReader(test))
	actual := make(tokenResult, 0)
	lexer.Lex(
		func(t Token) {
			actual = append(actual, t)
		})
	if ! expected.Equal(actual) {
		t.Errorf("ERROR %v\texpected\t%v\tgot\t%v", test, expected, actual)
	} else {
		t.Logf("OK %v \tlexed as %v", test, actual)
	}
}

func (r tokenResult) Equal(other tokenResult) bool {
	if len(r) != len(other) {
		return false
	}
	for i, v := range other {
		if r[i].Kind != v.Kind || !strings.EqualFold(r[i].Value, v.Value) {
			return false
		}

	}
	return true
}
