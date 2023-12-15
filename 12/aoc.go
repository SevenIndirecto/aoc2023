package main

import (
	"fmt"
	"io/ioutil"
	"strconv"
	"strings"
)

func loadLine(line string, expand bool) (string, []int) {
	s := strings.Split(line, " ")

	record := s[0]
	groupsString := s[1]

	if expand {
		for i := 0; i < 4; i++ {
			record += "?" + s[0]
			groupsString += "," + s[1]
		}
	}

	groups := make([]int, 0)
	groupsSplit := strings.Split(groupsString, ",")
	for _, c := range groupsSplit {
		num, _ := strconv.Atoi(c)
		groups = append(groups, num)
	}
	return record, groups
}

var cache = make(map[string]int)

func handleDot(record string, groups []int) int {
	// Skip dot and look for next group
	return getArrangementCount(record[1:], groups)
}

func handlePound(record string, groups []int) int {
	// First n characters must be damaged, n = group[0]
	currGroupSize := groups[0]
	if len(record) < currGroupSize {
		return 0
	}
	currGroup := strings.ReplaceAll(record[0:currGroupSize], "?", "#")

	if strings.Contains(currGroup, ".") || len(currGroup) < currGroupSize {
		// Can't fit current group in remaining record -> invalid
		return 0
	}

	// Rest of the record is just the last group
	if len(record) == currGroupSize {
		// This is the last group -> ok
		if len(groups) == 1 {
			return 1
		}
		// We have more groups than we can fit -> invalid
		return 0
	}

	// Next character must not be # so we can have a valid group
	if record[currGroupSize] == '#' {
		return 0
	}

	// Next character can be a separator so skip (+1) and reduce groups
	return getArrangementCount(record[currGroupSize+1:], groups[1:])
}

func getArrangementCount(record string, groups []int) int {
	if len(groups) < 1 {
		if strings.Contains(record, "#") {
			// More damaged springs than groups allow -> invalid
			return 0
		}
		return 1
	}

	if len(record) == 0 {
		// More groups, but ran out of record -> invalid
		return 0
	}

	cacheKey := record
	for _, n := range groups {
		cacheKey += "," + strconv.Itoa(n)
	}

	val, exists := cache[cacheKey]
	if exists {
		//fmt.Println("Cache hit", cacheKey, record, groups, "val", val)
		return val
	}

	nextChar := record[0]

	count := 0
	if nextChar == '.' {
		count = handleDot(record, groups)
	} else if nextChar == '#' {
		count = handlePound(record, groups)
	} else if nextChar == '?' {
		count = handleDot(record, groups) + handlePound(record, groups)
	} else {
		panic("Invalid char")
	}

	cache[cacheKey] = count
	return count
}

func PartOne(lines []string) int {
	// Reset cache
	cache = map[string]int{}

	sum := 0
	for _, l := range lines {
		record, groups := loadLine(l, false)
		sum += getArrangementCount(record, groups)
	}
	return sum
}

func PartTwo(lines []string) int {
	// Reset cache
	cache = map[string]int{}
	sum := 0
	for _, l := range lines {
		record, groups := loadLine(l, true)
		sum += getArrangementCount(record, groups)
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
