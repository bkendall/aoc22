package main

import (
	"flag"
	"fmt"
	"log"
	"math"
	"os"
	"strings"
	"sync"
)

type Node struct {
	c    string
	x, y int

	up, down, left, right *Node
}

func (n *Node) Val() int {
	if n.c == "E" {
		return int('z')
	} else if n.c == "S" {
		return int('a')
	}
	return int(n.c[0])
}

func (n *Node) Neighbors() []*Node {
	arr := []*Node{}
	if n.up != nil && n.canGoTo(n.up) {
		arr = append(arr, n.up)
	}
	if n.right != nil && n.canGoTo(n.right) {
		arr = append(arr, n.right)
	}
	if n.down != nil && n.canGoTo(n.down) {
		arr = append(arr, n.down)
	}
	if n.left != nil && n.canGoTo(n.left) {
		arr = append(arr, n.left)
	}
	return arr
}

func (n *Node) canGoTo(n2 *Node) bool {
	return n2.Val() <= n.Val()+1
}

func main() {
	sourcePtr := flag.String("source", "./sample.txt", "input file")
	flag.Parse()

	buff, err := os.ReadFile(*sourcePtr)
	if err != nil {
		log.Fatalf("Error: %v", err)
	}
	input := strings.Split(string(buff), "\n")

	aStarts := []*Node{}
	var start, end *Node
	var xPtr, yPtr *Node
	for y, line := range input {
		arr := strings.Split(line, "")
		var first *Node
		for x, v := range arr {
			n := &Node{c: v, x: x, y: y}
			if v == "S" {
				start = n
				aStarts = append(aStarts, n)
			} else if v == "E" {
				end = n
			} else if v == "a" {
				aStarts = append(aStarts, n)
			}

			if first == nil {
				first = n
			}

			if xPtr != nil {
				xPtr.right = n
				n.left = xPtr
				xPtr = n
			} else {
				xPtr = n
			}
			if yPtr != nil {
				yPtr.down = n
				n.up = yPtr
				yPtr = yPtr.right
			}
		}
		xPtr = nil
		yPtr = first
	}

	// Part 1
	path := AStar(start, end, func(n *Node) int { return 1 })
	fmt.Printf("Path of length %d steps found!\n", len(path)-1)

	// Part 2
	fmt.Printf("Part 2\n")
	fmt.Printf("%d choices to start from\n", len(aStarts))
	var wg sync.WaitGroup
	var mu sync.Mutex
	min := math.MaxInt
	for _, aStart := range aStarts {
		wg.Add(1)
		thisStart := aStart
		go func() {
			path := AStar(thisStart, end, func(n *Node) int { return 1 })
			mu.Lock()
			if len(path) > 0 && len(path) < min {
				min = len(path)
			}
			mu.Unlock()
			wg.Done()
		}()
	}
	wg.Wait()
	fmt.Printf("Most scenic route is %d steps!\n", min-1)
}

func rebuildPath(cameFrom map[*Node]*Node, current *Node) []*Node {
	totalPath := []*Node{current}
	_, ok := cameFrom[current]
	for ok {
		current = cameFrom[current]
		totalPath = append([]*Node{current}, totalPath...)
		_, ok = cameFrom[current]
	}
	return totalPath
}

func AStar(start, goal *Node, h func(n *Node) int) []*Node {
	// The set of discovered nodes that may need to be (re-)expanded.
	// Initially, only the start node is known.
	// This is usually implemented as a min-heap or priority queue rather than a hash-set.
	queue := Queue{}
	queue.Add(start)

	// For node n, cameFrom[n] is the node immediately preceding it on the cheapest path from start
	// to n currently known.
	cameFrom := map[*Node]*Node{}

	// For node n, gScore[n] is the cost of the cheapest path from start to n currently known.
	gScore := NewMapWithDefault(math.MaxInt)
	gScore.Set(start, 0)

	// For node n, fScore[n] := gScore[n] + h(n). fScore[n] represents our current best guess as to
	// how cheap a path could be from start to finish if it goes through n.
	fScore := NewMapWithDefault(math.MaxInt)
	fScore.Set(start, h(start))

	for queue.Length() > 0 {
		current := lowestScored(queue, fScore)
		if current == goal {
			return rebuildPath(cameFrom, current)
		}

		queue.Remove(current)
		for _, neighbor := range current.Neighbors() {
			// d(current,neighbor) is the weight of the edge from current to neighbor
			// tentative_gScore is the distance from start to the neighbor through current
			tentative_gScore := gScore.Get(current) + stepScore(current, neighbor)
			if tentative_gScore < gScore.Get(neighbor) {
				cameFrom[neighbor] = current
				gScore.Set(neighbor, tentative_gScore)
				fScore.Set(neighbor, tentative_gScore+h(neighbor))
				if !queue.Has(neighbor) {
					queue.Add(neighbor)
				}
			}
		}
	}

	return []*Node{}
}

type MapWithDefault struct {
	def int
	m   map[*Node]int
}

func NewMapWithDefault(def int) *MapWithDefault {
	return &MapWithDefault{
		def: def,
		m:   map[*Node]int{},
	}
}

func (m *MapWithDefault) Get(v *Node) int {
	if v, ok := m.m[v]; ok {
		return v
	}
	return m.def
}

func (m *MapWithDefault) Set(v *Node, i int) {
	m.m[v] = i
}

func lowestScored(q Queue, m *MapWithDefault) *Node {
	min := math.MaxInt
	var ret *Node
	for _, n := range q.Values() {
		score := m.Get(n)
		if score < min {
			min, ret = score, n
		}
	}
	return ret
}

// Returns a 1 if the step is valid, a large, large number otherwise.
func stepScore(from, to *Node) int {
	if to.Val() <= from.Val()+1 {
		return 1
	}
	return math.MaxInt / 4
}

type Queue struct {
	arr []*Node
}

func (q *Queue) Length() int {
	return len(q.arr)
}

func (q *Queue) Values() []*Node {
	return q.arr
}

func (q *Queue) Add(n *Node) {
	q.arr = append(q.arr, n)
}

func (q *Queue) Has(node *Node) bool {
	for _, n := range q.arr {
		if n == node {
			return true
		}
	}
	return false
}

func (q *Queue) Remove(node *Node) {
	i := -1
	for j := 0; j < len(q.arr); j++ {
		if q.arr[j] == node {
			i = j
		}
	}

	if i == -1 {
		return
	}
	if len(q.arr) == 1 {
		q.arr = []*Node{}
		return
	}
	copy(q.arr[i:], q.arr[i+1:])
	q.arr = q.arr[:len(q.arr)-1]
}

func max(a, b int) int {
	return int(math.Max(float64(a), float64(b)))
}

func min(a, b int) int {
	return int(math.Min(float64(a), float64(b)))
}

func abs(a, b int) int {
	return int(math.Abs(float64(a - b)))
}
