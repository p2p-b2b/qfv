package qfv

import (
	"reflect"
	"testing"
)

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
			expectedErr: &QFVFieldsError{Message: "empty input expression"},
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
			expectedErr: &QFVFieldsError{Field: "", Message: "empty field expression"},
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
			expectedErr: &QFVFieldsError{Field: "unknown", Message: "unknown field"},
		},
		{
			name:        "invalid field uppercase",
			input:       "NAME",
			expected:    FieldsNode{},
			expectedErr: &QFVFieldsError{Field: "NAME", Message: "unknown field"},
		},
		{
			name:        "empty part",
			input:       "name, ,age",
			expected:    FieldsNode{},
			expectedErr: &QFVFieldsError{Field: "", Message: "empty field expression"},
		},
		{
			name:        "single comma",
			input:       ",",
			expected:    FieldsNode{},
			expectedErr: &QFVFieldsError{Field: "", Message: "empty field expression"},
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
