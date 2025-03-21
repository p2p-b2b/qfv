package qfv

import (
	"errors"
	"fmt"
	"strings"
	"unicode"
)

// Token represents a lexical token
type Token struct {
	Type  string
	Value string
}

// Lexer breaks the input string into tokens
type Lexer struct {
	input  string
	pos    int
	tokens []Token
}

// newLexer creates a new lexer
func newLexer(input string) *Lexer {
	return &Lexer{
		input:  input,
		pos:    0,
		tokens: []Token{},
	}
}

// tokenize breaks the input into tokens
func (l *Lexer) tokenize() ([]Token, error) {
	for l.pos < len(l.input) {
		char := l.currentChar()

		switch {
		case unicode.IsSpace(char):
			l.skipWhitespace()
		case unicode.IsLetter(char):
			l.tokenizeIdentifier()
		case char == '\'':
			err := l.tokenizeString()
			if err != nil {
				return nil, err
			}
		case char == '=':
			l.tokens = append(l.tokens, Token{Type: "OPERATOR", Value: "="})
			l.pos++
		case char == '!' && l.peek() == '=':
			l.tokens = append(l.tokens, Token{Type: "OPERATOR", Value: "!="})
			l.pos += 2
		case char == '<':
			if l.peek() == '=' {
				l.tokens = append(l.tokens, Token{Type: "OPERATOR", Value: "<="})
				l.pos += 2
			} else {
				l.tokens = append(l.tokens, Token{Type: "OPERATOR", Value: "<"})
				l.pos++
			}
		case char == '>':
			if l.peek() == '=' {
				l.tokens = append(l.tokens, Token{Type: "OPERATOR", Value: ">="})
				l.pos += 2
			} else {
				l.tokens = append(l.tokens, Token{Type: "OPERATOR", Value: ">"})
				l.pos++
			}
		case char == '(':
			l.tokens = append(l.tokens, Token{Type: "LPAREN", Value: "("})
			l.pos++
		case char == ')':
			l.tokens = append(l.tokens, Token{Type: "RPAREN", Value: ")"})
			l.pos++

		default:
			return nil, fmt.Errorf("unexpected character: %c at position %d", char, l.pos)
		}
	}

	return l.tokens, nil
}

// currentChar returns the current character
func (l *Lexer) currentChar() rune {
	if l.pos >= len(l.input) {
		return 0
	}
	return rune(l.input[l.pos])
}

// peek returns the next character without advancing
func (l *Lexer) peek() rune {
	if l.pos+1 >= len(l.input) {
		return 0
	}
	return rune(l.input[l.pos+1])
}

// skipWhitespace skips all whitespace
func (l *Lexer) skipWhitespace() {
	for l.pos < len(l.input) && unicode.IsSpace(rune(l.input[l.pos])) {
		l.pos++
	}
}

// tokenizeIdentifier tokenizes an identifier or keyword
func (l *Lexer) tokenizeIdentifier() {
	start := l.pos
	for l.pos < len(l.input) && (unicode.IsLetter(rune(l.input[l.pos])) || unicode.IsDigit(rune(l.input[l.pos])) || rune(l.input[l.pos]) == '_') {
		l.pos++
	}

	value := l.input[start:l.pos]
	upper := strings.ToUpper(value)

	// Check if it's a logical operator
	switch upper {
	case "AND", "OR", "NOT":
		l.tokens = append(l.tokens, Token{Type: "LOGICAL_OP", Value: upper})
	case "ASC", "DESC":
		l.tokens = append(l.tokens, Token{Type: "SORT_DIR", Value: upper})
	default:
		l.tokens = append(l.tokens, Token{Type: "IDENTIFIER", Value: value})
	}
}

// tokenizeString tokenizes a string literal
func (l *Lexer) tokenizeString() error {
	l.pos++ // Skip opening quote
	start := l.pos

	for l.pos < len(l.input) && rune(l.input[l.pos]) != '\'' {
		// Handle escaped quotes
		if rune(l.input[l.pos]) == '\\' && l.pos+1 < len(l.input) && rune(l.input[l.pos+1]) == '\'' {
			l.pos += 2
			continue
		}
		l.pos++
	}

	if l.pos >= len(l.input) {
		return errors.New("unterminated string literal")
	}

	value := l.input[start:l.pos]
	l.tokens = append(l.tokens, Token{Type: "STRING", Value: value})
	l.pos++ // Skip closing quote
	return nil
}
