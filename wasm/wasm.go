package wasm

import (
	"syscall/js"
)

var clickCallback func(x, y int)
var mouseMoveCallback func(x, y int)
var frameCallback func(now float64)

var offSetLeft, offSetTop float64
var canvasWidth, canvasHeight float64

var mouseMoveEvt, renderFrameEvt, canvasClickEvt js.Func

// MousePos gives the x,y coordinates of the mouse over the canvas, these values will be constrained to the
// edges of the canvas
var mousePos [2]float64

// NewJsDoc returns a new JsDoc.
func NewJsDoc(canvasName string) JsDoc {
	doc := js.Global().Get("document")
	canvas := doc.Call("getElementById", canvasName)
	ctx := canvas.Call("getContext", "2d")

	offSetLeft = canvas.Get("offsetLeft").Float()
	offSetTop = canvas.Get("offsetTop").Float()

	jsDoc := JsDoc{
		document:   doc,
		canvasElem: canvas,
		TwoDCtx:    ctx,
	}
	return jsDoc
}

// StartAnimLoop starts the requestAnimationFrame loop.
func (d *JsDoc) StartAnimLoop(frame func(now float64)) {
	frameCallback = frame
	renderFrameEvt = js.FuncOf(renderFrame)
	js.Global().Call("requestAnimationFrame", renderFrameEvt)
}

// InitEvents will bind the click, mousemove
func (d *JsDoc) InitEvents(click func(x, y int), mouseMove func(x, y int)) {
	clickCallback = click

	mouseMoveCallback = mouseMove

	mouseMoveEvt = js.FuncOf(mouseMoveHandler)
	d.document.Call("addEventListener", "mousemove", mouseMoveEvt)

	canvasClickEvt = js.FuncOf(canvasClick)
	d.canvasElem.Call("addEventListener", "click", canvasClickEvt)
}

// ReleaseEvents removes the event listeners.
func ReleaseEvents() {
	mouseMoveEvt.Release()
	renderFrameEvt.Release()
	canvasClickEvt.Release()
}

func canvasClick(this js.Value, args []js.Value) interface{} {
	clickCallback(int(mousePos[0]), int(mousePos[1]))
	return nil
}

func renderFrame(this js.Value, args []js.Value) interface{} {
	frameCallback(args[0].Float())
	js.Global().Call("requestAnimationFrame", renderFrameEvt)
	return nil
}

func mouseMoveHandler(this js.Value, args []js.Value) interface{} {
	e := args[0]

	mousePos[0] = e.Get("clientX").Float() - offSetLeft
	mousePos[1] = e.Get("clientY").Float() - offSetTop

	if mousePos[0] < 0 {
		mousePos[0] = 0
	}
	if mousePos[1] < 0 {
		mousePos[1] = 0
	}

	if mousePos[0] > canvasWidth {
		mousePos[0] = canvasWidth
	}

	if mousePos[1] > canvasHeight {
		mousePos[1] = canvasHeight
	}

	mouseMoveCallback(int(mousePos[0]), int(mousePos[1]))
	return nil
}
