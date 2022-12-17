package main

import (
	"flag"
	"fmt"
	"log"
	"math"
	"os"
	"regexp"
	"strconv"
	"strings"
)

type Material int

const (
	Nothing Material = iota
	Sensor
	Beacon
)

type Spot struct {
	x, y          int
	thing         Material
	closestBeacon *Spot
}

func (s Spot) DistanceFromBeacon() int {
	return s.DistanceFrom(s.closestBeacon.x, s.closestBeacon.y)
}

func (s Spot) DistanceFrom(x, y int) int {
	return abs(s.x-x) + abs(s.y-y)
}

func (s Spot) IsInSensorZone(x, y int) bool {
	// fmt.Printf("SensorZone check: %d, %d from %d, %d\n", x, y, s.closestBeacon.x, s.closestBeacon.y)
	// fmt.Printf("SensorZone distance from sensor: %d\n", s.DistanceFrom(x, y))
	// fmt.Printf("SensorZone distance from beacon: %d\n", s.DistanceFromBeacon())
	return s.DistanceFrom(x, y) <= s.DistanceFromBeacon()
}

func (s Spot) FurthestRightFrom(x, y int) int {
	// fmt.Printf("SensorZone check: %d, %d from %d, %d\n", x, y, s.closestBeacon.x, s.closestBeacon.y)
	// fmt.Printf("SensorZone distance from sensor: %d\n", s.DistanceFrom(x, y))
	// fmt.Printf("SensorZone distance from beacon: %d\n", s.DistanceFromBeacon())
	totalDistance := s.DistanceFromBeacon()
	yUsed := abs(y - totalDistance)
	xDist := abs(totalDistance - yUsed)
	return x + xDist
}

type World struct {
	m map[int]map[int]*Spot
}

func (w *World) Add(x, y int, s *Spot) {
	if _, ok := w.m[x]; !ok {
		w.m[x] = map[int]*Spot{}
	}
	w.m[x][y] = s
}

func (w *World) Get(x, y int) *Spot {
	if _, ok := w.m[x]; ok {
		return w.m[x][y]
	}
	return nil
}

func main() {
	sourcePtr := flag.String("source", "./sample.txt", "input file")
	flag.Parse()

	buff, err := os.ReadFile(*sourcePtr)
	if err != nil {
		log.Fatalf("Error: %v", err)
	}
	input := strings.Split(string(buff), "\n")

	world := World{map[int]map[int]*Spot{}}
	lineRegex := regexp.MustCompile(`Sensor at x=(-?\d+), y=(-?\d+): closest beacon is at x=(-?\d+), y=(-?\d+)`)

	beacons := []*Spot{}
	sensors := []*Spot{}

	for _, line := range input {
		arr := lineRegex.FindStringSubmatch(line)
		if arr == nil {
			log.Fatalf("Could not match string: %q\n", line)
		}
		beacon := &Spot{x: toInt(arr[3]), y: toInt(arr[4]), thing: Beacon}
		sensor := &Spot{x: toInt(arr[1]), y: toInt(arr[2]), thing: Sensor, closestBeacon: beacon}
		world.Add(beacon.x, beacon.y, beacon)
		world.Add(sensor.x, sensor.y, sensor)
		beacons = append(beacons, beacon)
		sensors = append(sensors, sensor)
	}

	{
		// y := 10
		y := 2000000
		deadSpots := checkRowForCoverage(y, -10_000_000, 10_000_000, sensors, world, false)
		fmt.Printf("Dead spots, not a beacon: %d\n", len(deadSpots))
	}

	{
		// Part 2. Thank you, reddit.
		// size := 20
		size := 4_000_000

		dist := func(x1, y1, x2, y2 int) int {
			return abs(x2-x1) + abs(y2-y1)
		}

		aCoeffs, bCoeffs := map[int]any{}, map[int]any{}
		for _, sensor := range sensors {
			x, y, r := sensor.x, sensor.y, dist(sensor.x, sensor.y, sensor.closestBeacon.x, sensor.closestBeacon.y)
			aCoeffs[y-x+r+1] = true
			aCoeffs[y-x-r-1] = true
			bCoeffs[x+y+r+1] = true
			bCoeffs[x+y-r-1] = true
		}

		// Goal x=14, y=11
	loop:
		for a := range aCoeffs {
			for b := range bCoeffs {
				pX, pY := (b-a)/2, (a+b)/2
				if 0 < pX && pX < size && 0 < pY && pY < size {
					c := 0
					for _, sensor := range sensors {
						if dist(pX, pY, sensor.x, sensor.y) > sensor.DistanceFromBeacon() {
							c++
						}
					}
					if c == len(sensors) {
						fmt.Printf("Found it. %d, %d: %d\n", pX, pY, 4_000_000*pX+pY)
						break loop
					}
				}
			}
		}
	}
}

func checkRowForCoverage(y, xMin, xMax int, sensors []*Spot, world World, includeBeacons bool) []*Spot {
	arr := []*Spot{}
	for x := xMin; x <= xMax; x++ {
		for _, s := range sensors {
			if checkSpot := world.Get(x, y); !includeBeacons && checkSpot != nil && checkSpot.thing == Beacon {
				continue
			}
			if s.IsInSensorZone(x, y) {
				arr = append(arr, s)
				break
			}
		}
	}
	return arr
}

func toInt(s string) int {
	v, err := strconv.Atoi(s)
	if err != nil {
		log.Fatalf("failed to parse %q: %v\n", s, err)
	}
	return v
}

func max(a, b int) int {
	return int(math.Max(float64(a), float64(b)))
}

func min(a, b int) int {
	return int(math.Min(float64(a), float64(b)))
}

func abs(a int) int {
	return int(math.Abs(float64(a)))
}

func floor(a int) int {
	return int(math.Floor(float64(a)))
}
