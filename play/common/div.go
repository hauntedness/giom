package common

import (
	"gioui.org/layout"
	"gioui.org/widget"
)

type Div struct {
	widget.Border
	layout.Direction
	layout.Inset
	NoDirection bool
}

func (b Div) Layout(gtx layout.Context, w layout.Widget) layout.Dimensions {
	dim := b.Border.Layout(
		gtx,
		func(gtx layout.Context) layout.Dimensions {
			if b.NoDirection {
				return b.Inset.Layout(gtx, w)
			}
			return b.Inset.Layout(
				gtx,
				func(gtx layout.Context) layout.Dimensions {
					return b.Direction.Layout(gtx, w)
				},
			)
		},
	)
	return dim
}
