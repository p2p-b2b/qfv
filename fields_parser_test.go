package qfv

import (
	"fmt"
	"reflect"
	"testing"
)

func ErrEmptyFieldExpression(input string) error {
	return fmt.Errorf("empty field expression: %s", input)
}

func TestFieldsParser_Parse(t *testing.T) {
	allowedFields := []string{"name", "age", "city"}
	parser := NewFieldsParser(allowedFields)

	tests := []struct {
		name        string
		input       string
		expected    FieldsNode
		expectedErr error
	}{
		{
			name:        "empty input",
			input:       "",
			expected:    FieldsNode{},
			expectedErr: ErrEmptyInputExpression(),
		},
		{
			name:  "single field",
			input: "name",
			expected: FieldsNode{
				Fields: []string{"name"},
			},
			expectedErr: nil,
		},
		{
			name:        "single field, extra comma",
			input:       "name,",
			expected:    FieldsNode{},
			expectedErr: ErrEmptyFieldExpression("name,"),
		},
		{
			name:  "multiple fields",
			input: "name,age,city",
			expected: FieldsNode{
				Fields: []string{"name", "age", "city"},
			},
			expectedErr: nil,
		},
		{
			name:  "multiple fields with extra spaces",
			input: "name,   age, city",
			expected: FieldsNode{
				Fields: []string{"name", "age", "city"},
			},
			expectedErr: nil,
		},
		{
			name:        "invalid field",
			input:       "unknown",
			expected:    FieldsNode{},
			expectedErr: ErrUnknownField("unknown"),
		},
		{
			name:        "invalid field uppercase",
			input:       "NAME",
			expected:    FieldsNode{},
			expectedErr: ErrUnknownField("NAME"),
		},
		{
			name:        "empty part",
			input:       "name, ,age",
			expected:    FieldsNode{},
			expectedErr: ErrEmptyFieldExpression("name, ,age"),
		},
		{
			name:        "single comma",
			input:       ",",
			expected:    FieldsNode{},
			expectedErr: ErrEmptyFieldExpression(","),
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			actual, err := parser.Parse(tt.input)

			if tt.expectedErr != nil {
				if err == nil {
					t.Fatalf("expected error '%v', got nil", tt.expectedErr)
				}
				if err.Error() != tt.expectedErr.Error() {
					t.Fatalf("expected error '%v', got '%v'", tt.expectedErr, err)
				}
				return
			}

			if err != nil {
				t.Fatalf("unexpected error: %v", err)
			}

			if !reflect.DeepEqual(actual, tt.expected) {
				t.Errorf("expected '%v', got '%v'", tt.expected, actual)
			}
		})
	}
}

func TestFieldsParser_Validate(t *testing.T) {
	allowedFields := []string{"name", "age", "city"}
	parser := NewFieldsParser(allowedFields)

	tests := []struct {
		name        string
		input       string
		expected    FieldsNode
		expectedErr error
	}{
		{
			name:        "empty input",
			input:       "",
			expected:    FieldsNode{},
			expectedErr: ErrEmptyInputExpression(),
		},
		{
			name:  "single field",
			input: "name",
			expected: FieldsNode{
				Fields: []string{"name"},
			},
			expectedErr: nil,
		},
		{
			name:  "multiple fields",
			input: "name,age,city",
			expected: FieldsNode{
				Fields: []string{"name", "age", "city"},
			},
			expectedErr: nil,
		},
		{
			name:        "invalid field",
			input:       "unknown",
			expected:    FieldsNode{},
			expectedErr: ErrUnknownField("unknown"),
		},
		{
			name:        "empty part",
			input:       "name, ,age",
			expected:    FieldsNode{},
			expectedErr: ErrEmptyFieldExpression("name, ,age"),
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			actual, err := parser.Validate(tt.input)

			if tt.expectedErr != nil {
				if err == nil {
					t.Fatalf("expected error '%v', got nil", tt.expectedErr)
				}
				if err.Error() != tt.expectedErr.Error() {
					t.Fatalf("expected error '%v', got '%v'", tt.expectedErr, err)
				}
				return
			}

			if err != nil {
				t.Fatalf("unexpected error: %v", err)
			}

			if !reflect.DeepEqual(actual, tt.expected) {
				t.Errorf("expected '%v', got '%v'", tt.expected, actual)
			}
		})
	}
}
