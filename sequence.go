package diff

import "regexp"

var lines, words = regexp.MustCompile(`[^\n\r]+`), regexp.MustCompile(`[^ \t\n\r]+`)

type SequenceType int16

//TODO: something more sophisticated than chars vs words vs lines.
//e.g. ignore whitespace between words but still treat lines as the units of difference.

const (
	//treat chars as units of difference
	CHAR_SPLIT SequenceType = iota

	//treat words as units of difference
	WORD_SPLIT

	//treat lines as units of difference
	LINE_SPLIT
)

//Some string, sequenced into lines, words, or characters of text. Elements can be random-accessed at a given index.
//The whitespace (if any) that precedes the elements can be random-accessed at a given index.
type Sequence interface {
	Len() int
	// The content of this sequence from start through end (inclusive). Includes any gaps in between elements, but does
	// not include whitespace preceding the start element or following the end element.
	Range(start, end int) (result string, ok bool)
	At(i int) (result string, ok bool)
	// The space between each word or line can be accessed using Pre (for space preceding the word or line at index i).
	Pre(i int) (result string, ok bool)
	// The space following the last word or line can be accessed using Suffix()
	Suffix() string
}

func seq(s string, split SequenceType) Sequence {
	var rex *regexp.Regexp
	switch split {
	case CHAR_SPLIT:
		value := chars(s)
		return &value
	case WORD_SPLIT:
		rex = words
	case LINE_SPLIT:
		rex = lines
	}
	var matches [][]int
	matches = rex.FindAllStringIndex(s, -1)
	return &sequence {s, matches}
}

//Sequence data type for char sequences
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

//Sequence data type for word and line sequences
type sequence struct {
	raw string
	runs [][]int
}

func (a *sequence) Range(from, to int) (result string, ok bool) {
	if from < 0 || to < from || to >= a.Len() {
		return "", false
	}
	fStart, tEnd := a.runs[from][0], a.runs[to][1]
	return a.raw[fStart:tEnd], true
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
	return a.raw[start:end], true
}

func (a *sequence) Len() int {
	return len(a.runs)
}
