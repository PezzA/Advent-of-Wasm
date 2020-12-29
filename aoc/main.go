package main

import (
	"fmt"
	"image/color"

	"github.com/pezza/advent-of-code/common"
	"github.com/pezza/advent-of-wasm/wasm"
)

func getHex(size int) []common.Point {
	return []common.Point{
		{X: size, Y: 0},
		{X: 3 * size, Y: 0},
		{X: 4 * size, Y: 2 * size},
		{X: 3 * size, Y: 4 * size},
		{X: size, Y: 4 * size},
		{X: 0, Y: 2 * size},
		{X: size, Y: 0},
	}
}

func createHexBackGroundLayer(c config) (*wasm.JsCanvas, int, int) {
	doc := wasm.NewJsDoc()
	canvas := doc.GetOrCreateCanvas("hexBackGround", c.canvasWidth, c.canvasHeight, false, true, true)

	canvas.SetFont("14pt Arial")
	canvas.SetTextAlign(wasm.TextAlignCenter)
	canvas.SetTextBaseLine(wasm.TextBaseLineMiddle)
	hex := getHex(c.hexSize)

	canvas.SetStrokeStyle("white")

	hexTileWidth := 4 * c.hexSize
	halfWidth := 2 * c.hexSize
	threeQuarterWidth := (hexTileWidth / 4) * 3

	xCells := (c.canvasWidth / threeQuarterWidth) + 4
	yCells := c.canvasHeight / threeQuarterWidth

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

type tile struct {
	common.Point
	margin int
	color  color.RGBA
}
type state struct {
	tiles       []tile
	currentTime float64
}

type config struct {
	hexSize      int
	canvasWidth  int
	canvasHeight int
}

func main() {
	done := make(chan bool, 0)
	doc := wasm.NewJsDoc()

	canvasDrawWidth, canvasDrawHeight := doc.GetWindowSize()

	config := config{
		10, canvasDrawWidth, canvasDrawHeight,
	}

	hexDrop, _, _ := createHexBackGroundLayer(config)

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

	layers := make(map[string]*wasm.JsCanvas, 0)

	layers["ui"] = uiCanvas
	layers["draw"] = drawCanvas
	layers["background"] = hexDrop

	state := state{
		tiles: make([]tile, 0),
	}

	doc.StartAnimLoop(func(now float64) {
		delta := now - state.currentTime
		state.currentTime = now

		draw(delta, layers, config, state)
		state = update(delta, state)
	})

	<-done
}

func update(delta float64, s state) state {

	return s
}

func draw(delta float64, layers map[string]*wasm.JsCanvas, c config, s state) {
	layers["draw"].ClearFrame(0, 0, c.canvasWidth, c.canvasHeight)
	//layers["draw"].DrawCanvas(*layers["background"], 0, 0)
	layers["draw"].DrawBufferedFrame()

	layers["ui"].ClearFrame(0, 0, c.canvasWidth, c.canvasHeight)
	layers["ui"].DrawText(fmt.Sprintf("Delta: %v", int(delta)), 50, 50, true)
	layers["ui"].DrawBufferedFrame()
}

// hexColor returns an HTML hex-representation of c. The alpha channel is dropped
// and precision is truncated to 8 bits per channel
// https://www.reddit.com/r/golang/comments/adm6l5/do_i_really_need_a_third_party_library_to_get_hex/edibj7a?utm_source=share&utm_medium=web2x&context=3
func hexColor(c color.Color) string {
	rgba := color.RGBAModel.Convert(c).(color.RGBA)
	return fmt.Sprintf("#%.2x%.2x%.2x", rgba.R, rgba.G, rgba.B)
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
