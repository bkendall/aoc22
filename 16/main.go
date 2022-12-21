package main

import (
	"flag"
	"fmt"
	"log"
	"math"
	"os"
	"regexp"
	"strconv"
	"strings"
)

type Node struct {
	name        string
	open        bool
	flowRate    int
	connections []*Node
}

func main() {
	sourcePtr := flag.String("source", "./sample.txt", "input file")
	flag.Parse()

	buff, err := os.ReadFile(*sourcePtr)
	if err != nil {
		log.Fatalf("Error: %v", err)
	}
	input := strings.Split(string(buff), "\n")

	var world *Node
	worldMap := map[string]*Node{}

	lineRegex := regexp.MustCompile(`Valve (.+) has flow rate=(\d+); tunnel[s]? lead[s]? to valve[s]? (.+)`)

	for _, line := range input {
		arr := lineRegex.FindStringSubmatch(line)
		if arr == nil {
			log.Fatalf("Could not match string: %q\n", line)
		}
		name := arr[1]
		flowRate := toInt(arr[2])
		connections := strings.Split(arr[3], ", ")

		connectionsArr := []*Node{}
		for _, n := range connections {
			_, ok := worldMap[n]
			if ok {
				continue
			}
			worldMap[n] = &Node{name: n}
			connectionsArr = append(connectionsArr, worldMap[n])
		}

		node, ok := worldMap[name]
		if !ok {
			node = &Node{name: name}
		}
		node.flowRate = flowRate
		node.connections = connectionsArr
		if name == "AA" {
			world = node
		}
	}

	{
		// Part 1.
		// 30 minutes. 1m to move, 1m to open.
		maxMinutes := 30

		type Action int
		const (
			move Action = iota
			open
			wait
		)
		type Step struct {
			action                Action
			node                  *Node
			totalPressureReleased int
		}
		sumSteps := func(arr []*Step) int {
			sum := 0
			for i, s := range arr {
				switch s.action {
				case move:
					continue
				case open:
					// fmt.Printf("Valve %s opened on minute %d\n", s.node.name, i+1)
					// i is 0 index, minute 1. we have 30 minutes total, 29, 0-index.
					// We don't count _this_ minute, but we get the rest.
					score := s.node.flowRate * (maxMinutes - 1 - i)
					sum += score
				}
			}
			return sum
		}
		stepsHasOpenNode := func(arr []*Step, n *Node) bool {
			for _, s := range arr {
				if s.node == n && s.action == open {
					return true
				}
			}
			return false
		}

		var findPath func(seen []*Step, place *Node) int

		findPath = func(seen []*Step, place *Node) int {
			if len(seen) >= maxMinutes {
				// fmt.Printf("len of seen: %d\n", len(seen))
				// fmt.Println()
				// for i, s := range seen {
				// 	fmt.Printf("%d: %d @ %s\n", i, s.action, s.node.name)
				// }
				score := sumSteps(seen)
				// fmt.Printf("Score: %d\n", score)
				// fmt.Println()
				return score
			}

			// We can either open or not open a node, so we'll have to branch.
			maxScore := 0
			if !stepsHasOpenNode(seen, place) {
				seenOpenHere := append([]*Step{}, seen...)
				seenOpenHere = append(seenOpenHere, &Step{action: open, node: place})
				score := findPath(seenOpenHere, place)
				maxScore = max(maxScore, score)
			}
			// var maxArr []*Step
			for _, n := range place.connections {
				seenMove := append([]*Step{}, seen...)
				seenMove = append(seenMove, &Step{action: move, node: n})
				score := findPath(seenMove, n)
				maxScore = max(maxScore, score)
				// if score > maxScore {
				// 	maxScore = score
				// 	// maxArr = append([]*Step{}, seen...)
				// }
			}
			// Or maybe we just wait:
			stayHere := append([]*Step{}, seen...)
			stayHere = append(stayHere, &Step{action: wait, node: place})
			score := findPath(stayHere, place)
			maxScore = max(maxScore, score)
			return maxScore
		}

		score := findPath([]*Step{}, world)
		fmt.Printf("Found path with score: %d\n", score)
	}

	{
	}
}

func toInt(s string) int {
	v, err := strconv.Atoi(s)
	if err != nil {
		log.Fatalf("failed to parse %q: %v\n", s, err)
	}
	return v
}

func max(a, b int) int {
	return int(math.Max(float64(a), float64(b)))
}

func min(a, b int) int {
	return int(math.Min(float64(a), float64(b)))
}

func abs(a int) int {
	return int(math.Abs(float64(a)))
}

func floor(a int) int {
	return int(math.Floor(float64(a)))
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
