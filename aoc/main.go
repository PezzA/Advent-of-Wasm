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

func main() {
	done := make(chan bool, 0)
	doc := wasm.NewJsDoc()

	_, _, title := Day202024.Entry.Describe()

	canvasDrawWidth, canvasDrawHeight := doc.GetWindowSize()

	div := doc.CreateElement("div", "aocContainer")
	doc.AddClass(div, "control-panel")

	partOneButton := doc.CreateElement("button", "")
	doc.SetInnerHTML(partOneButton, "Part One")

	partTwoButton := doc.CreateElement("button", "")
	doc.SetInnerHTML(partTwoButton, "Part Two")

	stopButton := doc.CreateElement("button", "")
	doc.SetInnerHTML(stopButton, "StopAnimation")

	doc.Document.Get("body").Call("appendChild", div)

	panelTitle := doc.CreateElement("h2", "")
	doc.SetInnerHTML(panelTitle, title)

	div.Call("appendChild", panelTitle)
	div.Call("appendChild", partOneButton)
	div.Call("appendChild", partTwoButton)
	div.Call("appendChild", stopButton)

	config := getNewConfig(canvasDrawWidth, canvasDrawHeight, part1HexSize)

	hexDrop, _, _ := createHexBackGroundLayer(config)

	drawCanvas := doc.GetOrCreateCanvas("drawCanvas", canvasDrawWidth, canvasDrawHeight, true, true)
	drawCanvas.SetFillStyle("red")
	drawCanvas.SetStrokeStyle("blue")
	drawCanvas.SetFont("10pt Arial")
	drawCanvas.SetTextAlign(wasm.TextAlignCenter)
	drawCanvas.SetTextBaseLine(wasm.TextBaseLineMiddle)

	layers := make(map[string]*wasm.JsCanvas, 0)
	layers["draw"] = drawCanvas
	layers["background"] = hexDrop

	state := getNewState(false, 1)

	frameFunc := func(now float64) {
		delta := now - state.currentTime
		state.currentTime = now
		draw(delta, layers, config, state)
		state = update(delta, state, config)
	}

	doc.AddEventListener(stopButton, "click", js.FuncOf(func(this js.Value, args []js.Value) interface{} {
		if state.paused {
			doc.SetInnerHTML(stopButton, "Pause Animation")
			doc.StartAnimLoop(frameFunc)
		} else {
			doc.SetInnerHTML(stopButton, "Resume Animation")
			doc.CancelAnimLoop()
		}

		state.paused = !state.paused
		return false
	}))

	doc.AddEventListener(partOneButton, "click", js.FuncOf(func(this js.Value, args []js.Value) interface{} {
		doc.CancelAnimLoop()
		state = getNewState(true, 1)
		config.setSizes(part1HexSize)
		layers["background"], _, _ = createHexBackGroundLayer(config)
		doc.StartAnimLoop(frameFunc)
		return false
	}))

	doc.AddEventListener(partTwoButton, "click", js.FuncOf(func(this js.Value, args []js.Value) interface{} {
		doc.CancelAnimLoop()
		state = getNewState(true, 2)
		config.setSizes(part2HexSize)
		//layers["background"], _, _ = createHexBackGroundLayer(config)
		state.Floor.FlipAllTiles(Day202024.GetData(Day202024.Entry.PuzzleInput()))
		doc.StartAnimLoop(frameFunc)

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

	if s.puzzlePart == 2 && s.stateTick > c.state2TickPoll && s.automataIteration <= s.maxIteration {
		s.stateTick = 0
		s.Floor = s.Floor.Automata()
		s.automataIteration++
		return s
	}

	// if there is going to be no update, return current state
	if s.animTick < c.animTickPoll && s.stateTick < c.stateTickPoll {
		return s
	}

	newTiles := make([]tile, 0)
	newEndTiles := make([]tile, 0)

	// update any existing tiles
	if s.animTick > c.animTickPoll || s.stateTick > c.stateTickPoll {
		s.animTick = 0

		for i := range s.tiles {
			if s.tiles[i].margin < c.hexSize {
				s.tiles[i].margin++
				newTiles = append(newTiles, s.tiles[i])
			}
		}
	}

	// see if we have some new ones to add.
	if s.puzzlePart == 1 && s.stateTick > c.stateTickPoll {
		s.stateTick = 0

		// do part on update
		if s.puzzlePart == 1 {
			var t tile

			for i := 0; i < tilesToAdd; i++ {
				if s.running {
					s, t = getNextTile(s, c)
					if t.endTile {
						newEndTiles = append(newEndTiles, t)
					} else {
						newTiles = append(newTiles, t)
					}
				}
			}
		}
	}

	s.tiles = newTiles

	return s
}

func getNextTile(s state, c config) (state, tile) {
	var t tile
	endTile := false

	// inc route counter, when at end, inc route list counter
	s.subInsIndex++
	if s.subInsIndex == len(c.routeList[s.insIndex])-1 {
		endTile = true
		s.Floor.FlipTile(c.insList[s.insIndex])
	}

	if !endTile {
		t = tile{
			Point:   common.AxialToOffset(c.routeList[s.insIndex][s.subInsIndex]),
			endTile: endTile,
		}
	}

	if endTile {
		s.subInsIndex = 0
		s.insIndex++

		if s.insIndex == len(c.routeList) {
			s.running = false
		}
	}

	return s, t
}

func draw(delta float64, layers map[string]*wasm.JsCanvas, c config, s state) {
	layers["draw"].Clear()

	if s.puzzlePart == 1 {
		layers["draw"].SetFillStyle("#000000")
	}

	if s.puzzlePart == 2 {
		layers["draw"].SetFillStyle("#FFFFFF")
		layers["draw"].DrawRectInt(-c.canvasWidth/2, -c.canvasHeight/2, c.canvasWidth, c.canvasHeight, true)
	}

	for k, v := range s.Floor {
		if s.puzzlePart == 2 {
			if v > 0 {
				if v == 1 || v == 2 {
					layers["draw"].SetFillStyle("#dddddd")
				}

				if v == 3 || v == 4 {
					layers["draw"].SetFillStyle("#bbbbbb")
				}

				if v == 5 || v == 6 {
					layers["draw"].SetFillStyle("#999999")
				}

				if v == 7 || v == 8 {
					layers["draw"].SetFillStyle("#777777")
				}

				if v == 9 || v == 10 {
					layers["draw"].SetFillStyle("#555555")
				}

				if v == 11 || v == 12 {
					layers["draw"].SetFillStyle("#333333")
				}

				if v > 12 {
					layers["draw"].SetFillStyle("#000000")
				}
			}
		}
		layers["draw"].DrawPolyLine(getCanvasPoint(common.AxialToOffset(k), 0, c), c.hex, true)

	}

	if s.puzzlePart == 1 {
		for _, tile := range s.tiles {
			layers["draw"].DrawPolyLine(getCanvasPoint(tile.Point, tile.margin, c), getHex(c.hexSize-tile.margin), false)
		}
	}
}

// hexColor returns an HTML hex-representation of c. The alpha channel is dropped
// and precision is truncated to 8 bits per channel
// https://www.reddit.com/r/golang/comments/adm6l5/do_i_really_need_a_third_party_library_to_get_hex/edibj7a?utm_source=share&utm_medium=web2x&context=3
func hexColor(c color.Color) string {
	rgba := color.RGBAModel.Convert(c).(color.RGBA)
	return fmt.Sprintf("#%.2x%.2x%.2x", rgba.R, rgba.G, rgba.B)
}

func getCanvasPoint(p common.Point, margin int, c config) common.Point {
	canvasPoint := common.Point{
		X: p.X * c.hexThreeQSize,
		Y: p.Y * c.hexFullSize,
	}

	if p.X%2 == 0 {
		canvasPoint.Y += c.hexHalfSize
	}

	if margin == 0 {
		return canvasPoint
	}
	canvasPoint.X += 2 * margin
	canvasPoint.Y += 2 * margin

	return canvasPoint
}
