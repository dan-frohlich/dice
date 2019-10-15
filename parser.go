package dice

type diceParser struct {
	lexer      *lexer
	precedence map[byte]int
}

func (p *diceParser) init(data string) *diceParser {
	p.lexer = new(lexer).init(data)
	//p.precedence = make(map[byte]int)
	p.precedence = map[byte]int{
		'+': 1,
		'-': 1,
		'*': 2,
		'/': 2,
		'd': 4,
		'b': 3,
		'w': 3,
	}
	p.lexer.next()
	return p
}

func (p *diceParser) reset(data string) *diceParser {
	p.lexer = new(lexer).init(data)
	p.lexer.next()
	return p
}

//func (p *diceParser) AddOperator(op byte, precedence int) {
//	p.precedence[op] = precedence
//}

func (p *diceParser) parse() (node, bool) {
	lhs, ok := p.parsePrimary()
	if !ok {
		return nil, false
	}
	// starting with 1 instead of 0, because
	// map[*]int returns 0 for non-existant items
	node, ok := p.parseOperators(lhs, 1)
	if !ok {
		return nil, false
	}
	return node, true
}

func (p *diceParser) parsePrimary() (node, bool) {
	switch p.lexer.kind {
	case tokenFudge:
		node := &fudge{}
		p.lexer.next()
		return node, true
	case tokenNum:
		node := &leaf{value: p.lexer.num}
		p.lexer.next()
		return node, true
	case tokenLPar:
		p.lexer.next()
		node, ok := p.parse()
		if !ok {
			return nil, false
		}
		if p.lexer.kind == tokenRPar {
			p.lexer.next()
		}
		return node, true
	}
	return nil, false
}

func (p *diceParser) parseOperators(lhs node, min_precedence int) (node, bool) {
	var ok bool
	var rhs node
	for p.lexer.kind == tokenOp && p.precedence[p.lexer.oper] >= min_precedence {
		op := p.lexer.oper
		p.lexer.next()
		rhs, ok = p.parsePrimary()
		if !ok {
			return nil, false
		}
		for p.lexer.kind == tokenOp && p.precedence[p.lexer.oper] > p.precedence[op] {
			op2 := p.lexer.oper
			rhs, ok = p.parseOperators(rhs, p.precedence[op2])
			if !ok {
				return nil, false
			}
		}
		lhs = new(binary).init(op, lhs, rhs)
	}
	return lhs, p.lexer.kind != tokenErr
}
