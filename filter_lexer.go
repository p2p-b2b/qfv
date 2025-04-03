package qfv

import (
	"strings"
	"text/scanner"
)

// Lexer breaks the input string into tokens
type Lexer struct {
	s        scanner.Scanner
	input    string
	inputLen int
	pos      int
	tokens   []Token
}

// NewLexer creates a new lexer
func NewLexer(input string) *Lexer {
	var s scanner.Scanner
	s.Init(strings.NewReader(input))
	// Customize scanner: recognize identifiers, numbers, strings.
	s.Mode = scanner.ScanIdents | scanner.ScanFloats | scanner.ScanStrings
	s.Whitespace = 1<<'\t' | 1<<'\n' | 1<<'\r' | 1<<' ' // Define whitespace chars
	s.Error = func(*scanner.Scanner, string) {}         // Suppress default errors

	l := &Lexer{s: s, input: input, pos: -1} // Start at -1, first Next moves to 0
	l.inputLen = len(input)

	return l
}

// Parse reads all tokens from the scanner and buffers them.
func (l *Lexer) Parse() {
	for {
		scanTok := l.s.Scan()
		pos := l.s.Position
		lit := l.s.TokenText()
		var tok TokenType

		switch scanTok {
		case scanner.EOF:
			tok = TokenEOF
		case scanner.Ident:
			upperLit := strings.ToUpper(lit)
			switch upperLit {
			case "AND":
				tok = TokenOperatorAnd
			case "OR":
				tok = TokenOperatorOr
			case "NOT":
				tok = TokenOperatorNot
			case "LIKE": // , "IN", "BETWEEN", "DISTINCT", "IS NOT NULL", "IS NULL", "NOT LIKE", "NOT IN", "NOT BETWEEN", "NOT DISTINCT"
				tok = TokenOperatorLike
			case "IN":
				tok = TokenOperatorIn
			case "BETWEEN":
				tok = TokenOperatorBetween
			case "DISTINCT":
				tok = TokenOperatorDistinct
			case "IS":
				next := l.Next().Value
				if next == "NOT" {
					tok = TokenOperatorIsNotNull
				} else {
					l.Backup() // Go back to IS
					tok = TokenOperatorIsNull
				}
			case "TRUE", "FALSE", "YES", "NO":
				tok = TokenBoolean
			default:
				tok = TokenIdentifier
			}
		case scanner.Int:
			tok = TokenInt
		case scanner.Float:
			tok = TokenFloat
		case scanner.String: // double quotes string are not supported
			tok = TokenIllegal
		case '\'':
			var sb strings.Builder
			sb.WriteByte(byte(scanTok)) // Write the opening quote
			var invalid bool
			numQuotes := 1 // because we already have one

			for {
				char := l.s.Next()

				if char == scanner.EOF {
					break
				} else if char == '\\' {
					sb.WriteByte(byte(char)) // Write the escape character
					char = l.s.Next()        // Consume the escaped character
					sb.WriteByte(byte(char)) // Write the escaped character
					char = l.s.Next()        // Consume the escaped character
				} else if char == '\'' {
					if l.s.Peek() == ' ' || l.s.Peek() == ')' || l.s.Peek() == scanner.EOF {
						sb.WriteByte(byte(char)) // Write the closing quote
						numQuotes++
						break
					}
				}

				sb.WriteRune(char)
			}

			if numQuotes%2 != 0 {
				// Even number of quotes means the string is not closed
				invalid = true
			}

			if invalid {
				tok = TokenIllegal
			} else {
				tok = TokenString
			}

			lit = sb.String()

		case '(':
			tok = TokenLPAREN
		case ')':
			tok = TokenRPAREN
		case ',':
			tok = TokenComma
		case '=':
			tok = TokenOperatorEqual
		case '+':
		case '<':
			if l.s.Peek() == '=' {
				l.s.Scan()
				tok = TokenOperatorLessThanOrEqualTo
			} else if l.s.Peek() == '>' {
				l.s.Scan()
				tok = TokenOperatorNotEqual
			} else {
				tok = TokenOperatorLessThan
			}
		case '>':
			if l.s.Peek() == '=' {
				l.s.Scan()
				tok = TokenOperatorGreaterThanOrEqualTo
			} else {
				tok = TokenOperatorGreaterThan
			}
		case '!':
			if l.s.Peek() == '=' {
				l.s.Scan()
				tok = TokenOperatorNotEqualAlias
				lit = "!="
			} else {
				// Can ! be unary NOT? Let's assume keyword NOT for that.
				// If ! is encountered alone, treat as ILLEGAL or assign a specific token if needed.
				tok = TokenIllegal
				lit = "!" // Keep literal for error message
			}
		default:
			// Handle other single characters if necessary
			tok = TokenIllegal
			lit = string(scanTok) // Store the problematic character
		}

		l.tokens = append(l.tokens, Token{Pos: pos, Type: tok, Value: lit})

		if tok == TokenEOF {
			break
		}
	}
}

// Peek returns the next token without consuming it.
func (l *Lexer) Peek() Token {
	if l.pos+1 >= len(l.tokens) {
		return l.tokens[len(l.tokens)-1] // Return EOF
	}

	return l.tokens[l.pos+1]
}

// Next consumes and returns the next token.
func (l *Lexer) Next() Token {
	l.pos++
	if l.pos >= len(l.tokens) {
		return l.tokens[len(l.tokens)-1] // Return EOF repeatedly
	}

	return l.tokens[l.pos]
}

// Backup goes back one token. Useful for some parsing patterns.
func (l *Lexer) Backup() {
	if l.pos > -1 {
		l.pos--
	}
}

// Current returns the last token returned by Next().
func (l *Lexer) Current() Token {
	if l.pos < 0 || l.pos >= len(l.tokens) {
		// Return an initial dummy token or EOF if out of bounds
		if len(l.tokens) > 0 {
			return l.tokens[len(l.tokens)-1]
		} // EOF

		return Token{Type: TokenIllegal} // Should not happen if lexer ran
	}

	return l.tokens[l.pos]
}
