package ot

import (
	"log"
	"sync"

	"github.com/alimasry/gopad/internal/pkg/gapbuffer"
	"github.com/alimasry/gopad/internal/services/document"
)

type OTBufferMap map[string]*OTBuffer

type OTBuffer struct {
	UUID      string
	Version   int
	Pending   []OTransformation
	gapBuffer *gapbuffer.GapBuffer
	sync.RWMutex
}

// creates a new OTBuffer
func NewOTBuffer(documentUUID string) *OTBuffer {
	doc, err := document.GetDocumentFromCache(documentUUID)

	if err != nil {
		log.Println("Error occured", err.Error())
	}

	return &OTBuffer{
		UUID:      doc.UUID,
		Version:   doc.Version,
		gapBuffer: gapbuffer.NewGapBufferWithContent(doc.Content),
		Pending:   make([]OTransformation, 0),
	}
}

// push a transformation to the list of pending transformation so that it could be processed
func (otb *OTBuffer) PushTransformation(t OTransformation) {
	otb.Lock()
	defer otb.Unlock()

	for i := range otb.Pending {
		Transform(&t, &otb.Pending[i])
	}

	otb.Pending = append(otb.Pending, t)
}

// get current content of the document
func (otb *OTBuffer) Content() string {
	return otb.gapBuffer.String()
}

// save the document to the database
func (otb *OTBuffer) save() error {
	doc, err := document.GetDocument(otb.UUID)
	if err != nil {
		return err
	}

	doc.Version = otb.Version + 1
	doc.Content = otb.Content()

	err = document.SaveDocument(doc)

	if err != nil {
		return err
	}

	otb.Version++

	return nil
}

// process transformations related to that buffer
func (otb *OTBuffer) Process() {
	if len(otb.Pending) == 0 {
		return
	}

	otb.Lock()
	defer otb.Unlock()

	for _, t := range otb.Pending {
		if t.Delete > 0 {
			otb.gapBuffer.DeleteAt(t.Position, t.Delete)
		}
		if t.Insert != "" {
			otb.gapBuffer.InsertAt(t.Position, t.Insert)
		}
	}

	otb.Pending = make([]OTransformation, 0)

	if err := otb.save(); err != nil {
		log.Println("Error occured: ", err)
	}
}
