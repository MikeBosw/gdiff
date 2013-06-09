package gdiff

import "fmt"

type DiffAlgo string

const (
	MYERS DiffAlgo = "Myers"
)

type Diff struct {
	edits []*Edit
	a, b  Sequence
	split Sequencer
}

type Differ interface {
	Diff(as, bs string, split Sequencer) (diff *Diff)
    Algorithm() DiffAlgo
}

func (diff *Diff) Edits() []*Edit {
	return diff.edits
}

type Edit struct {
	Start, End int
	Type       editType
}

type editType rune

const (
	INSERT editType = 'i'
	DELETE editType = 'd'
)

//// see sequence.go for Sequence and SequenceType

func DifferUsing(algorithm DiffAlgo) Differ {
	switch algorithm {
	case MYERS:
		return MyersDiffer()
	default:
		panic(fmt.Sprintf("unrecognized algorithm: %x", string(algorithm)))
	}
	return nil
}
