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
	return string(sd)
}

// Operator represents the comparison operators in filter expressions
type Operator string

const (
	OperatorAnd                  Operator = "AND"
	OperatorOr                   Operator = "OR"
	OperatorNot                  Operator = "NOT"
	OperatorEqual                Operator = "="
	OperatorNotEqual             Operator = "<>"
	OperatorNotEqualAlias        Operator = "!="
	OperatorLessThan             Operator = "<"
	OperatorLessThanOrEqualTo    Operator = "<="
	OperatorGreaterThan          Operator = ">"
	OperatorGreaterThanOrEqualTo Operator = ">="
	// ----
	OperatorLike        Operator = "LIKE"
	OperatorNotLike     Operator = "NOT LIKE"
	OperatorIn          Operator = "IN"
	OperatorIsNotNull   Operator = "IS NOT NULL"
	OperatorIsNull      Operator = "IS NULL"
	OperatorNotIn       Operator = "NOT IN"
	OperatorBetween     Operator = "BETWEEN"
	OperatorNotBetween  Operator = "NOT BETWEEN"
	OperatorDistinct    Operator = "DISTINCT"
	OperatorNotDistinct Operator = "NOT DISTINCT"
)

// String returns the string representation of the Operator
func (co Operator) String() string {
	return string(co)
}
