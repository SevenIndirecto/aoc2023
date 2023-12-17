package main

import (
	"fmt"
	"io/ioutil"
	"strings"
)

const (
	empty = '.'
	round = 'O'
	cube  = '#'
)

func gridToString(g [][]rune) string {
	out := ""
	for y := range g {
		for _, c := range g[y] {
			out += string(c)
		}
		out += "\n"
	}
	return out
}

func loadGrid(lines []string) [][]rune {
	grid := make([][]rune, len(lines))
	for y, rowString := range lines {
		row := make([]rune, len(rowString))
		for x, c := range rowString {
			row[x] = c
		}
		grid[y] = row
	}
	return grid
}

func slideNorth(g [][]rune) {
	for x := range g[0] {
		for y := range g {
			if y == 0 || g[y][x] != round {
				continue
			}

			newY := 0
			for ny := y - 1; ny >= 0; ny-- {
				if g[ny][x] != empty {
					newY = ny + 1
					break
				}
			}
			g[y][x] = empty
			g[newY][x] = round
		}
	}
}

func slideSouth(g [][]rune) {
	for x := range g[0] {
		for y := len(g) - 2; y >= 0; y-- {
			if g[y][x] != round {
				continue
			}

			newY := len(g) - 1
			for ny := y + 1; ny <= len(g)-1; ny++ {
				if g[ny][x] != empty {
					newY = ny - 1
					break
				}
			}
			g[y][x] = empty
			g[newY][x] = round
		}
	}
}

func slideEast(g [][]rune) {
	for y := range g[0] {
		for x := len(g[y]) - 2; x >= 0; x-- {
			if g[y][x] != round {
				continue
			}

			newX := len(g[y]) - 1
			for nx := x + 1; nx <= len(g[y])-1; nx++ {
				if g[y][nx] != empty {
					newX = nx - 1
					break
				}
			}
			g[y][x] = empty
			g[y][newX] = round
		}
	}
}

func slideWest(g [][]rune) {
	for y := range g[0] {
		for x := range g {
			if x == 0 || g[y][x] != round {
				continue
			}

			newX := 0
			for nx := x - 1; nx >= 0; nx-- {
				if g[y][nx] != empty {
					newX = nx + 1
					break
				}
			}
			g[y][x] = empty
			g[y][newX] = round
		}
	}
}

func cycle(g [][]rune) {
	slideNorth(g)
	slideWest(g)
	slideSouth(g)
	slideEast(g)
}

func score(g [][]rune) int {
	sum := 0
	for y, row := range g {
		for _, c := range row {
			if c == round {
				sum += len(g) - y
			}
		}
	}
	return sum
}

func PartOne(lines []string) int {
	grid := loadGrid(lines)
	//fmt.Println(gridToString(grid))
	slideNorth(grid)
	//fmt.Println(gridToString(grid))
	return score(grid)
}

func PartTwo(lines []string) int {
	grid := loadGrid(lines)

	layoutCache := map[string]int{}
	targetCycles := 1000000000
	found := false
	for step := 1; step <= targetCycles; step++ {
		cycle(grid)
		key := gridToString(grid)
		prevStep, exists := layoutCache[key]

		if exists && !found {
			found = true
			offset := prevStep
			delta := step - prevStep
			requiredSteps := (targetCycles - offset) % delta
			step = targetCycles - requiredSteps
		} else {
			layoutCache[key] = step
		}
	}
	return score(grid)
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
