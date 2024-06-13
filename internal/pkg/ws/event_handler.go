package ws

import (
	"log"

	"github.com/alimasry/gopad/internal/pkg/ot"
	"github.com/alimasry/gopad/internal/services/editor"
)

var otBufferManager *ot.OTBufferManager = ot.GetOTBufferManager()

// handle incoming insert events
func handleInsert(client *Client, insertData InsertData) {
	otBuffer := otBufferManager.GetOTBuffer(client.documentUUID)
	otBuffer.PushTransformation(ot.OTransformation{
		Position:  insertData.Position,
		Insert:    insertData.String,
		Version:   client.ActiveVersion,
		ReplicaId: client.ReplicaId,
	})
}

// handle incoming delete events
func handleDelete(client *Client, deleteData DeleteData) {
	otBuffer := otBufferManager.GetOTBuffer(client.documentUUID)
	otBuffer.PushTransformation(ot.OTransformation{
		Position:  deleteData.Position,
		Delete:    deleteData.Delete,
		Version:   client.ActiveVersion,
		ReplicaId: client.ReplicaId,
	})
}

func handleUndo(client *Client) {
	err := editor.Undo(client.documentUUID)
	if err != nil {
		log.Println("Error occured: ", err)
	}
}

func handleRedo(client *Client) {
	err := editor.Redo(client.documentUUID)
	if err != nil {
		log.Println("Error occured: ", err)
	}
}
