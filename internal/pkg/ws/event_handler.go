package ws

import (
	"github.com/alimasry/gopad/internal/pkg/ot"
)

var otBufferManager *ot.OTBufferManager = ot.GetOTBufferManager()

// handle incoming insert events
func handleInsert(client *Client, insertData InsertData) {
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
	otBuffer := otBufferManager.GetOTBuffer(client.documentUUID)
	otBuffer.PushTransformation(ot.OTransformation{
		Position:  deleteData.Position,
		Delete:    deleteData.Delete,
		Version:   client.lastSyncedVersion,
		ReplicaId: client.ReplicaId,
	})
}
