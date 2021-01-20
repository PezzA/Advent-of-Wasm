package wasm

import (
	"syscall/js"

	"github.com/pezza/advent-of-code/common"
)

type JsCanvas struct {
	Canvas    js.Value
	Context   js.Value
	Width     int
	Height    int
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
		Width:     drawWidth,
		Height:    drawHeight,
		cartesian: cartesian,
	}
}

func (c *JsCanvas) GetImageData() []byte {
	jsImageData := c.Context.Call("getImageData", 0, 0, c.Width, c.Height)
	imageData := make([]byte, jsImageData.Get("height").Int()*jsImageData.Get("width").Int()*4)
	_ = js.CopyBytesToGo(imageData, jsImageData.Get("data"))
	return imageData
}

func (c *JsCanvas) GetBlankBytes() []byte {
	return make([]byte, c.Width*c.Height*4)
}

func (c *JsCanvas) ClearArea(x, y, w, h int) {
	c.Context.Call("clearRect", x, y, w, h)
}

func (c *JsCanvas) Clear() {
	x, y := 0, 0

	if c.cartesian {
		x, y = -c.Width/2, -c.Height/2
	}

	c.Context.Call("clearRect", x, y, c.Width, c.Height)
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

func (c *JsCanvas) PutImageData(data []byte) {
	dst := js.Global().Get("Uint8ClampedArray").New(len(data))
	_ = js.CopyBytesToJS(dst, data)
	id := js.Global().Get("ImageData").New(dst, c.Width, c.Height)
	c.Context.Call("putImageData", id, 0, 0)
}

func (c *JsCanvas) DrawRect(x, y, w, h float64, fill bool) {
	if fill {
		c.Context.Call("fillRect", x, y, w, h)
	} else {
		c.Context.Call("strokeRect", x, y, w, h)
	}
}

func (c *JsCanvas) DrawRectInt(x, y, w, h int, fill bool) {
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

func (c *JsCanvas) DrawTextOnly(text string, x int, y int, fill bool) {
	if fill {
		c.Context.Call("fillText", text, x, y)
	} else {
		c.Context.Call("strokeText", text, x, y)
	}
}

func (c *JsCanvas) DrawText(text string, x int, y int, fill bool, font string, style string, align TextAlign, base TextBaseLine) {
	c.SetFont(font)
	c.SetTextAlign(align)
	c.SetTextBaseLine(base)

	if fill {
		c.SetFillStyle(style)
		c.Context.Call("fillText", text, x, y)
	} else {
		c.SetStrokeStyle(style)
		c.Context.Call("strokeText", text, x, y)
	}
}

func (c *JsCanvas) CopyCanvas(t JsCanvas, x int, y int) {
	c.Context.Call("drawImage", t.Canvas, x, y)
}

func (c *JsCanvas) SetGlobalAlpha(alpha float64) {
	c.Context.Set("globalAlpha", alpha)
}

func (c *JsCanvas) DrawImage(img js.Value, sX, sY, sW, sH, dX, dY, dW, dH int) {
	c.Context.Call("drawImage", img, sX, sY, sW, sH, dX, dY, dW, dH)
}
