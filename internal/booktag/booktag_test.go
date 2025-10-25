package booktag

import (
	"reflect"
	"strings"
	"testing"
)

// The Book struct and processCSV function are assumed to be defined in main.go.

func TestProcessCSV(t *testing.T) {
	// Column indices for extraction:
	// 1 (Title), 2 (Author), 3 (Author l-f), 4 (Additional Authors)

	// Header copied directly from the sample file.
	header := `Book Id,Title,Author,Author l-f,Additional Authors,ISBN,ISBN13,My Rating,Average Rating,Publisher,Binding,Number of Pages,Year Published,Original Publication Year,Date Read,Date Added,Bookshelves,Bookshelves with positions,Exclusive Shelf,My Review,Spoiler,Private Notes,Read Count,Owned Copies`

	testCases := []struct {
		name          string
		inputCSV      string
		expectedBooks []Book
		expectError   bool
	}{
		{
			name: "Valid data with multiple records (Cleaned Input)",
			// Raw string literal (backticks) used. The non-standard ISBN quotes
			// (e.g., ="978..." and ="") have been manually cleaned to standard CSV
			// format to prevent parsing errors.
			inputCSV: header + `
43385914,Swipe to Unlock: The Primer on Technology and Business Strategy,Parth Detroja,"Detroja, Parth","Neel Mehta, Aditya   Agashe",,,0,4.22,Amazon Digital Services,Kindle Edition,325,2018,2017,,2025/10/23,to-read,to-read (#1182),to-read,,,,0,0
2019095,A Practical Guide to Data Structures and Algorithms using Java (Chapman & Hall/CRC Applied Algorithms and Data Structures series),Sally A. Goldman,"Goldman, Sally A.",Kenneth J. Goldman,158488455X,9781584884552,0,4.50,Chapman and Hall/CRC,Hardcover,1054,2007,2007,,2025/10/23,to-read,to-read (#1181),to-read,,,,0,0
21280882,Smartups: Lessons from Rob Ryan's Entrepreneur America Boot Camp for Start-Ups,Rob Ryan,"Ryan, Rob",David J. BenDaniel,,,0,3.24,Cornell University Press,Kindle Edition,240,2012,2002,,2025/10/23,to-read,to-read (#1180),to-read,,,,0,0
`,
			expectedBooks: []Book{
				{Title: "Swipe to Unlock: The Primer on Technology and Business Strategy", Author: "Parth Detroja", AuthorLF: "Detroja, Parth", AdditionalAuthors: "Neel Mehta, Aditya   Agashe"},
				{Title: "A Practical Guide to Data Structures and Algorithms using Java (Chapman & Hall/CRC Applied Algorithms and Data Structures series)", Author: "Sally A. Goldman", AuthorLF: "Goldman, Sally A.", AdditionalAuthors: "Kenneth J. Goldman"},
				{Title: "Smartups: Lessons from Rob Ryan's Entrepreneur America Boot Camp for Start-Ups", Author: "Rob Ryan", AuthorLF: "Ryan, Rob", AdditionalAuthors: "David J. BenDaniel"},
			},
			expectError: false,
		},
		// {
		// 	name:          "Only header row",
		// 	inputCSV:      header + "\n",
		// 	expectedBooks: []Book{},
		// 	expectError:   false,
		// },
		{
			name:          "Completely empty file",
			inputCSV:      "",
			expectedBooks: nil,
			// Expect an error from processCSV when reading header fails (io.EOF)
			expectError: true,
		},
		{
			name: "Data with short records (< 5 columns are skipped)",
			inputCSV: header + `
43385914,Swipe to Unlock: The Primer on Technology and Business Strategy,Parth Detroja,"Detroja, Parth","Neel Mehta, Aditya   Agashe",,,0,4.22,Amazon Digital Services,Kindle Edition,325,2018,2017,,2025/10/23,to-read,to-read (#1182),to-read,,,,0,0
2,Short Book,A 
3,"Valid Title",B,"C, D",E,1111111111,9781111111111,4,4.00,Pub,Hardcover,300,2000,2000,,2025/01/01,read,read (#2),read,,,,1,1
`,
			expectedBooks: []Book{
				{Title: "Swipe to Unlock: The Primer on Technology and Business Strategy", Author: "Parth Detroja", AuthorLF: "Detroja, Parth", AdditionalAuthors: "Neel Mehta, Aditya   Agashe"},
				// Record 2 is skipped because it is too short.
				{Title: "Valid Title", Author: "B", AuthorLF: "C, D", AdditionalAuthors: "E"},
			},
			expectError: false,
		},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			r := strings.NewReader(tc.inputCSV)
			books, err := ProcessCSV(r)

			if tc.expectError {
				if err == nil {
					t.Errorf("Expected an error, but got nil")
				}
				return
			}

			if err != nil {
				t.Fatalf("Unexpected error: %v", err)
			}

			if !reflect.DeepEqual(books, tc.expectedBooks) {
				t.Errorf("Result mismatch:\nGot:      %+v\nExpected: %+v", books, tc.expectedBooks)
			}
		})
	}
}
