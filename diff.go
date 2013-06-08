package gdiff

type Diff struct {
	edits []*edit
	a, b  Sequence
	split SequenceType
}

type Differ interface {
	Diff(as, bs string, split SequenceType) (diff *Diff)
}

func (diff *Diff) Edits() []*edit {
	return diff.edits
}

type edit struct {
	Start, End int
	Type       editType
}

type editType rune

const (
	INSERT editType = 'i'
	DELETE editType = 'd'
)

//// see sequence.go for Sequence and SequenceType
