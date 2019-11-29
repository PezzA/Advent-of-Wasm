package main

import (
	"fmt"

	puzzles "github.com/pezza/advent-of-code"
)

func runPuzzle(day, year int, updateChan chan []string) {
	puzzle, err := puzzles.GetPuzzle(day, year)

	if err != nil {
		fmt.Println(err)
	}

	go doPart(puzzle.PartTwo, puzzle.PuzzleInput(), updateChan)
}

func doPart(fn puzzles.PuzzlePart, inputData string, updateChan chan []string) {
	fmt.Println("dopartstart")
	fmt.Println(fn(inputData, updateChan))
	fmt.Println("dopartend")
}
