package qfv

// Validator is the interface for all validators
// It defines a method to validate the fields of a query
// and return a Node representing the parsed fields.
// The Validate method takes a string input representing the fields
// and returns a Node and an error.
type Validator interface {
	Validate(fields string) (Node, error)
}

// NodeType represents the type of AST node
// It is used to identify the kind of node in the AST
// and to provide type safety when working with different node types.
// Each node type corresponds to a specific part of the query
// (e.g., fields, sort, filter, etc.).
type NodeType int

const (
	NodeTypeLiteral NodeType = iota
	NodeTypeField
	NodeTypeSort
	NodeTypeSortField
	NodeTypeFilter
	NodeTypeLogicalOperator
	NodeTypeComparisonOperator
	NodeTypeString
	NodeTypeBoolean
	NodeTypeBinaryOperator
	NodeTypeNumber
	NodeTypeIdentifier
	NodeTypeFieldList
	NodeTypeFieldListItem
	NodeTypeBetween
	NodeTypeParentExpression
)

// String returns the string representation of the NodeType
func (nt NodeType) String() string {
	switch nt {
	case NodeTypeLiteral:
		return "Literal"
	case NodeTypeField:
		return "Field"
	case NodeTypeSort:
		return "Sort"
	case NodeTypeSortField:
		return "SortField"
	case NodeTypeFilter:
		return "Filter"
	case NodeTypeLogicalOperator:
		return "LogicalOperator"
	case NodeTypeComparisonOperator:
		return "ComparisonOperator"
	case NodeTypeString:
		return "String"
	case NodeTypeBoolean:
		return "Boolean"
	case NodeTypeBinaryOperator:
		return "BinaryOperator"
	case NodeTypeNumber:
		return "Number"
	case NodeTypeIdentifier:
		return "Identifier"
	case NodeTypeFieldList:
		return "FieldList"
	case NodeTypeFieldListItem:
		return "FieldListItem"
	case NodeTypeBetween:
		return "Between"
	case NodeTypeParentExpression:
		return "ParentExpression"
	default:
		return "Unknown"
	}
}

// Node is the base interface for all AST nodes
type Node interface {
	Type() NodeType
}

// FilterNode represents a filter expression
type FilterNode struct {
	Expression Node
}

func (n *FilterNode) Type() NodeType {
	return NodeTypeFilter
}

// LiteralNode represents a literal value (string, number, bool)
type LiteralNode struct {
	Value any
}

func (n *LiteralNode) Type() NodeType {
	return NodeTypeLiteral
}

// IdentifierNode represents a field name
type IdentifierNode struct {
	Name string
}

func (n *IdentifierNode) Type() NodeType {
	return NodeTypeIdentifier
}

// BinaryOperatorNode represents a binary operation (AND, OR)
type BinaryOperatorNode struct {
	Left     Node
	Right    Node
	Operator Operator // "AND" and "OR"
}

func (n *BinaryOperatorNode) Type() NodeType {
	return NodeTypeBinaryOperator
}

// UnaryOperatorNode represents a unary operation (NOT)
type UnaryOperatorNode struct {
	Expression Node
	Operator   Operator // "NOT"
}

func (n *UnaryOperatorNode) Type() NodeType {
	return NodeTypeBinaryOperator
}

// ComparisonNode represents a comparison operation (=, !=, etc.)
type ComparisonNode struct {
	Field    string
	Operator Operator
	Value    Node // StringNode or IdentifierNode
}

func (n ComparisonNode) Type() NodeType {
	return NodeTypeComparisonOperator
}

// BetweenNode represents a BETWEEN or NOT BETWEEN operation
type BetweenNode struct {
	Field *IdentifierNode
	Not   bool
	Lower Node
	Upper Node
}

func (n *BetweenNode) Type() NodeType {
	return NodeTypeBetween
}

// IsNullNode represents an IS NULL or IS NOT NULL operation
type IsNullNode struct {
	Field *IdentifierNode
	Not   bool
}

func (n *IsNullNode) Type() NodeType {
	return NodeTypeComparisonOperator
}

// DistinctNode represents a DISTINCT or NOT DISTINCT operation
type DistinctNode struct {
	Field *IdentifierNode
	Not   bool
}

func (n *DistinctNode) Type() NodeType {
	return NodeTypeComparisonOperator
}

// ParentExpressionNode represents a parenthesized expression
type ParentExpressionNode struct {
	Expr Node
}

func (n *ParentExpressionNode) Type() NodeType {
	return NodeTypeParentExpression
}
