package main

import (
	"github.com/pezza/advent-of-code/puzzles/2020/Day202024"
	"github.com/pezza/advent-of-code/puzzles/common"
	"github.com/pezza/wasm"
)

// size of grid
const part1HexSize = 8
const part2HexSize = 2
const tilesToAdd = 15

type config struct {
	animTickPoll   float64
	stateTickPoll  float64
	state2TickPoll float64
	hexSize        int
	canvasWidth    int
	canvasHeight   int
	hex            []wasm.Point
	routeList      [][]common.Point
	insList        [][]string
	hexHalfSize    int
	hexThreeQSize  int
	hexFullSize    int
}

func getNewConfig(w, h, size int) config {
	c := config{
		animTickPoll:   float64(10),
		stateTickPoll:  float64(20),
		state2TickPoll: float64(20),
		canvasWidth:    w,
		canvasHeight:   h,
		insList:        Day202024.GetData(Day202024.Entry.PuzzleInput())}

	c.routeList = make([][]common.Point, len(c.insList))

	for i := range c.insList {
		c.routeList[i] = Day202024.GetMoveList(c.insList[i])
	}

	c.setSizes(size)

	return c
}

func (c *config) setSizes(size int) {
	c.hexSize = size
	c.hexHalfSize = size * 2
	c.hexThreeQSize = size * 3
	c.hexFullSize = size * 4
	c.hex = getHex(size)
}
