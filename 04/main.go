package main

import (
	"flag"
	"fmt"
	"log"
	"os"
	"strconv"
	"strings"
)

type section struct {
	start, end int
}

type pair struct {
	one, two section
}

func (one section) Contains(second section) bool {
	startWithin := second.start >= one.start && second.start <= one.end
	endWithin := second.end <= one.end && second.end >= one.start
	return startWithin && endWithin
}

func (one section) Overlap(second section) bool {
	startWithin := second.start >= one.start && second.start <= one.end
	endWithin := second.end <= one.end && second.end >= one.start
	return startWithin || endWithin
}

func main() {
	sourcePtr := flag.String("source", "./sample.txt", "input file")
	flag.Parse()

	buff, err := os.ReadFile(*sourcePtr)
	if err != nil {
		log.Fatalf("Error: %v", err)
	}
	input := strings.Split(string(buff), "\n")

	pairs := []pair{}
	for _, line := range input {
		ranges := strings.Split(line, ",")
		p := pair{}
		oneVals := strings.Split(ranges[0], "-")
		twoVals := strings.Split(ranges[1], "-")
		var err error
		p.one.start, err = strconv.Atoi(oneVals[0])
		if err != nil {
			log.Fatalf("Failed to parse int: %v", err)
		}
		p.one.end, err = strconv.Atoi(oneVals[1])
		if err != nil {
			log.Fatalf("Failed to parse int: %v", err)
		}
		p.two.start, err = strconv.Atoi(twoVals[0])
		if err != nil {
			log.Fatalf("Failed to parse int: %v", err)
		}
		p.two.end, err = strconv.Atoi(twoVals[1])
		if err != nil {
			log.Fatalf("Failed to parse int: %v", err)
		}
		pairs = append(pairs, p)
	}

	{
		// Part 1.
		count := 0
		for _, p := range pairs {
			if p.one.Contains(p.two) || p.two.Contains(p.one) {
				count++
			}
		}

		fmt.Printf("Inclusive pairs: %d\n", count)
	}

	{
		// Part 2.
		count := 0
		for _, p := range pairs {
			if p.one.Overlap(p.two) || p.two.Overlap(p.one) {
				count++
			}
		}

		fmt.Printf("Overlaping pairs: %d\n", count)
	}
}
