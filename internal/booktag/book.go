package booktag

// Book struct holds the extracted book data.
type Book struct {
	Title             string `json:"title"`
	Author            string `json:"author"`
	AuthorLF          string `json:"authorLF"`          // Corresponds to "Author l-f"
	AdditionalAuthors string `json:"additionalAuthors"` // Corresponds to "Additional Authors"
}

// TaggedBook reflects the structure of the desired output, combining the Book data with the new Tag field.
type TaggedBook struct {
	Book
	Tag string `json:"tag"`
}
