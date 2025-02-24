package editor

type CollaborationStrategy interface {
	Initialize(documentUUID string) error
	ProcessEdit(edit Edit) error
	GetContent() string
	Close() error
}
