package main

import (
	"fmt"
	"math/rand"
	"strconv"
	"syscall/js"

	"github.com/pezza/wasm"
)

type flake struct {
	x     int
	y     int
	speed int
	style string
}

var doc wasm.DomDocument
var canvas wasm.Canvas

var flakes []flake

var canvasDrawWidth, canvasDrawHeight = 800, 600

func main() {
	done := make(chan bool, 0)

	doc = wasm.GetJSDocument()
	canvas = doc.GetCanvas("canv")

	doc.AddEventListener("flakecount", "input", js.FuncOf(countHandlerfunc))

	flakeCount := 250
	flakeCount, _ = strconv.Atoi(doc.GetElementInnerHTML("flakecount-value"))

	flakes = createFlakes(flakeCount)

	doc.StartAnimLoop(frame)

	<-done
}

func createFlakes(flakeCount int) []flake {
	flakeArray := make([]flake, flakeCount)

	for index := range flakeArray {
		flakeArray[index].x = rand.Intn(canvasDrawWidth)
		flakeArray[index].y = rand.Intn(canvasDrawHeight)

		speed := rand.Intn(5) + 1
		flakeArray[index].speed = speed

		style := "#000000"

		switch speed {
		case 1:
			style = "#777777"
		case 2:
			style = "#999999"
		case 3:
			style = "#aaaaaa"
		case 4:
			style = "#bbbbbb"
		case 5:
			style = "#eeeeee"
		}

		flakeArray[index].style = style
	}

	return flakeArray
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

func resize() {
	fmt.Println("resized!")
}

var currentTime float64

func frame(now float64) {
	delta := now - currentTime
	currentTime = now
	canvas.ClearFrame(0, 0, canvasDrawWidth, canvasDrawHeight)
	for i := range flakes {
		canvas.DrawRect(flakes[i].x, flakes[i].y, flakes[i].speed, flakes[i].speed, flakes[i].style)
		flakes[i].y += int(float64(flakes[i].speed) * (delta / 20))
		if flakes[i].y > canvasDrawHeight {
			flakes[i].y -= canvasDrawHeight
		}
	}
}
