package main

import (
	"github.com/pezza/advent-of-wasm/wasm"
	"strconv"
	"syscall/js"
)

var doc wasm.JsDoc
var canvas wasm.JsCanvas
var currentTime float64

var canvasDrawWidth, canvasDrawHeight int

var flakeCountSlider js.Value
var flakeCountLabel js.Value


func main() {
	done := make(chan bool, 0)

	doc = wasm.NewJsDoc()

	flakeCountSlider = doc.GetElementByID("flakecount")
	flakeCountLabel = doc.GetElementByID("flakecount-value")

	flakeCount := 500

	doc.SetInnerHTML(flakeCountLabel, strconv.Itoa(flakeCount))

	canvasDrawWidth, canvasDrawHeight = doc.GetWindowSize()
	canvasDrawWidth = canvasDrawWidth / 3
	canvasDrawHeight = canvasDrawHeight / 3

	flakes = createFlakes(flakeCount, canvasDrawWidth, canvasDrawHeight)
	canvas = doc.GetOrCreateCanvas("canv", canvasDrawWidth, canvasDrawHeight)

	doc.AddEventListener(flakeCountSlider, "input", js.FuncOf(countHandlerFunc))
	doc.StartAnimLoop(frame)

	<-done
}

func frame(now float64) {
	delta := now - currentTime
	currentTime = now

	update(delta)
	draw(delta)
}

func update(delta float64) {
	flakes.update(delta, canvasDrawHeight)
}

func draw(delta float64) {
	canvas.ClearFrame(0, 0, canvasDrawWidth, canvasDrawHeight)

	prevStyle :=""
	for i := range flakes {
		if prevStyle != flakes[i].style {
			canvas.SetFillStyle(flakes[i].style)
			prevStyle = flakes[i].style
		}
		canvas.DrawFilledRect(flakes[i].x, flakes[i].y, flakes[i].speed, flakes[i].speed-1)
	}

	canvas.DrawFrame()
}

func countHandlerFunc(this js.Value, args []js.Value) interface{} {
	newFlakeCount, _ := strconv.Atoi(this.Get("value").String())
	flakes = flakes.adjustFlakes(newFlakeCount)
	doc.SetInnerHTML(flakeCountLabel, strconv.Itoa(newFlakeCount))
	return false
}


