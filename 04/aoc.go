package main

import (
	"fmt"
	"io/ioutil"
	"math"
	"strconv"
	"strings"
)

type card struct {
	id        int
	count     int
	winners   map[int]bool
	myNumbers []int
}

func loadCards(lines []string) []card {
	cards := make([]card, len(lines))

	for i, l := range lines {
		s := strings.Split(l, ": ")
		s = strings.Split(s[1], " | ")
		ws := strings.Split(s[0], " ")
		ms := strings.Split(s[1], " ")

		winners := map[int]bool{}
		mine := make([]int, len(ms))
		for _, w := range ws {
			if w == "" {
				continue
			}
			wnum, _ := strconv.Atoi(w)
			winners[wnum] = true
		}
		for j, m := range ms {
			if m == "" {
				continue
			}
			mnum, _ := strconv.Atoi(m)
			mine[j] = mnum
		}

		cards[i] = card{id: i + 1, winners: winners, myNumbers: mine, count: 1}
	}

	return cards
}

func PartOne(lines []string) int {
	cards := loadCards(lines)

	sum := 0
	for _, c := range cards {
		hits := 0
		for _, mine := range c.myNumbers {
			if c.winners[mine] {
				hits++
			}
		}
		if hits > 0 {
			sum += int(math.Pow(2, float64(hits)-1))
		}
	}
	return sum
}

func PartTwo(lines []string) int {
	cards := loadCards(lines)

	sum := 0
	for i, c := range cards {
		sum += c.count
		hits := 0
		for _, mine := range c.myNumbers {
			if c.winners[mine] {
				hits++
			}
		}

		if hits > 0 {
			for delta := 1; delta <= hits; delta++ {
				cards[i+delta].count += c.count
			}
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
