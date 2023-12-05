package main

import (
	"fmt"
	"io/ioutil"
	"strconv"
	"strings"
)

type conversionMap struct {
	dst      string
	mappings []mapping
}

func (cm *conversionMap) srcToDst(input int) int {
	for _, m := range cm.mappings {
		if m.src <= input && input <= m.src+m.r-1 {
			delta := input - m.src
			return m.dst + delta
		}
	}

	return input
}

type mapping struct {
	dst int
	src int
	r   int
}

func loadConversionMapsAndSeeds(lines []string) (map[string]conversionMap, []int, [][]int) {
	seeds := make([]int, 0)
	seedsAndRanges := make([][]int, 0)
	cms := map[string]conversionMap{}
	currentSrc := ""

	for i, line := range lines {
		if len(line) == 0 {
			continue
		}

		if i == 0 {
			s := strings.Split(line, "seeds: ")
			nums := strings.Split(s[1], " ")
			inputAsInts := make([]int, 0)

			for _, n := range nums {
				seed, _ := strconv.Atoi(n)
				inputAsInts = append(inputAsInts, seed)
				seeds = append(seeds, seed)
			}

			for i := 0; i < len(inputAsInts); i += 2 {
				seedsAndRanges = append(seedsAndRanges, []int{inputAsInts[i], inputAsInts[i+1]})
			}
			continue
		}

		s := strings.Split(line, " map")
		if len(s) > 1 {
			// Start new mapping
			s = strings.Split(s[0], "-to-")
			currentSrc = s[0]
			cms[currentSrc] = conversionMap{dst: s[1], mappings: make([]mapping, 0)}
			continue
		}

		params := strings.Split(line, " ")
		dst, _ := strconv.Atoi(params[0])
		src, _ := strconv.Atoi(params[1])
		r, _ := strconv.Atoi(params[2])

		mappings := append(cms[currentSrc].mappings, mapping{dst: dst, src: src, r: r})
		cms[currentSrc] = conversionMap{dst: cms[currentSrc].dst, mappings: mappings}
	}

	return cms, seeds, seedsAndRanges
}

func PartOne(lines []string) int {
	cms, seeds, _ := loadConversionMapsAndSeeds(lines)

	bestResult := -1
	for _, seed := range seeds {
		srcType := "seed"
		runningConversion := seed

		for {
			cm := cms[srcType]
			runningConversion = cm.srcToDst(runningConversion)
			srcType = cm.dst

			if cm.dst == "location" {
				if bestResult == -1 || runningConversion < bestResult {
					bestResult = runningConversion
				}
				break
			}
		}
	}

	return bestResult
}

func PartTwo(lines []string) int {
	cms, _, seedsAndRanges := loadConversionMapsAndSeeds(lines)

	bestResult := -1
	// TODO: Not really a solution lol, just a brute force that works...
	for _, sr := range seedsAndRanges {
		for i := sr[0]; i < sr[0]+sr[1]; i++ {
			srcType := "seed"
			runningConversion := i

			for {
				cm := cms[srcType]
				runningConversion = cm.srcToDst(runningConversion)
				srcType = cm.dst

				if cm.dst == "location" {
					if bestResult == -1 || runningConversion < bestResult {
						bestResult = runningConversion
					}
					break
				}
			}
		}
	}

	return bestResult
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
