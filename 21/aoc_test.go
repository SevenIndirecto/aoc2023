package main

import "testing"

func TestPartOne(t *testing.T) {
	lines, _ := LoadLines("test.txt")

	expected := 16
	got := PartOne(lines, 6)
	if got != expected {
		t.Errorf("Part one failed, got %v, expected %v", got, expected)
	}
}

func TestPartTwo(t *testing.T) {
	lines, _ := LoadLines("test.txt")

	stepsToExpectedPlots := map[int]int{
		//6:    16,
		//10:   50,
		//50:   1594,
		//100:  6536,
		500: 167004,
		//1000: 668697,
		//5000: 16733044,
	}

	for steps, expected := range stepsToExpectedPlots {
		got := PartTwo(lines, steps)
		if got != expected {
			t.Errorf("Test failed for %v steps, got %v, expected %v", steps, got, expected)
		}
	}
}
