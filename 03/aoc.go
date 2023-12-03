package main

import (
	"fmt"
	"io/ioutil"
	"strconv"
	"strings"
	"unicode"
)

type variant int

const (
	empty variant = iota
	digit
	symbol
	gear
)

type number struct {
	number int
	x      int
	y      int
	len    int
}

func loadSchematic(lines []string) ([][]variant, []number) {
	rows := make([][]variant, len(lines))
	points := make([]number, 0)
	lineLength := len(lines[0])

	numBeingConstructed := ""
	pointBeingConstructed := number{}

	for y, line := range lines {
		row := make([]variant, lineLength)

		for x, c := range line {
			finishCurrentNumber := false
			if c == '.' {
				row[x] = empty
				finishCurrentNumber = true
			} else if unicode.IsDigit(c) {
				row[x] = digit
				if numBeingConstructed == "" {
					pointBeingConstructed = number{x: x, y: y}
				}
				numBeingConstructed += string(c)
			} else {
				if c == '*' {
					row[x] = gear
				} else {
					row[x] = symbol
				}
				finishCurrentNumber = true
			}

			if (x == len(line)-1 || finishCurrentNumber) && numBeingConstructed != "" {
				n, _ := strconv.Atoi(numBeingConstructed)
				pointBeingConstructed.number = n
				pointBeingConstructed.len = len(numBeingConstructed)
				numBeingConstructed = ""
				points = append(points, pointBeingConstructed)
			}
		}
		rows[y] = row
	}

	return rows, points
}

type coordinate struct {
	x int
	y int
}

func getCoordinatesToCheck(schematic [][]variant, n number) []coordinate {
	coordsToCheck := make([]coordinate, 2)
	coordsToCheck[0] = coordinate{x: n.x - 1, y: n.y}
	coordsToCheck[1] = coordinate{x: n.x + n.len, y: n.y}

	for x := n.x - 1; x <= n.x+n.len; x++ {
		coordsToCheck = append(coordsToCheck, coordinate{x: x, y: n.y - 1}, coordinate{x: x, y: n.y + 1})
	}

	filtered := make([]coordinate, 0)
	for _, c := range coordsToCheck {
		if c.y < 0 || c.x < 0 || c.y >= len(schematic) || c.x >= len(schematic[0]) {
			continue
		}
		filtered = append(filtered, c)
	}

	return filtered
}

func PartOne(lines []string) int {
	schematic, numbers := loadSchematic(lines)
	sum := 0

	for _, n := range numbers {
		coordsToCheck := getCoordinatesToCheck(schematic, n)

		for _, c := range coordsToCheck {
			if schematic[c.y][c.x] == symbol || schematic[c.y][c.x] == gear {
				sum += n.number
				break
			}
		}
	}

	return sum
}

func PartTwo(lines []string) int {
	schematic, numbers := loadSchematic(lines)
	gears := make(map[int][]int)

	for _, n := range numbers {
		coordsToCheck := getCoordinatesToCheck(schematic, n)

		for _, c := range coordsToCheck {
			if schematic[c.y][c.x] == gear {
				key := 1000*c.y + c.x
				g, exists := gears[key]
				if !exists {
					gears[key] = make([]int, 0)
				}
				gears[key] = append(g, n.number)
				break
			}
		}
	}

	sum := 0
	for _, g := range gears {
		if len(g) == 2 {
			sum += g[0] * g[1]
		}
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
