package qfv

import (
	"errors"
	"fmt"
)

// FilterNode represents the filter part of the query
type FilterNode struct {
	Root FilterExprNode
}

func (n FilterNode) Type() NodeType {
	return NodeTypeFilter
}

// FilterParser parses the query parameter for filtering
type FilterParser struct {
	allowedFields map[string]any // any because don't allocate memory for struct{}
}

// NewFilterParser creates a new parser with the allowed fields
func NewFilterParser(allowedFields []string) *FilterParser {
	filterFields := make(map[string]any, len(allowedFields))

	for _, f := range allowedFields {
		filterFields[f] = struct{}{}
	}

	return &FilterParser{
		allowedFields: filterFields,
	}
}

// Parse parses the filter parameter
func (p *FilterParser) Parse(input string) (FilterNode, error) {
	if input == "" {
		return FilterNode{}, nil
	}

	lex := newLexer(input)
	tokens, err := lex.tokenize()
	if err != nil {
		return FilterNode{}, err
	}

	parser := newFilterParser(tokens, p.allowedFields)
	expr, err := parser.parse()
	if err != nil {
		return FilterNode{}, err
	}

	return FilterNode{Root: expr}, nil
}

func (p *FilterParser) Validate(input string) (FilterNode, error) {
	if input == "" {
		return FilterNode{}, nil
	}

	node, err := p.Parse(input)
	if err != nil {
		return FilterNode{}, err
	}

	return node, nil
}

// filterParser is a parser for filter expressions
type filterParser struct {
	tokens        []Token
	pos           int
	allowedFields map[string]any
}

// newFilterParser creates a new filter parser
func newFilterParser(tokens []Token, allowedFields map[string]any) *filterParser {
	return &filterParser{
		tokens:        tokens,
		pos:           0,
		allowedFields: allowedFields,
	}
}

// parse parses the filter expression
func (p *filterParser) parse() (FilterExprNode, error) {
	return p.parseExpression()
}

// parseExpression parses a logical expression
func (p *filterParser) parseExpression() (FilterExprNode, error) {
	left, err := p.parseComparison()
	if err != nil {
		return nil, err
	}

	for p.pos < len(p.tokens) && p.tokens[p.pos].Type == "LOGICAL_OP" {
		op := p.tokens[p.pos].Value
		p.pos++

		if op == "NOT" {
			right, err := p.parseComparison()
			if err != nil {
				return nil, err
			}
			left = LogicalOpNode{
				Operator: LogicalOperator(op),
				Left:     right,
				Right:    nil,
			}
		} else {
			right, err := p.parseComparison()
			if err != nil {
				return nil, err
			}
			left = LogicalOpNode{
				Operator: LogicalOperator(op),
				Left:     left,
				Right:    right,
			}
		}
	}

	return left, nil
}

// parseComparison parses a comparison expression
func (p *filterParser) parseComparison() (FilterExprNode, error) {
	if p.pos >= len(p.tokens) {
		return nil, errors.New("unexpected end of input")
	}

	if p.tokens[p.pos].Type == "LPAREN" {
		p.pos++ // Skip '('
		expr, err := p.parseExpression()
		if err != nil {
			return nil, err
		}

		if p.pos >= len(p.tokens) || p.tokens[p.pos].Type != "RPAREN" {
			return nil, errors.New("missing closing parenthesis")
		}
		p.pos++ // Skip ')'
		return expr, nil
	}

	if p.tokens[p.pos].Type == "LOGICAL_OP" && p.tokens[p.pos].Value == "NOT" {
		p.pos++ // Skip NOT
		expr, err := p.parseComparison()
		if err != nil {
			return nil, err
		}
		return LogicalOpNode{
			Operator: OpNot,
			Left:     expr,
			Right:    nil,
		}, nil
	}

	if p.tokens[p.pos].Type != "IDENTIFIER" {
		return nil, fmt.Errorf("expected field name, got %s", p.tokens[p.pos].Value)
	}

	field := p.tokens[p.pos].Value
	if _, exists := p.allowedFields[field]; !exists {
		return nil, fmt.Errorf("unknown field: %s", field)
	}
	p.pos++

	if p.pos >= len(p.tokens) || p.tokens[p.pos].Type != "OPERATOR" {
		return nil, errors.New("expected comparison operator")
	}

	operator := ComparisonOperator(p.tokens[p.pos].Value)
	p.pos++

	if p.pos >= len(p.tokens) {
		return nil, errors.New("expected value after operator")
	}

	var value Node
	switch p.tokens[p.pos].Type {
	case "STRING":
		value = StringNode{Value: p.tokens[p.pos].Value}
	case "IDENTIFIER":
		value = IdentifierNode{Value: p.tokens[p.pos].Value}
	default:
		return nil, fmt.Errorf("expected string or identifier, got %s", p.tokens[p.pos].Type)
	}

	p.pos++

	return ComparisonNode{
		Field:    field,
		Operator: operator,
		Value:    value,
	}, nil
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
