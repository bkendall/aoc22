package main

import (
	"flag"
	"fmt"
	"log"
	"math"
	"os"
	"regexp"
	"sort"
	"strconv"
	"strings"
)

type Node struct {
	name        string
	open        bool
	flowRate    int
	connections []string
}

var worldMap = map[string]*Node{}

func main() {
	sourcePtr := flag.String("source", "./sample.txt", "input file")
	flag.Parse()

	buff, err := os.ReadFile(*sourcePtr)
	if err != nil {
		log.Fatalf("Error: %v", err)
	}
	input := strings.Split(strings.TrimSpace(string(buff)), "\n")

	lineRegex := regexp.MustCompile(`Valve (.+) has flow rate=(\d+); tunnel[s]? lead[s]? to valve[s]? (.+)`)

	for _, line := range input {
		arr := lineRegex.FindStringSubmatch(line)
		if arr == nil {
			log.Fatalf("Could not match string: %q\n", line)
		}
		name := arr[1]
		flowRate := toInt(arr[2])
		connections := strings.Split(arr[3], ", ")

		node := &Node{name: name, flowRate: flowRate, connections: connections}
		worldMap[name] = node
	}

	{
		// Part 1.
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

		score := maxRate("AA", make(map[string]bool), maxMinutes)
		fmt.Printf("Found path with score: %d\n", score)
	}

	{
	}
}

var cache = map[string]int{}

func hashArgs(curr string, opened map[string]bool, t int) string {
	arr := []string{}
	for k := range opened {
		arr = append(arr, k)
	}
	sort.Strings(arr)
	return fmt.Sprintf("%s@%d::%s", curr, t, strings.Join(arr, ","))
}

func maxRateForArgs(curr string, opened map[string]bool, t int) int {
	h := hashArgs(curr, opened, t)
	score, ok := cache[h]
	if ok {
		return score
	}
	calculatedScore := maxRate(curr, opened, t)
	cache[h] = calculatedScore
	return calculatedScore
}

func maxRate(curr string, opened map[string]bool, t int) int {
	// fmt.Printf("maxRate %s, %+v, %d\n", curr.name, opened, t)
	if t <= 0 {
		return 0
	}
	maxScore := 0
	val := (t - 1) * worldMap[curr].flowRate
	currentOpened := copyMap(opened)
	currentOpened[curr] = true
	for _, n := range worldMap[curr].connections {
		if !opened[curr] && val != 0 {
			maxScore = max(maxScore, val+maxRateForArgs(n, currentOpened, t-2))
		}
		maxScore = max(maxScore, maxRateForArgs(n, opened, t-1))
	}
	return maxScore
}

func copyMap(m map[string]bool) map[string]bool {
	n := make(map[string]bool)
	for k, v := range m {
		n[k] = v
	}
	return n
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
