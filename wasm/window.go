package wasm

func (d *JsDoc) GetWindowSize() (int, int) {
	width := d.Document.Call("getElementsByTagName", "html").Index(0).Get("clientWidth").Int()
	height := d.Document.Call("getElementsByTagName", "html").Index(0).Get("clientHeight").Int()

	return width, height
}