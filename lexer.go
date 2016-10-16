package func_purrser

import (
	"bufio"
	"bytes"
	"fmt"
	"io"
	"regexp"
	"unicode"
)

type tokentypetype int

const (
	TokenIdent tokentypetype = iota
	TokenRange
	TokenOpenParen
	TokenCloseParen
	TokenSeparator
)

type Token struct {
	ttype   tokentypetype
	content string
}

const rangePattern = `\A[a-zA-Z]{1,2}\d+(?:\:(?:\d+|[a-zA-Z]{1,2}\d*))?\z`

type runeSelector func(rune) bool

func readall(input *bufio.Reader, s runeSelector, init rune, buf *bytes.Buffer) error {
	_, err := buf.WriteRune(init)
	if err != nil {
		return err
	}
	c, _, err := input.ReadRune()
	if err == io.EOF {
		return nil
	} else if err != nil {
		return err
	}
	for s(c) {
		_, err = buf.WriteRune(c)
		if err != nil {
			return err
		}
		c, _, err = input.ReadRune()
		if err == io.EOF {
			return nil
		} else if err != nil {
			return err
		}
	}
	err = input.UnreadRune()
	return err
}

func identBodySelector(c rune) bool {
	return unicode.IsLetter(c) || unicode.IsDigit(c) || c == ':'
}

func Tokenize(f io.Reader) ([]Token, error) {
	reader := bufio.NewReader(f)
	identBuffer := bytes.NewBufferString("")
	var tokens []Token
	for {
		c, _, err := reader.ReadRune()
		if err == io.EOF {
			return tokens, nil
		} else if err != nil {
			return nil, err
		}
		switch {
		case unicode.IsLetter(c):
			identBuffer.Reset()
			err = readall(reader, identBodySelector, c, identBuffer)
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
			tokens = append(tokens, Token{ttype: tokenType, content: identBuffer.String()})
		case c == '(':
			tokens = append(tokens, Token{ttype: TokenOpenParen})
		case c == ')':
			tokens = append(tokens, Token{ttype: TokenCloseParen})
		case c == ',' || c == ';':
			tokens = append(tokens, Token{ttype: TokenSeparator})
		case unicode.IsSpace(c):
		default:
			return nil, fmt.Errorf("Character not recognized: %c", c)
		}
	}
	return tokens, nil
}
