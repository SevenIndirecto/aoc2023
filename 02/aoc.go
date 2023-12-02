package main

import (
	"fmt"
	"io/ioutil"
	"strconv"
	"strings"
)

type game struct {
	id   int
	sets []set
}

type set struct {
	r int
	g int
	b int
}

func (g *game) isPossible(p *params) bool {
	for _, s := range g.sets {
		if s.r > p.r || s.b > p.b || s.g > p.g {
			return false
		}
	}
	return true
}

func (g *game) calculatePower() int {
	minRed := 0
	minBlue := 0
	minGreen := 0

	for _, s := range g.sets {
		if s.r > minRed {
			minRed = s.r
		}
		if s.g > minGreen {
			minGreen = s.g
		}
		if s.b > minBlue {
			minBlue = s.b
		}
	}

	return minRed * minBlue * minGreen
}

type params struct {
	r int
	g int
	b int
}

func loadGame(gameLine string) game {
	s := strings.Split(gameLine, ": ")
	id, _ := strconv.Atoi(s[0][5:])

	setsAsStrings := strings.Split(s[1], "; ")
	sets := make([]set, len(setsAsStrings))

	for i, gameSet := range setsAsStrings {
		red := 0
		blue := 0
		green := 0

		cubes := strings.Split(gameSet, ", ")
		for _, cube := range cubes {
			numberColor := strings.Split(cube, " ")
			n, _ := strconv.Atoi(numberColor[0])

			switch numberColor[1] {
			case "red":
				red = n
				break
			case "blue":
				blue = n
				break
			case "green":
				green = n
				break
			default:
				panic(n)
			}
		}
		sets[i] = set{r: red, g: green, b: blue}
	}

	return game{id: id, sets: sets}
}

func PartOne(lines []string) int {
	p := params{r: 12, g: 13, b: 14}
	sum := 0

	for _, line := range lines {
		g := loadGame(line)
		if g.isPossible(&p) {
			sum += g.id
		}
	}

	return sum
}

func PartTwo(lines []string) int {
	sum := 0

	for _, line := range lines {
		g := loadGame(line)
		sum += g.calculatePower()
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
