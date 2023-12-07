package main

import (
	"fmt"
	"io/ioutil"
	"sort"
	"strconv"
	"strings"
)

const (
	highCard int = iota
	onePair
	twoPair
	threeOfAKind
	fullHouse
	fourOfAKind
	fiveOfAKind
)

type card int

type game struct {
	asString string
	hand     [5]card
	handMap  map[card]int
	bid      int
	handType int
}

func determineHandType(handMap map[card]int) int {
	switch len(handMap) {
	case 1:
		return fiveOfAKind
	case 2:
		for _, count := range handMap {
			if count == 4 || count == 1 {
				return fourOfAKind
			} else {
				return fullHouse
			}
		}
	case 3:
		for _, count := range handMap {
			if count == 2 {
				return twoPair
			} else if count == 3 {
				return threeOfAKind
			}
		}
	case 4:
		return onePair
	case 5:
		return highCard
	}
	fmt.Println(handMap)
	panic("Unsupported hand type")
}

const (
	wildcard card = 1
)

func loadGames(lines []string, isPartTwo bool) []game {
	games := make([]game, len(lines))
	powerMap := map[rune]card{
		'A': 14,
		'K': 13,
		'Q': 12,
		'J': 11,
		'T': 10,
		'9': 9,
		'8': 8,
		'7': 7,
		'6': 6,
		'5': 5,
		'4': 4,
		'3': 3,
		'2': 2,
	}

	for i, line := range lines {
		s := strings.Split(line, " ")
		suits := s[0]
		bid, _ := strconv.Atoi(s[1])

		hand := [5]card{}
		handMap := map[card]int{}

		for j, c := range suits {
			mappedCard := powerMap[c]

			if isPartTwo && c == 'J' {
				mappedCard = wildcard
			}

			hand[j] = mappedCard
			handMap[mappedCard]++
		}

		if isPartTwo {
			if handMap[wildcard] > 0 && handMap[wildcard] != 5 {
				// Merge wildcards into highest count
				var highestCountCard card
				highestCount := 0
				for c, count := range handMap {
					if c == wildcard {
						continue
					}
					if count > highestCount {
						highestCountCard = c
						highestCount = count
					}
				}
				handMap[highestCountCard] = highestCount + handMap[wildcard]
				delete(handMap, wildcard)
			}
		}

		games[i] = game{
			asString: suits,
			hand:     hand,
			bid:      bid,
			handType: determineHandType(handMap),
			handMap:  handMap,
		}
	}

	return games
}

func determineWinnings(games []game) int {
	sort.Slice(games, func(i, j int) bool {
		g1 := games[i]
		g2 := games[j]

		// Fallback
		if g1.handType == g2.handType {
			for z, _ := range g1.hand {
				if g1.hand[z] != g2.hand[z] {
					return g1.hand[z] < g2.hand[z]
				}
			}
			fmt.Println(g1, g2)
			panic("Could not determine stable ordering")
		}

		return g1.handType < g2.handType
	})

	winnings := 0
	for rank, game := range games {
		winnings += (rank + 1) * game.bid
	}

	return winnings
}

func PartOne(lines []string) int {
	games := loadGames(lines, false)
	return determineWinnings(games)
}

func PartTwo(lines []string) int {
	games := loadGames(lines, true)
	return determineWinnings(games)
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
