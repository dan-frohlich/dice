package dice

import (
	"github.com/dan-frohlich/dice/lex"
	"math/rand"
	"regexp"
	"strings"
	"time"
)

type Roller interface {
	Roll(input string) (result int, plan string, err error)
}

type roller struct {
	r *rand.Rand
}

func NewSeededRoller(seed int64) Roller {
	return roller{r: rand.New(rand.NewSource(seed))}
}

func NewRoller() Roller {
	return NewSeededRoller(time.Now().UnixNano())
}

func (r roller) Roll(input string) (result int, plan string, err error) {
	re := regexp.MustCompile(` |\t|\n`)
	sanitized := re.ReplaceAllString(input, "")
	p := lex.NewParser(strings.NewReader(sanitized))
	var ast lex.AST
	ast, err = p.Parse()

	if err == nil {
		result, _, err = ast.Evaluate(r.r)
		if err == nil {
			return int(result), ast.Plan(), nil
		}
	}
	return 0, "", err
}
