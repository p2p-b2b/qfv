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
	sortParser := qfv.NewSortParser(allowedSortFields)
	fieldsParser := qfv.NewFieldsParser(allowedFieldsFields)
	filterParser := qfv.NewFilterParser(allowedFilterFields)

	// Example inputs
	sortInput := "first_name ASC,created_at DESC"
	fieldsInput := "first_name, last_name, email"
	filterInput := "first_name = 'John' AND last_name = 'Doe' OR email = 'example@example.com' AND (created_at > '2023-01-01' OR updated_at < '2023-12-31')"
	fmt.Printf("Filter input: %s\n", filterInput)
	fmt.Printf("Sort input: %s\n", sortInput)
	fmt.Printf("Fields input: %s\n", fieldsInput)

	// Validate the inputs
	sortNode, err := sortParser.Parse(sortInput)
	if err != nil {
		log.Fatalf("Validation error: %v", err)
	}

	fieldsNode, err := fieldsParser.Parse(fieldsInput)
	if err != nil {
		log.Fatalf("Validation error: %v", err)
	}

	_, err = filterParser.Parse(filterInput)
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

	// Print the parsed filter
	fmt.Println("\nFilter is valid")
}
