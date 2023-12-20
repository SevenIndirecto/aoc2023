package main

import (
	"fmt"
	"io/ioutil"
	"strings"
)

type point struct {
	x, y int
}

type beam struct {
	loc       point
	direction int
}

const (
	up = iota
	right
	down
	left
)

func loadGrid(lines []string) [][]rune {
	grid := make([][]rune, len(lines))

	for y, line := range lines {
		row := make([]rune, len(line))
		for x, c := range line {
			row[x] = c
		}
		grid[y] = row
	}
	return grid
}

// If true returned, split beam
func (b *beam) move(grid [][]rune) bool {
	// Determine new loc. Tile type doesn't matter.
	if b.direction == up {
		b.loc.y--
	} else if b.direction == right {
		b.loc.x++
	} else if b.direction == down {
		b.loc.y++
	} else {
		b.loc.x--
	}

	if b.loc.x < 0 || b.loc.x >= len(grid[0]) || b.loc.y < 0 || b.loc.y >= len(grid) {
		// off-grid
		return false
	}

	tileType := grid[b.loc.y][b.loc.x]
	if tileType == '.' {
		return false
	} else if tileType == '/' {
		if b.direction == right {
			b.direction = up
		} else if b.direction == down {
			b.direction = left
		} else if b.direction == left {
			b.direction = down
		} else {
			b.direction = right
		}
		return false
	} else if tileType == '\\' {
		if b.direction == down {
			b.direction = right
		} else if b.direction == left {
			b.direction = up
		} else if b.direction == up {
			b.direction = left
		} else {
			b.direction = down
		}
		return false
	} else if tileType == '-' {
		if b.direction == right || b.direction == left {
			return false
		}
		b.direction = left // Other beam will be right
		return true

	} else if tileType == '|' {
		if b.direction == up || b.direction == down {
			return false
		}
		b.direction = up // Other beam will be down
		return true
	}

	panic("Invalid tile type")
}

func process(startBeam beam, grid [][]rune) int {
	beams := []beam{startBeam}
	beamVisits := map[beam]bool{}
	energizedPoints := map[point]bool{}

	for len(beams) > 0 {
		newBeams := make([]beam, 0)

		for i := range beams {
			shouldSplit := beams[i].move(grid)

			if shouldSplit {
				newDirection := down
				if beams[i].direction == left {
					newDirection = right
				}
				newSplitBeam := beam{loc: beams[i].loc, direction: newDirection}
				if !beamVisits[newSplitBeam] {
					newBeams = append(newBeams, newSplitBeam)
				}
			}
		}

		// After moving all the beams evaluate if any need to be deactivated
		for _, b := range beams {
			if b.loc.x < 0 || b.loc.x >= len(grid[0]) || b.loc.y < 0 || b.loc.y >= len(grid) {
				// 1. Deactivate any off grid beams
				continue
			}
			energizedPoints[b.loc] = true

			// 2. deactivate any beams with direction / point combo that has already been visited
			if beamVisits[b] {
				continue
			}

			newBeams = append(newBeams, b)
			beamVisits[b] = true
		}
		beams = newBeams
	}

	return len(energizedPoints)
}

func PartOne(lines []string) int {
	grid := loadGrid(lines)
	startBeam := beam{loc: point{x: -1, y: 0}, direction: right}
	return process(startBeam, grid)
}

func PartTwo(lines []string) int {
	grid := loadGrid(lines)
	maxEnergized := 0

	for x := range grid[0] {
		// Top row
		energized := process(beam{loc: point{x: x, y: -1}, direction: down}, grid)
		if energized > maxEnergized {
			maxEnergized = energized
		}
		// Bottom row
		energized = process(beam{loc: point{x: x, y: len(grid)}, direction: up}, grid)
		if energized > maxEnergized {
			maxEnergized = energized
		}
	}

	for y := range grid {
		// Left-most column
		energized := process(beam{loc: point{x: -1, y: y}, direction: right}, grid)
		if energized > maxEnergized {
			maxEnergized = energized
		}
		// Right-most column
		energized = process(beam{loc: point{x: len(grid[0]), y: y}, direction: left}, grid)
		if energized > maxEnergized {
			maxEnergized = energized
		}
	}

	return maxEnergized
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
