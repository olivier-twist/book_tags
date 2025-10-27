// Package booktag provides functionality to extract specific book information from CSV data.
package booktag

// Book struct holds the extracted book data.
import (
	"context"
	"database/sql"
	"encoding/csv"
	"fmt"
	"io"
	"log"
	"os"
	"strings"

	"github.com/google/generative-ai-go/genai"
	_ "github.com/mattn/go-sqlite3"
	"google.golang.org/api/option"
)

const dbFile = "db/taggedbooks.db"

// processCSV reads book data from an io.Reader (like a file or a string),
// extracts columns 2, 3, 4, and 5 (index 1, 2, 3, 4), and returns a slice of Book structs.
func ProcessCSV(r io.Reader) ([]Book, error) {
	reader := csv.NewReader(r)
	reader.FieldsPerRecord = -1
	reader.LazyQuotes = true

	// Read and discard the header row.
	_, err := reader.Read()
	if err == io.EOF {
		// Handle case where file is completely empty (no header)
		return nil, fmt.Errorf("CSV file is empty: %w", err)
	}
	if err != nil {
		return nil, fmt.Errorf("error reading CSV header: %w", err)
	}

	var books []Book

	for {
		// Read one record (row)
		record, err := reader.Read()
		if err == io.EOF {
			break // Reached end of file
		}
		if err != nil {
			return nil, fmt.Errorf("error reading CSV record: %w", err)
		}

		// The required columns are 2, 3, 4, 5, which are indices 1, 2, 3, 4.
		// A valid record must have at least 5 columns (index 0 to 4).
		if len(record) > 4 {
			// Extract and populate the Book struct
			book := Book{
				Title:             record[1],
				Author:            record[2],
				AuthorLF:          record[3],
				AdditionalAuthors: record[4],
			}

			books = append(books, book)
		}
		// Records with fewer than 5 columns are simply skipped (as per test cases).
	}

	return books, nil
}

// processBooksWithGemini takes a slice of books and calls the Gemini API for each
// sequentially (synchronous) to determine relevant tags, returning a new slice of TaggedBook elements.
func ProcessBooksWithGemini(inputBooks []Book) ([]TaggedBook, error) {
	if os.Getenv("GEMINI_API_KEY") == "" {
		return nil, fmt.Errorf("GEMINI_API_KEY environment variable is not set. Cannot call the API")
	}

	ctx := context.Background()

	// Initialize the Gemini client
	client, err := genai.NewClient(ctx, option.WithAPIKey(os.Getenv("GEMINI_API_KEY")))
	if err != nil {
		return nil, fmt.Errorf("error creating Gemini client: %w", err)
	}
	defer client.Close()

	var outputList []TaggedBook // List to accumulate all results

	// Process each book sequentially (synchronously)
	for _, book := range inputBooks {

		// Construct the prompt
		prompt := fmt.Sprintf(`
			Analyze the following book details and determine a concise list of primary, single-word genre or subject tags.
			Return ONLY the comma-separated list of tags (e.g., "Programming,InterviewPrep,Dystopian"). DO NOT include any other text, quotes, or formatting.
			
			Title: %s
			Author: %s
			Additional Authors: %s
		`, book.Title, book.Author, book.AdditionalAuthors)

		// Execute the synchronous API call
		resp, err := client.GenerativeModel("gemini-2.5-flash").
			GenerateContent(ctx, genai.Text(prompt))

		if err != nil {
			// Log the error and continue to the next book
			log.Printf("Error generating content for book '%s': %v", book.Title, err)
			continue
		}

		text := ""
		if len(resp.Candidates) > 0 && len(resp.Candidates[0].Content.Parts) > 0 {
			text = fmt.Sprintf("%s", resp.Candidates[0].Content.Parts[0])
		}

		// Process the model's response string
		tagStrings := strings.Split(text, ",")

		for _, t := range tagStrings {
			tag := strings.TrimSpace(t)
			if tag != "" {
				// Append results directly to the output list
				outputList = append(outputList, TaggedBook{
					Book: book,
					Tag:  tag,
				})
			}
		}
	}

	return outputList, nil
}

// insertTaggedBooks connects to the SQLite database and inserts all tagged books.
func InsertTaggedBooks(taggedBooks []TaggedBook) error {
	// --- New Step: Ensure the directory exists ---
	// Check if the directory "../db" exists, and create it if necessary (0755 permissions)
	if err := os.MkdirAll("../db", 0755); err != nil {
		return fmt.Errorf("failed to create directory '../db': %w", err)
	}

	// 1. Connect to the database (creates file if it doesn't exist inside the "db" folder)
	db, err := sql.Open("sqlite3", dbFile)
	if err != nil {
		return fmt.Errorf("failed to open database file '%s': %w", dbFile, err)
	}
	defer db.Close()

	// 3. Prepare the INSERT statement
	// SQLite uses '?' as the placeholder instead of '$1, $2, ...'
	stmt, err := db.PrepareContext(
		context.Background(),
		`INSERT INTO BOOK (title, author, authorLF, additionalAuthors, tag) 
		 VALUES (?, ?, ?, ?, lower(?))`,
	)
	if err != nil {
		return fmt.Errorf("failed to prepare statement: %w", err)
	}
	defer stmt.Close()

	// 4. Begin a transaction for atomic and efficient bulk insert
	// Note: Transactions are optional but highly recommended for bulk inserts.
	tx, err := db.BeginTx(context.Background(), nil)
	if err != nil {
		return fmt.Errorf("failed to start transaction: %w", err)
	}

	var insertedCount int

	// 5. Execute the statement for each TaggedBook within the transaction
	for _, tb := range taggedBooks {
		// Use ExecContext to execute the prepared statement within the transaction
		_, err := tx.Stmt(stmt).ExecContext(
			context.Background(),
			tb.Title,
			tb.Author,
			tb.AuthorLF,
			tb.AdditionalAuthors,
			tb.Tag,
		)
		if err != nil {
			// Roll back the entire transaction upon failure
			tx.Rollback()
			return fmt.Errorf("failed to execute insert for book '%s': %w", tb.Title, err)
		}
		insertedCount++
	}

	// 6. Commit the transaction if all inserts were successful
	err = tx.Commit()
	if err != nil {
		return fmt.Errorf("failed to commit transaction: %w", err)
	}

	log.Printf("Successfully inserted %d records into the BOOK table in %s.", insertedCount, dbFile)
	return nil
}
