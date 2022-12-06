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
	input := strings.Split(string(buff), "\n")[0]

	{
		for i := 0; i < len(input); i++ {
			if unique(input[i : i+4]) {
				fmt.Printf("Marker at: %d\n", i+4)
				break
			}
		}
	}

	{
		for i := 0; i < len(input); i++ {
			if unique(input[i : i+14]) {
				fmt.Printf("Marker at: %d\n", i+14)
				break
			}
		}
	}
}

func unique(s string) bool {
	arr := strings.Split(s, "")
	fmt.Printf("unique: %v\n", s)
	m := map[string]bool{}
	for _, v := range arr {
		if m[v] {
			return false
		}
		m[v] = true
	}
	return true
}
