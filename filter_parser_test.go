package qfv

// func TestFilterParser_Parse(t *testing.T) {
// 	allowedFields := []string{"name", "age"}
// 	parser := NewFilterParser(allowedFields)

// 	tests := []struct {
// 		name        string
// 		input       string
// 		expected    FilterNode
// 		expectedErr error
// 	}{
// 		// {
// 		// 	name:        "empty input",
// 		// 	input:       "",
// 		// 	expected:    FilterNode{},
// 		// 	expectedErr: ErrEmptyFilterExpression(),
// 		// },
// 		// {
// 		// 	name:  "simple comparison",
// 		// 	input: "name = 'test'",
// 		// 	expected: FilterNode{
// 		// 		Root: ComparisonNode{
// 		// 			Field:    "name",
// 		// 			Operator: OperatorEqual,
// 		// 			Value:    StringNode{Value: "test"},
// 		// 		},
// 		// 	},
// 		// 	expectedErr: nil,
// 		// },
// 		// {
// 		// 	name:  "simple comparison with identifier",
// 		// 	input: "name = age",
// 		// 	expected: FilterNode{
// 		// 		Root: ComparisonNode{
// 		// 			Field:    "name",
// 		// 			Operator: OperatorEqual,
// 		// 			Value:    IdentifierNode{Value: "age"},
// 		// 		},
// 		// 	},
// 		// 	expectedErr: nil,
// 		// },
// 		// {
// 		// 	name:        "unknown field",
// 		// 	input:       "unknown = test",
// 		// 	expected:    FilterNode{},
// 		// 	expectedErr: ErrUnknownField("unknown"),
// 		// },
// 		{
// 			name:  "logical AND",
// 			input: "name = 'test' AND age = 10",
// 			expected: FilterNode{
// 				Root: LogicalOperationNode{
// 					Operator: OperatorAnd,
// 					Left: ComparisonNode{
// 						Field:    "name",
// 						Operator: OperatorEqual,
// 						Value:    IdentifierNode{Value: "test"},
// 					},
// 					Right: ComparisonNode{
// 						Field:    "age",
// 						Operator: OperatorEqual,
// 						Value:    NumberNode{Value: 10},
// 					},
// 				},
// 			},
// 			expectedErr: nil,
// 		},
// 		// {
// 		// 	name:  "logical OR",
// 		// 	input: "name = test OR age = 10",
// 		// 	expected: FilterNode{
// 		// 		Root: LogicalOperationNode{
// 		// 			Operator: OperatorOr,
// 		// 			Left: ComparisonNode{
// 		// 				Field:    "name",
// 		// 				Operator: OperatorEqual,
// 		// 				Value:    StringNode{Value: "test"},
// 		// 			},
// 		// 			Right: ComparisonNode{
// 		// 				Field:    "age",
// 		// 				Operator: OperatorEqual,
// 		// 				Value:    StringNode{Value: "10"},
// 		// 			},
// 		// 		},
// 		// 	},
// 		// 	expectedErr: nil,
// 		// },
// 		// {
// 		// 	name:  "logical NOT",
// 		// 	input: "NOT name = test",
// 		// 	expected: FilterNode{
// 		// 		Root: LogicalOperationNode{
// 		// 			Operator: OperatorNot,
// 		// 			Left: ComparisonNode{
// 		// 				Field:    "name",
// 		// 				Operator: OperatorEqual,
// 		// 				Value:    StringNode{Value: "test"},
// 		// 			},
// 		// 			Right: nil,
// 		// 		},
// 		// 	},
// 		// 	expectedErr: nil,
// 		// },
// 		// {
// 		// 	name:  "parentheses",
// 		// 	input: "(name = test)",
// 		// 	expected: FilterNode{
// 		// 		Root: ComparisonNode{
// 		// 			Field:    "name",
// 		// 			Operator: OperatorEqual,
// 		// 			Value:    StringNode{Value: "test"},
// 		// 		},
// 		// 	},
// 		// 	expectedErr: nil,
// 		// },
// 		// {
// 		// 	name:  "nested parentheses",
// 		// 	input: "((name = test))",
// 		// 	expected: FilterNode{
// 		// 		Root: ComparisonNode{
// 		// 			Field:    "name",
// 		// 			Operator: OperatorEqual,
// 		// 			Value:    StringNode{Value: "test"},
// 		// 		},
// 		// 	},
// 		// 	expectedErr: nil,
// 		// },
// 		// {
// 		// 	name:  "parentheses with AND",
// 		// 	input: "(name = test) AND (age = 10)",
// 		// 	expected: FilterNode{
// 		// 		Root: LogicalOperationNode{
// 		// 			Operator: OperatorAnd,
// 		// 			Left: ComparisonNode{
// 		// 				Field:    "name",
// 		// 				Operator: OperatorEqual,
// 		// 				Value:    StringNode{Value: "test"},
// 		// 			},
// 		// 			Right: ComparisonNode{
// 		// 				Field:    "age",
// 		// 				Operator: OperatorEqual,
// 		// 				Value:    StringNode{Value: "10"},
// 		// 			},
// 		// 		},
// 		// 	},
// 		// 	expectedErr: nil,
// 		// },
// 		// {
// 		// 	name:  "complex expression",
// 		// 	input: "(name = test AND age = 10) OR NOT city = newyork",
// 		// 	expected: FilterNode{
// 		// 		Root: LogicalOperationNode{
// 		// 			Operator: OperatorOr,
// 		// 			Left: LogicalOperationNode{
// 		// 				Operator: OperatorAnd,
// 		// 				Left: ComparisonNode{
// 		// 					Field:    "name",
// 		// 					Operator: OperatorEqual,
// 		// 					Value:    StringNode{Value: "test"},
// 		// 				},
// 		// 				Right: ComparisonNode{
// 		// 					Field:    "age",
// 		// 					Operator: OperatorEqual,
// 		// 					Value:    StringNode{Value: "10"},
// 		// 				},
// 		// 			},
// 		// 			Right: LogicalOperationNode{
// 		// 				Operator: OperatorNot,
// 		// 				Left: ComparisonNode{
// 		// 					Field:    "city",
// 		// 					Operator: OperatorEqual,
// 		// 					Value:    StringNode{Value: "newyork"},
// 		// 				},
// 		// 				Right: nil,
// 		// 			},
// 		// 		},
// 		// 	},
// 		// 	expectedErr: nil,
// 		// },
// 	}

// 	for _, tt := range tests {
// 		t.Run(tt.name, func(t *testing.T) {
// 			actual, err := parser.Parse(tt.input)

// 			if tt.expectedErr != nil {
// 				if err == nil {
// 					t.Fatalf("expected error '%v', got nil", tt.expectedErr)
// 				}
// 				if err.Error() != tt.expectedErr.Error() {
// 					t.Fatalf("expected error '%v', got '%v'", tt.expectedErr, err)
// 				}
// 				return
// 			}

// 			if err != nil {
// 				t.Fatalf("unexpected error: %v", err)
// 			}

// 			if !reflect.DeepEqual(actual, tt.expected) {
// 				t.Errorf("expected '%v', got '%v'", tt.expected, actual)
// 			}
// 		})
// 	}
// }

// func ErrEmptyFilterExpression() error {
// 	return fmt.Errorf("empty filter expression")
// }
