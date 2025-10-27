package main

import (
	"encoding/json"
	"fmt"
	"log"
	"os"
	"path/filepath"

	"github.com/olivier-twist/book_tags/internal/booktag"
)

func findModuleRoot(startDir string) (string, error) {
	currentDir := startDir
	for {
		goModPath := filepath.Join(currentDir, "go.mod")
		if _, err := os.Stat(goModPath); err == nil {
			return currentDir, nil // Found go.mod, this is the module root
		}

		parentDir := filepath.Dir(currentDir)
		if parentDir == currentDir { // Reached the filesystem root
			return "", os.ErrNotExist
		}
		currentDir = parentDir
	}
}

func main() {
	cwd, err := os.Getwd()
	if err != nil {
		log.Fatalf("Error getting current working directory: %s", err)
	}

	rootDir, err := findModuleRoot(cwd)
	if err != nil {
		log.Fatalf("Error finding module root: %s", err)
	}

	err = os.Chdir(rootDir)
	if err != nil {
		log.Fatalf("Error changing directory to module root: %s", err)
	}

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
	_, err = json.MarshalIndent(books, "", "  ")
	if err != nil {
		log.Fatalf("Error marshaling JSON: %s", err)
	}

	// 4. Print the JSON to standard output
	//fmt.Println(string(jsonData))
	tagged_books, err := booktag.ProcessBooksWithGemini(books)
	if err != nil {
		log.Fatalf("Error processing books with Gemini: %s", err)
	}

	// Convert the tagged books slice to JSON
	taggedJsonData, err := json.MarshalIndent(tagged_books, "", "  ")
	if err != nil {
		log.Fatalf("Error marshaling tagged books JSON: %s", err)
	}

	// Print the tagged books JSON to standard output
	fmt.Println(string(taggedJsonData))

	err = booktag.InsertTaggedBooks(tagged_books)
	if err != nil {
		log.Fatalf("Error inserting tagged books into database: %s", err)
	}
}
