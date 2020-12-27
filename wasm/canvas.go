package wasm

import (
	"syscall/js"

	"github.com/pezza/advent-of-code/common"
)

type JsCanvas struct {
	Canvas         js.Value
	Canvas2        js.Value
	Context        js.Value
	Context2       js.Value
	width          int
	height         int
	doubleBuffered bool
}

func (d *JsDoc) GetOrCreateCanvas(name string, drawWidth int, drawHeight int, doubleBuffer bool, addToDom bool, cartesian bool) JsCanvas {
	var canvas, canvas2, ctx, ctx2 js.Value
	canvas = d.Document.Call("getElementById", name)

	if canvas.IsNull() {
		canvas = d.Document.Call("createElement", "canvas")
		canvas.Set("id", name)

		if addToDom {
			d.Document.Get("body").Call("appendChild", canvas)
		}
	}
	canvas.Set("width", drawWidth)
	canvas.Set("height", drawHeight)

	ctx = canvas.Call("getContext", "2d", "{ alpha: false }")

	if cartesian && !doubleBuffer {
		ctx.Call("translate", drawWidth/2, drawHeight/2)
	}

	if doubleBuffer {
		canvas2 = d.Document.Call("createElement", "canvas")
		canvas2.Set("id", name+"2")
		canvas2.Set("width", drawWidth)
		canvas2.Set("height", drawHeight)

		ctx2 = canvas2.Call("getContext", "2d", "{ alpha: false }")

		if cartesian {
			ctx2.Call("translate", drawWidth/2, drawHeight/2)
		}
	}

	return JsCanvas{
		Canvas:         canvas,
		Canvas2:        canvas2,
		Context:        ctx,
		Context2:       ctx2,
		width:          drawWidth,
		height:         drawHeight,
		doubleBuffered: doubleBuffer,
	}
}

// ClearFrame will draw a clear frame of the entire canvas
func (c *JsCanvas) ClearFrame(x, y, w, h int) {
	if c.doubleBuffered {
		c.Context2.Call("clearRect", x, y, w, h)
	}

	c.Context.Call("clearRect", x, y, w, h)
}

func (c *JsCanvas) getDrawContext() js.Value {
	if c.doubleBuffered {
		return c.Context2
	}
	return c.Context
}

func (c *JsCanvas) LineWidth(width int) {
	c.getDrawContext().Set("lineWidth", width)
}

func (c *JsCanvas) SetFillStyle(style string) {
	c.getDrawContext().Set("fillStyle", style)
}

func (c *JsCanvas) SetStrokeStyle(style string) {
	c.getDrawContext().Set("strokeStyle", style)
}

// DrawRect draws a filled rectangle to the canvas
func (c *JsCanvas) DrawFilledRect(x, y, w, h int) {
	c.getDrawContext().Call("fillRect", x, y, w, h)
}

func (c *JsCanvas) DrawPolyLine(start common.Point, points []common.Point, fill bool) {
	c.getDrawContext().Call("beginPath")
	c.getDrawContext().Call("moveTo", start.X+points[0].X, start.Y+points[0].Y)

	for i := 0; i < len(points); i++ {
		c.getDrawContext().Call("lineTo", points[i].X+start.X, points[i].Y+start.Y)
	}

	if fill {
		c.getDrawContext().Call("fill")
	} else {
		c.getDrawContext().Call("stroke")
	}
}

func (c *JsCanvas) SetFont(font string) {
	c.getDrawContext().Set("font", font)
}

type TextAlign string

const (
	TextAlignStart  TextAlign = "start"
	TextAlignEnd    TextAlign = "end"
	TextAlignLeft   TextAlign = "left"
	TextAlignRight  TextAlign = "right"
	TextAlignCenter TextAlign = "center"
)

func (c *JsCanvas) SetTextAlign(alignment TextAlign) {
	c.getDrawContext().Set("textAlign", string(alignment))
}

type TextBaseLine string

const (
	TextBaseLineTop        TextBaseLine = "top"
	TextBaseLineBottom     TextBaseLine = "bottom"
	TextBaseLineMiddle     TextBaseLine = "middle"
	TextBaseLineAlphabetic TextBaseLine = "alphabetic"
	TextBaseLineHanging    TextBaseLine = "hanging"
)

func (c *JsCanvas) SetTextBaseLine(baseline TextBaseLine) {
	c.getDrawContext().Set("textBaseline", string(baseline))
}

func (c *JsCanvas) DrawText(text string, x int, y int, fill bool) {
	if fill {
		c.getDrawContext().Call("fillText", text, x, y)
	} else {
		c.getDrawContext().Call("strokeText", text, x, y)
	}
}

func (c *JsCanvas) DrawCanvas(t JsCanvas, x int, y int) {
	c.Context.Call("drawImage", t.Canvas, x, y)
}

func (c *JsCanvas) DrawBufferedFrame() {
	c.Context.Call("drawImage", c.Canvas2, 0, 0)
}
