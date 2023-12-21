package main

import (
	"fmt"
	"io/ioutil"
	"strconv"
	"strings"
)

const (
	lessThanCheck = iota
	greaterThanCheck
	noCheck
)

type rule struct {
	checkType   int
	category    int
	value       int
	destination string
}

type workflow struct {
	id    string
	rules []rule
}

type part [4]int

func (r *rule) check(p part) (string, bool) {
	if r.checkType == noCheck {
		return r.destination, false
	} else if r.checkType == greaterThanCheck {
		if p[r.category] > r.value {
			return r.destination, false
		}
		return "", true
	} else {
		if p[r.category] < r.value {
			return r.destination, false
		}
	}
	return "", true
}

func loadPartsAndWorkflows(lines []string) ([]part, map[string]workflow) {
	categoryMapping := map[string]int{"x": 0, "m": 1, "a": 2, "s": 3}
	workflows := map[string]workflow{}

	// Load workflows
	i := 0
	for ; i < len(lines); i++ {
		workflowString := lines[i]
		if workflowString == "" {
			i++
			break
		}
		s := strings.Split(workflowString, "{")
		w := workflow{id: s[0], rules: make([]rule, 0)}

		ruleStrings := strings.Split(s[1][:len(s[1])-1], ",")
		for _, singleRuleString := range ruleStrings {
			r := rule{}

			var rs []string
			if strings.Contains(singleRuleString, "<") {
				r.checkType = lessThanCheck
				rs = strings.Split(singleRuleString, "<")
				r.category = categoryMapping[rs[0]]
			} else if strings.Contains(singleRuleString, ">") {
				r.checkType = greaterThanCheck
				rs = strings.Split(singleRuleString, ">")
				r.category = categoryMapping[rs[0]]
			} else {
				r.checkType = noCheck
				r.destination = singleRuleString
			}

			if r.checkType == lessThanCheck || r.checkType == greaterThanCheck {
				checkSplit := strings.Split(rs[1], ":")
				v, _ := strconv.Atoi(checkSplit[0])
				r.value = v
				r.destination = checkSplit[1]
			}
			w.rules = append(w.rules, r)
		}
		workflows[w.id] = w
	}

	// Load parts
	parts := make([]part, 0)
	for ; i < len(lines); i++ {
		partString := lines[i][1 : len(lines[i])-1]
		s := strings.Split(partString, ",")
		p := part{}
		for _, categoryString := range s {
			cs := strings.Split(categoryString, "=")
			v, _ := strconv.Atoi(cs[1])
			p[categoryMapping[cs[0]]] = v
		}
		parts = append(parts, p)
	}

	return parts, workflows
}

func evaluatePart(p part, workflows map[string]workflow) bool {
	currentWorkflow := "in"

	for {
		w := workflows[currentWorkflow]

		for _, r := range w.rules {
			destination, noMatch := r.check(p)
			if noMatch {
				// Try next rule
				continue
			} else {
				if destination == "A" {
					return true
				} else if destination == "R" {
					return false
				} else {
					currentWorkflow = destination
					break
				}
			}
		}

		_, exists := workflows[currentWorkflow]
		if !exists {
			fmt.Println("Testing part", p, "current wf target", currentWorkflow)
			panic("Unexpected end, ran out of workflows")
			return false
		}
	}
}

func result(parts []part) int {
	sum := 0
	for _, p := range parts {
		for _, v := range p {
			sum += v
		}
	}
	return sum
}

func PartOne(lines []string) int {
	parts, workflows := loadPartsAndWorkflows(lines)
	acceptedParts := make([]part, 0)

	for _, p := range parts {
		if evaluatePart(p, workflows) {
			acceptedParts = append(acceptedParts, p)
		}
	}

	return result(acceptedParts)
}

// Part two

// Each limit has [Lowest possible value, highest possible] inclusive
type condition struct {
	minMax          [4][2]int
	workflowToCheck string
}

func newCondition(workflowId string) condition {
	c := condition{workflowToCheck: workflowId}
	for i := 0; i < 4; i++ {
		c.minMax[i][floor] = 1
		c.minMax[i][ceil] = 4000
	}
	return c
}

const (
	floor = iota
	ceil
)

// returns (conditionForNewWorkflow, conditionForNextRule, isRejected, isAccepted)
func mergeRuleAndCondition(r rule, c condition) (condition, condition, bool, bool) {
	nextRuleCondition := c

	if r.checkType == noCheck {
		// No change to condition limits since no checks were made
		c.workflowToCheck = r.destination
		if r.destination == "R" {
			return c, c, true, false
		}
		return c, c, false, r.destination == "A"
	}

	catLimits := c.minMax[r.category]
	nextRuleCheckLimits := c.minMax[r.category]

	if r.checkType == greaterThanCheck {
		if r.value >= catLimits[floor] {
			catLimits[floor] = r.value + 1
		}
		// Update next check rule limits
		if nextRuleCheckLimits[ceil] > r.value {
			nextRuleCheckLimits[ceil] = r.value
		}
	} else if r.checkType == lessThanCheck {
		if r.value <= catLimits[ceil] {
			catLimits[ceil] = r.value - 1
		}
		// Update next check rule limits
		if nextRuleCheckLimits[floor] < r.value {
			nextRuleCheckLimits[floor] = r.value
		}
	} else {
		panic("Unexpected state")
	}

	nextRuleCondition.minMax[r.category] = nextRuleCheckLimits

	if catLimits[floor] > catLimits[ceil] {
		// reject
		return c, nextRuleCondition, true, false
	}

	c.minMax[r.category] = catLimits
	c.workflowToCheck = r.destination
	return c, nextRuleCondition, r.destination == "R", r.destination == "A"
}

func findCombinations(workflows map[string]workflow) []condition {
	acceptedConditions := make([]condition, 0)
	queue := []condition{newCondition("in")}

	for len(queue) > 0 {
		// pop condition
		c := queue[len(queue)-1]
		queue = queue[:len(queue)-1]

		// process possibilities
		for _, r := range workflows[c.workflowToCheck].rules {
			conditionForNewWorkflow, conditionForNextRule, isRejected, isAccepted := mergeRuleAndCondition(r, c)

			if isAccepted {
				acceptedConditions = append(acceptedConditions, conditionForNewWorkflow)
			} else if !isRejected {
				queue = append(queue, conditionForNewWorkflow)
			}
			c = conditionForNextRule
		}
	}

	//for _, c := range acceptedConditions {
	//	fmt.Println(c)
	//}
	return acceptedConditions
}

func scoreCombinations(combinations []condition) int64 {
	sum := int64(0)

	for _, c := range combinations {
		mul := int64(1)

		for _, minMax := range c.minMax {
			delta := minMax[ceil] - minMax[floor] + 1
			if delta < 0 {
				fmt.Println(minMax)
				panic("Unexpected minMax")
			}
			mul *= int64(delta)
		}
		sum += mul
	}
	return sum
}

func PartTwo(lines []string) int64 {
	_, workflows := loadPartsAndWorkflows(lines)
	combinations := findCombinations(workflows)
	return scoreCombinations(combinations)
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
