package view

import (
	"gioui.org/layout"
	"gioui.org/text"
	"gioui.org/unit"
	"gioui.org/widget/material"
	"github.com/hauntedness/giom/assets/fonts"
)

type Greetings struct {
	*material.Theme
}

func NewGreetings(theme *material.Theme) Greetings {
	return Greetings{Theme: theme}
}

func (cp *Greetings) Layout(gtx layout.Context) layout.Dimensions {
	if cp.Theme == nil {
		cp.Theme = fonts.NewTheme()
	}

	flex := layout.Flex{Axis: layout.Vertical, Spacing: layout.SpaceSides, Alignment: layout.Middle}
	gtx.Constraints.Min.Y = gtx.Constraints.Max.Y
	return flex.Layout(gtx,
		layout.Rigid(func(gtx layout.Context) layout.Dimensions {
			return DrawProtonetImageCenter(gtx, cp.Theme)
		}),
		layout.Rigid(layout.Spacer{Height: unit.Dp(16)}.Layout),
		layout.Rigid(func(gtx layout.Context) layout.Dimensions {
			return layout.Center.Layout(gtx, func(gtx layout.Context) layout.Dimensions {
				body := material.Body1(cp.Theme, "Welcome to Protonet !")
				body.Alignment = text.Middle
				body.Font.Weight = text.Black
				return body.Layout(gtx)
			})
		}),
		layout.Rigid(layout.Spacer{Height: unit.Dp(16)}.Layout),
	)
}
