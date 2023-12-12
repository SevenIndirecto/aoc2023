package main

import (
	"fmt"
	"io/ioutil"
	"strings"
)

const (
	ground   = iota
	vertical = iota
	horizontal
	ne
	nw
	sw
	se
	start
)

type point struct {
	x, y int
}

type pipe struct {
	prev, next *pipe
	loc        point
	t          int
	distance   int
}

func loadGrid(lines []string) ([][]int, point) {
	grid := make([][]int, len(lines))
	startPoint := point{}

	mapping := map[rune]int{
		'|': vertical,
		'-': horizontal,
		'L': ne,
		'J': nw,
		'7': sw,
		'F': se,
		'.': ground,
		'S': start,
	}

	for y, line := range lines {
		row := make([]int, len(line))

		for x, c := range line {
			row[x] = mapping[c]
			if row[x] == start {
				startPoint.x = x
				startPoint.y = y
			}
		}
		grid[y] = row
	}

	return grid, startPoint
}

func constructLoop(grid [][]int, startPoint point) *pipe {
	directionDeltas := []point{
		{x: -1, y: 0},
		{x: 1, y: 0},
		{x: 0, y: -1},
		{x: 0, y: 1},
	}

	pipeDirectionMapping := map[int]map[point]point{
		vertical: {
			// N -> S
			{x: 0, y: 1}: {x: 0, y: 1},
			// S -> N
			{x: 0, y: -1}: {x: 0, y: -1},
		},
		horizontal: {
			// W -> E
			{x: 1, y: 0}: {x: 1, y: 0},
			// E -> W
			{x: -1, y: 0}: {x: -1, y: 0},
		},
		ne: {
			// N -> E
			{x: 0, y: 1}: {x: 1, y: 0},
			// E -> N
			{x: -1, y: 0}: {x: 0, y: -1},
		},
		nw: {
			// N -> W
			{x: 0, y: 1}: {x: -1, y: 0},
			// W -> N
			{x: 1, y: 0}: {x: 0, y: -1},
		},
		sw: {
			// S -> W
			{x: 0, y: -1}: {x: -1, y: 0},
			// W -> S
			{x: 1, y: 0}: {x: 0, y: 1},
		},
		se: {
			// S -> E
			{x: 0, y: -1}: {x: 1, y: 0},
			// E -> S
			{x: -1, y: 0}: {x: 0, y: 1},
		},
	}

	// Find first viable next direction
	var next point
	var nextType int
	for _, delta := range directionDeltas {
		next = point{x: startPoint.x + delta.x, y: startPoint.y + delta.y}

		if next.y < 0 || next.x < 0 || next.y >= len(grid) || next.x >= len(grid[next.y]) {
			continue
		}

		nextType = grid[next.y][next.x]
		if (delta.y == 1 && (nextType == ne || nextType == nw || nextType == vertical)) ||
			(delta.y == -1 && (nextType == se || nextType == sw || nextType == vertical)) ||
			(delta.x == 1 && (nextType == nw || nextType == sw || nextType == horizontal)) ||
			(delta.x == -1 && (nextType == ne || nextType == se || nextType == horizontal)) {
			break
		}
	}

	prevPipe := &pipe{loc: startPoint, t: start, distance: 0}
	currentPipe := &pipe{loc: next, prev: prevPipe, t: nextType, distance: -1}
	prevPipe.next = currentPipe
	startPipe := prevPipe

	for {
		enterDelta := point{x: currentPipe.loc.x - prevPipe.loc.x, y: currentPipe.loc.y - prevPipe.loc.y}
		exitDelta := pipeDirectionMapping[currentPipe.t][enterDelta]
		nextPoint := point{x: currentPipe.loc.x + exitDelta.x, y: currentPipe.loc.y + exitDelta.y}

		if nextPoint.x == startPoint.x && nextPoint.y == startPoint.y {
			// Connect start pipe
			startPipe.prev = currentPipe
			currentPipe.next = startPipe
			break
		}

		currentPipe.next = &pipe{loc: nextPoint, prev: currentPipe, t: grid[nextPoint.y][nextPoint.x], distance: -1}
		prevPipe = currentPipe
		currentPipe = currentPipe.next
	}

	// Determine start pipe type
	np := startPipe.next.loc
	pp := startPipe.prev.loc
	sp := startPipe.loc

	if np.y == pp.y {
		startPipe.t = horizontal
	} else if np.x == pp.x {
		startPipe.t = vertical
	} else if np.x > sp.x && pp.y < sp.y || pp.x > sp.x && np.y < sp.y {
		startPipe.t = ne
	} else if np.x < sp.x && pp.y < sp.y || pp.x < sp.x && np.y < sp.y {
		startPipe.t = nw
	} else if pp.x > sp.x && np.y > sp.y || np.x > sp.x && pp.y > sp.y {
		startPipe.t = se
	} else {
		startPipe.t = sw
	}

	return startPipe
}

func PartOne(lines []string) int {
	grid, startPoint := loadGrid(lines)
	startPipe := constructLoop(grid, startPoint)
	currentPipe := startPipe

	distance := 0
	for {
		distance++
		currentPipe = currentPipe.next
		if currentPipe.loc == startPipe.loc {
			break
		}
		currentPipe.distance = distance
	}

	distance = 0
	for {
		distance++
		currentPipe = currentPipe.prev

		if currentPipe.loc == startPipe.loc || distance >= currentPipe.distance {
			break
		}
		currentPipe.distance = distance
	}

	currentPipe = startPipe.next
	maxDistance := 0
	for {
		if currentPipe.distance > maxDistance {
			maxDistance = currentPipe.distance
		}
		if currentPipe.loc == startPipe.loc {
			break
		}
		currentPipe = currentPipe.next
	}

	return maxDistance
}

func PartTwo(lines []string) int {
	grid, startPoint := loadGrid(lines)
	startPipe := constructLoop(grid, startPoint)
	loopPointMap := map[point]int{startPipe.loc: startPipe.t}

	current := startPipe.next
	for {
		if current.loc == startPipe.loc {
			break
		}
		loopPointMap[current.loc] = current.t
		current = current.next
	}

	insideCount := 0
	for y, row := range grid {
		for x, _ := range row {
			// Ray trace

			// exclude loop points and edges
			if x == 0 || y == 0 || x == len(row)-1 || y == len(grid) || loopPointMap[point{x: x, y: y}] > 0 {
				continue
			}

			intersectCount := 0
			for i := 0; i < x; i++ {
				p := point{x: i, y: y}
				loopType := loopPointMap[p]
				if loopType == vertical || loopType == se || loopType == sw {
					intersectCount++
				}
			}
			insideCount += intersectCount % 2
		}
	}

	return insideCount
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
