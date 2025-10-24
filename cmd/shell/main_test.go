package main

import (
	"reflect"
	"strings"
	"testing"
)

func TestProcessCSV(t *testing.T) {
	// The function under test is assumed to be:
	// func processCSV(r io.Reader) ([]Book, error)

	testCases := []struct {
		name          string
		inputCSV      string
		expectedBooks []Book
		expectError   bool
	}{
		{
			name: "Valid data with multiple records",
			inputCSV: `Book Id,Title,Author,Author l-f,Additional Authors,ISBN,ISBN13,My Rating
1,Dune,Frank Herbert,Herbert, Frank,,0441172717,978441172719,5
2,Project Hail Mary,Andy Weir,Weir, Andy,,"0593135202","9780593135204",4
3,Neuromancer,William Gibson,Gibson, William,Sterling, Bruce,0441569584,9780441569585,5
`,
			expectedBooks: []Book{
				{Title: "Dune", Author: "Frank Herbert", AuthorLF: "Herbert, Frank", AdditionalAuthors: ""},
				{Title: "Project Hail Mary", Author: "Andy Weir", AuthorLF: "Weir, Andy", AdditionalAuthors: ""},
				{Title: "Neuromancer", Author: "William Gibson", AuthorLF: "Gibson, William", AdditionalAuthors: "Sterling, Bruce"},
			},
			expectError: false,
		},
		{
			name: "Only header row",
			inputCSV: `Book Id,Title,Author,Author l-f,Additional Authors,ISBN,ISBN13,My Rating
`,
			expectedBooks: []Book{},
			expectError:   false,
		},
		{
			name:          "Completely empty file",
			inputCSV:      "",
			expectedBooks: nil,
			expectError:   true, // Expect an error when reading the header fails (EOF)
		},
		{
			name: "Data with short records (records with < 5 columns should be skipped)",
			// The 2nd record is too short (only 3 columns, index 0, 1, 2)
			inputCSV: `Book Id,Title,Author,Author l-f,Additional Authors,ISBN,ISBN13,My Rating
1,Dune,Frank Herbert,Herbert, Frank,,0441172717,978441172719,5
2,Short Book,A
3,Valid Again,C,D,E,F,G,H
`,
			expectedBooks: []Book{
				{Title: "Dune", Author: "Frank Herbert", AuthorLF: "Herbert, Frank", AdditionalAuthors: ""},
				// Record 2 is skipped because it doesn't have a value for Additional Authors (index 4)
				{Title: "Valid Again", Author: "C", AuthorLF: "D", AdditionalAuthors: "E"},
			},
			expectError: false,
		},
		{
			name: "Row with missing optional values",
			inputCSV: `Book Id,Title,Author,Author l-f,Additional Authors,ISBN,ISBN13,My Rating
1,Book Title,Main Author,,,"ISBN","ISBN13",1
`,
			expectedBooks: []Book{
				{Title: "Book Title", Author: "Main Author", AuthorLF: "", AdditionalAuthors: ""},
			},
			expectError: false,
		},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			// Create an io.Reader from the test case string
			r := strings.NewReader(tc.inputCSV)

			// Call the function under test
			books, err := processCSV(r)

			// Check for expected error state
			if tc.expectError {
				if err == nil {
					t.Errorf("Expected an error, but got nil")
				}
				return // Exit if error was expected and received
			}

			// If error was not expected, check if one occurred
			if err != nil {
				t.Fatalf("Unexpected error: %v", err)
			}

			// Check if the resulting list of books matches the expected list
			if !reflect.DeepEqual(books, tc.expectedBooks) {
				t.Errorf("Result mismatch:\nGot:      %+v\nExpected: %+v", books, tc.expectedBooks)
			}
		})
	}
}
