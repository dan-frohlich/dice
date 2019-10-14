package dice

import (
	"testing"
)

func Test_can_add(t *testing.T) {
	tests := map[string]int{
		"1+1":   2,
		"0+0":   0,
		"101+1": 102,
	}

	roller := NewRoller()

	for test, expected := range tests {
		actual, _, err := roller.Roll(test)
		if err != nil {
			t.Error("ERROR", test, err)
			continue
		}
		if actual != expected {
			t.Error("ERROR", test, "expected", expected, "got", actual)
			continue
		}
		t.Log("OK", test, "expected", expected, "got", actual)
	}
}

func Test_various_dice(t *testing.T) {
	tests := map[string][]int{
		"3d6":     {3, 18},
		"2d10":    {2, 20},
		"1d20":    {1, 20},
		"1d100":   {1, 100},
		"3d6+2":   {5, 20},
	}

	samples := 10000
	roller := NewRoller()

	for test, minMax := range tests {
		min := minMax[0]
		max := minMax[1]
		testOK := true
		actual := 0
		var err error
		for i := 0; i < samples; i++ {
			actual, _, err = roller.Roll(test)
			if err != nil {
				t.Error("ERROR", test, "error:", err)
				testOK = false
				break
			}
			if actual < min || actual > max {
				t.Error("ERROR", test, "expected result in [", min, ",", max, "] got", actual)
				testOK = false
				break
			}
		}
		if testOK {
			t.Logf("OK %s all rolls within [%d,%d] after %d attempts: %d", test, min, max, samples, actual)
		}
	}
}

func Test_various(t *testing.T) {
	tests := map[string]int{
		"(1+3)*7": 28, // 28, example from task description.
		"1+3*7":   22, // 22, shows operator precedence.
		"7":       7,  // 7, a single literal is a valid expression.
		"7/3":     2,  // eval only does integer math.
		"7.3":     7,  //decimals are read as int
		"7.3+1.9": 8,  //decimals are read as int
	}

	roller := NewRoller()

	for test, expected := range tests {
		actual, _, err := roller.Roll(test)
		if err != nil {
			t.Error(err)
			continue
		}
		if actual != expected {
			t.Error("ERROR", test, "expected", expected, "got", actual)
			continue
		}
		t.Log("OK", test, "expected", expected, "got", actual)
	}
}

func Test_various_neg(t *testing.T) {
	tests := []string{
		"7^3", // parses, but disallowed in eval.
		"go",  // a valid keyword, not valid in an expression.
		"3@7", // error message is "illegal character."
		"",    // EOF seems a reasonable error message.
	}
	roller := NewRoller()

	for _, test := range tests {
		actual, _, err := roller.Roll(test)
		if err != nil {
			t.Log("OK", test, "got", err)
			continue
		}
		t.Error("ERROR", test, "expected", "error", "got", actual)
	}

}
