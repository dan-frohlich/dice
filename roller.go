package dice

import (
	"fmt"
	"math/rand"
	"regexp"
	"time"
)

type Roller interface {
	Roll(input string) (result int, plan string, err error)
}

type roller struct {
	r *rand.Rand
	p *diceParser
}

func NewSeededRoller(seed int64) Roller {
	p := new(diceParser).init("")
	return roller{r: rand.New(rand.NewSource(seed)), p: p}
}

func NewRoller() Roller {
	return NewSeededRoller(time.Now().UnixNano())
}

func (r roller) Roll(input string) (result int, plan string, err error) {
	re := regexp.MustCompile(` |\t|\n`)
	sanitized := re.ReplaceAllString(input, "")
	node, parseOk := r.p.reset(sanitized).parse()

	if parseOk {
		result, ok := node.Eval(r.r)
		if ok {
			return int(result), node.String(), nil
		}
		return 0, "", fmt.Errorf("failed to evaluate %s", sanitized)
	}
	return 0, "", fmt.Errorf("failed to parse %s", sanitized)
}
