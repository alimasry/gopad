package ws

import (
	"log"

	"github.com/alimasry/gopad/internal/pkg/ot"
)

// handle incoming insert events
func handleInsert(client *Client, insertData InsertData) {
	log.Printf("client : %v", client)
	otBufferManager := ot.GetOTBufferManager()
	otBuffer := otBufferManager.GetOTBuffer(client.documentUUID)
	otBuffer.PushTransformation(ot.OTransformation{
		Position:  insertData.Position,
		Insert:    insertData.String,
		Version:   client.lastSyncedVersion,
		ReplicaId: client.ReplicaId,
	})
}

// handle incoming delete events
func handleDelete(client *Client, deleteData DeleteData) {
	otBufferManager := ot.GetOTBufferManager()
	otBuffer := otBufferManager.GetOTBuffer(client.documentUUID)
	otBuffer.PushTransformation(ot.OTransformation{
		Position:  deleteData.Position,
		Delete:    deleteData.Delete,
		Version:   client.lastSyncedVersion,
		ReplicaId: client.ReplicaId,
	})
}
