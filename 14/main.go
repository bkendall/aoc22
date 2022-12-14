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

type Material int

const (
	Air Material = iota
	Rock
	Sand
)

type World struct {
	m map[int]map[int]Material
}

func (w *World) Add(x, y int, m Material) {
	if _, ok := w.m[x]; !ok {
		w.m[x] = map[int]Material{}
	}
	w.m[x][y] = m
}

func (w *World) Get(x, y int) Material {
	if _, ok := w.m[x]; ok {
		return w.m[x][y]
	}
	return Air
}

func (w *World) XRange() (int, int) {
	minX, maxX := math.MaxInt, math.MinInt
	for x := range w.m {
		minX = min(minX, x)
		maxX = max(maxX, x)
	}
	return minX, maxX
}

func (w *World) YRange() (int, int) {
	minY, maxY := 0, math.MinInt
	for _, row := range w.m {
		for y := range row {
			minY = min(minY, y)
			maxY = max(maxY, y)
		}
	}
	return minY, maxY
}

func (w *World) String() string {
	s := ""
	xmin, xmax := w.XRange()
	ymin, ymax := w.YRange()
	for y := ymin; y <= ymax; y++ {
		for x := xmin; x <= xmax; x++ {
			switch w.Get(x, y) {
			case Air:
				s += "."
			case Rock:
				s += "#"
			case Sand:
				s += "+"
			}
		}
		s += "\n"
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

	world := World{map[int]map[int]Material{}}

	for _, line := range input {
		first := true
		x, y := 0, 0
		for _, pt := range strings.Split(line, " -> ") {
			arr := strings.Split(pt, ",")
			if first {
				first = false
				x, y = toInt(arr[0]), toInt(arr[1])
				world.Add(x, y, Rock)
				continue
			}
			toX, toY := toInt(arr[0]), toInt(arr[1])
			if x != toX {
				deltaX := int(math.Copysign(1, float64(toX-x)))
				for x != toX {
					world.Add(x, y, Rock)
					x += deltaX
				}
				world.Add(x, y, Rock)
			} else if y != toY {
				deltaY := int(math.Copysign(1, float64(toY-y)))
				for y != toY {
					world.Add(x, y, Rock)
					y += deltaY
				}
				world.Add(x, y, Rock)
			} else {
				log.Fatalf("Could not figure out which way to go.\n")
			}
		}
	}

	fmt.Printf("World:\n%s\n", world.String())

	_, maxY := world.YRange()
	sands := 0
	for {
		x, y := 500, 0
		sands++
		for {
			if y+1 >= maxY+2 {
				world.Add(x, y, Sand)
				break
			}
			if world.Get(x, y+1) == Air {
				y++
			} else if world.Get(x-1, y+1) == Air {
				x--
				y++
			} else if world.Get(x+1, y+1) == Air {
				x++
				y++
			} else {
				// We stop.
				world.Add(x, y, Sand)
				break
			}
		}
		if world.Get(500, 0) == Sand {
			break
		}
	}
	fmt.Printf("World:\n%s\n", world.String())
	fmt.Printf("Sands: %d\n", sands)
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

func abs(a, b int) int {
	return int(math.Abs(float64(a - b)))
}
