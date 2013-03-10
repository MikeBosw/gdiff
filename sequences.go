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

//A sequence whose items can be random-accessed at a given index.
type Sequence interface {
	Len() int
	At(i int) interface{}
	// If a sequence allows gaps in between its items, the gaps can be accessed using Pre (for gaps preceding an item)
	// and Post (for gaps following an item).
	Pre(i int) interface{}
	// If a sequence allows gaps in between its items, the gaps can be accessed using Pre (for gaps preceding an item)
	// and Post (for gaps following an item).
	Post(i int) interface{}
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
	return fromMatches(s, matches)
}

func fromMatches(s string, si [][]int) Sequence {
	blanks, bullets := make([]string, len(si)+1), make([]string, len(si))
	var prevEnd int
	for i, v := range si {
		start, end := v[0], v[1]
		if i == len(si) - 1 {
			if end < len(s) {
				blanks[len(si)] = s[end:len(s)]
			}
		}
		blanks[i] = s[prevEnd:start]
		bullets[i] = s[start:end]
		prevEnd = end
	}
	seq := &ammo {blanks, bullets}
	return seq
}

//data structure for char sequences
type chars string

func (c *chars) At(i int) interface{} {
	return string(*c)[i]
}

func (c *chars) Len() int {
	return len(string(*c))
}

func (c *chars) Pre(i int) interface{} {
	return nil
}

func (c *chars) Post(i int) interface{} {
	return nil
}

//data structure for word and line sequences
//todo: verify memory consumption of string pieces (vs. just indices)
type ammo struct {
	blanks []string
	bullets []string
}

func (a *ammo) Pre(i int) interface{} {
	return a.blanks[i]
}

func (a *ammo) Post(i int) interface{} {
	return a.blanks[i+1]
}

func (a *ammo) At(i int) interface{} {
	return a.bullets[i]
}

func (a *ammo) Len() int {
	return len(a.bullets)
}
