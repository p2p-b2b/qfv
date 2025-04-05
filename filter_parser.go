package qfv

import (
	"fmt"
	"reflect"
	"strconv"
	"strings"
	"text/scanner"
)

// FilterParser parses the query parameter for filtering
type FilterParser struct {
	allowedFields map[string]any // any because don't allocate memory for struct{}
	lexer         *Lexer
	currentToken  Token
	errors        []error
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

// Parse parses the filter query and returns the AST
func (p *FilterParser) Parse(input string) (Node, error) {
	p.lexer = NewLexer(input)
	p.lexer.Parse()
	p.errors = nil

	// Check for illegal tokens in the input
	for _, token := range p.lexer.tokens {
		if token.Type == TokenIllegal {
			p.addError(fmt.Errorf("illegal token: %s", token.Value))
		}
	}

	p.nextToken()

	if p.currentToken.Type == TokenEOF {
		return nil, fmt.Errorf("empty filter expression")
	}

	node := p.parseExpression()

	if len(p.errors) > 0 {
		return nil, fmt.Errorf("parse errors: %v", p.errors)
	}

	return node, nil
}

// nextToken advances to the next token
func (p *FilterParser) nextToken() {
	p.currentToken = p.lexer.Next()
}

// peekToken returns the next token without consuming it
func (p *FilterParser) peekToken() Token {
	return p.lexer.Peek()
}

// expect checks if the current token is of the expected type
func (p *FilterParser) expect(tokenType TokenType) bool {
	if p.currentToken.Type == tokenType {
		p.nextToken()
		return true
	}
	p.addError(fmt.Errorf("expected %s, got %s", tokenType, p.currentToken.Type))
	return false
}

// addError adds an error to the error list
func (p *FilterParser) addError(err error) {
	p.errors = append(p.errors, err)
}

// parseExpression parses an expression
func (p *FilterParser) parseExpression() Node {
	return p.parseLogicalOr()
}

// parseLogicalOr parses OR expressions
func (p *FilterParser) parseLogicalOr() Node {
	left := p.parseLogicalAnd()

	for p.currentToken.Type == TokenOperatorOr {
		pos := p.currentToken.Pos
		operator := p.currentToken.Type
		p.nextToken()
		right := p.parseLogicalAnd()
		left = &BinaryOperatorNode{
			baseNode: baseNode{pos: pos},
			Left:     left,
			Right:    right,
			Operator: operator,
		}
	}

	return left
}

// parseLogicalAnd parses AND expressions
func (p *FilterParser) parseLogicalAnd() Node {
	left := p.parseComparison()

	for p.currentToken.Type == TokenOperatorAnd {
		pos := p.currentToken.Pos
		operator := p.currentToken.Type
		p.nextToken()
		right := p.parseComparison()
		left = &BinaryOperatorNode{
			baseNode: baseNode{pos: pos},
			Left:     left,
			Right:    right,
			Operator: operator,
		}
	}

	return left
}

// parseComparison parses comparison expressions
func (p *FilterParser) parseComparison() Node {
	// Check for NOT operator
	if p.currentToken.Type == TokenOperatorNot {
		pos := p.currentToken.Pos
		p.nextToken()
		expr := p.parseComparison()
		return &UnaryOperatorNode{
			baseNode: baseNode{pos: pos},
			Operator: TokenOperatorNot,
			X:        expr,
		}
	}

	// Check for parenthesized expressions
	if p.currentToken.Type == TokenLPAREN {
		pos := p.currentToken.Pos
		p.nextToken()
		expr := p.parseExpression()
		if !p.expect(TokenRPAREN) {
			p.addError(fmt.Errorf("expected closing parenthesis"))
		}
		return &GroupNode{
			baseNode:   baseNode{pos: pos},
			Expression: expr,
		}
	}

	// Parse field comparison
	if p.currentToken.Type == TokenIdentifier {
		field := &IdentifierNode{
			baseNode: baseNode{pos: p.currentToken.Pos},
			Name:     p.currentToken.Value,
		}
		p.nextToken()

		// Check if field is allowed
		if _, ok := p.allowedFields[field.Name]; !ok {
			p.addError(fmt.Errorf("field %s is not allowed", field.Name))
		}

		// Handle different operators
		switch p.currentToken.Type {
		case TokenOperatorEqual, TokenOperatorNotEqual, TokenOperatorNotEqualAlias,
			TokenOperatorLessThan, TokenOperatorLessThanOrEqualTo,
			TokenOperatorGreaterThan, TokenOperatorGreaterThanOrEqualTo:
			return p.parseComparisonOperator(field)
		case TokenOperatorLike:
			return p.parseLikeOperator(field, false)
		case TokenOperatorIn:
			return p.parseInOperator(field, false)
		case TokenOperatorBetween:
			return p.parseBetweenOperator(field, false)
		case TokenOperatorIsNull:
			return p.parseIsNullOperator(field)
		case TokenOperatorDistinct:
			return p.parseDistinctOperator(field, false)
		case TokenOperatorNot:
			// Handle NOT operators (NOT IN, NOT BETWEEN, NOT LIKE)
			pos := p.currentToken.Pos
			p.nextToken() // Consume NOT

			switch p.currentToken.Type {
			case TokenOperatorIn:
				return p.parseInOperator(field, true)
			case TokenOperatorBetween:
				return p.parseBetweenOperator(field, true)
			case TokenOperatorLike:
				return p.parseLikeOperator(field, true)
			default:
				p.addError(fmt.Errorf("unexpected token after NOT: %s", p.currentToken.Type))
				return &UnaryOperatorNode{
					baseNode: baseNode{pos: pos},
					Operator: TokenOperatorNot,
					X:        field,
				}
			}
		default:
			p.addError(fmt.Errorf("unexpected token after field: %s", p.currentToken.Type))
			return field
		}
	}

	// Parse literal
	return p.parsePrimary()
}

// parseComparisonOperator parses comparison operators (=, <>, !=, <, <=, >, >=)
func (p *FilterParser) parseComparisonOperator(field Node) Node {
	pos := p.currentToken.Pos
	operator := p.currentToken.Type
	p.nextToken()
	right := p.parsePrimary()
	return &BinaryOperatorNode{
		baseNode: baseNode{pos: pos},
		Left:     field,
		Right:    right,
		Operator: operator,
	}
}

// parseLikeOperator parses LIKE operator
func (p *FilterParser) parseLikeOperator(field Node, isNot bool) Node {
	pos := p.currentToken.Pos
	p.nextToken()
	pattern := p.parsePrimary()
	operator := TokenOperatorLike
	if isNot {
		return &UnaryOperatorNode{
			baseNode: baseNode{pos: pos},
			Operator: TokenOperatorNot,
			X: &BinaryOperatorNode{
				baseNode: baseNode{pos: pos},
				Left:     field,
				Right:    pattern,
				Operator: operator,
			},
		}
	}
	return &BinaryOperatorNode{
		baseNode: baseNode{pos: pos},
		Left:     field,
		Right:    pattern,
		Operator: operator,
	}
}

// parseInOperator parses IN operator
func (p *FilterParser) parseInOperator(field Node, isNot bool) Node {
	pos := p.currentToken.Pos
	p.nextToken()

	if !p.expect(TokenLPAREN) {
		p.addError(fmt.Errorf("expected opening parenthesis after IN"))
		return field
	}

	var values []Node

	// Parse the first value
	values = append(values, p.parsePrimary())

	// Parse additional values
	for p.currentToken.Type == TokenComma {
		p.nextToken()
		values = append(values, p.parsePrimary())
	}

	if !p.expect(TokenRPAREN) {
		p.addError(fmt.Errorf("expected closing parenthesis after IN values"))
	}

	inNode := &InNode{
		baseNode: baseNode{pos: pos},
		Field:    field,
		IsNot:    false, // Always false, we'll wrap it in a UnaryOperatorNode if isNot is true
		Values:   values,
	}

	if isNot {
		return &UnaryOperatorNode{
			baseNode: baseNode{pos: pos},
			Operator: TokenOperatorNot,
			X:        inNode,
		}
	}

	return inNode
}

// parseBetweenOperator parses BETWEEN operator
func (p *FilterParser) parseBetweenOperator(field Node, isNot bool) Node {
	pos := p.currentToken.Pos
	p.nextToken()

	lower := p.parsePrimary()

	if !p.expect(TokenOperatorAnd) {
		p.addError(fmt.Errorf("expected AND in BETWEEN expression"))
		return field
	}

	upper := p.parsePrimary()

	betweenNode := &BetweenNode{
		baseNode: baseNode{pos: pos},
		Field:    field,
		Lower:    lower,
		Upper:    upper,
		IsNot:    false, // Always false, we'll wrap it in a UnaryOperatorNode if isNot is true
	}

	if isNot {
		return &UnaryOperatorNode{
			baseNode: baseNode{pos: pos},
			Operator: TokenOperatorNot,
			X:        betweenNode,
		}
	}

	return betweenNode
}

// parseIsNullOperator parses IS NULL operator
func (p *FilterParser) parseIsNullOperator(field Node) Node {
	pos := p.currentToken.Pos
	p.nextToken()

	isNot := false
	if p.currentToken.Type == TokenOperatorNot {
		isNot = true
		p.nextToken()
	}

	// Check for NULL
	if p.currentToken.Type == TokenIdentifier && strings.ToUpper(p.currentToken.Value) == "NULL" {
		p.nextToken()
		return &IsNullNode{
			baseNode: baseNode{pos: pos},
			Field:    field,
			IsNot:    isNot,
		}
	}

	p.addError(fmt.Errorf("expected NULL after IS"))
	return field
}

// parseDistinctOperator parses DISTINCT operator
func (p *FilterParser) parseDistinctOperator(field Node, isNot bool) Node {
	pos := p.currentToken.Pos
	p.nextToken()

	return &DistinctNode{
		baseNode: baseNode{pos: pos},
		Field:    field,
		IsNot:    isNot,
	}
}

// parsePrimary parses primary expressions (literals)
func (p *FilterParser) parsePrimary() Node {
	switch p.currentToken.Type {
	case TokenString:
		node := &LiteralNode{
			baseNode: baseNode{pos: p.currentToken.Pos},
			Value:    strings.Trim(p.currentToken.Value, "'"),
			Kind:     reflect.String,
			Text:     p.currentToken.Value,
		}
		p.nextToken()
		return node

	case TokenInt:
		val, err := strconv.ParseInt(p.currentToken.Value, 10, 64)
		if err != nil {
			p.addError(fmt.Errorf("invalid integer: %s", p.currentToken.Value))
		}
		node := &LiteralNode{
			baseNode: baseNode{pos: p.currentToken.Pos},
			Value:    val,
			Kind:     reflect.Int64,
			Text:     p.currentToken.Value,
		}
		p.nextToken()
		return node

	case TokenFloat:
		val, err := strconv.ParseFloat(p.currentToken.Value, 64)
		if err != nil {
			p.addError(fmt.Errorf("invalid float: %s", p.currentToken.Value))
		}
		node := &LiteralNode{
			baseNode: baseNode{pos: p.currentToken.Pos},
			Value:    val,
			Kind:     reflect.Float64,
			Text:     p.currentToken.Value,
		}
		p.nextToken()
		return node

	case TokenBoolean:
		val := strings.ToUpper(p.currentToken.Value) == "TRUE" || strings.ToUpper(p.currentToken.Value) == "YES"
		node := &LiteralNode{
			baseNode: baseNode{pos: p.currentToken.Pos},
			Value:    val,
			Kind:     reflect.Bool,
			Text:     p.currentToken.Value,
		}
		p.nextToken()
		return node

	default:
		p.addError(fmt.Errorf("unexpected token: %s", p.currentToken.Type))
		// Skip the token to avoid infinite loops
		p.nextToken()
		return &LiteralNode{
			baseNode: baseNode{pos: scanner.Position{}},
			Value:    nil,
			Kind:     0,
			Text:     "",
		}
	}
}
