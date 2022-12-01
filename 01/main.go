package main

import (
	"flag"
	"fmt"
	"os"
	"sort"
	"strconv"
	"strings"
)

func fatalf(format string, args ...any) {
	fmt.Printf(format, args...)
	os.Exit(1)
}

type entry struct{ i, v int }

func main() {
	sourcePtr := flag.String("source", "./sample.txt", "input file")
	flag.Parse()

	buff, err := os.ReadFile(*sourcePtr)
	if err != nil {
		fatalf("Error: %v", err)
	}
	input := string(buff)

	nums := strArrToNumArr(input)

	sums := []entry{}
	sum := 0
	i := 0
	for _, v := range nums {
		if v == -1 {
			sums = append(sums, entry{i, sum})
			sum = 0
			i++
			continue
		}
		sum += v
	}
	sums = append(sums, entry{i, sum})

	max := -1
	for _, v := range sums {
		if v.v > max {
			max = v.v
		}
	}

	fmt.Printf("Max val: %v\n", max)

	sort.Slice(sums, func(a, b int) bool { return sums[a].v > sums[b].v })

	fmt.Printf("Top 3 sum: %v\n", sums[0].v+sums[1].v+sums[2].v)
}

func strArrToNumArr(s string) []int {
	arr := []int{}
	for _, v := range strings.Split(s, "\n") {
		if v == "" {
			arr = append(arr, -1)
			continue
		}
		i, err := strconv.Atoi(v)
		if err != nil {
			panic(err)
		}
		arr = append(arr, i)
	}
	return arr
}
