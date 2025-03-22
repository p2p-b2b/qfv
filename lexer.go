package qfv

import (
	"errors"
	"fmt"
	"strings"
	"unicode"
)

// TokenType represents the type of token
// It is used to identify the type of token during parsing
// and to provide type safety when working with different token types.
// Each token type corresponds to a specific part of the query
// (e.g., identifier, operator, string, number, etc.).
type TokenType string

const (
	// Identifier represents a variable name or keyword
	TokenIdentifier TokenType = "IDENTIFIER"

	// Operator represents an operator (e.g., =, <, >, etc.)
	TokenOperator TokenType = "OPERATOR"

	// String represents a string literal
	TokenString TokenType = "STRING"

	// Number represents a numeric literal
	TokenNumber TokenType = "NUMBER"

	// Boolean represents a boolean literal (true/false)
	TokenBoolean TokenType = "BOOLEAN"

	// LogicalOperation represents a logical operation (AND, OR, NOT)
	TokenLogicalOperation TokenType = "LOGICAL_OPERATION"

	// SortOperation represents a sort operation (ASC, DESC)
	TokenSortOperation TokenType = "SORT_OPERATION"

	// Parenthesis represents a parenthesis ((), ))
	TokenLPAREN TokenType = "LPAREN"
	TokenRPAREN TokenType = "RPAREN"
)

func (t TokenType) String() string {
	return string(t)
}

// Token represents a lexical token
type Token struct {
	Type  TokenType
	Value string
}

// Lexer breaks the input string into tokens
type Lexer struct {
	input    string
	inputLen int
	pos      int
	tokens   []Token
}

// newLexer creates a new lexer
func newLexer(input string) *Lexer {
	return &Lexer{
		input:    input,
		inputLen: len(input),
		pos:      0,
		tokens:   []Token{},
	}
}

// tokenize breaks the input into tokens
func (l *Lexer) tokenize() ([]Token, error) {
	for l.pos < l.inputLen {
		char := l.currentChar()

		switch {
		case unicode.IsSpace(char):
			l.skipWhitespace()
		case unicode.IsLetter(char):
			l.tokenizeIdentifier()
		case unicode.IsDigit(char) || char == '-':
			err := l.tokenizeNumber()
			if err != nil {
				return nil, err
			}
		case char == '\'': // String literal
			err := l.tokenizeString()
			if err != nil {
				return nil, err
			}
		case char == '=':
			l.tokens = append(l.tokens, Token{Type: TokenOperator, Value: OperatorEqual.String()})
			l.pos++
		case char == '!' && l.peek() == '=':
			l.tokens = append(l.tokens, Token{Type: TokenOperator, Value: OperatorNotEqualAlias.String()})
			l.pos += 2
		case char == '<':
			if l.peek() == '=' {
				l.tokens = append(l.tokens, Token{Type: TokenOperator, Value: OperatorLessThanOrEqualTo.String()})
				l.pos += 2
			} else {
				l.tokens = append(l.tokens, Token{Type: TokenOperator, Value: OperatorLessThan.String()})
				l.pos++
			}
		case char == '>':
			if l.peek() == '=' {
				l.tokens = append(l.tokens, Token{Type: TokenOperator, Value: OperatorGreaterThanOrEqualTo.String()})
				l.pos += 2
			} else {
				l.tokens = append(l.tokens, Token{Type: TokenOperator, Value: OperatorGreaterThan.String()})
				l.pos++
			}
		case char == '(':
			l.tokens = append(l.tokens, Token{Type: TokenLPAREN, Value: "("})
			l.pos++
		case char == ')':
			l.tokens = append(l.tokens, Token{Type: TokenRPAREN, Value: ")"})
			l.pos++

		default:
			return nil, fmt.Errorf("unexpected character: %c at position %d", char, l.pos)
		}
	}

	return l.tokens, nil
}

// currentChar returns the current character
func (l *Lexer) currentChar() rune {
	if l.pos >= l.inputLen {
		return 0
	}

	return rune(l.input[l.pos])
}

// peek returns the next character without advancing
func (l *Lexer) peek() rune {
	if l.pos+1 >= l.inputLen {
		return 0
	}
	return rune(l.input[l.pos+1])
}

// skipWhitespace skips all whitespace
func (l *Lexer) skipWhitespace() {
	for l.pos < l.inputLen && unicode.IsSpace(rune(l.input[l.pos])) {
		l.pos++
	}
}

// tokenizeIdentifier tokenizes an identifier or keyword
func (l *Lexer) tokenizeIdentifier() {
	start := l.pos

	// until current character is  a letter, digit, underscore or dot
	// we keep moving forward
	for l.pos < l.inputLen &&
		(unicode.IsLetter(rune(l.input[l.pos])) ||
			unicode.IsDigit(rune(l.input[l.pos])) ||
			rune(l.input[l.pos]) == '_' ||
			rune(l.input[l.pos]) == '.') {
		l.pos++
	}

	value := l.input[start:l.pos]
	upperValue := strings.ToUpper(value)

	// Check if it's a logical operator
	switch upperValue {
	// single word operators
	case OperatorAnd.String(), OperatorOr.String(), OperatorNot.String():
		l.tokens = append(l.tokens, Token{Type: TokenLogicalOperation, Value: upperValue})
	case SortAsc.String(), SortDesc.String():
		l.tokens = append(l.tokens, Token{Type: TokenSortOperation, Value: upperValue})
	case OperatorLike.String():
		l.tokens = append(l.tokens, Token{Type: TokenOperator, Value: OperatorLike.String()})
	case OperatorIn.String():
		l.tokens = append(l.tokens, Token{Type: TokenOperator, Value: OperatorIn.String()})
	case OperatorBetween.String():
		l.tokens = append(l.tokens, Token{Type: TokenOperator, Value: OperatorBetween.String()})
	case OperatorDistinct.String():
		l.tokens = append(l.tokens, Token{Type: TokenOperator, Value: OperatorDistinct.String()})
	case OperatorIsNotNull.String():
		l.tokens = append(l.tokens, Token{Type: TokenOperator, Value: OperatorIsNotNull.String()})
	case OperatorIsNull.String():
		l.tokens = append(l.tokens, Token{Type: TokenOperator, Value: OperatorIsNull.String()})
	case OperatorNotLike.String():
		l.tokens = append(l.tokens, Token{Type: TokenOperator, Value: OperatorNotLike.String()})
	case OperatorNotIn.String():
		l.tokens = append(l.tokens, Token{Type: TokenOperator, Value: OperatorNotIn.String()})
	case OperatorNotBetween.String():
		l.tokens = append(l.tokens, Token{Type: TokenOperator, Value: OperatorNotBetween.String()})
	case OperatorNotDistinct.String():
		l.tokens = append(l.tokens, Token{Type: TokenOperator, Value: OperatorNotDistinct.String()})
	case "TRUE", "FALSE", "true", "false", "YES", "NO", "yes", "no":
		l.tokens = append(l.tokens, Token{Type: TokenBoolean, Value: upperValue})
	default:
		l.tokens = append(l.tokens, Token{Type: TokenIdentifier, Value: value})
	}
}

// tokenizeNumber tokenizes a number literal
// It supports integers and floating-point numbers
// It also supports numbers with a single dot (e.g., 3.14, 2.0)
// but not multiple dots (e.g., 3.14.15)
// It also supports numbers with leading zeros (e.g., 0123)
// It does not support scientific notation (e.g., 1e10)
// It does not support hexadecimal or octal numbers
// It does not support numbers with underscores (e.g., 1_000)
// It does not support numbers with commas (e.g., 1,000)
// It does not support numbers with currency symbols (e.g., $1.00)
func (l *Lexer) tokenizeNumber() error {
	start := l.pos

	numDots := 0
	for l.pos < l.inputLen &&
		(unicode.IsDigit(rune(l.input[l.pos])) || (rune(l.input[l.pos]) == '.' || (rune(l.input[l.pos]) == '-'))) {
		if rune(l.input[l.pos]) == '.' {
			numDots++
		}

		l.pos++
	}

	if numDots > 1 {
		return fmt.Errorf("invalid number format: %s", l.input[start:l.pos])
	}

	if numDots == 1 {
		// Check if the input is only a dot
		if start == l.inputLen-1 && l.input[start] == '.' {
			return fmt.Errorf("invalid number format: %s", l.input[start:l.pos])
		}

		// Check if the last character is a dot
		if l.pos == l.inputLen && l.input[l.pos-1] == '.' {
			return fmt.Errorf("invalid number format: %s", l.input[start:l.pos])
		}

		// Check if the first character is a dot
		if l.input[start] == '.' && l.pos-start > 1 {
			return fmt.Errorf("invalid number format: %s", l.input[start:l.pos])
		}
	}

	// Check if the number has leading zeros
	if l.input[start] == '0' {
		if l.pos-start > 1 && l.input[start+1] != '.' {
			return fmt.Errorf("invalid number format: %s", l.input[start:l.pos])
		}
	}

	// Check if is a negative number
	if l.input[start] == '-' {
		// check if is only a negative sign
		if start == l.inputLen-1 {
			return fmt.Errorf("invalid number format: %s", l.input[start:l.pos])
		}

		// check if is a negative number with leading zeros
		if l.input[start+1] == '0' {
			if l.pos-start > 2 && l.input[start+2] != '.' {
				return fmt.Errorf("invalid number format: %s", l.input[start:l.pos])
			}
		}

		// check if is a negative number with a dot
		if l.input[start+1] == '.' {
			if l.pos-start > 2 {
				return fmt.Errorf("invalid number format: %s", l.input[start:l.pos])
			}
		}
	}

	value := l.input[start:l.pos]
	l.tokens = append(l.tokens, Token{Type: TokenNumber, Value: value})

	return nil
}

// tokenizeString tokenizes a string literal
func (l *Lexer) tokenizeString() error {
	l.pos++ // Skip opening quote

	start := l.pos

	for l.pos < l.inputLen && rune(l.input[l.pos]) != '\'' {
		// Handle escaped quotes
		if rune(l.input[l.pos]) == '\\' && l.pos+1 < l.inputLen && rune(l.input[l.pos+1]) == '\'' {
			l.pos += 2
			continue
		}

		l.pos++
	}

	if l.pos >= l.inputLen {
		return errors.New("unterminated string literal")
	}

	value := l.input[start:l.pos]
	l.tokens = append(l.tokens, Token{Type: TokenString, Value: value})

	l.pos++ // Skip closing quote

	return nil
}
