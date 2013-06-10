package gdiff

import "fmt"

type DiffAlgo string

const (
	Myers DiffAlgo = "Myers"
)

type Diff interface {
	Edits() []*Edit
	A() Sequence
	B() Sequence
}

type diff struct {
	edits []*Edit
	a, b  Sequence
	split Sequencer
}

type Differ interface {
	Diff(as, bs string, split Sequencer) (diff Diff)
	Algorithm() DiffAlgo
}

func (diff *diff) Edits() []*Edit {
	return diff.edits
}

func (diff *diff) A() Sequence {
	return diff.a
}

func (diff *diff) B() Sequence {
	return diff.b
}

type Edit struct {
	Start, End int
	Type       editType
}

type editType rune

const (
	Insert editType = 'i'
	Delete editType = 'd'
)

//// see sequence.go for Sequence and Sequencer

func DifferUsing(algorithm DiffAlgo) Differ {
	switch algorithm {
	case Myers:
		return MyersDiffer()
	default:
		panic(fmt.Sprintf("unrecognized algorithm: %x", string(algorithm)))
	}
	return nil
}
