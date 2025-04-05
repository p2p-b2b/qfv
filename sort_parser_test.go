package qfv

import (
	"reflect"
	"testing"
)

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
			expectedErr: &QFVSortError{Message: "empty input expression"},
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
			expectedErr: &QFVSortError{Field: "name", Message: "missing sort direction after field"},
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
			expectedErr: &QFVSortError{Field: "name", Message: "invalid sort direction"},
		},
		{
			name:        "single field multiple direction",
			input:       "name asc DESC",
			expected:    SortNode{},
			expectedErr: &QFVSortError{Field: "name asc DESC", Message: "too many sort expressions"},
		},
		{
			name:        "single field descending with extra comma",
			input:       "name DESC,",
			expected:    SortNode{},
			expectedErr: &QFVSortError{Message: "empty sort expression"},
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
			expectedErr: &QFVSortError{Field: "unknown", Message: "field not allowed for sorting"},
		},
		{
			name:        "invalid direction",
			input:       "name invalid",
			expected:    SortNode{},
			expectedErr: &QFVSortError{Field: "name", Message: "invalid sort direction"},
		},
		{
			name:        "empty part",
			input:       "name, ,age",
			expected:    SortNode{},
			expectedErr: &QFVSortError{Field: "name", Message: "missing sort direction after field"},
		},
		{
			name:        "invalid sort expression",
			input:       " ",
			expected:    SortNode{},
			expectedErr: &QFVSortError{Message: "empty sort expression"},
		},
		{
			name:        "single comma",
			input:       ",",
			expected:    SortNode{},
			expectedErr: &QFVSortError{Message: "empty sort expression"},
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
