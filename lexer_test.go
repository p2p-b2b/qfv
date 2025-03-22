package qfv

import (
	"reflect"
	"testing"
)

func TestLexer_tokenizeIdentifier(t *testing.T) {
	tests := []struct {
		name   string
		input  string
		tokens []Token
	}{
		{
			name:  "simple identifier",
			input: "name",
			tokens: []Token{
				{Type: TokenIdentifier, Value: "name"},
			},
		},
		{
			name:  "identifier with underscore",
			input: "user_name",
			tokens: []Token{
				{Type: TokenIdentifier, Value: "user_name"},
			},
		},
		{
			name:  "identifier with digits",
			input: "product123",
			tokens: []Token{
				{Type: TokenIdentifier, Value: "product123"},
			},
		},
		{
			name:  "identifier with dot",
			input: "address.city",
			tokens: []Token{
				{Type: TokenIdentifier, Value: "address.city"},
			},
		},
		{
			name:  "logical AND operator",
			input: "AND",
			tokens: []Token{
				{Type: TokenLogicalOperation, Value: "AND"},
			},
		},
		{
			name:  "logical OR operator",
			input: "OR",
			tokens: []Token{
				{Type: TokenLogicalOperation, Value: "OR"},
			},
		},
		{
			name:  "logical NOT operator",
			input: "NOT",
			tokens: []Token{
				{Type: TokenLogicalOperation, Value: "NOT"},
			},
		},
		{
			name:  "sort ASC operator",
			input: "ASC",
			tokens: []Token{
				{Type: TokenSortOperation, Value: "ASC"},
			},
		},
		{
			name:  "sort DESC operator",
			input: "DESC",
			tokens: []Token{
				{Type: TokenSortOperation, Value: "DESC"},
			},
		},
		{
			name:  "LIKE operator",
			input: "LIKE",
			tokens: []Token{
				{Type: TokenOperator, Value: "LIKE"},
			},
		},
		{
			name:  "IN operator",
			input: "IN",
			tokens: []Token{
				{Type: TokenOperator, Value: "IN"},
			},
		},
		{
			name:  "BETWEEN operator",
			input: "BETWEEN",
			tokens: []Token{
				{Type: TokenOperator, Value: "BETWEEN"},
			},
		},
		{
			name:  "DISTINCT operator",
			input: "DISTINCT",
			tokens: []Token{
				{Type: TokenOperator, Value: "DISTINCT"},
			},
		},
		{
			name:  "mixed identifier and operator",
			input: "nameAND",
			tokens: []Token{
				{Type: TokenIdentifier, Value: "nameAND"},
			},
		},
		{
			name:  "identifier with special characters",
			input: "username",
			tokens: []Token{
				{Type: TokenIdentifier, Value: "username"},
			},
		},
		{
			name:  "identifier with special characters and digits",
			input: "user123_name",
			tokens: []Token{
				{Type: TokenIdentifier, Value: "user123_name"},
			},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			lexer := newLexer(tt.input)
			lexer.tokenizeIdentifier()

			if !reflect.DeepEqual(lexer.tokens, tt.tokens) {
				t.Errorf("expected '%v', got '%v'", tt.tokens, lexer.tokens)
			}
		})
	}
}

func TestLexer_tokenizeNumber(t *testing.T) {
	tests := []struct {
		name    string
		input   string
		tokens  []Token
		wantErr bool
	}{
		{
			name:  "simple integer",
			input: "123",
			tokens: []Token{
				{Type: TokenNumber, Value: "123"},
			},
			wantErr: false,
		},
		{
			name:  "simple negative integer",
			input: "-123",
			tokens: []Token{
				{Type: TokenNumber, Value: "-123"},
			},
			wantErr: false,
		},
		{
			name:  "simple negative float",
			input: "-123.45",
			tokens: []Token{
				{Type: TokenNumber, Value: "-123.45"},
			},
			wantErr: false,
		},
		{
			name:  "simple negative float with leading zero",
			input: "-0.12345",
			tokens: []Token{
				{Type: TokenNumber, Value: "-0.12345"},
			},
			wantErr: false,
		},
		{
			name:    "simple negative float with leading zero and no digits",
			input:   "-.12345",
			tokens:  []Token{},
			wantErr: true,
		},
		{
			name:  "decimal number",
			input: "3.14",
			tokens: []Token{
				{Type: TokenNumber, Value: "3.14"},
			},
			wantErr: false,
		},
		{
			name:  "number with leading zero",
			input: "0.5",
			tokens: []Token{
				{Type: TokenNumber, Value: "0.5"},
			},
			wantErr: false,
		},
		{
			name:    "number with two leading zero",
			input:   "00.5",
			tokens:  []Token{},
			wantErr: true,
		},
		{
			name:    "multiple dots",
			input:   "1.2.3",
			tokens:  []Token{},
			wantErr: true,
		},
		{
			name:    "just a dot",
			input:   ".",
			tokens:  []Token{},
			wantErr: true,
		},
		{
			name:    "dot at the end",
			input:   "123.",
			tokens:  []Token{},
			wantErr: true,
		},
		{
			name:    "dot at the start",
			input:   ".123",
			tokens:  []Token{},
			wantErr: true,
		},
		{
			name:    "leading zeros",
			input:   "0123",
			tokens:  []Token{},
			wantErr: true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			lexer := newLexer(tt.input)
			err := lexer.tokenizeNumber()
			if (err != nil) != tt.wantErr {
				t.Errorf("tokenizeNumber() error = %v, wantErr %v", err, tt.wantErr)
				return
			}

			if tt.wantErr {
				if len(lexer.tokens) != 0 {
					t.Errorf("expected no tokens, got '%v'", lexer.tokens)
				}
				return
			}

			if !reflect.DeepEqual(lexer.tokens, tt.tokens) {
				t.Errorf("expected '%v', got '%v'", tt.tokens, lexer.tokens)
			}
		})
	}
}

func TestLexer_skipWhitespace(t *testing.T) {
	tests := []struct {
		name           string
		input          string
		initialPos     int
		expectedPos    int
		expectedTokens []Token
	}{
		{
			name:           "no whitespace",
			input:          "name",
			initialPos:     0,
			expectedPos:    0,
			expectedTokens: []Token{},
		},
		{
			name:           "single space",
			input:          " name",
			initialPos:     0,
			expectedPos:    1,
			expectedTokens: []Token{},
		},
		{
			name:           "multiple spaces",
			input:          "   name",
			initialPos:     0,
			expectedPos:    3,
			expectedTokens: []Token{},
		},
		{
			name:           "tab",
			input:          "\tname",
			initialPos:     0,
			expectedPos:    1,
			expectedTokens: []Token{},
		},
		{
			name:           "newline",
			input:          "\nname",
			initialPos:     0,
			expectedPos:    1,
			expectedTokens: []Token{},
		},
		{
			name:           "carriage return",
			input:          "\rname",
			initialPos:     0,
			expectedPos:    1,
			expectedTokens: []Token{},
		},
		{
			name:           "mixed whitespace",
			input:          "  \t\n\rname",
			initialPos:     0,
			expectedPos:    5,
			expectedTokens: []Token{},
		},
		{
			name:           "whitespace in the middle",
			input:          "name  value",
			initialPos:     4,
			expectedPos:    6,
			expectedTokens: []Token{},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			lexer := newLexer(tt.input)
			lexer.pos = tt.initialPos
			lexer.skipWhitespace()

			if lexer.pos != tt.expectedPos {
				t.Errorf("expected position '%d', got '%d'", tt.expectedPos, lexer.pos)
			}

			if !reflect.DeepEqual(lexer.tokens, tt.expectedTokens) {
				t.Errorf("expected tokens '%v', got '%v'", tt.expectedTokens, lexer.tokens)
			}
		})
	}
}

func TestLexer_tokenizeString(t *testing.T) {
	tests := []struct {
		name    string
		input   string
		tokens  []Token
		wantErr bool
	}{
		{
			name:  "simple string",
			input: "'name'",
			tokens: []Token{
				{Type: TokenString, Value: "name"},
			},
			wantErr: false,
		},
		{
			name:  "string with spaces",
			input: "'name value'",
			tokens: []Token{
				{Type: TokenString, Value: "name value"},
			},
			wantErr: false,
		},
		{
			name:  "string with escaped single quote",
			input: "'name\\'value'",
			tokens: []Token{
				{Type: TokenString, Value: "name\\'value"},
			},
			wantErr: false,
		},
		{
			name:    "unterminated string",
			input:   "'name",
			tokens:  []Token{},
			wantErr: true,
		},
		{
			name:  "empty string",
			input: "''",
			tokens: []Token{
				{Type: TokenString, Value: ""},
			},
			wantErr: false,
		},
		{
			name:  "string with special characters",
			input: "'name!@#$%^&*()_+=-`~[]\\{}|;\\':\",./<>?'",
			tokens: []Token{
				{Type: TokenString, Value: "name!@#$%^&*()_+=-`~[]\\{}|;\\':\",./<>?"},
			},
			wantErr: false,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			lexer := newLexer(tt.input)
			err := lexer.tokenizeString()

			if (err != nil) != tt.wantErr {
				t.Errorf("tokenizeString() error = %v, wantErr %v", err, tt.wantErr)
				return
			}

			if tt.wantErr {
				if len(lexer.tokens) != 0 {
					t.Errorf("expected no tokens, got '%v'", lexer.tokens)
				}
				return
			}

			if !reflect.DeepEqual(lexer.tokens, tt.tokens) {
				t.Errorf("expected '%v', got '%v'", tt.tokens, lexer.tokens)
			}
		})
	}
}

func TestLexer_tokenize(t *testing.T) {
	tests := []struct {
		name    string
		input   string
		tokens  []Token
		wantErr bool
	}{
		{
			name:    "empty input",
			input:   "",
			tokens:  []Token{},
			wantErr: false,
		},
		{
			name:  "simple identifier",
			input: "name",
			tokens: []Token{
				{Type: TokenIdentifier, Value: "name"},
			},
			wantErr: false,
		},
		{
			name:  "simple string",
			input: "'name'",
			tokens: []Token{
				{Type: TokenString, Value: "name"},
			},
			wantErr: false,
		},
		{
			name:  "simple integer number",
			input: "123",
			tokens: []Token{
				{Type: TokenNumber, Value: "123"},
			},
			wantErr: false,
		},
		{
			name:  "simple negative integer number",
			input: "-123",
			tokens: []Token{
				{Type: TokenNumber, Value: "-123"},
			},
			wantErr: false,
		},
		{
			name:  "simple float number",
			input: "123.45",
			tokens: []Token{
				{Type: TokenNumber, Value: "123.45"},
			},
			wantErr: false,
		},
		{
			name:  "simple float with leading zero",
			input: "0.12345",
			tokens: []Token{
				{Type: TokenNumber, Value: "0.12345"},
			},
			wantErr: false,
		},
		{
			name:  "simple negative float number",
			input: "-123.45",
			tokens: []Token{
				{Type: TokenNumber, Value: "-123.45"},
			},
			wantErr: false,
		},
		{
			name:  "simple float with leading zero",
			input: "-0.123",
			tokens: []Token{
				{Type: TokenNumber, Value: "-0.123"},
			},
			wantErr: false,
		},
		{
			name:  "simple operator",
			input: "=",
			tokens: []Token{
				{Type: TokenOperator, Value: "="},
			},
			wantErr: false,
		},
		{
			name:  "identifier and operator",
			input: "name=value",
			tokens: []Token{
				{Type: TokenIdentifier, Value: "name"},
				{Type: TokenOperator, Value: "="},
				{Type: TokenIdentifier, Value: "value"},
			},
			wantErr: false,
		},
		{
			name:  "string and operator",
			input: "'name'='value'",
			tokens: []Token{
				{Type: TokenString, Value: "name"},
				{Type: TokenOperator, Value: "="},
				{Type: TokenString, Value: "value"},
			},
			wantErr: false,
		},
		{
			name:  "number and operator",
			input: "123=456",
			tokens: []Token{
				{Type: TokenNumber, Value: "123"},
				{Type: TokenOperator, Value: "="},
				{Type: TokenNumber, Value: "456"},
			},
			wantErr: false,
		},
		{
			name:  "float and operator",
			input: "123.45=678.90",
			tokens: []Token{
				{Type: TokenNumber, Value: "123.45"},
				{Type: TokenOperator, Value: "="},
				{Type: TokenNumber, Value: "678.90"},
			},
			wantErr: false,
		},
		{
			name:  "parentheses",
			input: "(name)",
			tokens: []Token{
				{Type: TokenLPAREN, Value: "("},
				{Type: TokenIdentifier, Value: "name"},
				{Type: TokenRPAREN, Value: ")"},
			},
			wantErr: false,
		},
		{
			name:  "complex expression",
			input: "(name='value' AND age>18)",
			tokens: []Token{
				{Type: TokenLPAREN, Value: "("},
				{Type: TokenIdentifier, Value: "name"},
				{Type: TokenOperator, Value: "="},
				{Type: TokenString, Value: "value"},
				{Type: TokenLogicalOperation, Value: "AND"},
				{Type: TokenIdentifier, Value: "age"},
				{Type: TokenOperator, Value: ">"},
				{Type: TokenNumber, Value: "18"},
				{Type: TokenRPAREN, Value: ")"},
			},
			wantErr: false,
		},
		{
			name:  "complex expression with parentheses",
			input: "name != 'test string' OR (name='value' AND (age>18 OR city='New York'))",
			tokens: []Token{
				{Type: TokenIdentifier, Value: "name"},
				{Type: TokenOperator, Value: "!="},
				{Type: TokenString, Value: "test string"},
				{Type: TokenLogicalOperation, Value: "OR"},
				{Type: TokenLPAREN, Value: "("},
				{Type: TokenIdentifier, Value: "name"},
				{Type: TokenOperator, Value: "="},
				{Type: TokenString, Value: "value"},
				{Type: TokenLogicalOperation, Value: "AND"},
				{Type: TokenLPAREN, Value: "("},
				{Type: TokenIdentifier, Value: "age"},
				{Type: TokenOperator, Value: ">"},
				{Type: TokenNumber, Value: "18"},
				{Type: TokenLogicalOperation, Value: "OR"},
				{Type: TokenIdentifier, Value: "city"},
				{Type: TokenOperator, Value: "="},
				{Type: TokenString, Value: "New York"},
				{Type: TokenRPAREN, Value: ")"},
				{Type: TokenRPAREN, Value: ")"},
			},
			wantErr: false,
		},
		{
			name:    "invalid number",
			input:   "1.2.3",
			tokens:  []Token{},
			wantErr: true,
		},
		{
			name:    "unterminated string",
			input:   "'name",
			tokens:  []Token{},
			wantErr: true,
		},
		{
			name:  "not equal operator",
			input: "name!=value",
			tokens: []Token{
				{Type: TokenIdentifier, Value: "name"},
				{Type: TokenOperator, Value: "!="},
				{Type: TokenIdentifier, Value: "value"},
			},
			wantErr: false,
		},
		{
			name:  "greater than operator",
			input: "age>18",
			tokens: []Token{
				{Type: TokenIdentifier, Value: "age"},
				{Type: TokenOperator, Value: ">"},
				{Type: TokenNumber, Value: "18"},
			},
			wantErr: false,
		},
		{
			name:  "less than operator",
			input: "age<18",
			tokens: []Token{
				{Type: TokenIdentifier, Value: "age"},
				{Type: TokenOperator, Value: "<"},
				{Type: TokenNumber, Value: "18"},
			},
			wantErr: false,
		},
		{
			name:  "greater than or equal operator",
			input: "age>=18",
			tokens: []Token{
				{Type: TokenIdentifier, Value: "age"},
				{Type: TokenOperator, Value: ">="},
				{Type: TokenNumber, Value: "18"},
			},
			wantErr: false,
		},
		{
			name:  "less than or equal operator",
			input: "age<=18",
			tokens: []Token{
				{Type: TokenIdentifier, Value: "age"},
				{Type: TokenOperator, Value: "<="},
				{Type: TokenNumber, Value: "18"},
			},
			wantErr: false,
		},
		{
			name:  "AND operator",
			input: "name AND age",
			tokens: []Token{
				{Type: TokenIdentifier, Value: "name"},
				{Type: TokenLogicalOperation, Value: "AND"},
				{Type: TokenIdentifier, Value: "age"},
			},
			wantErr: false,
		},
		{
			name:  "OR operator",
			input: "name OR age",
			tokens: []Token{
				{Type: TokenIdentifier, Value: "name"},
				{Type: TokenLogicalOperation, Value: "OR"},
				{Type: TokenIdentifier, Value: "age"},
			},
			wantErr: false,
		},
		{
			name:  "NOT operator",
			input: "NOT name",
			tokens: []Token{
				{Type: TokenLogicalOperation, Value: "NOT"},
				{Type: TokenIdentifier, Value: "name"},
			},
			wantErr: false,
		},
		{
			name:  "BETWEEN operator",
			input: "age BETWEEN 18 AND 30",
			tokens: []Token{
				{Type: TokenIdentifier, Value: "age"},
				{Type: TokenOperator, Value: "BETWEEN"},
				{Type: TokenNumber, Value: "18"},
				{Type: TokenLogicalOperation, Value: "AND"},
				{Type: TokenNumber, Value: "30"},
			},
			wantErr: false,
		},
		{
			name:  "IN operator",
			input: "age IN (18, 30)",
			tokens: []Token{
				{Type: TokenIdentifier, Value: "age"},
				{Type: TokenOperator, Value: "IN"},
				{Type: TokenLPAREN, Value: "("},
				{Type: TokenNumber, Value: "18"},
				{Type: TokenOperator, Value: ","},
				{Type: TokenNumber, Value: "30"},
				{Type: TokenRPAREN, Value: ")"},
			},
			wantErr: false,
		},
		{
			name:  "IS NULL operator",
			input: "age IS NULL",
			tokens: []Token{
				{Type: TokenIdentifier, Value: "age"},
				{Type: TokenOperator, Value: "IS"},
				{Type: TokenOperator, Value: "NULL"},
			},
			wantErr: false,
		},
		{
			name:  "IS NOT NULL operator",
			input: "age IS NOT NULL",
			tokens: []Token{
				{Type: TokenIdentifier, Value: "age"},
				{Type: TokenOperator, Value: "IS"},
				{Type: TokenOperator, Value: "NOT"},
				{Type: TokenOperator, Value: "NULL"},
			},
			wantErr: false,
		},
		{
			name:  "NOT IN operator",
			input: "age NOT IN (18, 30)",
			tokens: []Token{
				{Type: TokenIdentifier, Value: "age"},
				{Type: TokenOperator, Value: "NOT"},
				{Type: TokenOperator, Value: "IN"},
				{Type: TokenLPAREN, Value: "("},
				{Type: TokenNumber, Value: "18"},
				{Type: TokenOperator, Value: ","},
				{Type: TokenNumber, Value: "30"},
				{Type: TokenRPAREN, Value: ")"},
			},
			wantErr: false,
		},
		{
			name:  "NOT BETWEEN operator",
			input: "age NOT BETWEEN 18 AND 30",
			tokens: []Token{
				{Type: TokenIdentifier, Value: "age"},
				{Type: TokenOperator, Value: "NOT"},
				{Type: TokenOperator, Value: "BETWEEN"},
				{Type: TokenNumber, Value: "18"},
				{Type: TokenLogicalOperation, Value: "AND"},
				{Type: TokenNumber, Value: "30"},
			},
			wantErr: false,
		},
		{
			name:  "DISTINCT operator",
			input: "SELECT DISTINCT name FROM users",
			tokens: []Token{
				{Type: TokenIdentifier, Value: "SELECT"},
				{Type: TokenOperator, Value: "DISTINCT"},
				{Type: TokenIdentifier, Value: "name"},
				{Type: TokenIdentifier, Value: "FROM"},
				{Type: TokenIdentifier, Value: "users"},
			},
			wantErr: false,
		},
		{
			name:  "NOT DISTINCT operator",
			input: "SELECT NOT DISTINCT name FROM users",
			tokens: []Token{
				{Type: TokenIdentifier, Value: "SELECT"},
				{Type: TokenOperator, Value: "NOT"},
				{Type: TokenOperator, Value: "DISTINCT"},
				{Type: TokenIdentifier, Value: "name"},
				{Type: TokenIdentifier, Value: "FROM"},
				{Type: TokenIdentifier, Value: "users"},
			},
			wantErr: false,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			lexer := newLexer(tt.input)
			tokens, err := lexer.tokenize()

			if (err != nil) != tt.wantErr {
				t.Errorf("tokenize() error = %v, wantErr %v", err, tt.wantErr)
				return
			}

			if tt.wantErr {
				if len(tokens) != 0 {
					t.Errorf("expected no tokens, got '%v'", tokens)
				}
				return
			}

			if !reflect.DeepEqual(tokens, tt.tokens) {
				t.Errorf("expected '%v', got '%v'", tt.tokens, tokens)
			}
		})
	}
}
