package main

import (
	"fmt"
	"testing"
)

func TestPartOne(t *testing.T) {
	lines, _ := LoadLines("test.txt")
	expected := 32000000
	got := PartOne(lines)
	if got != expected {
		t.Errorf("Part one-1 failed, got %v, expected %v", got, expected)
	}

	fmt.Println("Test 1-2")
	lines, _ = LoadLines("test2.txt")
	expected = 11687500
	got = PartOne(lines)
	if got != expected {
		t.Errorf("Part one-2 failed, got %v, expected %v", got, expected)
	}
}

func TestPartTwo(t *testing.T) {
	lines, _ := LoadLines("test.txt")

	expected := -1
	got := PartTwo(lines)
	if got != expected {
		t.Errorf("Test failed, got %v, expected %v", got, expected)
	}
}
