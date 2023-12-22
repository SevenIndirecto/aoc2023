package main

import (
	"fmt"
	"io/ioutil"
	"strings"
)

const (
	flipFlip = iota
	conjunction
	broadcaster
)

const (
	lowPulse = iota
	highPulse
)

type queuedPulse struct {
	pulseType int
	sender    string
	target    string
}

type module struct {
	id         string
	moduleType int
	outputs    []string

	lastPulse   int
	isOn        bool
	inputMemory map[string]int
}

type queue struct {
	bucket []queuedPulse
}

func newQueue() queue {
	return queue{bucket: make([]queuedPulse, 0)}
}

func (q *queue) isEmpty() bool {
	return len(q.bucket) < 1
}

func (q *queue) deque() queuedPulse {
	qp := q.bucket[0]
	q.bucket = q.bucket[1:]
	return qp
}

func (q *queue) add(qp queuedPulse) {
	q.bucket = append(q.bucket, qp)
}

func sendPulse(qp queuedPulse, modules map[string]*module, q *queue) {
	t := modules[qp.target]
	if t == nil {
		return
	}
	_, exists := t.inputMemory[qp.sender]
	if exists {
		t.inputMemory[qp.sender] = qp.pulseType
	}

	shouldSendPulse := false
	pulseToSend := lowPulse

	if t.moduleType == flipFlip {
		if qp.pulseType == highPulse {
			// NOOP
			return
		}

		t.isOn = !t.isOn
		shouldSendPulse = true

		if t.isOn {
			pulseToSend = highPulse
		}
	} else if t.moduleType == conjunction {
		allHigh := true
		for _, pulseType := range t.inputMemory {
			if pulseType == lowPulse {
				allHigh = false
				break
			}
		}

		shouldSendPulse = true
		if allHigh {
			pulseToSend = lowPulse
		} else {
			pulseToSend = highPulse
		}
	} else if t.moduleType == broadcaster {
		pulseToSend = qp.pulseType
		shouldSendPulse = true
	} else {
		panic("Invalid target type")
	}

	if shouldSendPulse {
		for _, o := range t.outputs {
			q.add(queuedPulse{pulseType: pulseToSend, sender: t.id, target: o})
		}
	}
}

func loadModules(lines []string) map[string]*module {
	modules := map[string]*module{}
	inputs := map[string]map[string]bool{}

	for _, l := range lines {
		s := strings.Split(l, " -> ")
		m := module{}
		if s[0] == "broadcaster" {
			m.id = s[0]
			m.moduleType = broadcaster
		} else {
			if s[0][0] == '%' {
				m.moduleType = flipFlip
			} else if s[0][0] == '&' {
				m.moduleType = conjunction
			} else {
				panic("Invalid module type")
			}
			m.id = s[0][1:]
		}

		m.outputs = strings.Split(s[1], ", ")
		for _, outputTarget := range m.outputs {
			_, exists := inputs[outputTarget]
			if !exists {
				inputs[outputTarget] = map[string]bool{}
			}
			inputs[outputTarget][m.id] = true
		}
		modules[m.id] = &m
	}

	// Input init for conjunctions
	for moduleId, inputsToModule := range inputs {
		_, exists := modules[moduleId]
		if !exists {
			continue
		}
		moduleInputMemory := map[string]int{}
		for inputId := range inputsToModule {
			moduleInputMemory[inputId] = lowPulse
		}
		modules[moduleId].inputMemory = moduleInputMemory
	}
	return modules
}

func PartOne(lines []string) int {
	q := newQueue()
	modules := loadModules(lines)

	sentPulses := [2]int{0, 0}

	i := 0
	for {
		if i >= 1000 && q.isEmpty() {
			break
		}

		if q.isEmpty() {
			// push button
			q.add(queuedPulse{pulseType: lowPulse, sender: "button", target: "broadcaster"})
			i++
		}
		// process queue
		qp := q.deque()
		sentPulses[qp.pulseType]++
		if qp.target == "rx" && qp.pulseType == lowPulse {
			break
		}
		sendPulse(qp, modules, &q)
	}

	fmt.Println("Pulse count", sentPulses)
	return sentPulses[0] * sentPulses[1]
}

func PartTwo(lines []string) int {
	q := newQueue()
	modules := loadModules(lines)

	target := modules["gf"]
	lastHigh := map[string]int{}

	i := 0
	for {
		if q.isEmpty() {
			// push button
			q.add(queuedPulse{pulseType: lowPulse, sender: "button", target: "broadcaster"})
			i++
		}
		// process queue
		qp := q.deque()
		if qp.target == "rx" && qp.pulseType == lowPulse {
			break
		}

		for inputId, pulse := range target.inputMemory {
			if pulse == highPulse {
				delta := i - lastHigh[inputId]
				if delta == 0 {
					continue
				}
				lastHigh[inputId] = i
				fmt.Println(inputId, "delta:", delta)
			}
		}

		mul := 1
		index := 0
		repeats := make([]int, 4)
		for _, repeat := range lastHigh {
			mul *= repeat
			repeats[index] = repeat
			index++
		}
		if mul > 0 && index == 4 {
			fmt.Println(repeats)
			return LCM(repeats[0], repeats[1], repeats[2:]...)
		}

		sendPulse(qp, modules, &q)
	}

	return i
}

// GCD greatest common divisor (GCD) via Euclidean algorithm
func GCD(a, b int) int {
	for b != 0 {
		t := b
		b = a % b
		a = t
	}
	return a
}

// LCM find Least Common Multiple (LCM) via GCD
func LCM(a, b int, integers ...int) int {
	result := a * b / GCD(a, b)

	for i := 0; i < len(integers); i++ {
		result = LCM(result, integers[i])
	}

	return result
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
