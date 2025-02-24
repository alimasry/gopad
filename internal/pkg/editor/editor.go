package editor

type Edit struct {
	Position  int
	Delete    int
	Insert    string
	Version   int
	ReplicaId string
}

type Editor struct {
	documentUUID string
	strategy     CollaborationStrategy
}

func NewEditor(documentUUID string, strategy CollaborationStrategy) (*Editor, error) {
	if err := strategy.Initialize(documentUUID); err != nil {
		return nil, err
	}
	return &Editor{
		documentUUID: documentUUID,
		strategy:     strategy,
	}, nil
}

func (e *Editor) SetStrategy(strategy CollaborationStrategy) error {
	if err := strategy.Close(); err != nil {
		return err
	}

	if err := strategy.Initialize(e.documentUUID); err != nil {
		return err
	}

	e.strategy = strategy
	return nil
}

func (e *Editor) ProcessEdit(edit Edit) error {
	return e.strategy.ProcessEdit(edit)
}

func (e *Editor) GetContent() string {
	return e.strategy.GetContent()
}
