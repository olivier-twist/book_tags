package booktag

// Book struct holds the extracted book data.
type Book struct {
	Title             string `json:"title"`
	Author            string `json:"author"`
	AuthorLF          string `json:"authorLF"`          // Corresponds to "Author l-f"
	AdditionalAuthors string `json:"additionalAuthors"` // Corresponds to "Additional Authors"
}
