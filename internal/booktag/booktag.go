// Package booktag provides functionality to extract specific book information from CSV data.
package booktag

// Book struct holds the extracted book data.
import (
	"encoding/csv"
	"fmt"
	"io"
)

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
