package func_purrser

import (
	"fmt"
	"io"
)

type parser struct {
	errs  []error
	input []Token
	i     int
}

func (ps parser) nextToken() Token {
	if ps.i >= len(ps.input) {
		ps.reportErr("Input abruptly finished while parsing expression")
		return ps.input[len(ps.input)-1]
	}
	tok := ps.input[ps.i]
	ps.i = ps.i + 1
	return tok
}

func (ps parser) peekToken() Token {
	if ps.i >= len(ps.input) {
		ps.reportErr("Input abruptly finished while parsing expression")
		return ps.input[len(ps.input)-1]
	}
	return ps.input[ps.i]
}

func (ps parser) reportErr(msg string, a ...interface{}) {
	ps.errs = append(ps.errs, fmt.Errorf(msg, a))
}

func (ps parser) assertTType(t Token, expected tokentypetype) bool {
	if t.ttype != expected {
		ps.reportErr("Unexpected %s `%s' at column %d; expected %s", t.ttype.desc, t.content, t.column, expected.desc)
		return false
	}
	return true
}

func (ps parser) parseExpr() Expr {
	token := ps.nextToken()
	switch token.ttype {
	case TokenRange:
		return TabRange{token.content}
	case TokenIdent:
		ftnName := token
		expParen := ps.nextToken()
		if !ps.assertTType(expParen, TokenOpenParen) {
			return nil
		}
		nextToken := ps.peekToken()
		var args []Expr
		if nextToken.ttype != TokenCloseParen {
			expr := ps.parseExpr()
			if expr == nil {
				return nil
			}
			args = append(args, expr)
			nextToken = ps.peekToken()
			for nextToken.ttype != TokenCloseParen {
				if !ps.assertTType(nextToken, TokenSeparator) {
					return nil
				}
				ps.nextToken() // consume the separator
				expr = ps.parseExpr()
				if expr == nil {
					return nil
				}
				args = append(args, expr)
				nextToken = ps.peekToken()
			}
		}
		ps.nextToken() // consume next token
		if !ps.assertTType(nextToken, TokenCloseParen) {
			return nil
		}
		return FunctionCall{name: ftnName.content, args: args}
	case TokenOpenParen:
		subexpr := ps.parseExpr()
		nextToken := ps.nextToken()
		if !ps.assertTType(nextToken, TokenCloseParen) {
			return nil
		}
		return subexpr
	default:
		ps.reportErr(fmt.Sprintf("Unexpected token `%s' at column %d", token.content, token.column))
		return nil
	}
}

func (ps parser) Err() error {
	if len(ps.errs) == 0 {
		return nil
	}
	return newCombinedError(ps.errs)
}

func ParseExprTokens(input []Token) (Expr, error) {
	ps := parser{input: input, i: 0}
	e := ps.parseExpr()
	if err := ps.Err(); err != nil {
		return nil, err
	}
	return e, nil
}

func ParseExpr(f io.Reader) (Expr, error) {
	Tokens, err := Tokenize(f)
	if err != nil {
		return nil, err
	}
	node, err := ParseExprTokens(Tokens)
	if err != nil {
		return nil, err
	}
	return node, nil
}
