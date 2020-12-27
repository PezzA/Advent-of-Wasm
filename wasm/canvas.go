package wasm

import (
	"syscall/js"
)

type JsCanvas struct {
	Canvas  js.Value
	Canvas2  js.Value
	Context js.Value
	Context2 js.Value
	offsetLeft int
	offsetRight int
}

func (d *JsDoc) GetOrCreateCanvas(name string, drawWidth int, drawHeight int) JsCanvas {
	var canvas js.Value
	canvas = d.Document.Call("getElementById", name)

	if canvas.IsNull() {
		canvas = d.Document.Call("createElement", "canvas")
		canvas.Set("id", name)
		d.Document.Get("body").Call("appendChild", canvas)
	}

	ctx := canvas.Call("getContext", "2d", "{ alpha: false }")

	var secondaryCanvas = d.Document.Call("createElement", "canvas")

	ctx2 := secondaryCanvas.Call("getContext","2d","{ alpha: false }")

	offSetLeft := canvas.Get("offsetLeft").Float()
	offSetTop := canvas.Get("offsetTop").Float()

	canvas.Set("width", drawWidth)
	canvas.Set("height", drawHeight)

	secondaryCanvas.Set("width", drawWidth)
	secondaryCanvas.Set("height", drawHeight)

	return JsCanvas{ Canvas: canvas, Canvas2: secondaryCanvas, Context: ctx, Context2: ctx2, offsetLeft: int(offSetLeft), offsetRight: int(offSetTop)}
}

// ClearFrame will draw a clear frame of the entire canvas
func (c *JsCanvas) ClearFrame(x, y, w, h int) {
	c.Context2.Call("clearRect", x, y, w, h)
	c.Context.Call("clearRect", x, y, w, h)
}

func (c *JsCanvas) SetFillStyle(style string) {
	c.Context2.Set("fillStyle", style)
}


// DrawRect draws a filled rectangle to the canvas
func (c *JsCanvas) DrawFilledRect(x, y, w, h int) {
	c.Context2.Call("fillRect", x, y, w, h)
}

func (c *JsCanvas) DrawFrame() {
	c.Context.Call("drawImage", c.Canvas2,0,0)
}

