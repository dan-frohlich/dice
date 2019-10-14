package dice

type lexer struct {
	data string
	pos  int
	kind int
	num  number
	oper byte
}

const (
	tokenUnk  = iota // unknown
	tokenErr         // error
	tokenNum         // number
	tokenLPar        // left parenthesis
	tokenRPar        // right parenthesis
	tokenOp          // operator
)

func (lexer *lexer) init(data string) *lexer {
	lexer.data = data
	lexer.pos = 0
	lexer.kind = tokenUnk
	lexer.num = 0
	lexer.oper = 0
	return lexer
}

//TODO implement best of and worst of (dice poll) ops : 3b4d6 (3 best or 4d6)
func (l *lexer) next() int {
	n := len(l.data)
	l.kind = tokenUnk

	if l.pos < n {

		//x := fmt.Sprintf("char[%d]: %d: '%s'\n", l.pos, char, string([]byte{char}))
		//fmt.Println(x)
		switch char := l.data[l.pos]; char {
		case '+', '-', '*', '/', 'd', 'b', 'w':
			l.pos++
			l.kind = tokenOp
			l.oper = char
		case '(':
			l.pos++
			l.kind = tokenLPar
			l.oper = char
		case ')':
			l.pos++
			l.kind = tokenRPar
			l.oper = char
		case '0', '1', '2', '3', '4', '5', '6', '7', '8', '9', '.':
			var value number = 0
			var divisor number = 1
			for ; l.pos < n && '0' <= l.data[l.pos] && l.data[l.pos] <= '9'; l.pos++ {
				value = value*10 + number(l.data[l.pos]-'0')
			}
			if l.pos < n && l.data[l.pos] == '.' {
				l.pos++
				for ; l.pos < n && '0' <= l.data[l.pos] && l.data[l.pos] <= '9'; l.pos++ {
					value = value*10 + number(l.data[l.pos]-'0')
					divisor *= 10
				}
			}
			l.kind = tokenNum
			l.num = value / divisor
		default:
			l.kind = tokenErr
			break
		}
	}
	return l.kind
}
