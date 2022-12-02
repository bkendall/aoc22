package main

import (
	"flag"
	"fmt"
	"log"
	"os"
	"strings"
)

func fatalf(format string, args ...any) {
	fmt.Printf(format, args...)
	os.Exit(1)
}

type round struct{ me, opp string }

var shapeScore = map[string]int{
	"A": 1,
	"B": 2,
	"C": 3,
}

func strToPairs(s string) []round {
	arr := []round{}
	for _, v := range strings.Split(s, "\n") {
		vs := strings.Split(v, " ")
		arr = append(arr, round{opp: vs[0], me: vs[1]})
	}
	return arr
}

var (
	loss = 0
	draw = 3
	win  = 6
)

func roundOutcome(opp, me string) int {
	switch me {
	case "X": // Rock.
		switch opp {
		case "A": // Rock.
			return draw
		case "B": // Paper.
			return loss
		case "C": // Scisors.
			return win
		}
	case "Y": // Paper.
		switch opp {
		case "A": // Rock.
			return win
		case "B": // Paper.
			return draw
		case "C": // Scisors.
			return loss
		}
	case "Z": // Scisors.
		switch opp {
		case "A": // Rock.
			return loss
		case "B": // Paper.
			return win
		case "C": // Scisors.
			return draw
		}
	}
	log.Fatalf("Should never happen: %v vs %v", opp, me)
	return -1
}

func throwForResult(opp, me string) string {
	switch me {
	case "X": // Lose.
		switch opp {
		case "A": // Rock.
			return "Z"
		case "B": // Paper.
			return "X"
		case "C": // Scisors.
			return "Y"
		}
	case "Y": // Draw.
		switch opp {
		case "A": // Rock.
			return "X"
		case "B": // Paper.
			return "Y"
		case "C": // Scisors.
			return "Z"
		}
	case "Z": // Win.
		switch opp {
		case "A": // Rock.
			return "Y"
		case "B": // Paper.
			return "Z"
		case "C": // Scisors.
			return "X"
		}
	}
	log.Fatalf("Should never happen: %v vs %v", opp, me)
	return ""
}

func main() {
	sourcePtr := flag.String("source", "./sample.txt", "input file")
	flag.Parse()

	buff, err := os.ReadFile(*sourcePtr)
	if err != nil {
		fatalf("Error: %v", err)
	}
	input := string(buff)

	entries := strToPairs(input)
	fmt.Printf("entries: %+v\n", entries)

	sum := 0
	for _, e := range entries {
		s := 0
		toThrow := throwForResult(e.opp, e.me)
		switch toThrow {
		case "X":
			s = 1
		case "Y":
			s = 2
		case "Z":
			s = 3
		}
		o := roundOutcome(e.opp, toThrow)

		fmt.Printf("%d + %d = %d\n", s, o, s+o)
		sum = sum + s + o
	}
	fmt.Printf("Sum: %d\n", sum)
}
