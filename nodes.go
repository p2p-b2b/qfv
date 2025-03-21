package qfv

// NodeType represents the type of AST node
// It is used to identify the kind of node in the AST
// and to provide type safety when working with different node types.
// Each node type corresponds to a specific part of the query
// (e.g., fields, sort, filter, etc.).
type NodeType int

const (
	NodeTypeUnknown NodeType = iota
	NodeTypeField
	NodeTypeSort
	NodeTypeSortField
	NodeTypeFilter
	NodeTypeLogicalOp
	NodeTypeComparisonOp
	NodeTypeCheckerOp
	NodeTypeString
	NodeTypeIdentifier
	NodeTypeFieldList
)

// String returns the string representation of the NodeType
func (nt NodeType) String() string {
	switch nt {
	case NodeTypeUnknown:
		return "Unknown"
	case NodeTypeField:
		return "Field"
	case NodeTypeSort:
		return "Sort"
	case NodeTypeSortField:
		return "SortField"
	case NodeTypeFilter:
		return "Filter"
	case NodeTypeLogicalOp:
		return "LogicalOp"
	case NodeTypeComparisonOp:
		return "ComparisonOp"
	case NodeTypeCheckerOp:
		return "CheckerOp"
	case NodeTypeString:
		return "String"
	case NodeTypeIdentifier:
		return "Identifier"
	case NodeTypeFieldList:
		return "FieldList"
	default:
		return "Unknown"
	}
}

// Node is the base interface for all AST nodes
type Node interface {
	Type() NodeType
}

// FieldsNode represents the fields part of the query
type FieldsNode struct {
	Fields []string
}

func (n FieldsNode) Type() NodeType {
	return NodeTypeFieldList
}

// SortNode represents the sort part of the query
type SortNode struct {
	Fields []SortFieldNode
}

func (n SortNode) Type() NodeType {
	return NodeTypeSort
}

// SortFieldNode represents a single field in the sort expression
type SortFieldNode struct {
	Field     string
	Direction SortDirection
}

func (n SortFieldNode) Type() NodeType {
	return NodeTypeSortField
}

// FilterNode represents the filter part of the query
type FilterNode struct {
	Root FilterExprNode
}

func (n FilterNode) Type() NodeType {
	return NodeTypeFilter
}

// FilterExprNode is the interface for all filter expression nodes
type FilterExprNode interface {
	Node
}

// LogicalOpNode represents a logical operation (AND, OR, NOT)
type LogicalOpNode struct {
	Operator LogicalOperator
	Left     FilterExprNode
	Right    FilterExprNode // nil for NOT
}

func (n LogicalOpNode) Type() NodeType {
	return NodeTypeLogicalOp
}

// ComparisonNode represents a comparison operation (=, !=, etc.)
type ComparisonNode struct {
	Field    string
	Operator ComparisonOperator
	Value    Node // StringNode or IdentifierNode
}

func (n ComparisonNode) Type() NodeType {
	return NodeTypeComparisonOp
}

// ComparisonPredicatesNode represents a checker operation (LIKE, IN, etc.)
type ComparisonPredicatesNode struct {
	Field    string
	Operator ComparisonPredicatesOperator
	Value    Node // StringNode or IdentifierNode
}

func (n ComparisonPredicatesNode) Type() NodeType {
	return NodeTypeCheckerOp
}

// StringNode represents a string literal
type StringNode struct {
	Value string
}

func (n StringNode) Type() NodeType {
	return NodeTypeString
}

// IdentifierNode represents an identifier
type IdentifierNode struct {
	Value string
}

func (n IdentifierNode) Type() NodeType {
	return NodeTypeIdentifier
}
