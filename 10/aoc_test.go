package main

import "testing"

func TestPartOne(t *testing.T) {
	lines, _ := LoadLines("test1-1.txt")
	expected := 4
	got := PartOne(lines)
	if got != expected {
		t.Errorf("Part one failed, got %v, expected %v", got, expected)
	}

	lines, _ = LoadLines("test1-2.txt")
	expected = 8
	got = PartOne(lines)
	if got != expected {
		t.Errorf("Part one failed, got %v, expected %v", got, expected)
	}
}

func TestPartTwo(t *testing.T) {
	lines, _ := LoadLines("test2-1.txt")

	expected := 4
	got := PartTwo(lines)
	if got != expected {
		t.Errorf("Test failed, got %v, expected %v", got, expected)
	}

	lines, _ = LoadLines("test2-2.txt")
	expected = 10
	got = PartTwo(lines)
	if got != expected {
		t.Errorf("Test failed, got %v, expected %v", got, expected)
	}
}
