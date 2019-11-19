package main

import (
	"math/rand"

	"github.com/pezza/advent-of-wasm/wasm"
)

type flake struct {
	x     int
	y     int
	speed int
	style string
}

var doc wasm.JsDoc
var flakes []flake
var tick float64 = 1000 / 60
var currentTick float64
var time float64

var flakeCount = 250

var canvasDrawWidth, canvasDrawHeight = 800, 600

func main() {
	done := make(chan bool, 0)

	doc = wasm.NewJsDoc("canv")

	flakes = make([]flake, flakeCount)

	for index := range flakes {
		flakes[index].x = rand.Intn(canvasDrawWidth)
		flakes[index].y = rand.Intn(canvasDrawHeight)

		speed := rand.Intn(4) + 1
		flakes[index].speed = speed

		style := "#000000"

		switch speed {
		case 1:
			style = "#333333"
		case 2:
			style = "#777777"
		case 3:
			style = "#aaaaaa"
		case 4:
			style = "#cccccc"
		case 5:
			style = "#ffffff"
		}

		flakes[index].style = style
	}

	doc.StartAnimLoop(frame)
	<-done
}

func frame(now float64) {
	delta := now - time
	time = now
	currentTick += delta

	if currentTick < tick {
		return
	}

	currentTick = 0
	doc.ClearFrame(0, 0, canvasDrawWidth, canvasDrawHeight)

	for i := range flakes {
		doc.DrawRect(flakes[i].x, flakes[i].y, flakes[i].speed, flakes[i].speed, flakes[i].style)
		flakes[i].y += flakes[i].speed
		if flakes[i].y > canvasDrawHeight {
			flakes[i].y -= canvasDrawHeight
		}
	}
}
