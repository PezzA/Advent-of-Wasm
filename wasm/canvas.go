package wasm

import (
	"syscall/js"

	"github.com/pezza/advent-of-code/common"
)

type JsCanvas struct {
	Canvas    js.Value
	Context   js.Value
	width     int
	height    int
	cartesian bool
}

func (d *JsDoc) GetOrCreateCanvas(name string, drawWidth int, drawHeight int, addToDom bool, cartesian bool) *JsCanvas {
	var canvas, ctx js.Value
	canvas = d.Document.Call("getElementById", name)

	if canvas.IsNull() {
		canvas = d.CreateElement("canvas", name)

		if addToDom {
			d.Document.Get("body").Call("appendChild", canvas)
		}
	}
	canvas.Set("width", drawWidth)
	canvas.Set("height", drawHeight)

	ctx = canvas.Call("getContext", "2d", "{ alpha: false }")

	if cartesian {
		ctx.Call("translate", drawWidth/2, drawHeight/2)
	}

	return &JsCanvas{
		Canvas:    canvas,
		Context:   ctx,
		width:     drawWidth,
		height:    drawHeight,
		cartesian: cartesian,
	}
}

func (c *JsCanvas) ClearArea(x, y, w, h int) {
	c.Context.Call("clearRect", x, y, w, h)
}

func (c *JsCanvas) Clear() {
	x, y := 0, 0

	if c.cartesian {
		x, y = -c.width/2, -c.height/2
	}

	c.Context.Call("clearRect", x, y, c.width, c.height)
}

func (c *JsCanvas) LineWidth(width int) {
	c.Context.Set("lineWidth", width)
}

func (c *JsCanvas) SetFillStyle(style string) {
	c.Context.Set("fillStyle", style)
}

func (c *JsCanvas) SetStrokeStyle(style string) {
	c.Context.Set("strokeStyle", style)
}

func (c *JsCanvas) DrawRect(x, y, w, h int, fill bool) {
	if fill {
		c.Context.Call("fillRect", x, y, w, h)
	} else {
		c.Context.Call("strokeRect", x, y, w, h)
	}
}

func (c *JsCanvas) DrawPolyLine(start common.Point, points []common.Point, fill bool) {
	c.Context.Call("beginPath")
	c.Context.Call("moveTo", start.X+points[0].X, start.Y+points[0].Y)

	for i := 0; i < len(points); i++ {
		c.Context.Call("lineTo", points[i].X+start.X, points[i].Y+start.Y)
	}

	if fill {
		c.Context.Call("fill")
	} else {
		c.Context.Call("stroke")
	}
}

func (c *JsCanvas) SetFont(font string) {
	c.Context.Set("font", font)
}

func (c *JsCanvas) SetTextAlign(alignment TextAlign) {
	c.Context.Set("textAlign", string(alignment))
}

func (c *JsCanvas) SetTextBaseLine(baseline TextBaseLine) {
	c.Context.Set("textBaseline", string(baseline))
}

func (c *JsCanvas) DrawText(text string, x int, y int, fill bool) {
	if fill {
		c.Context.Call("fillText", text, x, y)
	} else {
		c.Context.Call("strokeText", text, x, y)
	}
}

func (c *JsCanvas) CopyCanvas(t JsCanvas, x int, y int) {
	c.Context.Call("drawImage", t.Canvas, x, y)
}
