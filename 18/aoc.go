package main

import (
	"fmt"
	"io/ioutil"
	"math"
	"strconv"
	"strings"
)

const (
	empty = iota
	lagoon
)

const (
	up = iota
	right
	down
	left
)

type point struct {
	x, y int
}

type digCommand struct {
	direction int
	distance  int
}

func executeDigCommandForFloodFill(command digCommand, from point, grid [][]int, virtualDig bool) ([][]int, point) {
	deltas := map[int]point{
		up:    {x: 0, y: -1},
		right: {x: 1, y: 0},
		down:  {x: 0, y: 1},
		left:  {x: -1, y: 0},
	}
	delta := deltas[command.direction]
	to := point{x: from.x + command.distance*delta.x, y: from.y + command.distance*delta.y}

	if virtualDig {
		return grid, to
	}

	// Otherwise paint map
	current := from
	for {
		current.x += delta.x
		current.y += delta.y
		grid[current.y][current.x] = lagoon

		if current == to {
			break
		}
	}

	return grid, to
}

func determineGridSizeAndStart(commands []digCommand) (int, int, point) {
	grid := make([][]int, 0)
	from := point{0, 0}

	minPoint := point{x: 0, y: 0}
	maxPoint := point{x: 0, y: 0}
	for _, cmd := range commands {
		grid, from = executeDigCommandForFloodFill(cmd, from, grid, true)
		if from.x < minPoint.x {
			minPoint.x = from.x
		}
		if from.y < minPoint.y {
			minPoint.y = from.y
		}
		if from.y > maxPoint.y {
			maxPoint.y = from.y
		}
		if from.x > maxPoint.x {
			maxPoint.x = from.x
		}
	}

	// Determine starting point offset
	from = point{x: 0 - minPoint.x, y: 0 - minPoint.y}
	width := int(math.Abs(float64(minPoint.x-maxPoint.x))) + 1
	height := int(math.Abs(float64(minPoint.y-maxPoint.y))) + 1

	return width, height, from
}

func digTrenchForFloodFill(commands []digCommand) [][]int {
	width, height, from := determineGridSizeAndStart(commands)

	// build empty grid with starting point
	grid := make([][]int, height)
	for y, _ := range grid {
		row := make([]int, width)
		for x := 0; x < width; x++ {
			row[x] = empty
		}
		grid[y] = row
	}
	grid[from.y][from.x] = lagoon

	for _, cmd := range commands {
		grid, from = executeDigCommandForFloodFill(cmd, from, grid, false)
	}

	return grid
}

func printGrid(grid [][]int) {
	cmap := map[int]string{empty: ".", lagoon: "#"}

	for _, row := range grid {
		for _, t := range row {
			fmt.Print(cmap[t])
		}
		fmt.Println()
	}
	fmt.Println()
}

func loadCommands(lines []string) []digCommand {
	directionMap := map[string]int{
		"R": right,
		"D": down,
		"L": left,
		"U": up,
	}

	commands := make([]digCommand, 0)
	for _, line := range lines {
		s := strings.Split(line, " ")
		distance, _ := strconv.Atoi(s[1])

		commands = append(commands, digCommand{direction: directionMap[s[0]], distance: distance})
	}
	return commands
}

// NOTE: This is not working ok...
func isInside(p point, grid [][]int, excludeEdges bool) bool {
	if grid[p.y][p.x] == lagoon {
		// This is just a helper for finding a starting point for floodFill... since can't get this to work yet
		if excludeEdges {
			return false
		}
		return true
	}

	if p.y == 0 || p.y == len(grid)-1 {
		return false
	}

	intersectCount := 0
	requireNewEmpty := false

	for x := p.x + 1; x < len(grid[p.y]); x++ {
		if grid[p.y][x] == lagoon && !requireNewEmpty && x != len(grid[p.y])-1 {
			requireNewEmpty = true
		}

		if grid[p.y][x] == empty && requireNewEmpty {
			intersectCount++
			requireNewEmpty = false
		}
	}

	// Ray trace polygon check, if intersects odd times point is inside the polygon
	return intersectCount%2 == 1
}

func floodFill(grid [][]int) [][]int {
	// Find first point inside
	found := false
	start := point{x: 0, y: 0}
	for y, row := range grid {
		for x, _ := range row {
			if isInside(point{x: x, y: y}, grid, true) {
				start.x = x
				start.y = y
				found = true
				break
			}
		}
		if found {
			break
		}
	}

	deltas := []point{{x: 0, y: -1}, {x: 1, y: 0}, {x: 0, y: 1}, {x: -1, y: 0}}
	queue := []point{start}

	for len(queue) > 0 {
		// pop
		p := queue[len(queue)-1]
		queue = queue[0 : len(queue)-1]
		// Paint
		grid[p.y][p.x] = lagoon

		// Add all empty adjacent
		for _, delta := range deltas {
			np := point{x: p.x + delta.x, y: p.y + delta.y}
			if np.x < 0 || np.y < 0 || np.x >= len(grid[0]) || np.y >= len(grid) {
				continue
			}

			if grid[np.y][np.x] == empty {
				queue = append(queue, np)
			}
		}
	}
	return grid
}

func getLagoonCount(grid [][]int) int {
	count := 0
	for _, row := range grid {
		for _, c := range row {
			if c == lagoon {
				count++
			}
		}
	}
	return count
}

func PartOne(lines []string) int {
	commands := loadCommands(lines)
	grid := digTrenchForFloodFill(commands)
	grid = floodFill(grid)
	return getLagoonCount(grid)
}

// Implemented for part 2
func loadCommandsPartTwo(lines []string) []digCommand {
	directionMap := map[string]int{
		"0": right,
		"1": down,
		"2": left,
		"3": up,
	}

	commands := make([]digCommand, 0)
	for _, line := range lines {
		s := strings.Split(line, " ")
		distanceHexStr := s[2][2:7]
		distance, _ := strconv.ParseInt(distanceHexStr, 16, 32)
		direction := directionMap[s[2][len(s[2])-2:len(s[2])-1]]

		commands = append(commands, digCommand{direction: direction, distance: int(distance)})
	}
	return commands
}

func digTrenchPolygon(commands []digCommand) ([]point, int) {
	from := point{x: 0, y: 0}
	vertices := []point{from}
	var newEdgeLength int
	edgeLengthSum := 0

	for _, cmd := range commands {
		from, newEdgeLength = executeDigCommandPolygon(cmd, from)
		edgeLengthSum += newEdgeLength
		vertices = append(vertices, from)
	}

	return vertices, edgeLengthSum
}

func executeDigCommandPolygon(command digCommand, from point) (point, int) {
	deltas := map[int]point{
		up:    {x: 0, y: -1},
		right: {x: 1, y: 0},
		down:  {x: 0, y: 1},
		left:  {x: -1, y: 0},
	}
	delta := deltas[command.direction]
	to := point{x: from.x + command.distance*delta.x, y: from.y + command.distance*delta.y}
	edgeLength := int(math.Abs(float64(from.x-to.x)) + math.Abs(float64(from.y-to.y)))
	return to, edgeLength
}

func areaByShoelace(vertices []point) int {
	sum := 0
	p0 := vertices[len(vertices)-1]
	for _, p1 := range vertices {
		sum += p0.y*p1.x - p0.x*p1.y
		p0 = p1
	}
	return int(math.Abs(float64(sum) / 2))
}

func PartTwo(lines []string) int {
	commands := loadCommandsPartTwo(lines)
	vertices, edgeLengthSum := digTrenchPolygon(commands)
	area := areaByShoelace(vertices)
	// Couldn't figure out we need to add half the count of trench tiles
	// to get the correct number. Had to look it up... still don't quite get it.
	total := area + edgeLengthSum/2 + 1
	return total
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
