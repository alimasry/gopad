package ot

type OTransformation struct {
	Position  int
	Delete    int
	Insert    string
	Version   int
	ReplicaId string
}

// transforms the first operation index so that it would resolve conflicts between updates from different clients
func Transform(newOT, oldOT *OTransformation) {
	if newOT.Position-newOT.Delete > oldOT.Position-oldOT.Delete ||
		(newOT.Position-newOT.Delete == oldOT.Position-oldOT.Delete && !order(*newOT, *oldOT)) {
		newOT.Position += len(oldOT.Insert) - oldOT.Delete
	}
}

// returns true if the order newOT, oldOT is correct
func order(newOT, oldOT OTransformation) bool {
	// if different versions older version is first
	if newOT.Version != oldOT.Version {
		return newOT.Version < oldOT.Version
	}
	// lower replica id takes precedence (just for convergence)
	// if same replica id newOT is first
	return newOT.ReplicaId <= oldOT.ReplicaId
}
