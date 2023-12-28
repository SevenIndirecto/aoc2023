package main

import (
	"fmt"
	"io/ioutil"
	"strings"
)

// Entire solution copied from reddit...
// The idea is to find the three edges which get visited most often when traversing the graph from every node using BFS/DFS,
// since these are choke points connecting the two larger graphs.
//
// So i used BFS on all some nodes as input and counted each occurrence of an edge, the one with the highest should be one
// of the three. Then i remove that edge and run the search again etc., until i have all three edges. I remove them from my
// nodes list, start a BFS from one node to get the node count, then subtract from the total for the second value and voila, get the result.
func PartOne(lines []string) int {
	nodes := getNodes(lines)
	edgeCount := countEdges(nodes)
	first := findAndRemoveMax(edgeCount)
	removeEdge(nodes, first)
	edgeCount = countEdges(nodes)
	second := findAndRemoveMax(edgeCount)
	removeEdge(nodes, second)
	edgeCount = countEdges(nodes)
	third := findAndRemoveMax(edgeCount)
	removeEdge(nodes, third)

	countA := countNodes(nodes, first.from)
	countB := len(nodes) - countA

	return countA * countB
}

type Edge struct {
	from, to string
}

func removeEdge(nodes map[string][]string, edge Edge) {
	new := []string{}
	for _, val := range nodes[edge.from] {
		if val != edge.to {
			new = append(new, val)
		}
	}
	nodes[edge.from] = new
	new = []string{}
	for _, val := range nodes[edge.to] {
		if val != edge.from {
			new = append(new, val)
		}
	}
	nodes[edge.to] = new
}

func findAndRemoveMax(edges map[Edge]int) Edge {
	max := 0
	var maxEdge Edge
	for key, val := range edges {
		if val > max {
			max = val
			maxEdge = key
		}
	}
	delete(edges, maxEdge)
	return maxEdge
}

func countEdges(nodes map[string][]string) map[Edge]int {
	encountered := map[Edge]int{}
	i := 0
	for from := range nodes {
		walkNodes(nodes, from, encountered)
		i++
		// if this doesn't work, increase value
		if i > 50 {
			break
		}
	}
	return encountered
}

func countNodes(nodes map[string][]string, start string) int {
	visited := map[string]bool{}
	queue := []string{start}

	for len(queue) > 0 {
		from := queue[0]
		queue = queue[1:]

		for _, to := range nodes[from] {
			if _, found := visited[to]; found {
				continue
			}
			queue = append(queue, to)
			visited[to] = true
		}
	}
	return len(visited)
}

func walkNodes(nodes map[string][]string, start string, encountered map[Edge]int) {
	visited := map[string]bool{}
	queue := []string{start}

	for len(queue) > 0 {
		from := queue[0]
		queue = queue[1:]

		for _, to := range nodes[from] {
			if _, found := visited[to]; found {
				continue
			}
			queue = append(queue, to)
			visited[to] = true
			var edge Edge
			if from < to {
				edge = Edge{from, to}
			} else {
				edge = Edge{to, from}
			}
			encountered[edge]++
		}
	}
}

func getNodes(lines []string) map[string][]string {
	nodes := map[string][]string{}
	for _, line := range lines {
		split := strings.Split(line, ": ")
		from := split[0]
		if _, inside := nodes[from]; !inside {
			nodes[from] = []string{}
		}

		to := strings.Fields(split[1])
		for _, target := range to {
			nodes[from] = append(nodes[from], target)
			if _, inside := nodes[target]; !inside {
				nodes[target] = []string{}
			}
			nodes[target] = append(nodes[target], from)
		}
	}
	return nodes
}

func PartTwo(lines []string) int {
	return 0
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
