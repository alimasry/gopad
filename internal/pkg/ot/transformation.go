package ot

type OTransformation struct {
	Position  int
	Delete    int
	Insert    string
	Version   int
	ReplicaId string
}

func Transform(o1, o2 *OTransformation) {
	if o1.Position-o1.Delete > o2.Position-o2.Delete ||
		(o1.Position-o1.Delete == o2.Position-o2.Delete && !order(*o1, *o2)) {
		o1.Position += len(o2.Insert) - o2.Delete
	}
}

/*
Returns true if the order o1, o2 is correct
*/
func order(o1, o2 OTransformation) bool {
	if o1.Version != o2.Version {
		return o1.Version < o2.Version
	}
	return o1.ReplicaId < o2.ReplicaId
}
