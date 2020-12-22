package main

import (
	"github.com/pezza/advent-of-wasm/wasm"
	"math/rand"
	"strconv"
	"syscall/js"
	"fmt"
)

type flake struct {
	x     int
	y     int
	speed int
	style string
}

var doc wasm.JsDoc
var flakes []flake

var canvasDrawWidth, canvasDrawHeight = 800, 600

func createFlakes(flakeCount int) []flake {
	flakeArray := make([]flake, flakeCount)

	for index := range flakeArray {
		flakeArray[index].x = rand.Intn(canvasDrawWidth)
		flakeArray[index].y = rand.Intn(canvasDrawHeight) - canvasDrawHeight

		speed := rand.Intn(5) +2
		flakeArray[index].speed = speed

		style := "#666666" // white

		switch speed {
		case 1:
			style = "#00FF00" // green
		case 2:
			style = "#888888" //blue
		case 3:
			style = "#BBBBBB" //yellow
		case 4:
			style = "#DDDDDD" // cyan
		case 5:
			style = "#EEEEEE" // magenta
		case 6:
			style = "#FFFFFF" // red
		}



		flakeArray[index].style = style
	}

	return flakeArray
}

func main() {
	done := make(chan bool, 0)

	doc = wasm.NewJsDoc("canv")

	doc.AddEventListener("flakecount", "input", js.FuncOf(countHandlerfunc))

	flakeCount := 1000

	flakes = createFlakes(flakeCount)

	doc.StartAnimLoop(frame)

	<-done
}

var currentTime float64

func frame(now float64) {

	delta := now - currentTime
	currentTime = now
	doc.ClearFrame(0, 0, canvasDrawWidth, canvasDrawHeight)
	for i := range flakes {
		doc.DrawRect(flakes[i].x, flakes[i].y, flakes[i].speed-1, flakes[i].speed-1, flakes[i].style)
		flakes[i].y += int(float64(flakes[i].speed) * (delta / 20))
		if flakes[i].y > canvasDrawHeight {
			flakes[i].y -= canvasDrawHeight
		}
	}
}
func countHandlerfunc(this js.Value, args []js.Value) interface{} {
	newFlakeCount, err := strconv.Atoi(this.Get("value").String())

	if err != nil {
		fmt.Println(err)
		return nil
	}

	flakes = adjustFlakes(newFlakeCount, flakes)

	doc.SetElementInnerHTML("flakecount-value", strconv.Itoa(newFlakeCount))

	return nil
}

func adjustFlakes(newCount int, current []flake) []flake {
	if newCount == len(current) {
		return current
	}

	if newCount > len(current) {
		newflakes := createFlakes(newCount - len(current))
		return append(current, newflakes...)
	}

	return flakes[0:newCount]
}