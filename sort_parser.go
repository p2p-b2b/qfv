package qfv

import (
	"fmt"
	"strings"
)

// SortFieldNode represents a single field in the sort expression
type SortFieldNode struct {
	Field     string
	Direction SortDirection
}

func (n SortFieldNode) Type() NodeType {
	return NodeTypeSortField
}

// SortNode represents the sort part of the query
type SortNode struct {
	Fields []SortFieldNode
}

func (n SortNode) Type() NodeType {
	return NodeTypeSort
}

// SortParser parses the query parameter for sorting
type SortParser struct {
	allowedFields map[string]any // any because don't allocate memory for struct{}
}

// NewSortParser creates a new parser with the allowed fields for sorting
func NewSortParser(allowedFields []string) *SortParser {
	sortFields := make(map[string]any, len(allowedFields))

	for _, f := range allowedFields {
		sortFields[f] = struct{}{}
	}

	return &SortParser{
		allowedFields: sortFields,
	}
}

// Parse parses the sort parameter
func (p *SortParser) Parse(input string) (SortNode, error) {
	if input == "" {
		return SortNode{}, fmt.Errorf("empty input expression")
	}

	parts := strings.Split(input, ",")
	fields := make([]SortFieldNode, 0, len(parts))

	for _, part := range parts {
		part = strings.TrimSpace(part)
		if part == "" {
			return SortNode{}, fmt.Errorf("empty sort expression: %s", input)
		}

		sortParts := strings.Fields(part)
		if len(sortParts) == 0 {
			return SortNode{}, fmt.Errorf("invalid sort expression: %s", part)
		}

		if len(sortParts) > 2 {
			return SortNode{}, fmt.Errorf("invalid sort expression: %s", part)
		}

		fieldName := sortParts[0]
		if _, exists := p.allowedFields[fieldName]; !exists {
			return SortNode{}, fmt.Errorf("unknown field: %s", fieldName)
		}

		direction := SortAsc
		if len(sortParts) == 1 {
			return SortNode{}, fmt.Errorf("missing sort direction for field: %s", fieldName)
		}

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

func (p *SortParser) Validate(input string) (SortNode, error) {
	if input == "" {
		return SortNode{}, fmt.Errorf("empty input expression")
	}

	node, err := p.Parse(input)
	if err != nil {
		return SortNode{}, err
	}

	return node, nil
}
