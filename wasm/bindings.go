// Package wasm is a wrapper for wasm that assumes the containing web page has a single canvas which is used specifically for the wasm file.  The main
// reason for having this in a separate package is to prevent editor fails when trying to understand the build targets for syscall/js.
package wasm

import (
"fmt"
"syscall/js"
)

// JsDoc is an object that holds the references to the document, canvas and drawing context.  There is global state, so creating more than one will have
// unpredictable events.
type JsDoc struct {
	Document   js.Value
	canvasElem js.Value
	TwoDCtx    js.Value
}

// Log will send some text to the console, probably don't need this, but it calls out that the syscall is handling the output
func (d *JsDoc) Log(text string) {
	fmt.Println(text)
}

// GetElementInnerHTML wraps loading an element and getting the current innerHTML
func (d *JsDoc) GetElementInnerHTML(elementID string) string {
	return d.Document.Call("getElementById", elementID).Get("innerHtml").String()
}

func (d *JsDoc) SetElementInnerHTML(elementID string, html string) {
	d.Document.Call("getElementById", elementID).Set("innerHTML", html)
}

func (d *JsDoc) SetEvent(elementID string, html string) {
	d.Document.Call("getElementById", elementID).Set("innerHTML", html)
}




// GetElementByID wraps getting an element from the DOM.  Hint, this could be an image element for use in DrawImage
func (d *JsDoc) GetElementByID(elementID string) js.Value {
	return d.Document.Call("getElementById", elementID)
}

// SetGlobalAlpha sets the canvas global alpha
func (d *JsDoc) SetGlobalAlpha(alpha float64) {
	d.TwoDCtx.Set("globalAlpha", alpha)
}

// SetCanvasSize sets the size
func (d *JsDoc) SetCanvasSize(x, y int) {
	d.canvasElem.Set("width", x)
	d.canvasElem.Set("height", y)

	canvasWidth = float64(x)
	canvasHeight = float64(y)
}


// DrawImage will draw an image to specified coordinates, wants sprite offsets by default
func (d *JsDoc) DrawImage(img js.Value, sX, sY, sW, sH, dX, dY, dW, dH int) {
	d.TwoDCtx.Call("drawImage", img, sX, sY, sW, sH, dX, dY, dW, dH)
}

// ClearFrame will draw a clear frame of the entire canvas
func (d *JsDoc) ClearFrame(x, y, w, h int) {
	d.TwoDCtx.Call("clearRect", x, y, w, h)
}

// DrawText draws text to the canvas
func (d *JsDoc) DrawText(text, font, fillStyle, textAlign, textBaseLine string, x, y int) {
	d.TwoDCtx.Set("font", font)
	d.TwoDCtx.Set("fillStyle", fillStyle)
	d.TwoDCtx.Set("textAlign", textAlign)
	d.TwoDCtx.Set("textBaseline ", textBaseLine)
	d.TwoDCtx.Call("fillText", text, x, y)
}

// DrawRect draws a filled rectangle to the canvas
func (d *JsDoc) DrawRect(x, y, w, h int, fillStyle string) {
	d.TwoDCtx.Set("fillStyle", fillStyle)
	d.TwoDCtx.Call("fillRect", x, y, w, h)
}

// StrokeRect draws an unfilled rectangle to the canvas
func (d *JsDoc) StrokeRect(x, y, w, h int, strokeStyle string) {
	d.TwoDCtx.Set("strokeStyle", strokeStyle)
	d.TwoDCtx.Call("strokeRect", x, y, w, h)
}

// PlaySound plays a sound from the beginning
func (d *JsDoc) PlaySound(sound js.Value) {
	sound.Set("currentTime", 0)
	sound.Call("play")
}

// SetVolume sets the volume for a sound
func (d *JsDoc) SetVolume(sound js.Value, vol float64) {
	sound.Set("volume", vol)
}