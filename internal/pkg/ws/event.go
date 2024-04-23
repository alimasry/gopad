package ws

import (
	"encoding/json"
	"log"
)

var (
	SyncEvent   = "sync_event"
	InsertEvent = "insert_event"
	DeleteEvent = "delete_event"
)

type Event struct {
	Command string          `json="command"`
	UUID    string          `json="uuid"`
	Data    json.RawMessage `json="data"`
}

type InsertData struct {
	Position int    `json="position"`
	String   string `json="string"`
}

type DeleteData struct {
	Position int `json="position"`
	Delete   int `json="delete"`
}

// route incoming events to their handler functions
func routeEvent(event Event) {
	switch event.Command {
	case InsertEvent:
		var insertData InsertData
		if err := json.Unmarshal(event.Data, &insertData); err != nil {
			log.Printf("error: %v", err)
		}
		handleInsert(event.UUID, insertData)
	case DeleteEvent:
		var deleteData DeleteData
		if err := json.Unmarshal(event.Data, &deleteData); err != nil {
			log.Printf("error: %v", err)
		}
		handleDelete(event.UUID, deleteData)
	}
}