package qfv

import (
	"fmt"
	"reflect"
	"testing"
)

func ErrUnknownField(field string) error {
	return fmt.Errorf("unknown field: %s", field)
}

func ErrInvalidSortDirection(direction string) error {
	return fmt.Errorf("invalid sort direction: %s", direction)
}

func ErrMissingSortDirection(field string) error {
	return fmt.Errorf("missing sort direction for field: %s", field)
}

func ErrInvalidSortExpression(expression string) error {
	return fmt.Errorf("invalid sort expression: %s", expression)
}

func ErrEmptySortExpression(expression string) error {
	return fmt.Errorf("empty sort expression: %s", expression)
}

func ErrEmptyInputExpression() error {
	return fmt.Errorf("empty input expression")
}

func TestSortParser_Parse(t *testing.T) {
	allowedFields := []string{"name", "age", "city"}
	parser := NewSortParser(allowedFields)

	tests := []struct {
		name        string
		input       string
		expected    SortNode
		expectedErr error
	}{
		{
			name:        "empty input",
			input:       "",
			expected:    SortNode{},
			expectedErr: ErrEmptyInputExpression(),
		},
		{
			name:  "single field ascending lowercase",
			input: "name asc",
			expected: SortNode{
				Fields: []SortFieldNode{
					{Field: "name", Direction: SortAsc},
				},
			},
			expectedErr: nil,
		},
		{
			name:  "single field ascending Uppercase",
			input: "name ASC",
			expected: SortNode{
				Fields: []SortFieldNode{
					{Field: "name", Direction: SortAsc},
				},
			},
			expectedErr: nil,
		},
		{
			name:        "single field missing direction",
			input:       "name",
			expected:    SortNode{},
			expectedErr: ErrMissingSortDirection("name"),
		},
		{
			name:  "single field descending uppercase",
			input: "name DESC",
			expected: SortNode{
				Fields: []SortFieldNode{
					{Field: "name", Direction: SortDesc},
				},
			},
			expectedErr: nil,
		},
		{
			name:  "single field descending lowercase",
			input: "name desc",
			expected: SortNode{
				Fields: []SortFieldNode{
					{Field: "name", Direction: SortDesc},
				},
			},
			expectedErr: nil,
		},
		{
			name:        "single field, wrong direction",
			input:       "name invalid",
			expected:    SortNode{},
			expectedErr: ErrInvalidSortDirection("INVALID"),
		},
		{
			name:        "single field multiple direction",
			input:       "name asc DESC",
			expected:    SortNode{},
			expectedErr: ErrInvalidSortExpression("name asc DESC"),
		},
		{
			name:        "single field descending with extra comma",
			input:       "name DESC,",
			expected:    SortNode{},
			expectedErr: ErrEmptySortExpression("name DESC,"),
		},
		{
			name:  "multiple fields, lowercase and uppercase",
			input: "name asc, age DESC, city desc",
			expected: SortNode{
				Fields: []SortFieldNode{
					{Field: "name", Direction: SortAsc},
					{Field: "age", Direction: SortDesc},
					{Field: "city", Direction: SortDesc},
				},
			},
			expectedErr: nil,
		},
		{
			name:        "invalid field",
			input:       "unknown",
			expected:    SortNode{},
			expectedErr: ErrUnknownField("unknown"),
		},
		{
			name:        "invalid direction",
			input:       "name invalid",
			expected:    SortNode{},
			expectedErr: ErrInvalidSortDirection("INVALID"),
		},
		{
			name:        "empty part",
			input:       "name, ,age",
			expected:    SortNode{},
			expectedErr: ErrMissingSortDirection("name"),
		},
		{
			name:        "invalid sort expression",
			input:       " ",
			expected:    SortNode{},
			expectedErr: ErrEmptySortExpression(" "),
		},
		{
			name:        "single comma",
			input:       ",",
			expected:    SortNode{},
			expectedErr: ErrEmptySortExpression(","),
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
