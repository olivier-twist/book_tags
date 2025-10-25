package main

import (
	"encoding/json"
	"fmt"
	"log"
	"os"

	"github.com/olivier-twist/book_tags/internal/booktag"
)

func main() {
	// 1. Open the CSV file
	file, err := os.Open("data/goodreads_library_export.csv")
	if err != nil {
		log.Fatalf("Error opening file: %s", err)
	}
	defer file.Close()

	// 2. Process the CSV data using the testable function
	books, err := booktag.ProcessCSV(file)
	if err != nil {
		log.Fatalf("Error processing CSV: %s", err)
	}

	// 3. Convert the slice of structs to JSON
	// Using MarshalIndent for clean, readable output.
	jsonData, err := json.MarshalIndent(books, "", "  ")
	if err != nil {
		log.Fatalf("Error marshaling JSON: %s", err)
	}

	// 4. Print the JSON to standard output
	fmt.Println(string(jsonData))
}
