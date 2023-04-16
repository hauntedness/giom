package common

import (
	"image"
	"image/color"

	"gioui.org/layout"
	"gioui.org/op/clip"
	"gioui.org/op/paint"
	"gioui.org/widget"
)

var (
	ColorBackground = color.NRGBA{R: 0xc0, G: 0xc0, B: 0xc0, A: 0xff}
	ColorRed        = color.NRGBA{R: 0xc0, G: 0x40, B: 0x40, A: 0xff}
	ColorGreen      = color.NRGBA{R: 0x40, G: 0xc0, B: 0x40, A: 0xff}
	ColorBlue       = color.NRGBA{R: 0x40, G: 0x40, B: 0xc0, A: 0xff}
)

func ColorBox(gtx layout.Context, size image.Point, color color.NRGBA) layout.Dimensions {
	defer clip.Rect{Max: size}.Push(gtx.Ops).Pop()
	paint.ColorOp{Color: color}.Add(gtx.Ops)
	paint.PaintOp{}.Add(gtx.Ops)
	return layout.Dimensions{Size: size}
}

type ActionIcon []byte

func MustIcon(icon ActionIcon) *widget.Icon {
	ic, err := widget.NewIcon(icon)
	if err != nil {
		panic(err)
	}
	return ic
}
