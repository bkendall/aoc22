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

type tree struct {
	height, score int
}

func main() {
	sourcePtr := flag.String("source", "./sample.txt", "input file")
	flag.Parse()

	buff, err := os.ReadFile(*sourcePtr)
	if err != nil {
		log.Fatalf("Error: %v", err)
	}
	input := strings.Split(string(buff), "\n")

	world := map[int]map[int]*tree{}
	for y, line := range input {
		for x, c := range strings.Split(line, "") {
			v, err := strconv.Atoi(c)
			if err != nil {
				log.Fatalf("Could not convert %q: %v", c, err)
			}
			if _, ok := world[y]; !ok {
				world[y] = map[int]*tree{}
			}
			world[y][x] = &tree{height: v}
		}
	}

	maxY := len(world) - 1
	maxX := len(world[0]) - 1

	{
		numVisible := 0
		for y := 0; y <= maxY; y++ {
			for x := 0; x <= maxX; x++ {
				visible := isVisible(world, x, y)
				if visible {
					numVisible++
				}
			}
		}

		fmt.Printf("Number visible: %d\n", numVisible)
	}

	{
		topScore := -1
		for y := 0; y <= maxY; y++ {
			for x := 0; x <= maxX; x++ {
				score := scenicScore(world, x, y)
				world[y][x].score = score
				topScore = max(topScore, score)
			}
		}

		fmt.Printf("Top score: %d\n", topScore)
	}

}

func max(a, b int) int {
	return int(math.Max(float64(a), float64(b)))
}

func isVisible(w map[int]map[int]*tree, x, y int) bool {
	maxY := len(w) - 1
	maxX := len(w[0]) - 1
	hiddenDirections := 0
	v := w[y][x].height
	for i := x - 1; i >= 0; i-- {
		if w[y][i].height >= v {
			hiddenDirections++
			break
		}
	}
	for i := x + 1; i <= maxX; i++ {
		if w[y][i].height >= v {
			hiddenDirections++
			break
		}
	}
	for i := y - 1; i >= 0; i-- {
		if w[i][x].height >= v {
			hiddenDirections++
			break
		}
	}
	for i := y + 1; i <= maxY; i++ {
		if w[i][x].height >= v {
			hiddenDirections++
			break
		}
	}
	return hiddenDirections != 4
}

func scenicScore(w map[int]map[int]*tree, x, y int) int {
	maxY := len(w) - 1
	maxX := len(w[0]) - 1
	v := w[y][x].height
	scores := []int{0, 0, 0, 0}
	for i := x - 1; i >= 0; i-- {
		scores[0]++
		if w[y][i].height >= v {
			break
		}
	}
	for i := x + 1; i <= maxX; i++ {
		scores[1]++
		if w[y][i].height >= v {
			break
		}
	}
	for i := y - 1; i >= 0; i-- {
		scores[2]++
		if w[i][x].height >= v {
			break
		}
	}
	for i := y + 1; i <= maxY; i++ {
		scores[3]++
		if w[i][x].height >= v {
			break
		}
	}
	return scores[0] * scores[1] * scores[2] * scores[3]
}
