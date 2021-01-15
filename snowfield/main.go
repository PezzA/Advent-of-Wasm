package main

import (
	"fmt"
	"math/rand"
	"strconv"
	"syscall/js"
	"time"

	"github.com/pezza/advent-of-wasm/wasm"
)

const frameTolerance = 100

func main() {
	done := make(chan bool, 0)

	rand.Seed(time.Now().UnixNano())

	doc := wasm.NewJsDoc()
	flakeCountSlider := doc.GetElementByID("flakecount")
	flakeCountLabel := doc.GetElementByID("flakecount-value")

	flakeSpeedSlider := doc.GetElementByID("flakespeed")
	flakeSpeedLabel := doc.GetElementByID("flakespeed-value")

	flakeCount := 500
	flakeSpeed := float64(2)

	doc.SetInnerHTML(flakeCountLabel, strconv.Itoa(flakeCount))
	doc.SetInnerHTML(flakeSpeedLabel, fmt.Sprintf("%.2f", flakeSpeed))

	canvasDrawWidth, canvasDrawHeight := doc.GetWindowSize()

	canvasDrawWidth, canvasDrawHeight = canvasDrawWidth/2, canvasDrawHeight/2

	fmt.Println(canvasDrawWidth, canvasDrawHeight)
	flakes = adjustFlakes(flakeCount, make(snowField, 0), canvasDrawWidth, canvasDrawHeight)
	canvas := doc.GetOrCreateCanvas("canv", canvasDrawWidth, canvasDrawHeight, true, false)

	doc.AddEventListener(flakeCountSlider, "input", js.FuncOf(func(this js.Value, args []js.Value) interface{} {
		newFlakeCount, _ := strconv.Atoi(this.Get("value").String())
		flakes = adjustFlakes(newFlakeCount, flakes, canvasDrawWidth, canvasDrawHeight)
		doc.SetInnerHTML(flakeCountLabel, strconv.Itoa(newFlakeCount))
		return false
	}))

	doc.AddEventListener(flakeSpeedSlider, "input", js.FuncOf(func(this js.Value, args []js.Value) interface{} {
		flakeSpeed, _ = strconv.ParseFloat(this.Get("value").String(), 64)
		doc.SetInnerHTML(flakeSpeedLabel, fmt.Sprintf("%.2f", flakeSpeed))
		return false
	}))

	var currentTime float64

	doc.StartAnimLoop(func(now float64) {
		delta := now - currentTime
		currentTime = now

		if delta > frameTolerance {
			delta = float64(int(delta) % frameTolerance)
		}

		update(delta, canvasDrawHeight, flakeSpeed)
		draw(delta, canvas)
	})

	<-done
}

func update(delta float64, canvasDrawHeight int, speed float64) {
	flakes.update(delta, canvasDrawHeight, speed)

}

func draw(delta float64, canvas *wasm.JsCanvas) {
	frame := canvas.GetBlankBytes()
	for i := range flakes {
		for x := float64(0); x < flakes[i].drawSize; x++ {
			fx := int(flakes[i].x + x)

			if fx >= canvas.Width {
				break
			}

			for y := float64(0); y < flakes[i].drawSize; y++ {
				fy := int(flakes[i].y + y)
				if fy >= canvas.Height {
					break
				}

				pix := 4 * (fy*canvas.Width + fx)

				frame[pix] = flakes[i].col.R
				frame[pix+1] = flakes[i].col.G
				frame[pix+2] = flakes[i].col.B
				frame[pix+3] = 255
			}
		}
	}

	canvas.PutImageData(frame)
}
