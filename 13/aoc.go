package main

import (
	"fmt"
	"io/ioutil"
	"strconv"
	"strings"
)

const (
	vertical = iota
	horizontal
)

type pattern struct {
	grid         []string
	variant      int
	reflectionAt int
	cleanRefAt   int
	cleanVariant int
	smudge       [2]int
}

func findGridReflection(grid []string, reflectionToIgnore int) int {
	for y := 1; y < len(grid); y++ {
		if y == reflectionToIgnore {
			continue
		}
		if grid[y] == grid[y-1] {
			// See if all lines match until we reach an edge
			foundReflection := true
			for dy := 1; y+dy < len(grid) && y-dy-1 >= 0; dy++ {
				if grid[y+dy] != grid[y-dy-1] {
					foundReflection = false
				}
			}

			if foundReflection {
				return y
			}
		}
	}
	return -1
}

func flipGrid(g []string) []string {
	flipped := make([]string, len(g[0]))

	for i := range g[0] {
		flipped[i] = ""
	}

	for y := len(g) - 1; y >= 0; y-- {
		row := g[y]

		for i, c := range row {
			flipped[i] += string(c)
		}
	}
	return flipped
}

func (p *pattern) String() string {
	out := "\n------------\n"
	out += "Reflection at index: " + strconv.Itoa(p.reflectionAt)
	if p.variant == vertical {
		out += " VERTICAL"
	} else {
		out += " HORIZONTAL"
	}

	for _, r := range p.grid {
		out += "\n" + r
	}
	out += fmt.Sprint("\nSmudge: ", p.smudge, " New reflection at ", p.cleanRefAt)
	out += "\n------------"

	return out
}

func (p *pattern) determineReflection() {
	// Try to get a horizontal match
	horizontalReflectionAt := findGridReflection(p.grid, -1)
	if horizontalReflectionAt > -1 {
		p.reflectionAt = horizontalReflectionAt
		p.variant = horizontal
		return
	}

	verticalReflectionAt := findGridReflection(flipGrid(p.grid), -1)
	if verticalReflectionAt == -1 {
		fmt.Println("Could not find reflection for")
		fmt.Println(p)
		p.reflectionAt = 0
		p.variant = -1
		return
	}
	p.reflectionAt = verticalReflectionAt
	p.variant = vertical
}

func (p *pattern) getValue() int {
	if p.variant == vertical {
		return p.reflectionAt
	}
	return p.reflectionAt * 100
}

func (p *pattern) getSmudgeFreeValue() int {
	if p.cleanVariant == vertical {
		return p.cleanRefAt
	}
	return p.cleanRefAt * 100
}

func loadPatterns(lines []string) []pattern {
	patterns := make([]pattern, 0)

	g := make([]string, 0)

	for i, r := range lines {
		if r == "" || i == len(lines)-1 {
			patterns = append(patterns, pattern{grid: g})
			g = make([]string, 0)
		} else {
			g = append(g, r)
		}
	}

	return patterns
}

func (p *pattern) fixSmudge() {
	for y, row := range p.grid {
		for x, c := range row {
			newChar := "."
			if c == '.' {
				newChar = "#"
			}

			newRow := row[:x] + newChar + row[x+1:]
			newGrid := make([]string, len(p.grid))
			copy(newGrid, p.grid)
			newGrid[y] = newRow

			reflectionToIgnore := -1
			newReflectionVariant := horizontal
			if p.variant == newReflectionVariant {
				reflectionToIgnore = p.reflectionAt
			}
			reflectionAt := findGridReflection(newGrid, reflectionToIgnore)

			if reflectionAt == -1 || (reflectionAt == p.reflectionAt && newReflectionVariant == p.variant) {
				newReflectionVariant = vertical
				reflectionToIgnore = -1
				if p.variant == newReflectionVariant {
					reflectionToIgnore = p.reflectionAt
				}
				reflectionAt = findGridReflection(flipGrid(newGrid), reflectionToIgnore)
			}

			if reflectionAt > -1 && (reflectionAt != p.reflectionAt || newReflectionVariant != p.variant) {
				p.smudge = [2]int{x, y}
				p.cleanRefAt = reflectionAt
				p.cleanVariant = newReflectionVariant
				return
			}
		}
	}
	fmt.Println(p.String())
	panic("Could not find reflection")
}

func PartOne(lines []string) int {
	patterns := loadPatterns(lines)

	sum := 0
	for _, p := range patterns {
		p.determineReflection()
		sum += p.getValue()
	}
	return sum
}

func PartTwo(lines []string) int {
	patterns := loadPatterns(lines)

	sum := 0
	for _, p := range patterns {
		p.determineReflection()
		p.fixSmudge()
		sum += p.getSmudgeFreeValue()
	}
	return sum
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
