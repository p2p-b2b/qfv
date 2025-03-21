package qfv

// SortDirection represents the sorting direction in sort expressions
type SortDirection string

const (
	// SortAsc represents ascending sort order
	SortAsc SortDirection = "ASC"

	// SortDesc represents descending sort order
	SortDesc SortDirection = "DESC"
)

// String returns the string representation of the SortDirection
func (sd SortDirection) String() string {
	switch sd {
	case SortAsc:
		return "ASC"
	case SortDesc:
		return "DESC"
	default:
		return "Unknown"
	}
}

// LogicalOperator represents the logical operators in filter expressions
type LogicalOperator string

const (
	// OpAnd represents the AND logical operator
	OpAnd LogicalOperator = "AND"

	// OpOr represents the OR logical operator
	OpOr LogicalOperator = "OR"

	// OpNot represents the NOT logical operator
	OpNot LogicalOperator = "NOT"
)

// String returns the string representation of the LogicalOperator
func (lo LogicalOperator) String() string {
	switch lo {
	case OpAnd:
		return "AND"
	case OpOr:
		return "OR"
	case OpNot:
		return "NOT"
	default:
		return "Unknown"
	}
}

// ComparisonOperator represents the comparison operators in filter expressions
type ComparisonOperator string

const (
	// OpEqual represents the equality operator
	OpEqual ComparisonOperator = "="

	// OpNotEqual represents the inequality operator
	OpNotEqual ComparisonOperator = "<>"

	// OpEqualAlias is an alias for the equality operator
	OpNotEqualAlias ComparisonOperator = "!="

	// OpLt represents the less than operator
	OpLt ComparisonOperator = "<"

	// OpLte represents the less than or equal to operator
	OpLte ComparisonOperator = "<="

	// OpGt represents the greater than operator
	OpGt ComparisonOperator = ">"

	// OpGte represents the greater than or equal to operator
	OpGte ComparisonOperator = ">="
)

// String returns the string representation of the ComparisonOperator
func (co ComparisonOperator) String() string {
	switch co {
	case OpEqual:
		return "="
	case OpNotEqual:
		return "<>"
	case OpNotEqualAlias:
		return "!="
	case OpLt:
		return "<"
	case OpLte:
		return "<="
	case OpGt:
		return ">"
	case OpGte:
		return ">="
	default:
		return "Unknown"
	}
}

// ComparisonPredicatesOperator represents the comparison predicates operators in filter expressions
type ComparisonPredicatesOperator string

const (
	// OpLike represents the LIKE operator
	OpLike ComparisonPredicatesOperator = "LIKE"

	// OpNotLike represents the NOT LIKE operator
	OpNotLike ComparisonPredicatesOperator = "NOT LIKE"

	// OpIsNull represents the IS NULL operator
	OpIn ComparisonPredicatesOperator = "IN"

	// OpIsNotNull represents the IS NOT NULL operator
	OpIsNotNull ComparisonPredicatesOperator = "IS NOT NULL"

	// OpIsNull represents the IS NULL operator
	OpIsNull ComparisonPredicatesOperator = "IS NULL"

	// OpNotIn represents the NOT IN operator
	OpNotIn ComparisonPredicatesOperator = "NOT IN"

	// OpBetween represents the BETWEEN operator
	OpBetween ComparisonPredicatesOperator = "BETWEEN"

	// OpNotBetween represents the NOT BETWEEN operator
	OpNotBetween ComparisonPredicatesOperator = "NOT BETWEEN"

	// OpDistinct represents the DISTINCT operator
	OpDistinct ComparisonPredicatesOperator = "DISTINCT"

	// OpNotDistinct represents the NOT DISTINCT operator
	OpNotDistinct ComparisonPredicatesOperator = "NOT DISTINCT"
)

// String returns the string representation of the ComparisonPredicatesOperator
func (cpo ComparisonPredicatesOperator) String() string {
	switch cpo {
	case OpLike:
		return "LIKE"
	case OpNotLike:
		return "NOT LIKE"
	case OpIn:
		return "IN"
	case OpIsNotNull:
		return "IS NOT NULL"
	case OpIsNull:
		return "IS NULL"
	case OpNotIn:
		return "NOT IN"
	case OpBetween:
		return "BETWEEN"
	case OpNotBetween:
		return "NOT BETWEEN"
	case OpDistinct:
		return "DISTINCT"
	case OpNotDistinct:
		return "NOT DISTINCT"
	default:
		return "Unknown"
	}
}
