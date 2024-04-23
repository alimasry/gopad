package ws

import "github.com/alimasry/gopad/internal/pkg/ot"

// handle incoming insert events
func handleInsert(uuid string, insertData InsertData) {
	otBufferManager := ot.GetOTBufferManager()
	otBuffer := otBufferManager.GetOTBuffer(uuid)
	otBuffer.PushTransformation(ot.OTransformation{
		Position: insertData.Position,
		Insert:   insertData.String,
	})
}

// handle incoming delete events
func handleDelete(uuid string, deleteData DeleteData) {
	otBufferManager := ot.GetOTBufferManager()
	otBuffer := otBufferManager.GetOTBuffer(uuid)
	otBuffer.PushTransformation(ot.OTransformation{
		Position: deleteData.Position,
		Delete:   deleteData.Delete,
	})
}
