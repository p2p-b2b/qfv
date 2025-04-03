package qfv

// // FilterParser parses the query parameter for filtering
// type FilterParser struct {
// 	allowedFields map[string]any // any because don't allocate memory for struct{}
// }

// // NewFilterParser creates a new parser with the allowed fields
// func NewFilterParser(allowedFields []string) *FilterParser {
// 	filterFields := make(map[string]any, len(allowedFields))

// 	for _, f := range allowedFields {
// 		filterFields[f] = struct{}{}
// 	}

// 	return &FilterParser{
// 		allowedFields: filterFields,
// 	}
// }

// // Parse parses the filter parameter
// func (p *FilterParser) Parse(input string) (FilterNode, error) {
// 	if input == "" {
// 		return FilterNode{}, fmt.Errorf("empty filter expression")
// 	}

// 	lexer := newLexer(input)
// 	tokens, err := lexer.tokenize()
// 	if err != nil {
// 		return FilterNode{}, err
// 	}

// 	parser := newFilterParser(tokens, p.allowedFields)
// 	expr, err := parser.parse()
// 	if err != nil {
// 		return FilterNode{}, err
// 	}

// 	return FilterNode{Expression: expr}, nil
// }

// func (p *FilterParser) Validate(input string) (FilterNode, error) {
// 	if input == "" {
// 		return FilterNode{}, nil
// 	}

// 	node, err := p.Parse(input)
// 	if err != nil {
// 		return FilterNode{}, err
// 	}

// 	return node, nil
// }

// // filterParser is a parser for filter expressions
// type filterParser struct {
// 	tokens        []Token
// 	lenTokens     int
// 	pos           int
// 	allowedFields map[string]any
// }

// // newFilterParser creates a new filter parser
// func newFilterParser(tokens []Token, allowedFields map[string]any) *filterParser {
// 	return &filterParser{
// 		tokens:        tokens,
// 		lenTokens:     len(tokens),
// 		pos:           0,
// 		allowedFields: allowedFields,
// 	}
// }

// // parse parses the filter expression
// func (p *filterParser) parse() (FilterNode, error) {
// 	return p.parseExpression()
// }

// // parseExpression parses a logical expression
// func (p *filterParser) parseExpression() (FilterNode, error) {
// 	left, err := p.parseComparison()
// 	if err != nil {
// 		return FilterNode{}, err
// 	}

// 	for p.pos < p.lenTokens && p.tokens[p.pos].Type == TokenLogicalOperation {
// 		op := p.tokens[p.pos].Value
// 		p.pos++

// 		if op == OperatorNot.String() {
// 			right, err := p.parseComparison()
// 			if err != nil {
// 				return FilterNode{}, err
// 			}

// 			left = BinaryOperatorNode{
// 				Operator: OperatorNot,
// 				Left:     right,
// 				Right:    nil,
// 			}
// 		} else {
// 			right, err := p.parseComparison()
// 			if err != nil {
// 				return FilterNode{}, err
// 			}

// 			left = BinaryOperatorNode{
// 				Operator: Operator(op),
// 				Left:     left,
// 				Right:    right,
// 			}
// 		}
// 	}

// 	return left, nil
// }

// // parseComparison parses a comparison expression
// func (p *filterParser) parseComparison() (FilterNode, error) {
// 	if p.pos >= p.lenTokens {
// 		return nil, errors.New("unexpected end of input")
// 	}

// 	if p.tokens[p.pos].Type == TokenLPAREN {
// 		p.pos++ // Skip '('
// 		expr, err := p.parseExpression()
// 		if err != nil {
// 			return nil, err
// 		}

// 		if p.pos >= p.lenTokens || p.tokens[p.pos].Type != TokenRPAREN {
// 			return nil, errors.New("missing closing parenthesis")
// 		}

// 		p.pos++ // Skip ')'
// 		return expr, nil
// 	}

// 	if p.tokens[p.pos].Type == TokenLogicalOperation && p.tokens[p.pos].Value == OperatorNot.String() {
// 		p.pos++ // Skip NOT
// 		expr, err := p.parseComparison()
// 		if err != nil {
// 			return nil, err
// 		}

// 		return LogicalOperationNode{
// 			Operator: OperatorNot,
// 			Left:     expr,
// 			Right:    nil,
// 		}, nil
// 	}

// 	if p.tokens[p.pos].Type != TokenIdentifier {
// 		return nil, fmt.Errorf("expected field name, got %s", p.tokens[p.pos].Value)
// 	}

// 	field := p.tokens[p.pos].Value
// 	if _, exists := p.allowedFields[field]; !exists {
// 		return nil, fmt.Errorf("unknown field: %s", field)
// 	}
// 	p.pos++

// 	if p.pos >= p.lenTokens || p.tokens[p.pos].Type != TokenOperator {
// 		return nil, errors.New("expected comparison operator")
// 	}

// 	operator := Operator(p.tokens[p.pos].Value)
// 	p.pos++

// 	if p.pos >= p.lenTokens {
// 		return nil, errors.New("expected value after operator")
// 	}

// 	var value Node
// 	switch p.tokens[p.pos].Type {
// 	case TokenString:
// 		value = StringNode{Value: p.tokens[p.pos].Value}
// 	case TokenIdentifier:
// 		value = IdentifierNode{Value: p.tokens[p.pos].Value}
// 	case TokenBoolean:
// 		value = BooleanNode{Value: p.tokens[p.pos].Value == "true"}
// 	case TokenNumber:
// 		v, err := strconv.ParseFloat(p.tokens[p.pos].Value, 64)
// 		if err != nil {
// 			return nil, fmt.Errorf("invalid number: %s", p.tokens[p.pos].Value)
// 		}
// 		value = NumberNode{Value: v}

// 	default:
// 		return nil, fmt.Errorf("expected string or identifier, got %s", p.tokens[p.pos].Type)
// 	}

// 	p.pos++

// 	return ComparisonNode{
// 		Field:    field,
// 		Operator: operator,
// 		Value:    value,
// 	}, nil
// }
