package ot

type OTransformation struct {
	Position  int
	Delete    int
	Insert    string
	Version   int
	ReplicaId string
}

// Transforms the first operation index so that it would resolve conflicts between updates from different clients

func Transform(ot1, ot2 *OTransformation) {
	if ot1.Position-ot1.Delete > ot2.Position-ot2.Delete ||
		(ot1.Position-ot1.Delete == ot2.Position-ot2.Delete && !order(*ot1, *ot2)) {
		ot1.Position += len(ot2.Insert) - ot2.Delete
	}
}

// Returns true if the order o1, o2 is correct
func order(ot1, ot2 OTransformation) bool {
	if ot1.Version != ot2.Version {
		return ot1.Version < ot2.Version
	}
	return ot1.ReplicaId < ot2.ReplicaId
}
