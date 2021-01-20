package wasm

import (
	"syscall/js"
)

func (d *JsDoc) GetWindowSize() (int, int) {
	width := d.Document.Call("getElementsByTagName", "html").Index(0).Get("clientWidth").Int()
	height := d.Document.Call("getElementsByTagName", "html").Index(0).Get("clientHeight").Int()

	return width, height
}

// PlaySound plays a sound from the beginning
func PlaySound(sound js.Value) {
	sound.Set("currentTime", 0)
	sound.Call("play")
}

// SetVolume sets the volume for a sound
func SetVolume(sound js.Value, vol float64) {
	sound.Set("volume", vol)
}

// GetCanvasImage will load a slice of bytes and dimensions that have been extracted from
// imageExtractor and return a canvas element that can be used to draw it.
func (d JsDoc) GetCanvasImage(data []byte, width int, height int) js.Value {
	elem := d.Document.Call("createElement", "canvas")
	elem.Set("width", width)
	elem.Set("height", height)
	dst := js.Global().Get("Uint8ClampedArray").New(len(data))
	_ = js.CopyBytesToJS(dst, data)
	imageData := js.Global().Get("ImageData").New(dst, width, height)
	elem.Call("getContext", "2d").Call("putImageData", imageData, 0, 0)
	return elem
}
