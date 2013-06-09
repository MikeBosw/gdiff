package gdiff

import "regexp"

var lines, words = regexp.MustCompile(`[^\n\r]+`), regexp.MustCompile(`[^ \t\n\r]+`)

type Sequencer interface {
	Split(s string) Sequence
}

type sequenceType int8

//TODO: something more sophisticated than chars vs words vs lines.
//e.g. ignore whitespace between words but still treat lines as the units of difference.

const (
	//treat lines as units of difference
	LineSplit sequenceType = iota

	//treat words as units of difference
	WordSplit

	//treat chars as units of difference
	CharSplit
)

//Some string, sequenced into lines, words, or characters of text. Elements can be random-accessed at a given index.
//The whitespace (if any) that precedes the elements can be random-accessed at a given index.
type Sequence interface {
	Len() int
	// The content of this sequence from start through end (inclusive). Includes any gaps in between elements, but does
	// not include whitespace preceding the start element or following the end element.
	Range(start, end int) (result string, ok bool)
	// The content of this sequence from start through end (inclusive). Includes any gaps in between elements, as well
	// as any whitespace following the end element. Does not include whitespace preceding the start element.
	RangeWithTail(start, end int) (result string, ok bool)
	At(i int) (result string, ok bool)
	// The space following the word or line at index i
	Tail(i int) (result string, ok bool)
	// The space preceding the first word or line
	Head() string
}

func (seq sequenceType) Split(s string) Sequence {
	var rex *regexp.Regexp
	switch seq {
	case CharSplit:
		value := chars(s)
		return &value
	case WordSplit:
		rex = words
	case LineSplit:
		rex = lines
	}
	var matches [][]int
	matches = rex.FindAllStringIndex(s, -1)
	return &sequence{s, matches}
}

//Sequence data type for char sequences
type chars string

func (c *chars) Range(from, to int) (result string, ok bool) {
	if from < 0 || to < from || to >= len(string(*c)) {
		return "", false
	}
	return string(*c)[from:to], true
}

func (c *chars) RangeWithTail(from, to int) (result string, ok bool) {
	return c.Range(from, to)
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

func (c *chars) Tail(i int) (result string, ok bool) {
	if i < 0 || i >= len(string(*c)) {
		return "", false
	}
	return "", true
}

func (c *chars) Head() string {
	return ""
}

//Sequence data type for word and line sequences
type sequence struct {
	raw  string
	runs [][]int
}

func (a *sequence) Range(from, to int) (result string, ok bool) {
	if from < 0 || to < from || to >= a.Len() {
		return "", false
	}
	fStart, tEnd := a.runs[from][0], a.runs[to][1]
	return a.raw[fStart:tEnd], true
}

func (a *sequence) RangeWithTail(from, to int) (result string, ok bool) {
	if from < 0 || to < from || to >= a.Len() {
		return "", false
	}
	fStart, tEnd := a.runs[from][0], a.runs[to][1]
	tail, _ := a.Tail(to)
	return a.raw[fStart:tEnd] + tail, true
}

func (a *sequence) Tail(i int) (prefix string, ok bool) {
	if a.Len() <= i {
		return "", false
	}
	run := a.runs[i]
	iStart := run[1] //end of the ith word (index of its last char, + 1)
	pStart := iStart //beginning of the i+1th word (beginning of the suffix)
	if i+1 < a.Len() {
		pStart = a.runs[i+1][0]
	}
	return a.raw[iStart:pStart], true
}

func (a *sequence) Head() string {
	if len(a.runs) == 0 {
		return a.raw
	}
	start := 0
	end := a.runs[0][0]
	return a.raw[start:end]
}

func (a *sequence) At(i int) (result string, ok bool) {
	if i < 0 || i >= a.Len() {
		return "", false
	}
	run := a.runs[i]
	start, end := run[0], run[1]
	return a.raw[start:end], true
}

func (a *sequence) Len() int {
	return len(a.runs)
}
