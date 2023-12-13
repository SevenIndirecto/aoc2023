package main

import (
	"fmt"
	"io/ioutil"
	"math"
	"strings"
)

type point struct {
	x, y int
}

const (
	empty = iota
	galaxyType
)

func loadGalaxies(lines []string, isPartTwo bool) []point {
	universe := make([][]int, len(lines))
	galaxies := make([]point, 0)

	rowsWithGalaxies := map[int]bool{}
	colsWithGalaxies := map[int]bool{}

	for y, lineRow := range lines {
		row := make([]int, len(lineRow))

		for x, c := range lineRow {
			t := empty
			if c == '#' {
				t = galaxyType
				galaxies = append(galaxies, point{x: x, y: y})
				rowsWithGalaxies[y] = true
				colsWithGalaxies[x] = true
			}
			row[x] = t
		}
		universe[y] = row
	}

	bump := 1
	if isPartTwo {
		bump = 1000*1000 - 1
	}

	// Expand universe
	offset := 0
	for x, _ := range universe[0] {
		if !colsWithGalaxies[x] {
			// Column doesn't have galaxy, expand
			bumpColGalaxies(&galaxies, offset, x, bump)
			offset += bump
		}
	}

	offset = 0
	for y, _ := range universe {
		if !rowsWithGalaxies[y] {
			// Row  doesn't have galaxy, expand
			bumpRowGalaxies(&galaxies, offset, y, bump)
			offset += bump
		}
	}

	return galaxies
}

func bumpColGalaxies(galaxies *[]point, offset int, threshold int, bump int) {
	for i, _ := range *galaxies {
		if (*galaxies)[i].x < threshold+offset {
			continue
		}
		(*galaxies)[i].x += bump
	}
}

// Look... a duplicate is fine too :D
func bumpRowGalaxies(galaxies *[]point, offset int, threshold int, bump int) {
	for i, _ := range *galaxies {
		if (*galaxies)[i].y < threshold+offset {
			continue
		}
		(*galaxies)[i].y += bump
	}
}

func drawUniverse(galaxies []point) {
	mapOfGalaxies := map[point]bool{}
	for _, g := range galaxies {
		mapOfGalaxies[g] = true
	}

	// determine bounds
	maxX := 0
	maxY := 0
	for _, g := range galaxies {
		if g.x > maxX {
			maxX = g.x
		}
		if g.y > maxY {
			maxY = g.y
		}
	}

	for y := 0; y <= maxY; y++ {
		row := ""
		for x := 0; x <= maxX; x++ {
			if mapOfGalaxies[point{x: x, y: y}] {
				row += "#"
			} else {
				row += "."
			}
		}
		fmt.Println(row)
	}
}

func manhattanDistance(a, b point) int {
	return int(math.Abs(float64(a.x-b.x)) + math.Abs(float64(a.y-b.y)))
}

type shortestDistance struct {
	pointA, pointB point
	distance       int
}

func calcShortestDistancesSum(galaxies []point) int {
	shortestDistances := make([]shortestDistance, 0)
	sum := 0
	for i := 0; i < len(galaxies)-1; i++ {
		for j := i + 1; j < len(galaxies); j++ {
			d := manhattanDistance(galaxies[i], galaxies[j])
			shortestDistances = append(shortestDistances, shortestDistance{
				pointA:   galaxies[i],
				pointB:   galaxies[j],
				distance: d,
			})
			sum += d
		}
	}
	return sum
}

func PartOne(lines []string) int {
	galaxies := loadGalaxies(lines, false)
	drawUniverse(galaxies)
	return calcShortestDistancesSum(galaxies)
}

func PartTwo(lines []string) int {
	galaxies := loadGalaxies(lines, true)
	return calcShortestDistancesSum(galaxies)
}

func LoadLines(path string) ([]string, error) {
	dat, err := ioutil.ReadFile(path)
	if err != nil {
		return nil, err
	}
	txt := string(dat)
	lines := strings.Split(txt, "\n")
	return lines[:len(lines)-1], nil
}

func main() {
	lines, _ := LoadLines("input.txt")
	fmt.Printf("Part one %v\n", PartOne(lines))
	fmt.Printf("Part two %v\n", PartTwo(lines))
}
