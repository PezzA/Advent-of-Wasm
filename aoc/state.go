package main

import (
	"github.com/pezza/advent-of-code/puzzles/2020/Day202024"
	"github.com/pezza/advent-of-code/puzzles/common"
)

type state struct {
	tiles             []tile
	currentTime       float64
	animTick          float64
	stateTick         float64
	insIndex          int
	subInsIndex       int
	running           bool
	paused            bool
	puzzlePart        int
	automataIteration int
	maxIteration      int
	lastLen           int
	Day202024.Floor
}

type tile struct {
	common.Point
	margin  int
	endTile bool
}

func getNewState(running bool, puzzlePart int) state {
	return state{
		tiles:             make([]tile, 0),
		animTick:          0,
		stateTick:         0,
		currentTime:       0,
		insIndex:          0,
		subInsIndex:       0,
		running:           running,
		paused:            false,
		puzzlePart:        puzzlePart,
		automataIteration: 0,
		maxIteration:      100,
		lastLen:           0,
		Floor:             make(Day202024.Floor, 0),
	}
}
