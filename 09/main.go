package main

import (
	"flag"
	"fmt"
	"log"
	"math"
	"os"
	"strconv"
	"strings"
)

type move struct {
	dir  string
	dist int
}

type pointer struct {
	x, y int
}

type coordinateMap struct {
	m map[int]map[int]bool
}

func (m *coordinateMap) Add(p pointer) {
	if _, ok := m.m[p.x]; !ok {
		m.m[p.x] = map[int]bool{}
	}
	m.m[p.x][p.y] = true
}

func (m *coordinateMap) NumVisited() int {
	s := 0
	for _, r := range m.m {
		s += len(r)
	}
	return s
}

func main() {
	sourcePtr := flag.String("source", "./sample.txt", "input file")
	flag.Parse()

	buff, err := os.ReadFile(*sourcePtr)
	if err != nil {
		log.Fatalf("Error: %v", err)
	}
	input := strings.Split(string(buff), "\n")

	moves := []move{}
	for _, line := range input {
		arr := strings.Split(line, " ")
		n, err := strconv.Atoi(arr[1])
		if err != nil {
			log.Fatalf("Failed to parse %q: %v", arr[1], err)
		}
		moves = append(moves, move{dir: arr[0], dist: n})
	}

	{
		visited := coordinateMap{map[int]map[int]bool{}}
		knots := make([]pointer, 10)
		head := &knots[0]

		for _, m := range moves {
			fmt.Printf("== %s %d\n", m.dir, m.dist)
			for i := 0; i < m.dist; i++ {
				switch m.dir {
				case "R":
					head.x++
				case "L":
					head.x--
				case "U":
					head.y++
				case "D":
					head.y--
				}

			loop:
				for k := 1; k < len(knots); k++ {
					vectorX := knots[k-1].x - knots[k].x
					vectorY := knots[k-1].y - knots[k].y
					deltaX := abs(knots[k].x, knots[k-1].x)
					deltaY := abs(knots[k].y, knots[k-1].y)
					switch {
					case deltaX == 0 && deltaY == 0,
						deltaX == 1 && deltaY == 1,
						deltaX == 0 && deltaY == 1,
						deltaX == 1 && deltaY == 0:
						// Nothing to do in this case. We're done.
						break loop
					case deltaX == 1 && deltaY == 2,
						deltaX == 2 && deltaY == 1,
						deltaX == 2 && deltaY == 2:
						knots[k].y += int(math.Copysign(1, float64(vectorY)))
						knots[k].x += int(math.Copysign(1, float64(vectorX)))
					case deltaX == 0 && deltaY == 2:
						knots[k].y += int(math.Copysign(1, float64(vectorY)))
					case deltaX == 2 && deltaY == 0:
						knots[k].x += int(math.Copysign(1, float64(vectorX)))
					default:
						log.Fatalf("Deal with this: %d, %d", deltaX, deltaY)
					}
				}

				visited.Add(knots[len(knots)-1])
			}

			printPts(knots)
			fmt.Println()
			fmt.Println()
		}
		fmt.Printf("Num Visited: %d\n", visited.NumVisited())
	}

}

func max(a, b int) int {
	return int(math.Max(float64(a), float64(b)))
}

func abs(a, b int) int {
	return int(math.Abs(float64(a - b)))
}

func printPts(pts []pointer) {
	m := map[int]map[int]string{}
	for i, p := range pts {
		if _, ok := m[p.x]; !ok {
			m[p.x] = map[int]string{}
		}
		if i == 0 {
			m[p.x][p.y] = "H"
		} else {
			m[p.x][p.y] = fmt.Sprintf("%d", i)
		}
	}
	s := ""
	for y := 20; y >= -20; y-- {
		for x := -20; x <= 20; x++ {
			if _, ok := m[x]; ok {
				if v, ok := m[x][y]; ok {
					s += v
					continue
				}
			}
			if x == 0 && y == 0 {
				s += "s"
				continue
			}
			s += "."
		}
		s += "\n"
	}
	fmt.Printf("%s", s)
}
