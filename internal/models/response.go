package models

type GetDocumentResponse struct {
	UUID    string `json:"uuid"`
	Title   string `json:"title"`
	Content string `json:"content"`
	Size    int    `json:"size"`
	Version int    `json:"version_nr"`
}
