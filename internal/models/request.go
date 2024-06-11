package models

// CreateDocumentRequest represent request body for POST /documents
type CreateDocumentRequest struct {
	Title   string `json:"title"`
	Content string `json:"content"`
}
