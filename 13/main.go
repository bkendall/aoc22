package main

import (
	"flag"
	"fmt"
	"log"
	"math"
	"os"
	"reflect"
	"sort"
	"strconv"
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

	partTwoArr := [][]any{
		{[]any{2}},
		{[]any{6}},
	}

	sum, idx := 0, 1
	for i := 0; i < len(input); i += 3 {
		// for i := 0; i < 3; i += 3 {
		rawLeft := input[i]
		rawRight := input[i+1]

		left := parseString(rawLeft)
		right := parseString(rawRight)
		partTwoArr = append(partTwoArr, left, right)

		res, ok := listInOrder(left, right)
		fmt.Printf("Left: %v\nRight: %v\nIn Order: %v, %v\n\n", left, right, res, ok)

		if res {
			sum += idx
		}
		idx++
	}

	fmt.Printf("Sum: %d\n", sum)

	sort.Slice(partTwoArr, func(i, j int) bool {
		res, ok := listInOrder(partTwoArr[i], partTwoArr[j])
		if !ok {
			return true
		}
		return res
	})

	product := 1
	for i, a := range partTwoArr {
		str := fmt.Sprintf("%v", a)
		if str == "[[2]]" || str == "[[6]]" {
			product *= (i + 1)
		}
		fmt.Printf("%v\n", a)
	}
	fmt.Printf("Product: %d\n", product)
}

func listInOrder(l, r []any) (bool, bool) {
	for i := 0; i < min(len(l), len(r)); i++ {
		left := l[i]
		right := r[i]
		leftType := reflect.TypeOf(left).Kind()
		rightType := reflect.TypeOf(right).Kind()
		// fmt.Printf("%T, %T\n", left, right)
		if leftType == reflect.Int && rightType == reflect.Int {
			leftValue := left.(int)
			rightValue := right.(int)
			if leftValue > rightValue {
				return false, true
			} else if leftValue < rightValue {
				return true, true
			}
		} else if leftType == reflect.Slice && rightType == reflect.Slice {
			// fmt.Printf("Both are slices: %T, %T\n", left, right)
			res, ok := listInOrder(left.([]any), right.([]any))
			if ok {
				return res, true
			}
		} else if leftType == reflect.Int && rightType == reflect.Slice {
			newLeft := []any{left}
			res, ok := listInOrder(newLeft, right.([]any))
			if ok {
				return res, true
			}
		} else if leftType == reflect.Slice && rightType == reflect.Int {
			newRight := []any{right}
			res, ok := listInOrder(left.([]any), newRight)
			if ok {
				return res, true
			}
		} else {
			fmt.Printf("Failed to check types %T and %T\n", left, right)
		}
	}
	if len(r) < len(l) {
		return false, true
	} else if len(l) < len(r) {
		return true, true
	}
	fmt.Printf("Confused: not sure why this didn't pass: %v, %v\n", l, r)
	return false, false
}

func parseString(str string) []any {
	str = str[1 : len(str)-1]
	// fmt.Printf("str: %q\n", str)
	list := []any{}
	i := 0
	for i < len(str) {
		// fmt.Printf("%d:\n", i)
		c := str[i]
		switch {
		case c == '[':
			start := i
			brackets := 1
			for brackets > 0 {
				// fmt.Printf("bracket %d\n", brackets)
				i++
				switch str[i] {
				case '[':
					brackets++
				case ']':
					brackets--
				}
			}
			i++
			newList := parseString(str[start:i])
			list = append(list, newList)
		case isInt(string(c)):
			start := i
			i++
			for i < len(str) && isInt(string(str[i])) {
				i++
			}
			s := str[start:i]
			v, err := strconv.Atoi(s)
			if err != nil {
				log.Fatalf("Failed to parse %q: %v", s, err)
			}
			list = append(list, v)
		}
		i++
	}
	return list
}

func isInt(s string) bool {
	_, err := strconv.Atoi(s)
	return err == nil
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
