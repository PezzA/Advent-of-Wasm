package wasm

import "syscall/js"

type JsDoc struct {
	Document    js.Value
	animationID js.Value
}

func NewJsDoc() JsDoc {
	doc := js.Global().Get("document")

	return JsDoc{
		Document: doc,
	}
}

func (d *JsDoc) StartAnimLoop(frame func(now float64)) {
	var renderFrameEvent js.Func

	renderFrameEvent = js.FuncOf(func(this js.Value, args []js.Value) interface{} {
		frame(args[0].Float())
		d.animationID = js.Global().Call("requestAnimationFrame", renderFrameEvent)
		return nil
	})

	d.animationID = js.Global().Call("requestAnimationFrame", renderFrameEvent)
}

func (d *JsDoc) CancelAnimLoop() {
	js.Global().Call("cancelAnimationFrame", d.animationID)
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

func (d *JsDoc) AddClass(element js.Value, className string) {
	classList := element.Get("classList")
	classList.Call("add", className)
}

func (d *JsDoc) RemoveClass(element js.Value, className string) {
	classList := element.Get("classList")
	classList.Call("remove", className)
}

func (d *JsDoc) GetElementByID(elementID string) js.Value {
	return d.Document.Call("getElementById", elementID)
}

func (d *JsDoc) CreateElement(tagName string, id string) js.Value {
	elem := d.Document.Call("createElement", tagName)

	if id != "" {
		elem.Set("id", id)
	}

	return elem
}


