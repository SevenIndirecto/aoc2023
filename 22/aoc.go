package main

import (
	"fmt"
	"io/ioutil"
	"slices"
	"strconv"
	"strings"
)

type brick struct {
	// Start always has a Z lower than end
	start, end  [3]int
	axis        int
	pointsSlice [][3]int
	settled     bool
	id          int
}

const (
	x = iota
	y
	z
)

func makeBrick(start [3]int, end [3]int, id int) brick {
	b := brick{start: start, end: end, id: id}
	pointsSlice := make([][3]int, 0)
	if start == end {
		pointsSlice = append(pointsSlice, start)
	} else {
		for i := 0; i < 3; i++ {
			if b.start[i] != b.end[i] {
				for j := b.start[i]; j <= b.end[i]; j++ {
					p := start
					p[i] = j
					pointsSlice = append(pointsSlice, p)
				}
			}
		}
	}
	b.pointsSlice = pointsSlice
	return b
}

func loadBricks(lines []string) []brick {
	bricks := make([]brick, 0)

	for id, l := range lines {
		s := strings.Split(l, "~")
		startSplit := strings.Split(s[0], ",")
		start := [3]int{}
		for i, vs := range startSplit {
			v, _ := strconv.Atoi(vs)
			start[i] = v
		}
		endSplit := strings.Split(s[1], ",")
		end := [3]int{}
		for i, vs := range endSplit {
			v, _ := strconv.Atoi(vs)
			end[i] = v
		}
		if start[z] > end[z] {
			tmp := start
			start = end
			end = tmp
		}
		bricks = append(bricks, makeBrick(start, end, id+1))
	}

	return bricks
}

func settleTryTwo(bricks []brick, failFast bool) int {
	// Sort bricks by starting Z asc
	slices.SortFunc(bricks, func(a, b brick) int {
		return a.start[z] - b.start[z]
	})
	numUpdates := 0

	// Create map of all occupied cubes and link them to brick IDs for easier lookup
	cubeMap := map[[3]int]int{}
	for _, b := range bricks {
		for _, p := range b.pointsSlice {
			cubeMap[p] = b.id
		}
	}

	// Since bricks are ordered by start Z, move each down as far as we can and update cubeMap if needed
	for i, b := range bricks {
		newPoints := b.pointsSlice
		shouldUpdate := false
		for {
			collided := false
			newPointsCandidate := make([][3]int, 0)

			for _, p := range newPoints {
				// Lower all points by 1
				np := p
				np[z]--

				if np[z] < 1 || (cubeMap[np] > 0 && cubeMap[np] != b.id) {
					// Floor | Collision with another brick -> we're settled
					collided = true
					break
				}
				newPointsCandidate = append(newPointsCandidate, np)
			}

			if collided {
				break
			} else {
				newPoints = newPointsCandidate
				shouldUpdate = true
			}
		}

		if shouldUpdate {
			numUpdates++
			if failFast {
				return 1
			}
			// Clear old points
			for _, p := range bricks[i].pointsSlice {
				delete(cubeMap, p)
			}
			bricks[i] = makeBrick(newPoints[0], newPoints[len(newPoints)-1], bricks[i].id)
			for _, p := range newPoints {
				cubeMap[p] = bricks[i].id
			}
		}
	}
	return numUpdates
}

func disintegrate(brickToRemove int, bricks []brick) []brick {
	newBricks := make([]brick, 0)
	for _, b := range bricks {
		if b.id == brickToRemove {
			continue
		}
		newBricks = append(newBricks, makeBrick(b.start, b.end, b.id))
	}
	return newBricks
}

func PartOne(lines []string) int {
	bricks := loadBricks(lines)
	settleTryTwo(bricks, false)

	// Remove brick 1 by 1 by marking it as disintegrated and see if anything is not unsettled
	bricksThatCanBeDisintegrated := make([]int, 0)
	for _, b := range bricks {
		newMoves := settleTryTwo(disintegrate(b.id, bricks), true)
		if newMoves == 0 {
			bricksThatCanBeDisintegrated = append(bricksThatCanBeDisintegrated, b.id)
		}
	}
	return len(bricksThatCanBeDisintegrated)
}

func PartTwo(lines []string) int {
	bricks := loadBricks(lines)
	settleTryTwo(bricks, false)

	sum := 0
	for _, b := range bricks {
		sum += settleTryTwo(disintegrate(b.id, bricks), false)
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
