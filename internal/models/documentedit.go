package models

type EditType int

const (
	MoveCursor = iota
	Insert
	Delete
)

type DocumentEdit struct {
	Type EditType `json:"type"`
	Data []byte   `json:"data"`
}

type MoveCursorData struct {
	Position int `json:"position"`
}

type InsertData struct {
	Text string `json:"text"`
}

type DeleteData struct {
	Count int `json:"count"`
}
