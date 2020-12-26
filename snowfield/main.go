package main

import (
	"github.com/pezza/advent-of-wasm/wasm"
	"math/rand"
	"sort"
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
var canvas wasm.JsCanvas
var currentTime float64
var flakes []flake
var canvasDrawWidth, canvasDrawHeight int



func main() {
	done := make(chan bool, 0)

	doc = wasm.NewJsDoc()
	canvasDrawWidth, canvasDrawHeight = doc.GetWindowSize()

	canvasDrawWidth = canvasDrawWidth / 3
	canvasDrawHeight = canvasDrawHeight / 3

	flakeCount := 500
	flakes = createFlakes(flakeCount)

	canvas = doc.GetOrCreateCanvas("canv", canvasDrawWidth, canvasDrawHeight)

	doc.AddEventListener("flakecount", "input", js.FuncOf(countHandlerfunc))
	doc.StartAnimLoop(frame)

	<-done
}

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

	sort.Slice(flakeArray, func(i, j int) bool {
		return flakeArray[i].speed < flakeArray[j].speed
	})


	return flakeArray
}


func frame(now float64) {

	delta := now - currentTime
	currentTime = now
	canvas.ClearFrame(0, 0, canvasDrawWidth, canvasDrawHeight)

	prevStyle :=""
	for i := range flakes {

		if prevStyle != flakes[i].style {
			canvas.SetFillStyle(flakes[i].style)
			prevStyle = flakes[i].style
		}

		canvas.DrawFilledRect(flakes[i].x, flakes[i].y, flakes[i].speed, flakes[i].speed-1)

		yUpdate := int(float64(flakes[i].speed) * (delta / 30))

		if yUpdate < 1 {
			yUpdate = 1
		}
		flakes[i].y +=yUpdate
		if flakes[i].y > canvasDrawHeight {
			flakes[i].y -= canvasDrawHeight
		}
	}

	canvas.DrawFrame()
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

	return createFlakes(newCount)
}

