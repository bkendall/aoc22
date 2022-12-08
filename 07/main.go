package main

import (
	"flag"
	"fmt"
	"log"
	"os"
	"strconv"
	"strings"
)

type dir struct {
	dirs  map[string]dir
	files map[string]int
}

type file struct {
	name string
	size int
}

func newDir() dir {
	return dir{
		dirs:  make(map[string]dir),
		files: make(map[string]int),
	}
}

type node struct {
	name   string
	files  []file
	dirs   []*node
	parent *node

	totalSize int
	sizeDone  bool
}

func main() {
	sourcePtr := flag.String("source", "./sample.txt", "input file")
	flag.Parse()

	buff, err := os.ReadFile(*sourcePtr)
	if err != nil {
		log.Fatalf("Error: %v", err)
	}
	input := strings.Split(string(buff), "\n")

	// d := []string{""}
	// dirs := newDir()
	// markder := &dirs

	// files := map[string]int{}
	// dirs := map[string][]file{}
	home := &node{name: "/"}
	marker := home

	{
		i := 0
		for {
			fmt.Printf("i: %d\n", i)
			if i >= len(input) {
				break
			}
			l := strings.TrimSpace(input[i])
			arr := strings.Split(l, " ")
			if arr[0] == "$" {
				switch arr[1] {
				case "cd":
					if arr[2] == ".." {
						if marker.parent != nil {
							marker = marker.parent
						}
					} else if arr[2] == "/" {
						// d = []string{""}
						for marker.parent != nil {
							marker = marker.parent
						}
					} else {
						// d = append(d, arr[2])
						// marker.dirs = append(marker.dirs, &node{name: arr[2]})
						// marker = marker.dirs[len(marker.dirs)-1]
						for _, d := range marker.dirs {
							if d.name == arr[2] {
								marker = d
								break
							}
						}
					}
					fmt.Printf("dir: %v\n", marker)
				case "ls":
					fls := []file{}
					for {
						fmt.Printf("ls: i: %d\n", i)
						i++
						if i >= len(input) {
							break
						}
						ll := strings.TrimSpace(input[i])
						if strings.HasPrefix(ll, "$") {
							i--
							break // We're done.
						} else if strings.HasPrefix(ll, "dir") {
							aa := strings.Split(ll, " ")
							marker.dirs = append(marker.dirs, &node{name: aa[1], parent: marker})
							// Not important, maybe.
						} else {
							aa := strings.Split(ll, " ")
							size, err := strconv.Atoi(aa[0])
							if err != nil {
								log.Fatalf("Failed to parse size %q: %v", aa[0], err)
							}
							// name := strings.Join(append(d, aa[1]), "/")
							// files[name] = size
							fls = append(fls, file{name: aa[1], size: size})
						}
						// log.Fatal("STOP")
					}
					marker.files = append(marker.files, fls...)
					// dirName := strings.Join(d, "/")
					// if dirName == "" {
					// 	dirName = "/"
					// }
					// dirs[dirName] = fls
				}
			}
			i++
		}

		fmt.Printf("Files:\n")
		calculateSize(home)
		printNode(home, 0)

		lessThan100k := []*node{}
		forEveryNode(home, func(nn *node) {
			if nn.totalSize <= 100_000 {
				lessThan100k = append(lessThan100k, nn)
			}
		})
		// fmt.Printf("lessThan100k: %v\n", lessThan100k)
		sum := 0
		for _, nn := range lessThan100k {
			fmt.Printf("%s: %d\n", nn.name, nn.totalSize)
			sum += nn.totalSize
		}
		fmt.Printf("Total Size: %d\n", sum)

		// sizes := map[string]int{}
		// for d, fs := range dirs {
		// 	sum := 0
		// 	for _, f := range fs {
		// 		sum += f.size
		// 	}
		// 	fmt.Printf("%s: %d\n", d, sum)
		// }
	}

	{
		totalSize := 70000000
		totalUsed := home.totalSize
		totalNeeded := 30000000
		toDelete := totalNeeded - (totalSize - totalUsed)
		fmt.Printf("toDelete: %d\n", toDelete)

		closest := home.totalSize
		forEveryNode(home, func(nn *node) {
			if nn.totalSize > toDelete && nn.totalSize < closest {
				closest = nn.totalSize
			}
		})
		fmt.Printf("closest: %d\n", closest)
	}
}

func forEveryNode(n *node, fn func(nn *node)) {
	fn(n)
	for _, d := range n.dirs {
		forEveryNode(d, fn)
	}
}

func calculateSize(n *node) int {
	if n.sizeDone {
		return n.totalSize
	}
	dirSizes := 0
	for _, d := range n.dirs {
		dirSizes += calculateSize(d)
	}
	fileSizes := 0
	for _, f := range n.files {
		fileSizes += f.size
	}
	n.totalSize = dirSizes + fileSizes
	n.sizeDone = true
	return n.totalSize
}

func printNode(n *node, depth int) {
	fmt.Printf("%s- %s (%d)\n", strings.Repeat("\t", depth), n.name, n.totalSize)
	for _, d := range n.dirs {
		printNode(d, depth+1)
	}
	for _, f := range n.files {
		fmt.Printf("%s- %s (%d)\n", strings.Repeat("\t", depth+1), f.name, f.size)
	}
}
