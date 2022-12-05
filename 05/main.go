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

type stack struct {
	arr []string
}

func (s *stack) Pop() (string, bool) {
	if len(s.arr) == 0 {
		return "", false
	}
	last := len(s.arr) - 1
	v := s.arr[last]
	s.arr = s.arr[:last]
	return v, true
}

func (s *stack) PopN(n int) ([]string, bool) {
	if n == 0 {
		return []string{}, false
	}
	last := len(s.arr)
	first := max(last-n, 0)
	v := s.arr[first:last]
	if first == 0 {
		s.arr = []string{}
	} else {
		s.arr = s.arr[:first]
	}
	return v, true
}

func (s *stack) Push(v string) {
	s.arr = append(s.arr, v)
}

func (s *stack) PushArr(v []string) {
	s.arr = append(s.arr, v...)
}

func max(a, b int) int {
	return int(math.Max(float64(a), float64(b)))
}

type instruction struct {
	count, from, to int
}

func atoi(s string) int {
	i, err := strconv.Atoi(s)
	if err != nil {
		log.Fatalf("Failed to convert %q to int: %v", s, err)
	}
	return i
}

func main() {
	sourcePtr := flag.String("source", "./sample.txt", "input file")
	flag.Parse()

	buff, err := os.ReadFile(*sourcePtr)
	if err != nil {
		log.Fatalf("Error: %v", err)
	}
	input := strings.Split(string(buff), "\n")

	world := map[int]*stack{}
	world2 := map[int]*stack{}
	for i, line := range input {
		if line == "" {
			break
		}
		s := &stack{}
		s2 := &stack{}
		vs := strings.Split(strings.Split(line, ":")[1], ",")
		for _, v := range vs {
			s.Push(v)
			s2.Push(v)
		}
		world[i+1] = s
		world2[i+1] = s2
	}

	insts := []instruction{}
	input = input[len(world)+1:]
	for _, line := range input {
		vs := strings.Split(line, " ")
		insts = append(insts, instruction{
			count: atoi(vs[1]),
			from:  atoi(vs[3]),
			to:    atoi(vs[5]),
		})
	}

	{
		// Part 1.
		for _, inst := range insts {
			for i := 1; i <= inst.count; i++ {
				v, ok := world[inst.from].Pop()
				if ok {
					world[inst.to].Push(v)
				}
			}
		}

		for i := 1; i <= len(world); i++ {
			v, ok := world[i].Pop()
			if ok {
				fmt.Printf("%v", v)
			}
		}
		fmt.Println()
	}

	{
		// Part 2.
		world = world2
		for _, inst := range insts {
			v, ok := world[inst.from].PopN(inst.count)
			if ok {
				world[inst.to].PushArr(v)
			}
		}

		for i := 1; i <= len(world); i++ {
			v, ok := world[i].Pop()
			if ok {
				fmt.Printf("%v", v)
			}
		}
		fmt.Println()
	}
}
