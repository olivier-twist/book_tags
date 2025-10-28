package booktag

type OllamaRequest struct {
	Model  string `json:"model"`
	Prompt string `json:"prompt"`
	Stream bool   `json:"stream"` // Set to false for a single, complete response
}

// OllamaResponse defines the relevant fields from the Ollama API response.
type OllamaResponse struct {
	Model     string `json:"model"`
	CreatedAt string `json:"created_at"`
	Response  string `json:"response"` // This contains the generated text/tags
	Done      bool   `json:"done"`
}
