package wasm

import "syscall/js"

type JsDoc struct {
	Document   js.Value
}

var frameCallback func(now float64)
var renderFrameEvt js.Func

func NewJsDoc() JsDoc {
	doc := js.Global().Get("document")

	return JsDoc{
		Document: doc,
	}
}

// StartAnimLoop starts the requestAnimationFrame loop.
func (d *JsDoc) StartAnimLoop(frame func(now float64)) {
	frameCallback = frame
	renderFrameEvt = js.FuncOf(renderFrame)
	js.Global().Call("requestAnimationFrame", renderFrameEvt)
}

func renderFrame(this js.Value, args []js.Value) interface{} {
	frameCallback(args[0].Float())
	js.Global().Call("requestAnimationFrame", renderFrameEvt)
	return nil
}

// StartAnimLoop starts the requestAnimationFrame loop.
func (d *JsDoc) AddEventListener(element string, event string, handlerFunc js.Func) {
	elem:= d.GetElementByID(element)
	elem.Call("addEventListener", event, handlerFunc)
}

// GetElementInnerHTML wraps loading an element and getting the current innerHTML
func (d *JsDoc) GetElementInnerHTML(elementID string) string {
	return d.Document.Call("getElementById", elementID).Get("innerHtml").String()
}

func (d *JsDoc) SetElementInnerHTML(elementID string, html string) {
	d.Document.Call("getElementById", elementID).Set("innerHTML", html)
}

// GetElementByID wraps getting an element from the DOM.  Hint, this could be an image element for use in DrawImage
func (d *JsDoc) GetElementByID(elementID string) js.Value {
	return d.Document.Call("getElementById", elementID)
}
