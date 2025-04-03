package qfv

import (
	"testing"
	"text/scanner"
)

func TestLexer_Navigation(t *testing.T) {
	tests := []struct {
		name     string
		input    string
		expected []Token
	}{
		// {
		// 	name:  "double quoted string, 1",
		// 	input: `"comment = 'This is a string'"`,
		// 	expected: []Token{
		// 		{Pos: scanner.Position{Line: 1, Column: 1}, Type: TokenIllegal, Value: `"comment = 'This is a string'"`},
		// 		{Pos: scanner.Position{Line: 1, Column: 34}, Type: TokenEOF, Value: ""},
		// 	},
		// },
		// {
		// 	name:  "double quoted string, 2",
		// 	input: `comment = "'This is a string'"`,
		// 	expected: []Token{
		// 		{Pos: scanner.Position{Line: 1, Column: 1}, Type: TokenIdentifier, Value: "comment"},
		// 		{Pos: scanner.Position{Line: 1, Column: 9}, Type: TokenOperatorEqual, Value: "="},
		// 		{Pos: scanner.Position{Line: 1, Column: 1}, Type: TokenIllegal, Value: `"'This is a string'"`},
		// 		{Pos: scanner.Position{Line: 1, Column: 34}, Type: TokenEOF, Value: ""},
		// 	},
		// },
		// {
		// 	name:  "bad double quoted string, missing closing quote",
		// 	input: `comment = "'This is a bad string`,
		// 	expected: []Token{
		// 		{Pos: scanner.Position{Line: 1, Column: 1}, Type: TokenIdentifier, Value: "comment"},
		// 		{Pos: scanner.Position{Line: 1, Column: 9}, Type: TokenOperatorEqual, Value: "="},
		// 		{Pos: scanner.Position{Line: 1, Column: 11}, Type: TokenIllegal, Value: `"'This is a bad string`},
		// 		{Pos: scanner.Position{Line: 1, Column: 34}, Type: TokenEOF, Value: ""},
		// 	},
		// },
		{
			name:  "bad double quoted string, missing opening quote",
			input: `comment = This is a bad string'"`,
			expected: []Token{
				{Pos: scanner.Position{Line: 1, Column: 1}, Type: TokenIdentifier, Value: "comment"},
				{Pos: scanner.Position{Line: 1, Column: 9}, Type: TokenOperatorEqual, Value: "="},
				{Pos: scanner.Position{Line: 1, Column: 11}, Type: TokenIdentifier, Value: "This"},
				{Pos: scanner.Position{Line: 1, Column: 16}, Type: TokenIdentifier, Value: "is"},
				{Pos: scanner.Position{Line: 1, Column: 19}, Type: TokenIdentifier, Value: "a"},
				{Pos: scanner.Position{Line: 1, Column: 21}, Type: TokenIdentifier, Value: "bad"},
				{Pos: scanner.Position{Line: 1, Column: 25}, Type: TokenIdentifier, Value: "string"},
				{Pos: scanner.Position{Line: 1, Column: 31}, Type: TokenIllegal, Value: `'"`},
				{Pos: scanner.Position{Line: 1, Column: 34}, Type: TokenEOF, Value: ""},
			},
		},
		// {
		// 	name:  "bad quoted string, missing closing quote",
		// 	input: "comment = 'This is a bad string",
		// 	expected: []Token{
		// 		{Pos: scanner.Position{Line: 1, Column: 1}, Type: TokenIdentifier, Value: "comment"},
		// 		{Pos: scanner.Position{Line: 1, Column: 9}, Type: TokenOperatorEqual, Value: "="},
		// 		{Pos: scanner.Position{Line: 1, Column: 11}, Type: TokenIllegal, Value: `'This is a bad string`},
		// 		{Pos: scanner.Position{Line: 1, Column: 34}, Type: TokenEOF, Value: ""},
		// 	},
		// },
		// {
		// 	name:  "bad quoted string, missing opening quote",
		// 	input: "comment = This is a bad string'",
		// 	expected: []Token{
		// 		{Pos: scanner.Position{Line: 1, Column: 1}, Type: TokenIdentifier, Value: "comment"},
		// 		{Pos: scanner.Position{Line: 1, Column: 9}, Type: TokenOperatorEqual, Value: "="},
		// 		{Pos: scanner.Position{Line: 1, Column: 11}, Type: TokenIdentifier, Value: "This"},
		// 		{Pos: scanner.Position{Line: 1, Column: 16}, Type: TokenIdentifier, Value: "is"},
		// 		{Pos: scanner.Position{Line: 1, Column: 19}, Type: TokenIdentifier, Value: "a"},
		// 		{Pos: scanner.Position{Line: 1, Column: 21}, Type: TokenIdentifier, Value: "bad"},
		// 		{Pos: scanner.Position{Line: 1, Column: 25}, Type: TokenIdentifier, Value: "string"},
		// 		{Pos: scanner.Position{Line: 1, Column: 31}, Type: TokenIllegal, Value: "'"},
		// 		{Pos: scanner.Position{Line: 1, Column: 34}, Type: TokenEOF, Value: ""},
		// 	},
		// },
		// {
		// 	name:  "string literal with spaces and scaped quotes",
		// 	input: "comment = 'This is a \\' scaped single quote' AND age > 18",
		// 	expected: []Token{
		// 		{Pos: scanner.Position{Line: 1, Column: 1}, Type: TokenIdentifier, Value: "comment"},
		// 		{Pos: scanner.Position{Line: 1, Column: 9}, Type: TokenOperatorEqual, Value: "="},
		// 		{Pos: scanner.Position{Line: 1, Column: 11}, Type: TokenString, Value: `'This is a \' scaped single quote'`},
		// 		{Pos: scanner.Position{Line: 1, Column: 56}, Type: TokenOperatorAnd, Value: "AND"},
		// 		{Pos: scanner.Position{Line: 1, Column: 60}, Type: TokenIdentifier, Value: "age"},
		// 		{Pos: scanner.Position{Line: 1, Column: 63}, Type: TokenOperatorGreaterThan, Value: ">"},
		// 		{Pos: scanner.Position{Line: 1, Column: 65}, Type: TokenInt, Value: "18"},
		// 		{Pos: scanner.Position{Line: 1, Column: 67}, Type: TokenEOF, Value: ""},
		// 	},
		// },
		// {
		// 	name:  "multiple operators",
		// 	input: "age > 18 AND name = 'John' OR status != 'active'",
		// 	expected: []Token{
		// 		{Pos: scanner.Position{Line: 1, Column: 1}, Type: TokenIdentifier, Value: "age"},
		// 		{Pos: scanner.Position{Line: 1, Column: 4}, Type: TokenOperatorGreaterThan, Value: ">"},
		// 		{Pos: scanner.Position{Line: 1, Column: 6}, Type: TokenInt, Value: "18"},
		// 		{Pos: scanner.Position{Line: 1, Column: 8}, Type: TokenOperatorAnd, Value: "AND"},
		// 		{Pos: scanner.Position{Line: 1, Column: 12}, Type: TokenIdentifier, Value: "name"},
		// 		{Pos: scanner.Position{Line: 1, Column: 16}, Type: TokenOperatorEqual, Value: "="},
		// 		{Pos: scanner.Position{Line: 1, Column: 18}, Type: TokenString, Value: `'John'`},
		// 		{Pos: scanner.Position{Line: 1, Column: 24}, Type: TokenOperatorOr, Value: "OR"},
		// 		{Pos: scanner.Position{Line: 1, Column: 27}, Type: TokenIdentifier, Value: "status"},
		// 		{Pos: scanner.Position{Line: 1, Column: 33}, Type: TokenOperatorNotEqualAlias, Value: "!="},
		// 		{Pos: scanner.Position{Line: 1, Column: 35}, Type: TokenString, Value: `'active'`},
		// 		{Pos: scanner.Position{Line: 1, Column: 42}, Type: TokenEOF, Value: ""},
		// 	},
		// },
		// {
		// 	name:  "multiple operators, with sub-grouping",
		// 	input: "age > 18 AND (name = 'John' OR status != 'active')",
		// 	expected: []Token{
		// 		{Pos: scanner.Position{Line: 1, Column: 1}, Type: TokenIdentifier, Value: "age"},
		// 		{Pos: scanner.Position{Line: 1, Column: 4}, Type: TokenOperatorGreaterThan, Value: ">"},
		// 		{Pos: scanner.Position{Line: 1, Column: 6}, Type: TokenInt, Value: "18"},
		// 		{Pos: scanner.Position{Line: 1, Column: 8}, Type: TokenOperatorAnd, Value: "AND"},
		// 		{Pos: scanner.Position{Line: 1, Column: 12}, Type: TokenLPAREN, Value: "("},
		// 		{Pos: scanner.Position{Line: 1, Column: 13}, Type: TokenIdentifier, Value: "name"},
		// 		{Pos: scanner.Position{Line: 1, Column: 17}, Type: TokenOperatorEqual, Value: "="},
		// 		{Pos: scanner.Position{Line: 1, Column: 19}, Type: TokenString, Value: `'John'`},
		// 		{Pos: scanner.Position{Line: 1, Column: 25}, Type: TokenOperatorOr, Value: "OR"},
		// 		{Pos: scanner.Position{Line: 1, Column: 27}, Type: TokenIdentifier, Value: "status"},
		// 		{Pos: scanner.Position{Line: 1, Column: 33}, Type: TokenOperatorNotEqualAlias, Value: "!="},
		// 		{Pos: scanner.Position{Line: 1, Column: 35}, Type: TokenString, Value: `'active'`},
		// 		{Pos: scanner.Position{Line: 1, Column: 42}, Type: TokenRPAREN, Value: ")"},
		// 		{Pos: scanner.Position{Line: 1, Column: 43}, Type: TokenEOF, Value: ""},
		// 	},
		// },
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			lexer := NewLexer(tt.input)
			lexer.Parse()

			t.Logf("Lexer tokens: %v", lexer.tokens)

			for i, expected := range tt.expected {
				if len(lexer.tokens) <= i {
					t.Fatalf("expected token %d, got none", i)
				}

				token := lexer.tokens[i]
				if token.Type != expected.Type {
					t.Errorf("expected token type -->%s<--, got -->%s<--", expected.Type, token.Type)
				}

				if token.Value != expected.Value {
					t.Errorf("expected token value -->%s<--, got -->%s<--", expected.Value, token.Value)
				}
			}
		})
	}
}
