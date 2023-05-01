package fonts

import (
	"embed"
	"image/color"

	"gioui.org/font"
	"gioui.org/font/opentype"
	"gioui.org/widget/material"
)

var _ embed.FS

//go:embed CandyQelling.ttf
var ttf []byte

//go:embed HYMiaoCaoW.ttf
var regular []byte

//go:embed HYShuFang_55W.ttf
var song []byte

var (
	blackFont   = font.Font{Style: font.Regular, Weight: font.Black}
	regularFont = font.Font{Style: font.Regular, Weight: font.Bold}
)
var AppColor = color.NRGBA{R: 102, G: 117, B: 127, A: 255}

func NewTheme() *material.Theme {
	candyFace, _ := opentype.Parse(ttf)
	chinese, err := opentype.Parse(regular)
	if err != nil {
		panic(err)
	}
	song, err := opentype.Parse(song)
	if err != nil {
		panic(err)
	}
	collection := []font.FontFace{
		{
			Font: regularFont,
			Face: chinese,
		},
		{
			Font: regularFont,
			Face: candyFace,
		},
		{
			Font: blackFont,
			Face: song,
		},
	}
	th := material.NewTheme(collection)
	th.Bg.R = 245
	th.Bg.G = 245
	th.Bg.B = 255
	return th
}
