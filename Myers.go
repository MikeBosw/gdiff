package diff

import (
	/*"github.com/mikebosw/ifo"*/
	"fmt"
)

type edit struct {
	start, end int
	et editType
}

type editType rune

const (
	INSERT editType = 'i'
	DELETE editType = 'd'
)

//ABC and DEF

func toString(path []*vertex) string {
	s := "["
	for _, v := range path {
		s += fmt.Sprintf("(%d,%d) ", v.ai, v.bi)
	}
	s += "]"
	return s
}

func MyersDiff(a, b string) (edits []*edit) {
	if a == b {
		return
	}

	m, n := len(a), len(b)

	kLines := make([]int, (m + n)*2 + 1)

	follow := func(x, y int) (int, int) {
		for ; x < m && y < n;  {
			if a[x] != b[y] {
				break
			}
			x, y = x + 1, y + 1
		}
		return x, y
	}


	paths := make([][]*vertex, (m+n)*2 + 1)

outer:
	for d := 0; d <= m + n; d++ {
		for k := -d; k <= d; k += 2 {
			/* first, establish what our X, Y is, based on how we had to have gotten here.
			 * we must have gotten here in one of two ways: either down from the k+1 line, or right from the k-1 line.
			 * we pick whichever gets us further towards {m, n}. */
			var x, y int
			/* if we're at the bottom, leftmost k-line, the only way to have gotten here is down (from k+1). */
			bottomLeft := k == -d
			/* if we're at the top, rightmost k-line, the only way to have gotten here is right (from k-1). */
			topRight := k == d
			/* index into a k-line array for this k-line */
			ki := m + n + k

			var xs, ys int //, o_k, o_ki int
			var origins []*vertex
			if bottomLeft || (!topRight && kLines[ki - 1] < kLines[ki + 1]) {
				x = kLines[ki + 1]
				xs = x
				ys = xs - (k + 1)
				origins = paths[ki + 1]
			} else {
				x = kLines[ki - 1] + 1
				xs = x - 1
				ys = xs - (k - 1)
				origins = paths[ki - 1]
			}

			y = x - k

			o_copy := make([]*vertex, len(origins))
			copy(o_copy, origins)
			origins = o_copy

			if len(origins) == 0 {
				/*fmt.Printf("no previous vertex stored on path for k-line %d\n", o_k)*/
			} else {
				prev := origins[len(origins) - 1]
				/*fmt.Printf("previous vertex on path was %d, %d (from k-line %d)\n", prev.ai, prev.bi, o_k)*/
				if prev.ai != xs || prev.bi != ys {
					panic(fmt.Errorf("data integrity compromised: previous vertex (%d, %d) should be: (%d, %d)",
						prev.ai, prev.bi, xs, ys))
				}
			}
			/*fmt.Printf("printing path thus far (from k-line %d [%p]):\n", o_k, origins)
			for _, v := range origins {
				fmt.Printf("%d, %d\n", v.ai, v.bi)
			}
			fmt.Printf("printing complete.\n")
			fmt.Printf("adding %d, %d after origin path [%p] to create new path for k-line %d \n", x, y, origins, k)*/
			cursor := &vertex{x, y}
			path := append(origins, cursor)
			/*fmt.Printf("%p <- extended from k%d [%p] \n", path, o_k, origins)
			fmt.Printf("new path [%p] created for k-line %d from origin path [%p] (from k-line %d) \n", path, k, origins, o_k)*/

			x, y = follow(x, y)

			if x != cursor.ai || y != cursor.bi {
				/*fmt.Printf("followed diagonal from %d, %d to %d, %d \nadding %d, %d to the path [%p] for k-line %d\n",
					cursor.ai, cursor.bi, x, y, x, y, path, k)
				preAppend := path*/
				path = append(path, &vertex{x, y})
				/*fmt.Printf("%p <- extended from itself [%p]\n", path, preAppend)*/
			}

			/*fmt.Printf("writing path [%p] for k-line %d to index %d\n", path, k, ki)*/
			paths[ki] = path
			/*fmt.Printf("k%d <- %p %s", k, path, toString(path))*/

			//fmt.Printf("from %d, %d got to %d, %d (k-line %d) \n", xs, ys, x, y, k)

			kLines[ki] = x

			if x >= m && y >= n {
				fmt.Printf("edits required: %d\n", d)
				fmt.Printf("path for k-line %d: \n", k)
				for _, v := range paths[ki] {
					fmt.Printf("%d, %d\n", v.ai, v.bi)
				}
				edits = toEditsV(paths[ki])
				for _, v := range edits {
					fmt.Printf("edit: %c from %d to %d\n", v.et, v.start, v.end)
				}
				break outer
			}
		}
	}
	return
}

func toEditsV(path []*vertex) (edits []*edit) {
	edits = make([]*edit, 0)
	var x, y int
	flew := false
	for _, v := range path {
		fmt.Printf("{%d,%d} ", v.ai, v.bi)
		if v.ai == x && v.bi == y {
			continue
		}
		if v.ai > x && v.bi > y {
			fmt.Printf("(match) \n")
			x, y, flew = v.ai, v.bi, true
			continue
		}
		flew = false
		var es, ee int
		var et editType
		if v.ai > x {
			fmt.Printf("(deletion) \n")
			et = DELETE
			es = x
			ee = v.ai
		} else if v.bi > y {
			fmt.Printf("(insertion) \n")
			et = INSERT
			es = y
			ee = v.bi
		} else {
			panic(fmt.Errorf("impossible"))
		}
		var e *edit
		if len(edits) == 0 || flew {
			e = &edit{es, ee, et}
			edits = append(edits, e)
		} else {
			e = edits[len(edits)-1]
		}
		if e.et != et {
			edits = append(edits, &edit{es, ee, et})
		} else {
			e.end = ee
		}
		x, y = v.ai, v.bi
	}
	return
}

func toEditsN(path []*node) (edits []*edit) {
	var origin *node
	edits = make([]*edit, 0)
	for i, v := range path {
		if i == 0 {
			origin = v
			continue
		}
		if origin.aChar - v.aChar == 0 && origin.bChar - v.bChar == 0 {
			origin = v
			continue
		}
		//...
	}
	return
}


func route(origin *node, graph *node) (path []*node) {
	return
}

/*func Graph(a, b string) (graph *node) {
	populate := func (v *node) {
		neighbors := make([]*node, 0)
		var nextA, nextB int
		if nextA = v.aChar + 1; v.aChar < len(a) {
			neighbors = append(neighbors, &node{nextA, v.bChar, nil})
		}
		if nextB = v.bChar + 1; v.bChar < len(b) {
			neighbors = append(neighbors, &node{v.aChar, nextB, nil})
		}
		if v.aChar < len(a) && v.bChar < len(b) {
			if a[nextA - 1] == b[nextB - 1] {
				neighbors = append(neighbors, &node{nextA, nextB, nil})
			}
		}
		v.neighbors = neighbors
	}
	stack := ifo.NewStack()
	graph = &node {0, 0, nil}
	stack.Push(graph)
	for cursor := stack.Pop(); cursor != nil; {
		v := cursor.(*node)
		//fmt.Printf("node %x\n", v)
		populate(v)
		for _, n := range v.neighbors {
			//fmt.Printf("found neighbor %x for node %x\n", n, v)
			stack.Push(n)
		}
		cursor = stack.Pop()
	}
	return
}*/

/*func (v *node) toMatrix(a, b string, matrix [][]string) {
	var aChar, bChar string
	if indexA := v.aChar - 1; indexA >= 0 && indexA < len(a) {
		aChar = a[indexA]
	} else {
		aChar = ""
	}
	if indexB := v.bChar - 1; indexB >= 0 && indexB < len(b) {
		bChar = b[indexB]
	} else {
		bChar = ""
	}
	for _, n := range v.neighbors {
		n.toMatrix(a, b, matrix)
	}
}*/

func (v *node) printDFS(head, tail, a, b string) {
	if v == nil {
		return
	}
	for _, n := range v.neighbors {
		//fmt.Printf("node %x has neighbor %x \n", v, n)
		var step uint8
		newTail, newHead := tail, head
		switch {
		case n.aChar > v.aChar && n.bChar > v.bChar:
			step = a[n.aChar-1]
			if len(newTail) > 1 {
				newTail = newTail[1:]
			} else {
				newTail = ""
			}
			fmt.Printf("from %s[%c]%s, match on %c -> \"%s\" \n", newHead, step, newTail, step, newHead + newTail)
			newHead = fmt.Sprintf("%s%c", newHead, step)
		case n.aChar > v.aChar:
			step = a[n.aChar - 1]
			if len(newTail) > 1 {
				newTail = newTail[1:]
			} else {
				newTail = ""
			}
			fmt.Printf("from %s[%c]%s, remove %c -> \"%s\" \n", newHead, step, newTail, step, newHead + newTail)
		case n.bChar > v.bChar:
			step = b[n.bChar - 1]
			var cursor string
			if len(newTail) <= 1 {
				cursor = fmt.Sprintf("[%s]", newTail)
			} else {
				cursor = fmt.Sprintf("[%s]%s", newTail[0:1], newTail[1:])
			}
			newHead = fmt.Sprintf("%s%c", newHead, step)
			fmt.Printf("from %s%s, insert %c -> \"%s\" \n", head, cursor, step, newHead + newTail)
		default:
			panic(fmt.Errorf("node has a path to itself or to an earlier node: %x", n))
		}
		n.printDFS(newHead, newTail, a, b)
	}
	if len(v.neighbors) == 0 {
		fmt.Println()
	}
}

/*func PrintDFS(a, b string) {
	v := Graph(a, b)
	v.printDFS("", a, a, b)
}*/

/*func PrintMatrix(a, b string) {
	origin := Graph(a, b)
	d := math.Max(len(a), len(b)) + 1
	matrix := make([][]string, d)
	for i := 0; i < d; i++ {
		matrix[i] := make([][]string, d)
	}
	origin.toMatrix(a, b, matrix)

}*/

type node struct {
	/* 1-based index into a string (subtract 1 when indexing into the string) */
	aChar, bChar int
	neighbors []*node
}

type vertex struct {
	/* 1-based index into a string (subtract 1 when indexing into the string) */
	ai, bi int
}
