package main

import (
	"flag"
	"fmt"
	"log"
	"os"
	"strings"
)

func main() {
	sourcePtr := flag.String("source", "./sample.txt", "input file")
	flag.Parse()

	buff, err := os.ReadFile(*sourcePtr)
	if err != nil {
		log.Fatalf("Error: %v", err)
	}
	input := strings.Split(string(buff), "\n")

	sum := 0
	for i := 0; i < len(input)-1; i += 3 {
		one := input[i]
		two := input[i+1]
		three := input[i+2]

		m := map[rune]int{}
		for _, c := range one {
			m[c] = 1
		}
		for _, c := range two {
			if m[c] == 1 {
				m[c]++
			}
		}
		for _, c := range three {
			if m[c] == 2 {
				m[c]++
				v := int(c) - 96
				if v < 0 {
					v += 58
				}
				sum += v
			}
		}
	}

	fmt.Printf("Sum: %d\n", sum)
}
