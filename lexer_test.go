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
	assert.Equal(t, 1, len(tokens))
	assert.Equal(t, Token{ttype: TokenEOF, content: "EOF", column: 0}, tokens[0])
}

func TestRange(t *testing.T) {
	tokens, err := tokensForString("A2:A")
	assert.NoError(t, err)
	assert.Equal(t, 2, len(tokens))
	assert.Equal(t, Token{ttype: TokenRange, content: "A2:A", column: 0}, tokens[0])
	assert.Equal(t, Token{ttype: TokenEOF, content: "EOF", column: 4}, tokens[1])
}

func TestFunctionCall(t *testing.T) {
	tokens, err := tokensForString("Max(A2:B )")
	assert.NoError(t, err)
	assert.Equal(t, 5, len(tokens))
	assert.Equal(t, Token{ttype: TokenIdent, content: "Max", column: 0}, tokens[0])
	assert.Equal(t, Token{ttype: TokenOpenParen, content: "(", column: 3}, tokens[1])
	assert.Equal(t, Token{ttype: TokenRange, content: "A2:B", column: 4}, tokens[2])
	assert.Equal(t, Token{ttype: TokenCloseParen, content: ")", column: 9}, tokens[3])
	assert.Equal(t, Token{ttype: TokenEOF, content: "EOF", column: 10}, tokens[4])
}

func TestWhitespaceBreaksFtnName(t *testing.T) {
	tokens, err := tokensForString("Max Max()")
	assert.NoError(t, err)
	assert.Equal(t, 5, len(tokens))
	assert.Equal(t, Token{ttype: TokenIdent, content: "Max", column: 0}, tokens[0])
	assert.Equal(t, Token{ttype: TokenIdent, content: "Max", column: 4}, tokens[1])
	assert.Equal(t, Token{ttype: TokenOpenParen, content: "(", column: 7}, tokens[2])
	assert.Equal(t, Token{ttype: TokenCloseParen, content: ")", column: 8}, tokens[3])
	assert.Equal(t, Token{ttype: TokenEOF, content: "EOF", column: 9}, tokens[4])
}

func TestFunctionArgs(t *testing.T) {
	tokens, err := tokensForString("Max(A2, B2)")
	assert.NoError(t, err)
	assert.Equal(t, 7, len(tokens))
	assert.Equal(t, Token{ttype: TokenIdent, content: "Max", column: 0}, tokens[0])
	assert.Equal(t, Token{ttype: TokenOpenParen, content: "(", column: 3}, tokens[1])
	assert.Equal(t, Token{ttype: TokenRange, content: "A2", column: 4}, tokens[2])
	assert.Equal(t, Token{ttype: TokenSeparator, content: ",", column: 6}, tokens[3])
	assert.Equal(t, Token{ttype: TokenRange, content: "B2", column: 8}, tokens[4])
	assert.Equal(t, Token{ttype: TokenCloseParen, content: ")", column: 10}, tokens[5])
	assert.Equal(t, Token{ttype: TokenEOF, content: "EOF", column: 11}, tokens[6])
}
