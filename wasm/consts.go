package wasm

type TextBaseLine string

const (
	TextBaseLineTop        TextBaseLine = "top"
	TextBaseLineBottom     TextBaseLine = "bottom"
	TextBaseLineMiddle     TextBaseLine = "middle"
	TextBaseLineAlphabetic TextBaseLine = "alphabetic"
	TextBaseLineHanging    TextBaseLine = "hanging"
)

type TextAlign string

const (
	TextAlignStart  TextAlign = "start"
	TextAlignEnd    TextAlign = "end"
	TextAlignLeft   TextAlign = "left"
	TextAlignRight  TextAlign = "right"
	TextAlignCenter TextAlign = "center"
)
