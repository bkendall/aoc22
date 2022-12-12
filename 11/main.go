package main

import (
	"fmt"
	"math"
	"sort"
)

type monkey struct {
	// Worry for each item being held.
	items []int

	// shows how your worry level changes as that monkey inspects an item.
	operationFn func(int) int

	// shows how the monkey uses your worry level to decide where to throw an item next.
	testFn func(int) int

	inspectionCount int
}

func main() {
	sampleMonkeys := []*monkey{
		{
			items:       []int{79, 98},
			operationFn: func(w int) int { return w * 19 },
			testFn: func(w int) int {
				if w%23 == 0 {
					return 2
				}
				return 3
			},
		},
		{
			items:       []int{54, 65, 75, 74},
			operationFn: func(w int) int { return w + 6 },
			testFn: func(w int) int {
				if w%19 == 0 {
					return 2
				}
				return 0
			},
		},
		{
			items:       []int{79, 60, 97},
			operationFn: func(w int) int { return w * w },
			testFn: func(w int) int {
				if w%13 == 0 {
					return 1
				}
				return 3
			},
		},
		{
			items:       []int{74},
			operationFn: func(w int) int { return w + 3 },
			testFn: func(w int) int {
				if w%17 == 0 {
					return 0
				}
				return 1
			},
		},
	}
	// sampleMod := 23 * 19 * 13 * 17

	realMonkeys := []*monkey{
		{ // 0
			items:       []int{64},
			operationFn: func(w int) int { return w * 7 },
			testFn: func(w int) int {
				if w%13 == 0 {
					return 1
				}
				return 3
			},
		},
		{ // 1
			items:       []int{60, 84, 84, 65},
			operationFn: func(w int) int { return w + 7 },
			testFn: func(w int) int {
				if w%19 == 0 {
					return 2
				}
				return 7
			},
		},
		{ // 2
			items:       []int{52, 67, 74, 88, 51, 61},
			operationFn: func(w int) int { return w * 3 },
			testFn: func(w int) int {
				if w%5 == 0 {
					return 5
				}
				return 7
			},
		},
		{ // 3
			items:       []int{67, 72},
			operationFn: func(w int) int { return w + 3 },
			testFn: func(w int) int {
				if w%2 == 0 {
					return 1
				}
				return 2
			},
		},
		{ // 4
			items:       []int{80, 79, 58, 77, 68, 74, 98, 64},
			operationFn: func(w int) int { return w * w },
			testFn: func(w int) int {
				if w%17 == 0 {
					return 6
				}
				return 0
			},
		},
		{ // 5
			items:       []int{62, 53, 61, 89, 86},
			operationFn: func(w int) int { return w + 8 },
			testFn: func(w int) int {
				if w%11 == 0 {
					return 4
				}
				return 6
			},
		},
		{ // 6
			items:       []int{86, 89, 82},
			operationFn: func(w int) int { return w + 2 },
			testFn: func(w int) int {
				if w%7 == 0 {
					return 3
				}
				return 0
			},
		},
		{ // 7
			items:       []int{92, 81, 70, 96, 69, 84, 83},
			operationFn: func(w int) int { return w + 4 },
			testFn: func(w int) int {
				if w%3 == 0 {
					return 4
				}
				return 5
			},
		},
	}
	realMod := 13 * 19 * 5 * 2 * 17 * 11 * 7 * 3

	// {
	// 	monkeys := sampleMonkeys
	// 	_ = realMonkeys
	// 	for i := 0; i < 20; i++ {
	// 		for mI, monkey := range monkeys {
	// 			fmt.Printf("Monkey %d:\n", mI)
	// 			items := append([]int{}, monkey.items...)
	// 			monkey.items = []int{}
	// 			for _, item := range items {
	// 				worry := item
	// 				fmt.Printf("\tMonkey inspects an item with a worry level of %d.\n", worry)
	// 				monkey.inspectionCount++
	// 				worry = monkey.operationFn(worry)
	// 				fmt.Printf("\t\tWorry level is... %d.\n", worry)
	// 				worry /= 3
	// 				fmt.Printf("\t\tMonkey gets bored with item. Worry level is divided by 3 to %d.\n", worry)
	// 				targetMonkey := monkey.testFn(worry)
	// 				fmt.Printf("\t\tMonkey runs the test...\n")
	// 				monkeys[targetMonkey].items = append(monkeys[targetMonkey].items, worry)
	// 				fmt.Printf("\t\tItem with worry level %d is thrown to monkey %d.\n", worry, targetMonkey)
	// 			}
	// 		}
	// 	}

	// 	h := maxHeap{maxN: 2}
	// 	for mI, monkey := range monkeys {
	// 		fmt.Printf("Monkey %d inspected items %d times.\n", mI, monkey.inspectionCount)
	// 		h.add(monkey.inspectionCount)
	// 	}

	// 	fmt.Printf("Answer: %d\n", h.values[0]*h.values[1])
	// }

	{
		monkeys := sampleMonkeys
		monkeys = realMonkeys
		for i := 0; i < 10_000; i++ {
			for _, monkey := range monkeys {
				// fmt.Printf("Monkey %d:\n", mI)
				items := append([]int{}, monkey.items...)
				monkey.items = []int{}
				for _, item := range items {
					worry := item
					// fmt.Printf("\tMonkey inspects an item with a worry level of %d.\n", worry)
					monkey.inspectionCount++
					worry = monkey.operationFn(worry)
					// fmt.Printf("\t\tWorry level is... %d.\n", worry)
					// worry %= sampleMod
					worry %= realMod
					// realMod
					// fmt.Printf("\t\tMonkey gets bored with item. Worry level is divided by 3 to %d.\n", worry)
					targetMonkey := monkey.testFn(worry)
					// fmt.Printf("\t\tMonkey runs the test...\n")
					monkeys[targetMonkey].items = append(monkeys[targetMonkey].items, worry)
					// fmt.Printf("\t\tItem with worry level %d is thrown to monkey %d.\n", worry, targetMonkey)
				}
			}
		}

		h := maxHeap{maxN: 2}
		for mI, monkey := range monkeys {
			fmt.Printf("Monkey %d inspected items %d times.\n", mI, monkey.inspectionCount)
			h.add(monkey.inspectionCount)
		}

		fmt.Printf("Answer: %d\n", h.values[0]*h.values[1])
	}

}

type maxHeap struct {
	maxN   int
	values []int
}

func (h *maxHeap) add(v int) {
	if len(h.values) >= h.maxN {
		if h.values[len(h.values)-1] < v {
			h.values[len(h.values)-1] = v
		}
	} else {
		h.values = append(h.values, v)
	}
	sort.Slice(h.values, func(i, j int) bool { return h.values[i] > h.values[j] })
}

func max(a, b int) int {
	return int(math.Max(float64(a), float64(b)))
}

func abs(a, b int) int {
	return int(math.Abs(float64(a - b)))
}
