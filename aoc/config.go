package main

import (
	"github.com/pezza/advent-of-code/2020/Day202024"
	"github.com/pezza/advent-of-code/common"
)

// size of grid
const part1HexSize = 8
const part2HexSize = 2

// animation speed
const animTickPoll = float64(10)

// how often to update state
const stateTickPoll = float64(20)

const state2TickPoll = float64(20)

const tilesToAdd = 15

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

type tile struct {
	common.Point
	margin  int
	endTile bool
}
