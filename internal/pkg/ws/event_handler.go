package ws

import (
	"log"
	"sync"

	"github.com/alimasry/gopad/internal/pkg/editor"
)

var (
	editors     = make(map[string]*editor.Editor)
	editorMutex sync.RWMutex
)

// getOrCreateEditor returns an existing editor or creates a new one
func getOrCreateEditor(documentUUID string) (*editor.Editor, error) {
	editorMutex.Lock()
	defer editorMutex.Unlock()

	if e, exists := editors[documentUUID]; exists {
		return e, nil
	}

	strategy := editor.NewOTStrategy()
	e, err := editor.NewEditor(documentUUID, strategy)
	if err != nil {
		return nil, err
	}

	editors[documentUUID] = e
	return e, nil
}

// handle incoming insert events
func handleInsert(client *Client, insertData InsertData) {
	e, err := getOrCreateEditor(client.documentUUID)
	if err != nil {
		log.Printf("Error getting editor: %v", err)
		return
	}

	err = e.ProcessEdit(editor.Edit{
		Position:  insertData.Position,
		Insert:    insertData.String,
		Version:   client.ActiveVersion,
		ReplicaId: client.ReplicaId,
	})

	if err != nil {
		log.Printf("Error processing edit: %v", err)
	}
}

// handle incoming delete events
func handleDelete(client *Client, deleteData DeleteData) {
	e, err := getOrCreateEditor(client.documentUUID)
	if err != nil {
		log.Printf("Error getting editor: %v", err)
		return
	}

	err = e.ProcessEdit(editor.Edit{
		Position:  deleteData.Position,
		Delete:    deleteData.Delete,
		Version:   client.ActiveVersion,
		ReplicaId: client.ReplicaId,
	})

	if err != nil {
		log.Printf("Error processing edit: %v", err)
	}
}
