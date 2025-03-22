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

// Validator is the interface for all validators
// It defines a method to validate the fields of a query
// and return a Node representing the parsed fields.
// The Validate method takes a string input representing the fields
// and returns a Node and an error.
type Validator interface {
	Validate(fields string) (Node, error)
}
