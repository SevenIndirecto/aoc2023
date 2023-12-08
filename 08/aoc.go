package main

import (
	"fmt"
	"io/ioutil"
	"strings"
)

const (
	left = iota
	right
)

type node struct {
	id         string
	directions [2]string
}

func loadNodesAndDirections(lines []string) (map[string]node, []int) {
	instructions := make([]int, len(lines[0]))
	for i, c := range lines[0] {
		if c == 'L' {
			instructions[i] = left
		} else {
			instructions[i] = right
		}
	}

	nodes := make(map[string]node)
	for _, line := range lines[2:] {
		s := strings.Split(line, " = ")
		n := node{id: s[0], directions: [2]string{"", ""}}

		s = strings.Split(s[1][1:len(s[1])-1], ", ")
		n.directions[left] = s[0]
		n.directions[right] = s[1]
		nodes[n.id] = n
	}

	return nodes, instructions
}

func PartOne(lines []string) int {
	nodes, instructions := loadNodesAndDirections(lines)
	target := "ZZZ"
	current := "AAA"
	ip := 0
	steps := 1

	for {
		current = nodes[current].directions[instructions[ip]]

		if current == target {
			return steps
		}

		ip = (ip + 1) % len(instructions)
		steps++
	}
}

func PartTwo(lines []string) int {
	nodes, instructions := loadNodesAndDirections(lines)

	currentNodes := make([]string, 0)
	for id, _ := range nodes {
		if id[2] == 'A' {
			currentNodes = append(currentNodes, id)
		}
	}

	ip := 0
	steps := 1

	gapBetweenZ := make([]map[int]bool, len(currentNodes))
	for i, _ := range gapBetweenZ {
		gapBetweenZ[i] = map[int]bool{}
	}

	stepsSinceLastZ := make([]int, len(currentNodes))

	for {
		nextNodes := make([]string, len(currentNodes))
		allEndWithZ := true

		for i, id := range currentNodes {
			nextNodes[i] = nodes[id].directions[instructions[ip]]
			if nextNodes[i][2] != 'Z' {
				allEndWithZ = false
				stepsSinceLastZ[i]++
			} else {
				gapBetweenZ[i][stepsSinceLastZ[i]+1] = true
				stepsSinceLastZ[i] = 0
			}
		}

		repeats := make([]int, 0)
		for i, _ := range gapBetweenZ {
			for k, _ := range gapBetweenZ[i] {
				repeats = append(repeats, k)
			}
		}

		if len(repeats) == len(currentNodes) {
			return LCM(repeats[0], repeats[1], repeats[2:]...)
		}

		currentNodes = nextNodes

		if allEndWithZ {
			return steps
		}

		ip = (ip + 1) % len(instructions)
		steps++
	}
}

// GCD greatest common divisor (GCD) via Euclidean algorithm
func GCD(a, b int) int {
	for b != 0 {
		t := b
		b = a % b
		a = t
	}
	return a
}

// LCM find Least Common Multiple (LCM) via GCD
func LCM(a, b int, integers ...int) int {
	result := a * b / GCD(a, b)

	for i := 0; i < len(integers); i++ {
		result = LCM(result, integers[i])
	}

	return result
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
