package diff

import "io"

func (d *Diff) Unified(w io.Writer) bool {
	lastEnd := 0
	a, b := d.a, d.b
	for _, v := range d.edits {
		switch v.Type {
		case DELETE:
			if lastEnd < v.Start {
				if value, ok := a.Range(0, v.Start-1); ok {
					w.Write([]byte(value))
				} else {
					return false
				}
			}
			lastEnd = v.End + 1
			w.Write([]byte("-"))
			if value, ok := a.Range(v.Start, v.End); ok {
				w.Write([]byte(value))
			} else {
				return false
			}
		case INSERT:
			w.Write([]byte("+"))
			if value, ok := b.Range(v.Start, v.End); ok {
				w.Write([]byte(value))
			} else {
				return false
			}
		}
	}
	if lastEnd < a.Len() {
		if value, ok := a.Range(lastEnd, a.Len()-1); ok {
			w.Write([]byte(value))
		} else {
			return false
		}
	}
	return true
}
