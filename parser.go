package qfv

import (
	"errors"
	"fmt"
	"strings"
)

// Parser parses the query parameters and returns the AST
type Parser struct {
	allowedSortFields   map[string]any // Map of allowed field names for sorting
	allowedFieldsFields map[string]any // Map of allowed field names for filters
	allowedFilterFields map[string]any // Map of allowed field names for filters
}

// NewParser creates a new parser with the allowed fields
func NewParser(allowedSortFields []string, allowedFieldsFields []string, allowedFilterFields []string) *Parser {
	sortFields := make(map[string]any, len(allowedSortFields))
	for _, f := range allowedSortFields {
		sortFields[f] = struct{}{}
	}

	fieldsFields := make(map[string]any, len(allowedFieldsFields))
	for _, f := range allowedFieldsFields {
		fieldsFields[f] = struct{}{}
	}

	filterFields := make(map[string]any, len(allowedFilterFields))
	for _, f := range allowedFilterFields {
		filterFields[f] = struct{}{}
	}

	return &Parser{
		allowedSortFields:   sortFields,
		allowedFieldsFields: fieldsFields,
		allowedFilterFields: filterFields,
	}
}

// ParseSort parses the sort parameter
func (p *Parser) ParseSort(input string) (SortNode, error) {
	if input == "" {
		return SortNode{}, nil
	}

	parts := strings.Split(input, ",")
	fields := make([]SortFieldNode, 0, len(parts))

	for _, part := range parts {
		part = strings.TrimSpace(part)
		if part == "" {
			continue
		}

		sortParts := strings.Fields(part)
		if len(sortParts) == 0 {
			return SortNode{}, fmt.Errorf("invalid sort expression: %s", part)
		}

		fieldName := sortParts[0]
		if _, exists := p.allowedSortFields[fieldName]; !exists {
			return SortNode{}, fmt.Errorf("unknown field: %s", fieldName)
		}

		direction := SortAsc
		if len(sortParts) > 1 {
			dirStr := strings.ToUpper(sortParts[1])

			switch dirStr {
			case SortDesc.String():
				direction = SortDesc
			case SortAsc.String():
				direction = SortAsc
			default:
				return SortNode{}, fmt.Errorf("invalid sort direction: %s", dirStr)
			}
		}

		fields = append(fields, SortFieldNode{
			Field:     fieldName,
			Direction: direction,
		})
	}

	return SortNode{Fields: fields}, nil
}

// ParseFields parses the fields parameter
func (p *Parser) ParseFields(input string) (FieldsNode, error) {
	if input == "" {
		return FieldsNode{}, nil
	}

	parts := strings.Split(input, ",")
	fields := make([]string, 0, len(parts))

	for _, part := range parts {
		part = strings.TrimSpace(part)
		if part == "" {
			continue
		}

		if _, exists := p.allowedFieldsFields[part]; !exists {
			return FieldsNode{}, fmt.Errorf("unknown field: %s", part)
		}

		fields = append(fields, part)
	}

	return FieldsNode{Fields: fields}, nil
}

// ParseFilter parses the filter parameter
func (p *Parser) ParseFilter(input string) (FilterNode, error) {
	if input == "" {
		return FilterNode{}, nil
	}

	lex := newLexer(input)
	tokens, err := lex.tokenize()
	if err != nil {
		return FilterNode{}, err
	}

	parser := newFilterParser(tokens, p.allowedFilterFields)
	expr, err := parser.parse()
	if err != nil {
		return FilterNode{}, err
	}

	return FilterNode{Root: expr}, nil
}

// Validate validates the query parameters
func (p *Parser) Validate(sort, fields, filter string) (SortNode, FieldsNode, FilterNode, error) {
	sortNode, err := p.ParseSort(sort)
	if err != nil {
		return SortNode{}, FieldsNode{}, FilterNode{}, fmt.Errorf("sort: %v", err)
	}

	fieldsNode, err := p.ParseFields(fields)
	if err != nil {
		return SortNode{}, FieldsNode{}, FilterNode{}, fmt.Errorf("fields: %v", err)
	}

	filterNode, err := p.ParseFilter(filter)
	if err != nil {
		return SortNode{}, FieldsNode{}, FilterNode{}, fmt.Errorf("filter: %v", err)
	}

	return sortNode, fieldsNode, filterNode, nil
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
