package main

import (
	"fmt"
	"io/ioutil"
	"math"
	"regexp"
	"strconv"
	"strings"
)

type race struct {
	time     int
	distance int
}

func loadRaces(lines []string) []race {
	races := make([]race, 0)
	r := regexp.MustCompile(`\w+`)
	times := r.FindAllStringSubmatch(lines[0], -1)
	distances := r.FindAllStringSubmatch(lines[1], -1)
	for i := 1; i < len(times); i++ {
		t, _ := strconv.Atoi(times[i][0])
		d, _ := strconv.Atoi(distances[i][0])
		races = append(races, race{time: t, distance: d})
	}
	return races
}

func PartOne(lines []string) int {
	races := loadRaces(lines)
	result := 1

	for _, race := range races {
		minTime := int(math.Ceil(float64(race.distance / race.time)))
		hits := 0

		for t := minTime; t < race.time; t++ {
			v := t
			d := (race.time - t) * v
			if d > race.distance {
				hits++
			}
		}
		result *= hits
	}

	return result
}

func PartTwo(lines []string) int {
	races := loadRaces(lines)

	timeStr := ""
	distanceStr := ""
	for _, race := range races {
		timeStr += strconv.Itoa(race.time)
		distanceStr += strconv.Itoa(race.distance)
	}

	time, _ := strconv.Atoi(timeStr)
	distance, _ := strconv.Atoi(distanceStr)

	minTime := int(math.Ceil(float64(distance / time)))
	hits := 0

	for t := minTime; t < time; t++ {
		v := t
		d := (time - t) * v
		if d > distance {
			hits++
		}
	}

	return hits
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
