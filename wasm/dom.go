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

func (d *JsDoc) AddEventListener(element js.Value, event string, handlerFunc js.Func) js.Value {
	return element.Call("addEventListener", event, handlerFunc)
}

func (d *JsDoc) GetInnerHTMLById(elementID string) string {
	return d.Document.Call("getElementById", elementID).Get("innerHtml").String()
}

func (d *JsDoc) GetElementInnerHTML(elementID string) string {
	return d.Document.Call("getElementById", elementID).Get("innerHtml").String()
}

func (d *JsDoc) SetInnerHTMLById(elementID string, html string) {
	d.Document.Call("getElementById", elementID).Set("innerHTML", html)
}

func (d *JsDoc) SetValue(element js.Value, val interface{}) {
	element.Set("value", val)
}

func (d *JsDoc) SetInnerHTML(element js.Value, html string) {
	element.Set("innerHTML", html)
}

func (d *JsDoc) GetElementByID(elementID string) js.Value {
	return d.Document.Call("getElementById", elementID)
}
