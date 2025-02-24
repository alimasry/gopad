package ws

import (
	"encoding/json"
	"log"
)

var (
	SyncEvent   = "sync_event"
	InsertEvent = "insert_event"
	DeleteEvent = "delete_event"
	UndoEvent   = "undo_event"
	RedoEvent   = "redo_event"
)

type Event struct {
	client  *Client         `json="client"`
	Command string          `json="command"`
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
			log.Println("Error occured", err.Error())
		}
		handleInsert(event.client, insertData)
	case DeleteEvent:
		var deleteData DeleteData
		if err := json.Unmarshal(event.Data, &deleteData); err != nil {
			log.Println("Error occured", err.Error())
		}
		handleDelete(event.client, deleteData)
	// case UndoEvent:
	// 	handleUndo(event.client)
	// case RedoEvent:
	// 	handleRedo(event.client)
	}
}
