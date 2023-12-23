package main

import (
	"fmt"
	"io/ioutil"
	"strings"
)

type point struct {
	x, y int
}

type multiGridPoint struct {
	x, y   int // in-grid coordinates
	gx, gy int // Which grid are we on?
}

func loadGrid(lines []string) ([][]rune, point) {
	g := make([][]rune, len(lines))

	var start point

	for y, l := range lines {
		row := make([]rune, len(l))
		for x, c := range l {
			if c == 'S' {
				start = point{x: x, y: y}
				row[x] = '.'
			} else {
				row[x] = c
			}
		}
		g[y] = row
	}

	return g, start
}

func find(start point, steps int, grid [][]rune) map[point]bool {
	pointsToCheckForActiveStep := map[point]bool{start: true}
	directions := [4]point{{0, -1}, {1, 0}, {0, 1}, {-1, 0}}

	for i := 0; i < steps; i++ {
		newPointsToCheck := map[point]bool{}
		for p := range pointsToCheckForActiveStep {
			for _, d := range directions {
				np := point{p.x + d.x, p.y + d.y}
				if np.x < 0 || np.y < 0 || np.x >= len(grid[0]) || np.y >= len(grid) || grid[np.y][np.x] == '#' {
					continue
				}
				newPointsToCheck[np] = true
			}
		}
		pointsToCheckForActiveStep = newPointsToCheck
	}

	return pointsToCheckForActiveStep
}

func draw(grid [][]rune, currentStepPoints map[point]bool) {
	for y, row := range grid {
		for x, c := range row {
			if currentStepPoints[point{x, y}] {
				fmt.Print("O")
			} else {
				fmt.Print(string(c))
			}
		}
		fmt.Println()
	}
	fmt.Println()
}

func PartOne(lines []string, steps int) int {
	grid, start := loadGrid(lines)
	currentStepPoints := find(start, steps, grid)
	draw(grid, currentStepPoints)

	return len(currentStepPoints)
}

func findMultiGrid(start point, steps int, grid [][]rune) int {
	pointsToCheckForActiveStep := map[multiGridPoint]bool{multiGridPoint{x: start.x, y: start.y, gx: 0, gy: 0}: true}
	directions := [4]point{{0, -1}, {1, 0}, {0, 1}, {-1, 0}}

	gridWidth := len(grid[0])
	gridHeight := len(grid)
	combinationsAtStep := map[int]int{}
	offset := gridWidth / 2

	// We learned from Reddit...... that we're dealing with points that lie on a quadratic equation (ax^2 + bx + c = y)
	// 1) Our shape expands as a diamond on each step, so we need 65 (gridWidth / 2 rounded down) steps to make it from
	// the center to the edge, which is our first X on the quadratic curve.
	// 2) Then we need to find two more points, which occur every next 131 steps as the diamond expands and reaches another edge(gridWith)
	//
	// So we need to calculate at least 3 points to determine the a,b and c of the quadratic equation.
	numStepsWeNeedToSimulate := offset + 2*gridWidth + 1 // (+1 for good measure :D)

	for i := 0; i < numStepsWeNeedToSimulate; i++ {
		newPointsToCheck := map[multiGridPoint]bool{}
		for p := range pointsToCheckForActiveStep {
			for _, d := range directions {
				np := multiGridPoint{x: p.x + d.x, y: p.y + d.y, gx: p.gx, gy: p.gy}

				// Jump to other grids if needed
				if np.x < 0 {
					// New grid to te left
					np.gx -= 1
					np.x = gridWidth - 1
				} else if np.x >= gridWidth {
					np.gx += 1
					np.x = 0
				} else if np.y < 0 {
					np.gy -= 1
					np.y = gridHeight - 1
				} else if np.y >= gridHeight {
					np.gy += 1
					np.y = 0
				}

				if grid[np.y][np.x] == '.' {
					newPointsToCheck[np] = true
				}
			}
		}
		pointsToCheckForActiveStep = newPointsToCheck
		combinationsAtStep[i] = len(newPointsToCheck)
	}

	y1 := combinationsAtStep[offset-1]
	y2 := combinationsAtStep[offset+gridWidth-1]
	y3 := combinationsAtStep[offset+gridWidth*2-1]
	fmt.Println("y1", y1, "y2", y2, "y3", y3)
	// Now we know how many viable spots (y) are at each of the three points (x) that we're interested in (x1=0, x2=1 and x3=2)
	// This is derived on paper, standard stuff and input here
	a := (y3 - 2*y2 + y1) / 2
	b := y2 - y1 - a
	c := y1

	xStep := (steps - offset) / gridWidth
	return a*xStep*xStep + b*xStep + c
}

func PartTwo(lines []string, steps int) int {
	// NOTE: Annoyingly the solution doesn't work on the test case... but those numbers are small enough to simulate.
	grid, start := loadGrid(lines)
	return findMultiGrid(start, steps, grid)
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
	fmt.Printf("Part one %v\n", PartOne(lines, 64))
	fmt.Printf("Part two %v\n", PartTwo(lines, 26501365))
}
