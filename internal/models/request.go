package models

import "github.com/google/uuid"

type CreateDocumentRequest struct {
	Title   string `json:"title"`
	Content string `json:"content"`
}

type UpdateDocumentRequest struct {
	UUID    uuid.UUID `json:"uuid"`
	Title   *string   `json:"title"`
	Content *string   `json:"content"`
}

type EditDocumentRequest struct {
	Type     string `json:"type"`
	Position int32  `json:"position"`
	Value    rune   `json:"value"`
}
