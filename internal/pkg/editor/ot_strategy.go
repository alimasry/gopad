package editor

import "github.com/alimasry/gopad/internal/pkg/ot"

type OTStrategy struct {
	buffer       *ot.OTBuffer
	documentUUID string
}

func NewOTStrategy() *OTStrategy {
	return &OTStrategy{}
}

func (s *OTStrategy) Initialize(documentUUID string) error {
	s.documentUUID = documentUUID
	s.buffer = ot.GetOTBufferManager().GetOTBuffer(documentUUID)
	return nil
}

func (s *OTStrategy) ProcessEdit(edit Edit) error {
	s.buffer.PushTransformation(ot.OTransformation{
		Position:  edit.Position,
		Delete:    edit.Delete,
		Insert:    edit.Insert,
		Version:   edit.Version,
		ReplicaId: edit.ReplicaId,
	})
	return nil
}

func (s *OTStrategy) GetContent() string {
	return s.buffer.Content()
}

func (s *OTStrategy) Close() error {
  // TODO: save the file before closing
	return nil
}
