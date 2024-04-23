package models

type CreateDocumentRequest struct {
	Title   string `json:"title"`
	Content string `json:"content"`
}
