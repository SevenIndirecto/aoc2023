package main

import (
	"fmt"
	"io/ioutil"
	"strings"
)

type point [2]int

type visit struct {
	p       point
	steps   int
	history map[point]bool
}

func loadGrid(lines []string) [][]rune {
	grid := make([][]rune, len(lines))

	for y, l := range lines {
		row := make([]rune, len(l))
		for x, c := range l {
			row[x] = c
		}
		grid[y] = row
	}
	return grid
}

const (
	x = iota
	y
)

func copyHistory(m map[point]bool) map[point]bool {
	newMap := map[point]bool{}
	for k, v := range m {
		newMap[k] = v
	}
	return newMap
}

var directions = [4]point{{0, -1}, {-1, 0}, {0, 1}, {1, 0}}

// For part two-only, so ignore slopes
// NOTE: Takes several minutes, as I should have collapsed everything to a graph, with junctions for vertices... but ok.
func dfsRecurse(p, end point, grid [][]rune, visited map[point]bool, steps int, maxSteps int) int {
	if p == end {
		if steps > maxSteps {
			return steps
		}
		return maxSteps
	}

	visited[p] = true
	// Find all possible neighbors
	for _, d := range directions {
		np := point{p[x] + d[x], p[y] + d[y]}
		if np[x] < 0 || np[y] < 0 || np[x] >= len(grid[0]) || np[y] >= len(grid) || grid[np[y]][np[x]] == '#' {
			continue // avoid off-grid and forest
		}
		if !visited[np] {
			maxSteps = dfsRecurse(np, end, grid, visited, steps+1, maxSteps)
		}
	}
	visited[p] = false
	return maxSteps
}

func dfsLongestHike(start point, end point, grid [][]rune, ignoreSlopes bool) int {
	stack := []visit{{p: start, steps: 0, history: map[point]bool{start: true}}}
	longestPaths := map[point]int{}
	slopes := map[rune]point{'>': {1, 0}, 'v': {0, 1}, '<': {-1, 0}, '^': {0, -1}}
	mostSteps := 0

	for len(stack) > 0 {
		// pop
		cv := stack[len(stack)-1]
		stack = stack[0 : len(stack)-1]

		for _, d := range directions {
			nv := visit{p: point{cv.p[x] + d[x], cv.p[y] + d[y]}, steps: cv.steps + 1}
			if nv.p[x] < 0 || nv.p[y] < 0 || nv.p[x] >= len(grid[0]) || nv.p[y] >= len(grid) {
				continue // avoid off-grid
			}
			nv.history = copyHistory(cv.history)

			if !ignoreSlopes {
				nextType := grid[nv.p[y]][nv.p[x]]
				slopeDelta, exists := slopes[nextType]
				if exists {
					// Next visit is a slope, so automatically move in the mandatory position and increase step count
					// NOTE: Assume slope doesn't lead off-grid or to a forest.
					nv.history[nv.p] = true
					nv.p[x] += slopeDelta[x]
					nv.p[y] += slopeDelta[y]
					nv.steps++
				}
			}

			if longestPaths[nv.p] >= nv.steps || // avoid shorter hikes
				cv.history[nv.p] || // avoid loops
				grid[nv.p[y]][nv.p[x]] == '#' { // avoid forest
				continue
			}

			// Valid visit
			nv.history[nv.p] = true
			longestPaths[nv.p] = nv.steps

			if nv.p != end {
				stack = append(stack, nv)
			} else {
				if nv.steps > mostSteps {
					mostSteps = nv.steps
				}
			}
		}
	}

	return mostSteps
}

func PartOne(lines []string) int {
	grid := loadGrid(lines)
	start := point{1, 0}
	end := point{len(grid[0]) - 2, len(grid) - 1}
	steps := dfsLongestHike(start, end, grid, false)
	return steps
}

func PartTwo(lines []string) int {
	grid := loadGrid(lines)
	start := point{1, 0}
	end := point{len(grid[0]) - 2, len(grid) - 1}

	visited := map[point]bool{}
	maxSteps := dfsRecurse(start, end, grid, visited, 0, 0)
	return maxSteps
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
	// 4858 too low
	fmt.Printf("Part two %v\n", PartTwo(lines))
}
