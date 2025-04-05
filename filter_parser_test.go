package qfv

import (
	"reflect"
	"testing"
)

func TestFilterParser_Parse(t *testing.T) {
	tests := []struct {
		name          string
		input         string
		allowedFields []string
		wantErr       bool
		checkNode     func(t *testing.T, node Node)
	}{
		{
			name:          "empty input",
			input:         "",
			allowedFields: []string{"name", "age", "status"},
			wantErr:       true,
		},
		{
			name:          "simple equality",
			input:         "name = 'John'",
			allowedFields: []string{"name", "age", "status"},
			wantErr:       false,
			checkNode: func(t *testing.T, node Node) {
				binOp, ok := node.(*BinaryOperatorNode)
				if !ok {
					t.Fatalf("expected BinaryOperatorNode, got %T", node)
				}
				if binOp.Operator != TokenOperatorEqual {
					t.Errorf("expected = operator, got %s", binOp.Operator)
				}

				field, ok := binOp.Left.(*IdentifierNode)
				if !ok {
					t.Fatalf("expected IdentifierNode for left operand, got %T", binOp.Left)
				}
				if field.Name != "name" {
					t.Errorf("expected field name 'name', got %s", field.Name)
				}

				val, ok := binOp.Right.(*LiteralNode)
				if !ok {
					t.Fatalf("expected LiteralNode for right operand, got %T", binOp.Right)
				}
				if val.Value != "John" {
					t.Errorf("expected value 'John', got %v", val.Value)
				}
			},
		},
		{
			name:          "logical AND",
			input:         "name = 'John' AND age > 30",
			allowedFields: []string{"name", "age", "status"},
			wantErr:       false,
			checkNode: func(t *testing.T, node Node) {
				binOp, ok := node.(*BinaryOperatorNode)
				if !ok {
					t.Fatalf("expected BinaryOperatorNode, got %T", node)
				}
				if binOp.Operator != TokenOperatorAnd {
					t.Errorf("expected AND operator, got %s", binOp.Operator)
				}

				// Check left operand (name = 'John')
				leftBinOp, ok := binOp.Left.(*BinaryOperatorNode)
				if !ok {
					t.Fatalf("expected BinaryOperatorNode for left operand, got %T", binOp.Left)
				}
				if leftBinOp.Operator != TokenOperatorEqual {
					t.Errorf("expected = operator for left operand, got %s", leftBinOp.Operator)
				}

				// Check right operand (age > 30)
				rightBinOp, ok := binOp.Right.(*BinaryOperatorNode)
				if !ok {
					t.Fatalf("expected BinaryOperatorNode for right operand, got %T", binOp.Right)
				}
				if rightBinOp.Operator != TokenOperatorGreaterThan {
					t.Errorf("expected > operator for right operand, got %s", rightBinOp.Operator)
				}
			},
		},
		{
			name:          "logical OR",
			input:         "name = 'John' OR status = 'active'",
			allowedFields: []string{"name", "age", "status"},
			wantErr:       false,
			checkNode: func(t *testing.T, node Node) {
				binOp, ok := node.(*BinaryOperatorNode)
				if !ok {
					t.Fatalf("expected BinaryOperatorNode, got %T", node)
				}
				if binOp.Operator != TokenOperatorOr {
					t.Errorf("expected OR operator, got %s", binOp.Operator)
				}
			},
		},
		{
			name:          "parenthesized expression",
			input:         "(name = 'John' OR name = 'Jane') AND age > 30",
			allowedFields: []string{"name", "age", "status"},
			wantErr:       false,
			checkNode: func(t *testing.T, node Node) {
				binOp, ok := node.(*BinaryOperatorNode)
				if !ok {
					t.Fatalf("expected BinaryOperatorNode, got %T", node)
				}
				if binOp.Operator != TokenOperatorAnd {
					t.Errorf("expected AND operator, got %s", binOp.Operator)
				}

				// Check left operand (group)
				group, ok := binOp.Left.(*GroupNode)
				if !ok {
					t.Fatalf("expected GroupNode for left operand, got %T", binOp.Left)
				}

				// Check group expression (name = 'John' OR name = 'Jane')
				groupExpr, ok := group.Expression.(*BinaryOperatorNode)
				if !ok {
					t.Fatalf("expected BinaryOperatorNode for group expression, got %T", group.Expression)
				}
				if groupExpr.Operator != TokenOperatorOr {
					t.Errorf("expected OR operator for group expression, got %s", groupExpr.Operator)
				}
			},
		},
		{
			name:          "IS NULL operator",
			input:         "name IS NULL",
			allowedFields: []string{"name", "age", "status"},
			wantErr:       false,
			checkNode: func(t *testing.T, node Node) {
				isNull, ok := node.(*IsNullNode)
				if !ok {
					t.Fatalf("expected IsNullNode, got %T", node)
				}
				if isNull.IsNot {
					t.Errorf("expected IsNot to be false")
				}

				field, ok := isNull.Field.(*IdentifierNode)
				if !ok {
					t.Fatalf("expected IdentifierNode for field, got %T", isNull.Field)
				}
				if field.Name != "name" {
					t.Errorf("expected field name 'name', got %s", field.Name)
				}
			},
		},
		{
			name:          "IS NOT NULL operator",
			input:         "name IS NOT NULL",
			allowedFields: []string{"name", "age", "status"},
			wantErr:       false,
			checkNode: func(t *testing.T, node Node) {
				isNull, ok := node.(*IsNullNode)
				if !ok {
					t.Fatalf("expected IsNullNode, got %T", node)
				}
				if !isNull.IsNot {
					t.Errorf("expected IsNot to be true")
				}
			},
		},
		{
			name:          "IN operator",
			input:         "name IN ('John', 'Jane', 'Bob')",
			allowedFields: []string{"name", "age", "status"},
			wantErr:       false,
			checkNode: func(t *testing.T, node Node) {
				inNode, ok := node.(*InNode)
				if !ok {
					t.Fatalf("expected InNode, got %T", node)
				}
				if inNode.IsNot {
					t.Errorf("expected IsNot to be false")
				}

				if len(inNode.Values) != 3 {
					t.Fatalf("expected 3 values, got %d", len(inNode.Values))
				}

				// Check first value
				val1, ok := inNode.Values[0].(*LiteralNode)
				if !ok {
					t.Fatalf("expected LiteralNode for first value, got %T", inNode.Values[0])
				}
				if val1.Value != "John" {
					t.Errorf("expected first value 'John', got %v", val1.Value)
				}
			},
		},
		{
			name:          "NOT IN operator",
			input:         "name NOT IN ('John', 'Jane')",
			allowedFields: []string{"name", "age", "status"},
			wantErr:       false,
			checkNode: func(t *testing.T, node Node) {
				unaryOp, ok := node.(*UnaryOperatorNode)
				if !ok {
					t.Fatalf("expected UnaryOperatorNode, got %T", node)
				}
				if unaryOp.Operator != TokenOperatorNot {
					t.Errorf("expected NOT operator, got %s", unaryOp.Operator)
				}

				inExpr, ok := unaryOp.X.(*InNode)
				if !ok {
					t.Fatalf("expected InNode for NOT operand, got %T", unaryOp.X)
				}
				if inExpr.IsNot {
					t.Errorf("expected IsNot to be false in the InNode")
				}
			},
		},
		{
			name:          "BETWEEN operator",
			input:         "age BETWEEN 20 AND 30",
			allowedFields: []string{"name", "age", "status"},
			wantErr:       false,
			checkNode: func(t *testing.T, node Node) {
				between, ok := node.(*BetweenNode)
				if !ok {
					t.Fatalf("expected BetweenNode, got %T", node)
				}
				if between.IsNot {
					t.Errorf("expected IsNot to be false")
				}

				field, ok := between.Field.(*IdentifierNode)
				if !ok {
					t.Fatalf("expected IdentifierNode for field, got %T", between.Field)
				}
				if field.Name != "age" {
					t.Errorf("expected field name 'age', got %s", field.Name)
				}

				lower, ok := between.Lower.(*LiteralNode)
				if !ok {
					t.Fatalf("expected LiteralNode for lower bound, got %T", between.Lower)
				}
				if lower.Kind != reflect.Int64 {
					t.Errorf("expected lower bound to be int64, got %v", lower.Kind)
				}

				upper, ok := between.Upper.(*LiteralNode)
				if !ok {
					t.Fatalf("expected LiteralNode for upper bound, got %T", between.Upper)
				}
				if upper.Kind != reflect.Int64 {
					t.Errorf("expected upper bound to be int64, got %v", upper.Kind)
				}
			},
		},
		{
			name:          "NOT BETWEEN operator",
			input:         "age NOT BETWEEN 20 AND 30",
			allowedFields: []string{"name", "age", "status"},
			wantErr:       false,
			checkNode: func(t *testing.T, node Node) {
				unaryOp, ok := node.(*UnaryOperatorNode)
				if !ok {
					t.Fatalf("expected UnaryOperatorNode, got %T", node)
				}
				if unaryOp.Operator != TokenOperatorNot {
					t.Errorf("expected NOT operator, got %s", unaryOp.Operator)
				}

				betweenExpr, ok := unaryOp.X.(*BetweenNode)
				if !ok {
					t.Fatalf("expected BetweenNode for NOT operand, got %T", unaryOp.X)
				}
				if betweenExpr.IsNot {
					t.Errorf("expected IsNot to be false in the BetweenNode")
				}
			},
		},
		{
			name:          "LIKE operator",
			input:         "name LIKE '%John%'",
			allowedFields: []string{"name", "age", "status"},
			wantErr:       false,
			checkNode: func(t *testing.T, node Node) {
				binOp, ok := node.(*BinaryOperatorNode)
				if !ok {
					t.Fatalf("expected BinaryOperatorNode, got %T", node)
				}
				if binOp.Operator != TokenOperatorLike {
					t.Errorf("expected LIKE operator, got %s", binOp.Operator)
				}

				field, ok := binOp.Left.(*IdentifierNode)
				if !ok {
					t.Fatalf("expected IdentifierNode for left operand, got %T", binOp.Left)
				}
				if field.Name != "name" {
					t.Errorf("expected field name 'name', got %s", field.Name)
				}

				pattern, ok := binOp.Right.(*LiteralNode)
				if !ok {
					t.Fatalf("expected LiteralNode for right operand, got %T", binOp.Right)
				}
				if pattern.Value != "%John%" {
					t.Errorf("expected pattern '%%John%%', got %v", pattern.Value)
				}
			},
		},
		{
			name:          "NOT LIKE operator",
			input:         "name NOT LIKE '%John%'",
			allowedFields: []string{"name", "age", "status"},
			wantErr:       false,
			checkNode: func(t *testing.T, node Node) {
				unaryOp, ok := node.(*UnaryOperatorNode)
				if !ok {
					t.Fatalf("expected UnaryOperatorNode, got %T", node)
				}
				if unaryOp.Operator != TokenOperatorNot {
					t.Errorf("expected NOT operator, got %s", unaryOp.Operator)
				}

				likeExpr, ok := unaryOp.X.(*BinaryOperatorNode)
				if !ok {
					t.Fatalf("expected BinaryOperatorNode for NOT operand, got %T", unaryOp.X)
				}
				if likeExpr.Operator != TokenOperatorLike {
					t.Errorf("expected LIKE operator, got %s", likeExpr.Operator)
				}
			},
		},
		{
			name:          "complex expression",
			input:         "(name = 'John' OR name = 'Jane') AND (age > 30 OR status = 'active')",
			allowedFields: []string{"name", "age", "status"},
			wantErr:       false,
			checkNode: func(t *testing.T, node Node) {
				binOp, ok := node.(*BinaryOperatorNode)
				if !ok {
					t.Fatalf("expected BinaryOperatorNode, got %T", node)
				}
				if binOp.Operator != TokenOperatorAnd {
					t.Errorf("expected AND operator, got %s", binOp.Operator)
				}

				// Check left operand (group)
				leftGroup, ok := binOp.Left.(*GroupNode)
				if !ok {
					t.Fatalf("expected GroupNode for left operand, got %T", binOp.Left)
				}

				// Check right operand (group)
				rightGroup, ok := binOp.Right.(*GroupNode)
				if !ok {
					t.Fatalf("expected GroupNode for right operand, got %T", binOp.Right)
				}

				// Check left group expression (name = 'John' OR name = 'Jane')
				leftGroupExpr, ok := leftGroup.Expression.(*BinaryOperatorNode)
				if !ok {
					t.Fatalf("expected BinaryOperatorNode for left group expression, got %T", leftGroup.Expression)
				}
				if leftGroupExpr.Operator != TokenOperatorOr {
					t.Errorf("expected OR operator for left group expression, got %s", leftGroupExpr.Operator)
				}

				// Check right group expression (age > 30 OR status = 'active')
				rightGroupExpr, ok := rightGroup.Expression.(*BinaryOperatorNode)
				if !ok {
					t.Fatalf("expected BinaryOperatorNode for right group expression, got %T", rightGroup.Expression)
				}
				if rightGroupExpr.Operator != TokenOperatorOr {
					t.Errorf("expected OR operator for right group expression, got %s", rightGroupExpr.Operator)
				}
			},
		},
		{
			name:          "disallowed field",
			input:         "email = 'john@example.com'",
			allowedFields: []string{"name", "age", "status"},
			wantErr:       true,
		},
		{
			name:          "syntax error - missing closing parenthesis",
			input:         "(name = 'John' AND age > 30",
			allowedFields: []string{"name", "age", "status"},
			wantErr:       true,
		},
		{
			name:          "syntax error - missing value after operator",
			input:         "name =",
			allowedFields: []string{"name", "age", "status"},
			wantErr:       true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			p := NewFilterParser(tt.allowedFields)
			node, err := p.Parse(tt.input)

			if (err != nil) != tt.wantErr {
				t.Errorf("Parse() error = %v, wantErr %v", err, tt.wantErr)
				return
			}

			if err == nil && tt.checkNode != nil {
				tt.checkNode(t, node)
			}
		})
	}
}
