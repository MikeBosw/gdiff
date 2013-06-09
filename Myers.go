package gdiff

import "fmt"

type vertex struct {
	/* 1-based index into a string (subtract 1 when indexing into the string) */
	ai, bi int
	adj    *vertex
}

type myersDiffer struct{}

var singleton *myersDiffer = &myersDiffer{}

func MyersDiffer() Differ {
	return singleton
}

func (md *myersDiffer) Algorithm() DiffAlgo {
	return MYERS
}

func (md *myersDiffer) Diff(as, bs string, seq Sequencer) (diff *Diff) {
	diff = new(Diff)
	diff.a, diff.b = seq.Split(as), seq.Split(bs)

	if as == bs {
		return
	}

	m, n := diff.a.Len(), diff.b.Len()
	kLines := make([]int, (m+n)*2+1)
	breadcrumbs := make([]*vertex, (m+n)*2+1)

	for d := 0; d <= m+n; d++ {
		for k := -d; k <= d; k += 2 {
			ki := m + n + k //this k-line's index in the k-line array

			/* first, establish what our new X, Y is, based on how we had to have gotten here.
			 * we can have gotten here in one of two ways: either down, from the k+1 line, or right, from the k-1 line.
			 * we pick whichever gets us further towards {m, n}. */

			var origin *vertex
			var x, y int
			{
				// if we're at the NE (top, rightmost) k-line, the only way to have gotten here is right, from k-1.
				// if we're at the SW k-line, the only way to have gotten here is down, from k+1.
				isSW, isNE := k == -d, k == d
				if isSW || (!isNE && kLines[ki-1] < kLines[ki+1]) {
					x = kLines[ki+1]
					origin = breadcrumbs[ki+1]
				} else {
					x = kLines[ki-1] + 1
					origin = breadcrumbs[ki-1]
				}
				y = x - k
			}

			/** second, follow the k-line we're on, as far as possible **/

			cursor := &vertex{x, y, origin}
			{
				x, y = follow(diff.a, diff.b, x, y)
				if x != cursor.ai || y != cursor.bi {
					cursor = &vertex{x, y, cursor}
				}
				breadcrumbs[ki] = cursor
				kLines[ki] = x
			}

			/** third, check if we're at the end, and if so, construct our edit path **/

			if x >= m && y >= n {
				c := 0
				for v := cursor; v != nil; v = v.adj {
					c++
				}
				path := make([]*vertex, c)
				for i, v := c-1, cursor; v != nil; i, v = i-1, v.adj {
					path[i] = v
				}
				diff.edits = toEdits(path)
				return
			}
		}
	}
	return
}

func follow(a, b Sequence, x, y int) (int, int) {
	for x < a.Len() && y < b.Len() {
		atX, _ := a.At(x)
		atY, _ := b.At(y)
		if atX != atY {
			break
		}
		x, y = x+1, y+1
	}
	return x, y
}

func toEdits(path []*vertex) (edits []*Edit) {
	edits = make([]*Edit, 0)
	var x, y int
	flew := false
	for _, v := range path {
		if v.ai == x && v.bi == y {
			continue
		}
		if v.ai > x && v.bi > y {
			//match
			x, y, flew = v.ai, v.bi, true
			continue
		}
		flew = false
		var et editType
		//edit start, edit end
		var es, ee int
		if v.ai > x {
			//deletion
			et = DELETE
			es = x
			ee = v.ai - 1
		} else if v.bi > y {
			//insertion
			et = INSERT
			es = y
			ee = v.bi - 1
		} else {
			panic(fmt.Errorf("impossible"))
		}
		var e *Edit
		if len(edits) == 0 || flew {
			e = &Edit{es, ee, et}
			edits = append(edits, e)
		} else {
			e = edits[len(edits)-1]
		}
		if e.Type != et {
			edits = append(edits, &Edit{es, ee, et})
		} else {
			e.End = ee
		}
		x, y = v.ai, v.bi
	}
	return
}

func toString(path []*vertex) string {
	s := "["
	for _, v := range path {
		s += fmt.Sprintf("(%d,%d) ", v.ai, v.bi)
	}
	s += "]"
	return s
}
