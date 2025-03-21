package main

import (
	"fmt"
	"log"

	qfv "github.com/p2p-b2b/query-ast-validator"
)

func main() {
	// Define the allowed fields for your API
	allowedSortFields := []string{"first_name", "last_name", "email", "created_at", "updated_at"}
	fmt.Printf("Allowed Sort fields: %v\n", allowedSortFields)

	allowedFieldsFields := []string{"first_name", "last_name", "email", "created_at", "updated_at"}
	fmt.Printf("Allowed Fields fields: %v\n", allowedFieldsFields)

	allowedFilterFields := []string{"first_name", "last_name", "email", "created_at", "updated_at"}
	fmt.Printf("Allowed Filter fields: %v\n", allowedFilterFields)

	// Create a new parser with the allowed fields
	p := qfv.NewParser(allowedSortFields, allowedFieldsFields, allowedFilterFields)

	// Example inputs
	sortInput := "first_name ASC,created_at DESC"
	fieldsInput := "first_name, last_name, email"
	filterInput := "first_name = 'John' AND last_name = 'Doe' OR email = 'example@example.com' AND (created_at > '2023-01-01' OR updated_at < '2023-12-31')"
	fmt.Printf("Filter input: %s\n", filterInput)
	fmt.Printf("Sort input: %s\n", sortInput)
	fmt.Printf("Fields input: %s\n", fieldsInput)

	// Validate the inputs
	sortNode, fieldsNode, filterNode, err := p.Validate(sortInput, fieldsInput, filterInput)
	if err != nil {
		log.Fatalf("Validation error: %v", err)
	}

	// Print the parsed sort fields
	fmt.Println("Sort fields:")
	for _, field := range sortNode.Fields {
		fmt.Printf("  %s %s\n", field.Field, field.Direction)
	}

	// Print the parsed fields
	fmt.Println("\nRequested fields:")
	for _, field := range fieldsNode.Fields {
		fmt.Printf("  %s\n", field)
	}

	// Print the filter (this is more complex, we'll just confirm it's valid)
	fmt.Println("\nFilter is valid")

	// You could implement a function to print the filter AST for debugging
	printFilterTree(filterNode.Root, 0)
}

// Helper function to print the filter AST
func printFilterTree(node qfv.FilterExprNode, indent int) {
	indentStr := ""
	for range indent {
		indentStr += "  "
	}

	switch n := node.(type) {
	case qfv.LogicalOpNode:
		fmt.Printf("%s%s\n", indentStr, n.Operator)
		fmt.Printf("%sLeft:\n", indentStr)
		printFilterTree(n.Left, indent+1)
		if n.Right != nil {
			fmt.Printf("%sRight:\n", indentStr)
			printFilterTree(n.Right, indent+1)
		}
	case qfv.ComparisonNode:
		fmt.Printf("%s%s %s ", indentStr, n.Field, n.Operator)
		switch v := n.Value.(type) {
		case qfv.StringNode:
			fmt.Printf("'%s'\n", v.Value)
		case qfv.IdentifierNode:
			fmt.Printf("%s\n", v.Value)
		}
	case qfv.ComparisonPredicatesNode:
		fmt.Printf("%s%s %s ", indentStr, n.Field, n.Operator)
		switch v := n.Value.(type) {
		case qfv.StringNode:
			fmt.Printf("'%s'\n", v.Value)
		case qfv.IdentifierNode:
			fmt.Printf("%s\n", v.Value)
		}
	case qfv.IdentifierNode:
		fmt.Printf("%sIdentifier: %s\n", indentStr, n.Value)
	case qfv.StringNode:
		fmt.Printf("%sString: '%s'\n", indentStr, n.Value)
	default:
		fmt.Printf("%sUnknown node type: %T\n", indentStr, n)
	}
}
