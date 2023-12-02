package main

import (
	"fmt"
	"io/ioutil"
	"strconv"
	"strings"
	"unicode"
)

func PartOne(lines []string) int {
	sum := 0
	for _, line := range lines {
		first := ""
		last := ""

		for _, c := range line {
			if unicode.IsDigit(c) {
				if first == "" {
					first = string(c)
				}
				last = string(c)
			}
		}

		if last == "" {
			continue
		}

		num, _ := strconv.Atoi(first + last)
		sum += num
	}

	return sum
}

func PartTwo(lines []string) int {
	sum := 0
	needles := [9]string{"one", "two", "three", "four", "five", "six", "seven", "eight", "nine"}

	for _, line := range lines {
		first := ""
		last := ""
		firstIndex := len(line)
		lastIndex := -1

		for i, digitAsString := range needles {
			position := strings.Index(line, digitAsString)
			lastPosition := strings.LastIndex(line, digitAsString)

			if position < 0 {
				continue
			}

			if position < firstIndex {
				firstIndex = position
				first = strconv.Itoa(i + 1)
			}

			if lastPosition > lastIndex {
				lastIndex = lastPosition
				last = strconv.Itoa(i + 1)
			}
		}

		for i, c := range line {
			if !unicode.IsDigit(c) {
				continue
			}

			if i < firstIndex {
				firstIndex = i
				first = string(c)
			}

			if i > lastIndex {
				lastIndex = i
				last = string(c)
			}
		}

		if last == "" {
			continue
		}

		num, _ := strconv.Atoi(first + last)
		fmt.Println(line, first+last, num)
		sum += num
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
