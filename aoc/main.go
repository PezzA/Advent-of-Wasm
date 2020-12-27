package main

import (
	"fmt"

	"github.com/pezza/advent-of-code/common"
	"github.com/pezza/advent-of-wasm/wasm"
)

var currentTime float64

func getHex(size int) []common.Point {
	return []common.Point{
		{X: 1 * size, Y: 0 * size},
		{X: 3 * size, Y: 0 * size},
		{X: 4 * size, Y: 2 * size},
		{X: 3 * size, Y: 4 * size},
		{X: 1 * size, Y: 4 * size},
		{X: 0 * size, Y: 2 * size},
		{X: 1 * size, Y: size * 0},
	}
}

func createHexBackGroundLayer(width, height, size int) (wasm.JsCanvas, int, int) {
	doc := wasm.NewJsDoc()
	canvas := doc.GetOrCreateCanvas("hexBackGround", width, height, false, false, true)

	canvas.SetFont("14pt Arial")
	canvas.SetTextAlign(wasm.TextAlignCenter)
	canvas.SetTextBaseLine(wasm.TextBaseLineMiddle)
	hex := getHex(size)

	canvas.SetStrokeStyle("white")

	hexTileWidth := 4 * size
	halfWidth := 2 * size
	threeQuarterWidth := (hexTileWidth / 4) * 3

	xCells := (width / threeQuarterWidth) + 4
	yCells := height / threeQuarterWidth

	xCellsHalf, yCellsHalf := xCells/2, yCells/2

	for x := -xCellsHalf; x < xCellsHalf; x++ {
		for y := -yCellsHalf; y < yCellsHalf; y++ {
			plotPoint := common.Point{X: x, Y: y}

			canvasPoint := common.Point{
				X: plotPoint.X * threeQuarterWidth,
				Y: plotPoint.Y * hexTileWidth,
			}

			if x%2 == 0 {
				canvasPoint.Y += halfWidth
			}

			canvas.DrawPolyLine(canvasPoint, hex, false)
		}
	}

	return canvas, xCells, yCells
}

func main() {
	done := make(chan bool, 0)
	doc := wasm.NewJsDoc()

	canvasDrawWidth, canvasDrawHeight := doc.GetWindowSize()

	size := 10

	hexDrop, _, _ := createHexBackGroundLayer(canvasDrawWidth, canvasDrawHeight, size)

	drawCanvas := doc.GetOrCreateCanvas(
		"drawCanvas",
		canvasDrawWidth,
		canvasDrawHeight,
		true,
		true,
		true,
	)

	uiCanvas := doc.GetOrCreateCanvas(
		"uiCanvas",
		canvasDrawWidth,
		canvasDrawHeight,
		true,
		true,
		false,
	)

	uiCanvas.SetFillStyle("blue")
	uiCanvas.SetFont("18px Consolas")

	drawCanvas.SetFillStyle("red")
	drawCanvas.SetStrokeStyle("red")
	drawCanvas.SetFont("10pt Arial")
	drawCanvas.SetTextAlign(wasm.TextAlignCenter)
	drawCanvas.SetTextBaseLine(wasm.TextBaseLineMiddle)

	var s state
	doc.StartAnimLoop(func(now float64) {
		delta := now - currentTime
		currentTime = now

		draw(delta, drawCanvas, hexDrop, uiCanvas, s, canvasDrawWidth, canvasDrawHeight)
		s = update(delta, s)
	})

	<-done
}

func update(delta float64, s state) state {

	return s
}

func draw(delta float64, drawCanvas wasm.JsCanvas, hexDrop wasm.JsCanvas, uiCanvas wasm.JsCanvas, s state, width int, height int) {
	drawCanvas.ClearFrame(0, 0, width, height)
	drawCanvas.DrawCanvas(hexDrop, 0, 0)
	drawCanvas.DrawBufferedFrame()

	uiCanvas.ClearFrame(0, 0, width, height)
	uiCanvas.DrawText(fmt.Sprintf("Delta: %v", int(delta)), 50, 50, true)
	uiCanvas.DrawBufferedFrame()
}

type state struct {
	Current      []common.Point
	OpeningTiles []common.Point
}

func getCanvasPoint(p common.Point, margin int, size int) common.Point {
	hexTileWidth := 4 * size

	threeQuarterWidth := (hexTileWidth / 4) * 3

	canvasPoint := common.Point{
		X: p.X * threeQuarterWidth,
		Y: p.Y * hexTileWidth,
	}

	if p.X%2 == 0 {
		canvasPoint.Y += size * 2
	}

	canvasPoint.X += 2 * margin
	canvasPoint.Y += 2 * margin

	return canvasPoint
}
