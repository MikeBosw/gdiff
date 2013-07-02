package gdiff

import "math"

type ScoreAlgo string

const (
	simple ScoreAlgo = "simple"
)

type Comparator interface {
	Score(diff Diff) (score float64)
	Algorithm() ScoreAlgo
}

type simpleComparator struct {}

var simpleton *simpleComparator = &simpleComparator{}

func SimpleComparator() Comparator {
	return simpleton
}

func (*simpleComparator) Algorithm() ScoreAlgo {
	return simple
}

type strings []string

func (s strings) Len() int {
	return len([]string(s))
}

func (s strings) Less(i, j int) bool {
	slice := []string(s)
	return slice[i] < slice[j]
}

func (s strings) Swap(i, j int) {
	slice := []string(s)
	slice[i], slice[j] = slice[j], slice[i]
}

func (*simpleComparator) Score(diff Diff) (score float64) {
	edits := diff.Edits()
	max := math.Max(float64(diff.A().Len()), float64(diff.B().Len()))
	common := max
	for _, edit := range edits {
		common -= float64(edit.End + 1 - edit.Start)
	}
	return common * 100 / max
}

