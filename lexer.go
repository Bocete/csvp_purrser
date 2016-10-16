package func_purrser

import (
	"bufio"
	"bytes"
	"fmt"
	"io"
	"regexp"
	"unicode"
)

type tokentypetype struct {
	id   int
	desc string
}

var (
	TokenIdent      = tokentypetype{1, "identifier"}
	TokenRange      = tokentypetype{2, "table range"}
	TokenOpenParen  = tokentypetype{3, "open parenthesis"}
	TokenCloseParen = tokentypetype{4, "close parenthesis"}
	TokenSeparator  = tokentypetype{5, "argument separator"}
	TokenEOF        = tokentypetype{6, "end of input"}
)

type Token struct {
	ttype   tokentypetype
	column  int
	content string
}

const rangePattern = `\A[a-zA-Z]{1,2}\d+(?:\:(?:\d+|[a-zA-Z]{1,2}\d*))?\z`

type runeSelector func(rune) bool

func readall(input *bufio.Reader, s runeSelector, init rune, buf *bytes.Buffer) (int, error) {
	_, err := buf.WriteRune(init)
	if err != nil {
		return 0, err
	}
	c, _, err := input.ReadRune()
	if err == io.EOF {
		return 0, nil
	} else if err != nil {
		return 0, err
	}
	offset := 0
	for s(c) {
		_, err = buf.WriteRune(c)
		if err != nil {
			return offset, err
		}
		c, _, err = input.ReadRune()
		offset = offset + 1
		if err == io.EOF {
			return offset, nil
		} else if err != nil {
			return offset, err
		}
	}
	err = input.UnreadRune()
	return offset, err
}

func identBodySelector(c rune) bool {
	return unicode.IsLetter(c) || unicode.IsDigit(c) || c == ':'
}

func Tokenize(f io.Reader) ([]Token, error) {
	reader := bufio.NewReader(f)
	identBuffer := bytes.NewBufferString("")
	column := 0
	var tokens []Token
	for {
		c, _, err := reader.ReadRune()
		if err == io.EOF {
			tokens = append(tokens, Token{ttype: TokenEOF, content: "EOF", column: column})
			return tokens, nil
		} else if err != nil {
			return nil, err
		}
		switch {
		case unicode.IsLetter(c):
			identBuffer.Reset()
			offset, err := readall(reader, identBodySelector, c, identBuffer)
			if err != nil {
				return nil, err
			}
			content := identBuffer.String()
			identIsRange, err := regexp.MatchString(rangePattern, content)
			if err != nil {
				return nil, err
			}
			var tokenType tokentypetype
			if identIsRange {
				tokenType = TokenRange
			} else {
				tokenType = TokenIdent
			}
			tokens = append(tokens, Token{
				ttype:   tokenType,
				content: identBuffer.String(),
				column:  column,
			})
			column = column + offset
		case c == '(':
			tokens = append(tokens, Token{
				ttype:   TokenOpenParen,
				content: "(",
				column:  column,
			})
		case c == ')':
			tokens = append(tokens, Token{
				ttype:   TokenCloseParen,
				content: ")",
				column:  column,
			})
		case c == ',' || c == ';':
			tokens = append(tokens, Token{
				ttype:   TokenSeparator,
				content: string(c),
				column:  column,
			})
		case unicode.IsSpace(c):
		default:
			return nil, fmt.Errorf("Character %c not recognized at column %d", c, column)
		}
		column = column + 1
	}
	return tokens, nil
}
