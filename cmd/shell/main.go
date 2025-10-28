package main

import (
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

	tagged_books, err := booktag.ProcessBooksWithLocalLLM(books, "llama3")
	if err != nil {
		log.Fatalf("Error processing books with local LLM: %s", err)
	}

	err = booktag.InsertTaggedBooks(tagged_books)
	if err != nil {
		log.Fatalf("Error inserting tagged books into database: %s", err)
	}

	log.Printf("Successfully processed and inserted %d tagged books into the database.", len(tagged_books))
}
