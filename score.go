package gdiff

import "math"

type ScoreAlgo string

const (
	simple ScoreAlgo = "simple"
)

// produces a score from 0 to 100 indicating the similarity between a Diff's two sides. 100 means identical.
type Matcher interface {
	Score(diff Diff) (score float64)
	Algorithm() ScoreAlgo
}

type simpleMatcher int

func SimpleMatcher() *simpleMatcher {
	matcher := simpleMatcher(0)
	return &matcher
}

func (*simpleMatcher) Algorithm() ScoreAlgo {
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

func (*simpleMatcher) Score(diff Diff) (score float64) {
	edits := diff.Edits()
	max := math.Max(float64(diff.A().Len()), float64(diff.B().Len()))
	common := max
	for _, edit := range edits {
		common -= float64(edit.End + 1 - edit.Start)
	}
	return common * 100 / max
}
