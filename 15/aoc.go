package main

import (
	"fmt"
	"io/ioutil"
	"strconv"
	"strings"
)

func loadSteps(line string) []string {
	return strings.Split(line, ",")
}

func hash(s string) int {
	cv := 0
	for _, c := range s {
		cv += int(c)
		cv *= 17
		cv %= 256
	}
	return cv
}

func PartOne(lines []string) int {
	steps := loadSteps(lines[0])
	sum := 0
	for _, s := range steps {
		h := hash(s)
		sum += h
	}
	return sum
}

type lense struct {
	label string
	focus int
}

const (
	remove = iota
	insert
)

func PartTwo(lines []string) int {
	boxes := map[int][]lense{}
	steps := loadSteps(lines[0])
	for _, s := range steps {
		var label string
		var op int
		var focus int
		if strings.Contains(s, "-") {
			label = s[0 : len(s)-1]
			op = remove
		} else {
			label = s[0 : len(s)-2]
			op = insert
			f, _ := strconv.Atoi(s[len(s)-1:])
			focus = f
		}

		boxId := hash(label)
		if op == remove {
			boxLenses, exists := boxes[boxId]
			if !exists {
				continue
			}

			newLenses := make([]lense, 0)
			for _, l := range boxLenses {
				if l.label != label {
					newLenses = append(newLenses, l)
				}
			}
			boxes[boxId] = newLenses
		} else {
			// If already exists, replace
			alreadyExists := false
			for i, l := range boxes[boxId] {
				if l.label == label {
					boxes[boxId][i].focus = focus
					alreadyExists = true
					break
				}
			}
			if alreadyExists {
				continue
			}

			// Otherwise insert at the end
			boxes[boxId] = append(boxes[boxId], lense{label: label, focus: focus})
		}
	}

	score := 0
	for boxNumber, lenses := range boxes {
		for i, l := range lenses {
			mul := boxNumber + 1
			mul *= (i + 1) * l.focus
			score += mul
		}
	}
	return score
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
