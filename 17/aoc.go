package main

import (
	"fmt"
	"io/ioutil"
	"strconv"
	"strings"
)

const (
	up = iota
	right
	down
	left
)

type point struct {
	x, y int
}

type visit struct {
	from        int
	straightFor int
	p           point
}

func pop(queue []visit, cost map[visit]int) ([]visit, visit) {
	bestIndex := 0
	shortestDistance := cost[queue[0]]

	for i, v := range queue {
		if cost[v] < shortestDistance {
			bestIndex = i
			shortestDistance = cost[v]
		}
	}

	bestVisit := queue[bestIndex]
	return append(queue[0:bestIndex], queue[bestIndex+1:]...), bestVisit
}

func findLowestCost(to point, grid [][]int, useUltraCrucibles bool) int {
	start := visit{from: -1, straightFor: 0, p: point{x: 0, y: 0}}
	r := visit{from: left, straightFor: 1, p: point{x: 1, y: 0}}
	d := visit{from: up, straightFor: 1, p: point{x: 0, y: 1}}
	costs := map[visit]int{start: 0, r: grid[r.p.y][r.p.x], d: grid[r.p.y][r.p.x]}
	queue := []visit{r, d}

	deltas := map[int]point{
		up:    {x: 0, y: -1},
		right: {x: 1, y: 0},
		down:  {x: 0, y: 1},
		left:  {x: -1, y: 0},
	}

	var v visit
	for len(queue) > 0 {
		queue, v = pop(queue, costs)

		// For any not off-grid move left, right or forward (if 3 or less forward travels in sequence)
		for i := 1; i <= 3; i++ {
			newDirection := (v.from + i) % 4
			delta := deltas[newDirection]
			newPoint := point{x: v.p.x + delta.x, y: v.p.y + delta.y}
			// Off-grid check
			if newPoint.x < 0 || newPoint.y < 0 || newPoint.x >= len(grid[0]) || newPoint.y >= len(grid) {
				continue
			}

			newStraightFor := 1
			if useUltraCrucibles {
				if i == 2 {
					newStraightFor = v.straightFor + 1
				}
				if v.straightFor < 4 && i != 2 || newStraightFor > 10 {
					// Can't turn until we moved at least 4 straight
					// And can't move more than 10 straight
					continue
				}
			} else {
				// Make sure we only ever do max 3 forward movements for normal crucibles
				if i == 2 {
					newStraightFor = v.straightFor + 1
					if newStraightFor > 3 {
						continue
					}
				}
			}

			// Add new visit to queue if cost smaller or not already visited
			from := (newDirection + 2) % 4
			potentialVisit := visit{from: from, straightFor: newStraightFor, p: newPoint}
			newCost := costs[v] + grid[newPoint.y][newPoint.x]
			oldCost, alreadyVisited := costs[potentialVisit]

			if !alreadyVisited || newCost < oldCost {
				costs[potentialVisit] = newCost
				queue = append(queue, potentialVisit)
			}
		}
	}

	// Get all
	lowestCost := -1
	for vk, cost := range costs {
		if vk.p != to || useUltraCrucibles && vk.straightFor < 4 {
			continue
		}
		if lowestCost == -1 || cost < lowestCost {
			lowestCost = cost
		}
	}

	return lowestCost
}

func loadGrid(lines []string) [][]int {
	grid := make([][]int, len(lines))
	for y, l := range lines {
		row := make([]int, len(l))
		for x, c := range l {
			n, _ := strconv.Atoi(string(c))
			row[x] = n
		}
		grid[y] = row
	}

	return grid
}

func PartOne(lines []string) int {
	grid := loadGrid(lines)
	return findLowestCost(point{x: len(grid[0]) - 1, y: len(grid) - 1}, grid, false)
}

func PartTwo(lines []string) int {
	grid := loadGrid(lines)
	return findLowestCost(point{x: len(grid[0]) - 1, y: len(grid) - 1}, grid, true)
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
