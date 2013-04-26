package diff

import "io"

func Unified(a, b Sequence, edits []*edit, w io.Writer) {
	for _, v := range edits {
		switch v.Type {
		case DELETE:
			w.Write([]byte("-"))
			w.Write([]byte(a.Range(v.Start, v.End)))
		case INSERT:
			w.Write([]byte("+"))
			w.Write([]byte(b.Range(v.Start, v.End)))
		}
	}
}
