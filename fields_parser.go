package qfv

import (
	"fmt"
	"strings"
)

// FieldsNode represents the fields part of the query
type FieldsNode struct {
	Fields []string
}

func (n FieldsNode) Type() NodeType {
	return NodeTypeFieldList
}

// FieldsParser parses the query parameter for fields
type FieldsParser struct {
	allowedFieldsFields map[string]any // any because don't allocate memory for struct{}
}

// NewFieldsParser creates a new parser with the allowed fields
func NewFieldsParser(allowedFields []string) *FieldsParser {
	fieldsFields := make(map[string]any, len(allowedFields))

	for _, f := range allowedFields {
		fieldsFields[f] = struct{}{}
	}

	return &FieldsParser{
		allowedFieldsFields: fieldsFields,
	}
}

// Parse parses the fields parameter
func (p *FieldsParser) Parse(input string) (FieldsNode, error) {
	if input == "" {
		return FieldsNode{}, fmt.Errorf("empty input expression")
	}

	parts := strings.Split(input, ",")
	fields := make([]string, 0, len(parts))

	for _, part := range parts {
		part = strings.TrimSpace(part)
		if part == "" {
			return FieldsNode{}, fmt.Errorf("empty field expression: %s", input)
		}

		if _, exists := p.allowedFieldsFields[part]; !exists {
			return FieldsNode{}, fmt.Errorf("unknown field: %s", part)
		}

		fields = append(fields, part)
	}

	return FieldsNode{Fields: fields}, nil
}
