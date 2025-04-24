package qfv

import (
	"fmt"
	"reflect"
	"strconv"
	"strings"
	"text/scanner"
)

type QFVFilterError struct {
	Field   string
	Message string
}

func (e *QFVFilterError) Error() string {
	if e.Field != "" {
		return fmt.Sprintf("error on field '%s': %s", e.Field, e.Message)
	}

	return fmt.Sprintf("error: %s", e.Message)
}

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
			p.addError(&QFVFilterError{Field: token.Value, Message: "illegal token"})
		}
	}

	p.nextToken()

	if p.currentToken.Type == TokenEOF {
		return nil, &QFVFilterError{Message: "empty filter expression"}
	}

	node := p.parseExpression()

	if len(p.errors) > 0 {
		return nil, &QFVFilterError{Message: fmt.Sprintf("parsing errors: %v", p.errors)}
	}

	return node, nil
}

// nextToken advances to the next token
func (p *FilterParser) nextToken() {
	p.currentToken = p.lexer.Next()
}

// expect checks if the current token is of the expected type
func (p *FilterParser) expect(tokenType TokenType) bool {
	if p.currentToken.Type == tokenType {
		p.nextToken()
		return true
	}

	p.addError(&QFVFilterError{Message: fmt.Sprintf("expected %s, got %s", tokenType, p.currentToken.Type)})
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
			p.addError(&QFVFilterError{Message: "expected closing parenthesis"})
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
			p.addError(&QFVFilterError{Field: field.Name, Message: "field not allowed"})
		}

		// Handle different operators
		switch p.currentToken.Type {
		case TokenOperatorEqual, TokenOperatorNotEqual, TokenOperatorNotEqualAlias,
			TokenOperatorLessThan, TokenOperatorLessThanOrEqualTo,
			TokenOperatorGreaterThan, TokenOperatorGreaterThanOrEqualTo:
			return p.parseComparisonOperator(field)
		case TokenOperatorLike:
			p.nextToken() // Consume LIKE
			return p.parseLikeOperator(field)
		case TokenOperatorIn:
			p.nextToken() // Consume IN
			return p.parseInOperator(field)
		case TokenOperatorBetween:
			p.nextToken() // Consume BETWEEN
			return p.parseBetweenOperator(field)
		case TokenOperatorIsNull:
			p.nextToken() // Consume IS
			return p.parseIsNullOperator(field)
		case TokenOperatorDistinct:
			p.nextToken() // Consume DISTINCT
			return p.parseDistinctOperator(field)
		case TokenOperatorSimilarTo:
			p.nextToken() // Consume SIMILAR
			return p.parseSimilarToOperator(field)
		case TokenOperatorRegexMatchCS, TokenOperatorNotRegexMatchCS, TokenOperatorRegexMatchCI, TokenOperatorNotRegexMatchCI:
			opToken := p.currentToken
			p.nextToken() // Consume regex operator
			patternNode := p.parsePrimary()

			// Check if the pattern is a string literal
			patternLiteral, ok := patternNode.(*LiteralNode)
			if !ok || patternLiteral.Kind != reflect.String {
				p.addError(&QFVFilterError{Message: fmt.Sprintf("expected string pattern for regex operator %s, got %s", opToken.Type, patternNode.Type())})
				// Return the field node or the invalid pattern node on error
				// Returning the pattern node might give slightly better context
				return patternNode
			}

			return &RegexMatchNode{
				baseNode:          baseNode{pos: opToken.Pos},
				Field:             field,
				Pattern:           patternNode, // Use the parsed node
				IsNot:             opToken.Type == TokenOperatorNotRegexMatchCS || opToken.Type == TokenOperatorNotRegexMatchCI,
				IsCaseInsensitive: opToken.Type == TokenOperatorRegexMatchCI || opToken.Type == TokenOperatorNotRegexMatchCI,
			}
		case TokenOperatorNot:
			// Handle NOT operators (NOT IN, NOT BETWEEN, NOT LIKE, NOT SIMILAR TO, IS NOT NULL, NOT DISTINCT FROM)
			notPos := p.currentToken.Pos
			p.nextToken() // Consume NOT

			var notExpr Node
			switch p.currentToken.Type {
			case TokenOperatorIn:
				p.nextToken() // Consume IN
				notExpr = p.parseInOperator(field)
			case TokenOperatorBetween:
				p.nextToken() // Consume BETWEEN
				notExpr = p.parseBetweenOperator(field)
			case TokenOperatorLike:
				p.nextToken() // Consume LIKE
				notExpr = p.parseLikeOperator(field)
			case TokenOperatorSimilarTo:
				p.nextToken()                             // Consume SIMILAR
				notExpr = p.parseSimilarToOperator(field) // Expects TO next
			case TokenOperatorIsNull: // Handle IS NOT NULL here
				p.nextToken() // Consume IS
				// parseIsNullOperator handles the NOT internally now based on token sequence
				notExpr = p.parseIsNullOperator(field)
				// Check if parseIsNullOperator correctly identified IS NOT NULL
				if isNullNode, ok := notExpr.(*IsNullNode); !ok || !isNullNode.IsNot {
					// If parseIsNullOperator didn't return an IsNullNode with IsNot=true,
					// it means the sequence wasn't "IS NOT NULL".
					// The error would have been added inside parseIsNullOperator.
					// We might return the field or a generic error node, but returning
					// the result from parseIsNullOperator (which might be 'field') is consistent.
					return notExpr // Return whatever parseIsNullOperator returned on error
				}
				// If it was IS NOT NULL, we don't need to wrap it again
				return notExpr
			case TokenOperatorDistinct: // Handle NOT DISTINCT FROM here
				p.nextToken() // Consume DISTINCT
				// parseDistinctOperator handles the FROM internally
				notExpr = p.parseDistinctOperator(field)
				// Similar to IS NOT NULL, check if parseDistinctOperator failed
				if _, ok := notExpr.(*DistinctNode); !ok {
					return notExpr // Return error node or field
				}
			default:
				p.addError(&QFVFilterError{Message: fmt.Sprintf("unexpected token after NOT: %s", p.currentToken.Type)})
				// If NOT is followed by something unexpected, return a unary NOT node with the field
				// This might not be the most robust error handling, but fits the previous pattern.
				return &UnaryOperatorNode{
					baseNode: baseNode{pos: notPos},
					Operator: TokenOperatorNot,
					X:        field, // Apply NOT to the field itself? Or error?
				}
			}

			// Wrap the parsed expression (LIKE, IN, BETWEEN, SIMILAR TO, DISTINCT) in a UnaryOperatorNode(NOT)
			// Skip wrapping if it was handled internally (IS NOT NULL)
			if _, isIsNull := notExpr.(*IsNullNode); !isIsNull {
				return &UnaryOperatorNode{
					baseNode: baseNode{pos: notPos},
					Operator: TokenOperatorNot,
					X:        notExpr,
				}
			}
			// For IS NOT NULL, return the node directly as IsNot is set inside
			return notExpr

		default:
			p.addError(&QFVFilterError{Field: field.Name, Message: "unexpected token after field"})
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

// parseSimilarToOperator parses SIMILAR TO operator
// Expects the current token to be TO after SIMILAR was consumed.
func (p *FilterParser) parseSimilarToOperator(field Node) Node {
	pos := p.lexer.Current().Pos // Use position of SIMILAR token (already consumed)
	if p.currentToken.Type != TokenIdentifier || strings.ToUpper(p.currentToken.Value) != "TO" {
		p.addError(&QFVFilterError{Message: "expected TO after SIMILAR"})
		return field // Return field on error
	}
	p.nextToken() // Consume TO
	pattern := p.parsePrimary()
	return &SimilarToNode{
		baseNode: baseNode{pos: pos},
		Field:    field,
		Pattern:  pattern,
		IsNot:    false, // NOT is handled by parseComparison
	}
}

// parseLikeOperator parses LIKE operator
// Expects the current token to be the pattern after LIKE was consumed.
func (p *FilterParser) parseLikeOperator(field Node) Node {
	pos := p.lexer.Current().Pos // Use position of LIKE token (already consumed)
	pattern := p.parsePrimary()
	return &BinaryOperatorNode{
		baseNode: baseNode{pos: pos},
		Left:     field,
		Right:    pattern,
		Operator: TokenOperatorLike,
	}
}

// parseInOperator parses IN operator
// Expects the current token to be LPAREN after IN was consumed.
func (p *FilterParser) parseInOperator(field Node) Node {
	pos := p.lexer.Current().Pos // Use position of IN token (already consumed)
	if !p.expect(TokenLPAREN) {
		p.addError(&QFVFilterError{Message: "expected opening parenthesis after IN"})
		return field
	}

	var values []Node
	// Parse the first value
	if p.currentToken.Type == TokenRPAREN {
		p.addError(&QFVFilterError{Message: "expected at least one value after IN ("})
	} else {
		values = append(values, p.parsePrimary())
	}

	// Parse additional values
	for p.currentToken.Type == TokenComma {
		p.nextToken()
		if p.currentToken.Type == TokenRPAREN { // Handle trailing comma
			p.addError(&QFVFilterError{Message: "unexpected closing parenthesis after comma in IN list"})
			break
		}
		values = append(values, p.parsePrimary())
	}

	if !p.expect(TokenRPAREN) {
		p.addError(&QFVFilterError{Message: "expected closing parenthesis after IN values"})
	}

	return &InNode{
		baseNode: baseNode{pos: pos},
		Field:    field,
		IsNot:    false, // NOT is handled by parseComparison
		Values:   values,
	}
}

// parseBetweenOperator parses BETWEEN operator
// Expects the current token to be the lower bound after BETWEEN was consumed.
func (p *FilterParser) parseBetweenOperator(field Node) Node {
	pos := p.lexer.Current().Pos // Use position of BETWEEN token (already consumed)
	lower := p.parsePrimary()

	if !p.expect(TokenOperatorAnd) {
		p.addError(&QFVFilterError{Message: "expected AND in BETWEEN expression"})
		return field
	}

	upper := p.parsePrimary()

	return &BetweenNode{
		baseNode: baseNode{pos: pos},
		Field:    field,
		Lower:    lower,
		Upper:    upper,
		IsNot:    false, // NOT is handled by parseComparison
	}
}

// parseIsNullOperator parses IS [NOT] NULL operator
// Expects the current token to be NOT or NULL after IS was consumed.
func (p *FilterParser) parseIsNullOperator(field Node) Node {
	pos := p.lexer.Current().Pos // Use position of IS token (already consumed)
	isNot := false
	if p.currentToken.Type == TokenOperatorNot {
		isNot = true
		p.nextToken() // Consume NOT
	}

	// Check for NULL
	if p.currentToken.Type == TokenIdentifier && strings.ToUpper(p.currentToken.Value) == "NULL" {
		p.nextToken() // Consume NULL
		return &IsNullNode{
			baseNode: baseNode{pos: pos},
			Field:    field,
			IsNot:    isNot,
		}
	}

	if isNot {
		p.addError(&QFVFilterError{Message: "expected NULL after IS NOT"})
	} else {
		p.addError(&QFVFilterError{Message: "expected NULL or NOT NULL after IS"})
	}
	return field // Return field on error
}

// parseDistinctOperator parses DISTINCT FROM operator
// Expects the current token to be FROM after DISTINCT was consumed.
func (p *FilterParser) parseDistinctOperator(field Node) Node {
	pos := p.lexer.Current().Pos // Use position of DISTINCT token (already consumed)
	// Expect FROM (treated as identifier by lexer)
	if p.currentToken.Type != TokenIdentifier || strings.ToUpper(p.currentToken.Value) != "FROM" {
		p.addError(&QFVFilterError{Message: "expected FROM after DISTINCT"})
		return field // Return field on error
	}
	p.nextToken() // Consume FROM
	// Parse the value being compared
	_ = p.parsePrimary() // Consume the value but ignore it for now

	// Note: The DistinctNode currently only holds the field and IsNot.
	// We might need to adjust the AST node or create a new one if
	// the `FROM value` part is significant for evaluation.
	// For now, just return a basic DistinctNode.
	return &DistinctNode{
		baseNode: baseNode{pos: pos},
		Field:    field,
		IsNot:    false, // NOT is handled by parseComparison
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
			p.addError(&QFVFilterError{Message: fmt.Sprintf("invalid integer: %s", p.currentToken.Value)})
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
			p.addError(&QFVFilterError{Message: fmt.Sprintf("invalid float: %s", p.currentToken.Value)})
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
		p.addError(&QFVFilterError{Message: fmt.Sprintf("unexpected token: %s", p.currentToken.Type)})
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
