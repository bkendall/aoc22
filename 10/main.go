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

type cpu struct {
	x int
}

type instruction struct {
	cmd string
	v   int
}

func main() {
	sourcePtr := flag.String("source", "./sample.txt", "input file")
	flag.Parse()

	buff, err := os.ReadFile(*sourcePtr)
	if err != nil {
		log.Fatalf("Error: %v", err)
	}
	input := strings.Split(string(buff), "\n")

	insts := []instruction{}
	for _, line := range input {
		arr := strings.Split(line, " ")
		v := 0
		if len(arr) > 1 {
			var err error
			v, err = strconv.Atoi(arr[1])
			if err != nil {
				log.Fatalf("Failed to parse %q: %v", arr[1], err)
			}
		}
		insts = append(insts, instruction{cmd: arr[0], v: v})
	}

	valueByCycle := map[int]int{}
	{
		cycle := 0
		c := cpu{x: 1}
		for _, inst := range insts {
			cycle++
			valueByCycle[cycle] = c.x
			switch inst.cmd {
			case "noop":
				continue
			case "addx":
				cycle += 1
				c.x += inst.v
			}
		}

		fmt.Printf("Cycles: %d\nCPU: %+v\n", cycle, c)

		sum := 0
		marks := []int{20, 60, 100, 140, 180, 220}
		for _, m := range marks {
			v := findCycleValue(valueByCycle, m)
			sum += v * m
		}
		fmt.Printf("Sum: %d\n", sum)
	}

	{
		s := ""
		cycle := 1
		pos := findCycleValue(valueByCycle, cycle)
		for row := 0; row < 6; row++ {
			for col := 0; col < 40; col++ {
				if pos-1 == col || col == pos || col == pos+1 {
					s += "#"
				} else {
					s += "."
				}
				cycle++
				pos = findCycleValue(valueByCycle, cycle)
			}

			fmt.Printf("%s\n", s)
			s = ""
		}
	}

}

func findCycleValue(m map[int]int, v int) int {
	for i := v; i >= 0; i-- {
		if vv, ok := m[i]; ok {
			return vv
		}
	}
	log.Fatalf("Could not find a value for %d in %+v", v, m)
	return 0
}

func max(a, b int) int {
	return int(math.Max(float64(a), float64(b)))
}

func abs(a, b int) int {
	return int(math.Abs(float64(a - b)))
}
