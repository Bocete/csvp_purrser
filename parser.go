package func_purrser

import (
	"io"
)

type parserState struct {
	errs  []error
	input []Token
	i     int
}

func (ps parserState) Err() error {
	if len(ps.errs) == 0 {
		return nil
	}
	return newCombinedError(ps.errs)
}

func ParseExprTokens(input []Token) (*Node, error) {
	ps := parserState{input: input, i: 0}
	return nil, nil
}

func ParseExpr(f io.Reader) (*Node, error) {
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
