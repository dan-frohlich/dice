package dice

type lexer struct {
	data string
	pos  int
	kind int
	num  number
	oper byte
}

const (
	tokenUnk   = iota // unknown
	tokenErr          // error
	tokenNum          // number
	tokenLPar         // left parenthesis
	tokenRPar         // right parenthesis
	tokenOp           // operator
	tokenFudge        // FudgeDie
)

func (l *lexer) init(data string) *lexer {
	l.data = data
	l.pos = 0
	l.kind = tokenUnk
	l.num = 0
	l.oper = 0
	return l
}

func (l *lexer) next() int {
	n := len(l.data)
	l.kind = tokenUnk

	if l.pos < n {

		//x := fmt.Sprintf("char[%d]: %d: '%s'\n", l.pos, char, string([]byte{char}))
		//fmt.Println(x)
		switch char := l.data[l.pos]; char {
		case '*', '+', '-', '/', 'b', 'd', 'w':
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
		case 'F':
			l.kind = tokenFudge
			l.num = 3
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
		case ' ':
			//no-op
		default:
			l.kind = tokenErr
			break
		}
	}
	return l.kind
}
