package diff

import "regexp"

type mode int16

const (
	toChars mode = iota
	toWords
	toLines
)

var toc, tow, tol = toChars, toWords, toLines

var (
	CHAR = &toc
	WORD = &tow
	LINE = &tol
)

//A sequence whose items (strings) can be random-accessed at a given index.
type Sequence interface {
	Len() int
	// The content of this sequence from start through end (inclusive). Includes any gaps in between items, but does
	// not include gaps preceding the start or following the end.
	Range(start, end int) (result string, ok bool)
	At(i int) (result string, ok bool)
	// If a sequence allows gaps in between its items, the gaps can be accessed using Pre (for gaps preceding an item).
	Pre(i int) (result string, ok bool)
	// If a sequence allows gaps in between its items, the gap following a sequence can be accessed using Suffix().
	Suffix() string
}

type Sequencer interface {
	Sequence(s string) Sequence
}

var (
	words = regexp.MustCompile(`[^ \t\n\r]+`)
	lines = regexp.MustCompile(`[^\n\r]+`)
)

func (m *mode) Sequence(s string) Sequence {
	var matches [][]int
	switch *m {
	case toChars:
		seq := chars(s)
		return &seq
	case toWords:
		matches = words.FindAllStringIndex(s, -1)
	case toLines:
		matches = lines.FindAllStringIndex(s, -1)
	}
	return &sequence {s, matches}
}

//data structure for char sequences
type chars string

func (c *chars) Range(from, to int) (result string, ok bool) {
	if from < 0 || to < from || to >= len(string(*c)) {
		return "", false
	}
	return string(*c)[from:to], true
}

func (c *chars) At(i int) (result string, ok bool) {
	if i < 0 || i >= len(string(*c)) {
		return "", false
	}
	return string(string(*c)[i]), true
}

func (c *chars) Len() int {
	return len(string(*c))
}

func (c *chars) Pre(i int) (result string, ok bool) {
	if i < 0 || i >= len(string(*c)) {
		return "", false
	}
	return "", true
}

func (c *chars) Suffix() string {
	return ""
}

//data structure for word and line sequences
type sequence struct {
	raw string
	runs [][]int
}

func (a *sequence) Range(from, to int) (result string, ok bool) {
	if from < 0 || to < from || to >= a.Len() {
		return "", false
	}
	fStart, tEnd := a.runs[from][0], a.runs[to][1]
	return a.raw[fStart:tEnd+1], true
}

func (a *sequence) Pre(i int) (prefix string, ok bool) {
	if a.Len() <= i {
		return "", false
	}
	run := a.runs[i]
	iStart := run[0] //beginning of the ith word (1 + end of the prefix)
	pStart := 0		 //beginning of the i-1th word (beginning of the prefix)
	if i > 0 {
		pStart = a.runs[i-1][1]
	}
	return a.raw[pStart:iStart], true
}

func (a *sequence) Suffix() string {
	if len(a.runs) == 0 {
		return a.raw
	}
	run := a.runs[len(a.runs)-1]
	end := run[1]
	if end + 1 >= len(a.raw) {
		return ""
	}
	return a.raw[end+1:len(a.raw)]
}

func (a *sequence) At(i int) (result string, ok bool) {
	if i < 0 || i >= a.Len() {
		return "", false
	}
	run := a.runs[i]
	start, end := run[0], run[1]
	return a.raw[start:end+1], true
}

func (a *sequence) Len() int {
	return len(a.runs)
}
