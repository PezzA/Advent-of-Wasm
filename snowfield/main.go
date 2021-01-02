package main

import (
	"math/rand"
	"strconv"
	"syscall/js"
	"time"

	"github.com/pezza/advent-of-wasm/wasm"
)

func main() {
	done := make(chan bool, 0)

	rand.Seed(time.Now().UnixNano())

	doc := wasm.NewJsDoc()
	flakeCountSlider := doc.GetElementByID("flakecount")
	flakeCountLabel := doc.GetElementByID("flakecount-value")

	flakeCount := 500

	doc.SetInnerHTML(flakeCountLabel, strconv.Itoa(flakeCount))

	canvasDrawWidth, canvasDrawHeight := doc.GetWindowSize()

	flakes = createFlakes(flakeCount, canvasDrawWidth, canvasDrawHeight)
	canvas := doc.GetOrCreateCanvas("canv", canvasDrawWidth, canvasDrawHeight, true, false)

	doc.AddEventListener(flakeCountSlider, "input", js.FuncOf(func(this js.Value, args []js.Value) interface{} {
		newFlakeCount, _ := strconv.Atoi(this.Get("value").String())
		flakes = flakes.adjustFlakes(newFlakeCount, canvasDrawWidth, canvasDrawHeight)
		doc.SetInnerHTML(flakeCountLabel, strconv.Itoa(newFlakeCount))
		return false
	}))

	var currentTime float64

	doc.StartAnimLoop(func(now float64) {
		delta := now - currentTime
		currentTime = now

		update(delta, canvasDrawHeight)
		draw(delta, canvas)
	})

	<-done
}

func update(delta float64, canvasDrawHeight int) {
	flakes.update(delta, canvasDrawHeight)
}

func draw(delta float64, canvas *wasm.JsCanvas) {
	canvas.Clear()
	prevStyle := ""
	for i := range flakes {
		if prevStyle != flakes[i].style {
			canvas.SetFillStyle(flakes[i].style)
			prevStyle = flakes[i].style
		}
		canvas.DrawRect(flakes[i].x, flakes[i].y, flakes[i].size*2, flakes[i].size*2, true)
	}
}
