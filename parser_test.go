package dice

import (
	"math/rand"
	"testing"
)

func Test_best_worst(t *testing.T) {

	tests := map[string]*binary{
		"1w3d6": {
			op:   'w',
			left: &leaf{value: 1},
			right: &binary{
				op:     'd',
				left:   &leaf{value: 3},
				right:  &leaf{value: 6},
				value:  4,
				values: []number{1, 1, 2},
			},
			value:  1,
			values: []number{1},
		},
		"3b4d6": {
			op:   'b',
			left: &leaf{value: 3},
			right: &binary{
				op:     'd',
				left:   &leaf{value: 4},
				right:  &leaf{value: 6},
				value:  9,
				values: []number{1, 1, 2, 5},
			},
			value:  8,
			values: []number{1, 2, 5},
		},
		"1w(2b4d6)": {
			op:   'w',
			left: &leaf{value: 1},
			right: &binary{
				op:   'b',
				left: &leaf{value: 2},
				right: &binary{
					op:     'd',
					left:   &leaf{value: 4},
					right:  &leaf{value: 6},
					value:  9,
					values: []number{1, 1, 2, 5},
				},
				value:  8,
				values: []number{2, 5},
			},
			value:  2,
			values: []number{2},
		},
	}

	for test, expected := range tests {
		p := new(diceParser).init(test)

		failed := false

		n, ok := p.parse()
		if !ok {
			t.Errorf("failed to parse %v", test)
			failed = true
		}
		if ! n.Equals(expected) {
			t.Fatalf("%v expected to parse as %v but parsed as %v", test, expected, n)
			failed = true
		}

		r := rand.New(rand.NewSource(0))
		actual, ok := n.Eval(r)
		if !ok {
			t.Errorf("failed to evaluate %v", n)
			failed = true
		}

		r2 := rand.New(rand.NewSource(0))

		expectedVal, ok := expected.Eval(r2)
		if !ok {
			t.Errorf("failed to evaluate %v", expected)
			failed = true
		}

		if actual != expectedVal {
			t.Errorf("expected %v for %v but got %v", expectedVal, test, actual)
			failed = true
		}

		if ! failed {
			t.Logf("OK: evaluated %s as: %v", test, expected)
		}
	}
}

func Test_best_worst_neg(t *testing.T) {
	tests := map[string]*binary{
		"1w2b4d6": {
			op: 'b',
			left: &binary{
				op: 'w',
				left: &leaf{
					value: 1,
				},
				right: &leaf{
					value: 2,
				},
			},
			right: &binary{
				op: 'd',
				left: &leaf{
					value: 4,
				},
				right: &leaf{
					value: 6,
				},
			},
		},
	}

	for test, expected := range tests {
		p := new(diceParser).init(test)

		failed := false
		n, ok := p.parse()
		if !ok {
			t.Errorf("failed to parse %v", test)
			failed = true
		}

		if ! n.Equals(expected) {
			t.Errorf("%v did not parse to %v as expecetd", test, expected)
			failed = true
		}

		r := rand.New(rand.NewSource(0))
		val, ok := n.Eval(r)
		if ok {
			t.Errorf("should not be able to evaluate %v but got: %v", test, val)
			failed = true
		}

		if ! failed {
			t.Logf("OK: parsed (but faield to evaluate) %s as: %v", test, n)
		}
	}
}
