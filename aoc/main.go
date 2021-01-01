package main

import (
	"fmt"
	"image/color"
	"syscall/js"

	"github.com/pezza/advent-of-code/2020/Day202024"

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
	canvas := doc.GetOrCreateCanvas("hexBackGround", c.canvasWidth, c.canvasHeight, true, true)

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
	margin  int
	color   color.RGBA
	remove  bool
	endTile bool
}

const animTickPoll = 20
const stateTickPoll = 80

type state struct {
	tiles       []tile
	currentTime float64
	animTick    float64
	stateTick   float64
	insIndex    int
	subInsIndex int
	running     bool
}

func getNewState(running bool) state {
	return state{
		tiles:       make([]tile, 0),
		animTick:    0,
		stateTick:   0,
		currentTime: 0,
		insIndex:    0,
		subInsIndex: 0,
		running:     running,
	}
}

type config struct {
	hexSize      int
	canvasWidth  int
	canvasHeight int
	hex          []common.Point
	insList      [][]string
}

func main() {
	done := make(chan bool, 0)
	doc := wasm.NewJsDoc()

	_, _, title := Day202024.Entry.Describe()

	canvasDrawWidth, canvasDrawHeight := doc.GetWindowSize()

	div := doc.CreateElement("div", "aocContainer")
	doc.AddClass(div, "control-panel")

	partOneButton := doc.CreateElement("button", "")
	doc.SetInnerHTML(partOneButton, "Part One")

	doc.Document.Get("body").Call("appendChild", div)

	panelTitle := doc.CreateElement("h2", "")
	doc.SetInnerHTML(panelTitle, title)

	div.Call("appendChild", panelTitle)
	div.Call("appendChild", partOneButton)

	size := 8
	config := config{
		size,
		canvasDrawWidth,
		canvasDrawHeight,
		getHex(size),
		Day202024.GetData(Day202024.Entry.PuzzleInput()),
	}

	hexDrop, _, _ := createHexBackGroundLayer(config)

	drawCanvas := doc.GetOrCreateCanvas(
		"drawCanvas",
		canvasDrawWidth,
		canvasDrawHeight,
		true,
		true,
	)

	uiCanvas := doc.GetOrCreateCanvas(
		"uiCanvas",
		canvasDrawWidth,
		canvasDrawHeight,
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

	state := getNewState(false)

	doc.StartAnimLoop(func(now float64) {
		delta := now - state.currentTime
		state.currentTime = now

		draw(delta, layers, config, state)
		state = update(delta, state, config)
	})

	doc.AddEventListener(partOneButton, "click", js.FuncOf(func(this js.Value, args []js.Value) interface{} {
		state = getNewState(true)
		return false
	}))
	<-done
}

func update(delta float64, s state, c config) state {

	if !s.running {
		return s
	}

	s.animTick += delta
	s.stateTick += delta

	newTiles := make([]tile, 0)

	// update any existing tiles

	stateUpdate := false
	if s.animTick > animTickPoll || s.stateTick > stateTickPoll {
		s.animTick = 0
		stateUpdate = true

		for i := range s.tiles {
			if !s.tiles[i].endTile {
				if s.tiles[i].margin < c.hexSize {
					s.tiles[i].margin++
				}
				if s.tiles[i].margin < c.hexSize {
					newTiles = append(newTiles, s.tiles[i])
				}
			} else {
				newTiles = append(newTiles, s.tiles[i])
			}
		}
	}

	// see if we have some new ones to add.
	if s.stateTick > stateTickPoll {
		s.stateTick = 0
		stateUpdate = true

		if s.insIndex < len(c.insList) {
			tilesToAdd := Day202024.GetMoveList(c.insList[s.insIndex])

			r, g, b := uint8(0), uint8(0), uint8(0)

			if s.insIndex%3 == 0 {
				r = 255
			}

			if s.insIndex%3 == 1 {
				g = 255
			}

			if s.insIndex%3 == 2 {
				b = 255
			}

			for index, move := range tilesToAdd {
				endTile := index == len(tilesToAdd)-1

				newCol := color.RGBA{R: r, G: g, B: b}
				margin := 0
				if endTile {
					newCol = color.RGBA{B: 0}
					margin = 0
				}

				newTiles = append(newTiles, tile{
					Point:   common.AxialToOffset(move),
					remove:  false,
					endTile: endTile,
					margin:  margin,
					color:   newCol,
				})
			}

			s.insIndex++
		}
	}

	if stateUpdate {
		s.tiles = newTiles
	}
	return s
}

func draw(delta float64, layers map[string]*wasm.JsCanvas, c config, s state) {
	layers["draw"].Clear()

	for _, tile := range s.tiles {
		layers["draw"].SetFillStyle(hexColor(tile.color))
		layers["draw"].DrawPolyLine(getCanvasPoint(tile.Point, tile.margin, c.hexSize), getHex(c.hexSize-tile.margin), true)
	}

	layers["ui"].Clear()
	layers["ui"].DrawText(fmt.Sprintf("Delta: %v", int(delta)), 50, 50, true)
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
