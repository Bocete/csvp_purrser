package func_purrser

import (
	"github.com/stretchr/testify/assert"
	"strings"
	"testing"
)

func tokensForString(s string) ([]Token, error) {
	reader := strings.NewReader(s)
	return Tokenize(reader)
}

func TestEmptyString(t *testing.T) {
	tokens, err := tokensForString("")
	assert.NoError(t, err)
	assert.Equal(t, 0, len(tokens))
}

func TestRange(t *testing.T) {
	tokens, err := tokensForString("A2:A")
	assert.NoError(t, err)
	assert.Equal(t, 1, len(tokens))
	assert.Equal(t, Token{ttype: TokenRange, content: "A2:A"}, tokens[0])
}

func TestFunctionCall(t *testing.T) {
	tokens, err := tokensForString("Max(A2:B )")
	assert.NoError(t, err)
	assert.Equal(t, 4, len(tokens))
	assert.Equal(t, Token{ttype: TokenIdent, content: "Max"}, tokens[0])
	assert.Equal(t, Token{ttype: TokenOpenParen}, tokens[1])
	assert.Equal(t, Token{ttype: TokenRange, content: "A2:B"}, tokens[2])
	assert.Equal(t, Token{ttype: TokenCloseParen}, tokens[3])
}

func TestWhitespaceBreaksFtnName(t *testing.T) {
	tokens, err := tokensForString("Max Max()")
	assert.NoError(t, err)
	assert.Equal(t, 4, len(tokens))
	assert.Equal(t, Token{ttype: TokenIdent, content: "Max"}, tokens[0])
	assert.Equal(t, Token{ttype: TokenIdent, content: "Max"}, tokens[1])
	assert.Equal(t, Token{ttype: TokenOpenParen}, tokens[2])
	assert.Equal(t, Token{ttype: TokenCloseParen}, tokens[3])
}

func TestFunctionArgs(t *testing.T) {
	tokens, err := tokensForString("Max(A2, B2)")
	assert.NoError(t, err)
	assert.Equal(t, 6, len(tokens))
	assert.Equal(t, Token{ttype: TokenIdent, content: "Max"}, tokens[0])
	assert.Equal(t, Token{ttype: TokenOpenParen}, tokens[1])
	assert.Equal(t, Token{ttype: TokenRange, content: "A2"}, tokens[2])
	assert.Equal(t, Token{ttype: TokenSeparator}, tokens[3])
	assert.Equal(t, Token{ttype: TokenRange, content: "B2"}, tokens[4])
	assert.Equal(t, Token{ttype: TokenCloseParen}, tokens[5])
}
