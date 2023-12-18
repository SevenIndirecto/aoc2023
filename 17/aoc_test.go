package main

import "testing"

func TestPartOne(t *testing.T) {
	lines, _ := LoadLines("test.txt")

	expected := 102
	got := PartOne(lines)
	if got != expected {
		t.Errorf("Part one failed, got %v, expected %v", got, expected)
	}
}

func TestPartTwo(t *testing.T) {
	lines, _ := LoadLines("test2.txt")

	expected := 71
	got := PartTwo(lines)
	if got != expected {
		t.Errorf("Test failed, got %v, expected %v", got, expected)
	}
}
