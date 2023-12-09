package main

import (
	"fmt"
	"io/ioutil"
	"strconv"
	"strings"
)

func calcNext(history []int, isPartTwo bool) int {
	subs := make([][]int, 0)

	current := history
	subs = append(subs, current)

	// Build required sub sequences
	for {
		allZero := true
		for _, n := range current {
			if n != 0 {
				allZero = false
			}
		}
		if allZero {
			current = append(current, 0)
			break
		}

		sub := makeSubSequence(current)
		subs = append(subs, sub)
		current = sub
	}

	if isPartTwo {
		appendedValues := map[int]int{len(current) - 1: 0}
		for i := len(subs) - 2; i >= 0; i-- {
			appendedValues[i] = subs[i][0] - appendedValues[i+1]
		}
		return appendedValues[0]
	}

	// Part One
	// Fill up from next to last sub-sequence
	for i := len(subs) - 2; i >= 0; i-- {
		next := subs[i+1][len(subs[i+1])-1] + subs[i][len(subs[i])-1]
		subs[i] = append(subs[i], next)
	}

	return subs[0][len(subs[0])-1]
}

func makeSubSequence(history []int) []int {
	sub := make([]int, 0)
	for i := 0; i < len(history)-1; i++ {
		sub = append(sub, history[i+1]-history[i])
	}
	return sub
}

func loadHistories(lines []string) [][]int {
	histories := make([][]int, len(lines))
	for i, line := range lines {
		history := make([]int, 0)
		s := strings.Split(line, " ")
		for _, numAsStr := range s {
			n, _ := strconv.Atoi(numAsStr)
			history = append(history, n)
		}
		histories[i] = history
	}

	return histories
}

func PartOne(lines []string) int {
	histories := loadHistories(lines)
	sum := 0
	for _, history := range histories {
		sum += calcNext(history, false)
	}
	return sum
}

func PartTwo(lines []string) int {
	histories := loadHistories(lines)
	sum := 0
	for _, history := range histories {
		sum += calcNext(history, true)
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
