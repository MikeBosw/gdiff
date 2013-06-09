package gdiff

import "io"

type DiffFormatter interface {
	Print(d *Diff, w io.Writer) bool
}

type unifiedDiff string

var ud unifiedDiff = "unified"

func Unified() DiffFormatter {
	return &ud
}

func (ud *unifiedDiff) Print(d *Diff, w io.Writer) bool {
	lastEnd := 0
	a, b := d.a, d.b
	for _, v := range d.edits {
		switch v.Type {
		case Delete:
			if lastEnd < v.Start {
				if value, ok := a.RangeWithTail(0, v.Start-1); ok {
					w.Write([]byte(value))
				} else {
					return false
				}
			}
			lastEnd = v.End + 1
			w.Write([]byte("-"))
			if value, ok := a.RangeWithTail(v.Start, v.End); ok {
				w.Write([]byte(value))
			} else {
				return false
			}
		case Insert:
			w.Write([]byte("+"))
			if value, ok := b.RangeWithTail(v.Start, v.End); ok {
				w.Write([]byte(value))
			} else {
				return false
			}
		}
	}
	if lastEnd < a.Len() {
		if value, ok := a.RangeWithTail(lastEnd, a.Len()-1); ok {
			w.Write([]byte(value))
		} else {
			return false
		}
	}
	return true
}
