package view

import (
	"image/color"

	"gioui.org/layout"
	"gioui.org/text"
	"gioui.org/unit"
	"gioui.org/widget"
	"gioui.org/widget/material"
	"gioui.org/x/component"
	"github.com/hauntedness/giom/ui/api"
	"golang.org/x/exp/shiny/materialdesign/colornames"
)

// DrawAppBarLayout reusable function to draw consistent AppBar
func DrawAppBarLayout(gtx layout.Context, th *material.Theme, widget layout.Widget) layout.Dimensions {
	gtx.Constraints.Max.Y = gtx.Dp(56)
	component.Rect{Size: gtx.Constraints.Max, Color: th.ContrastBg}.Layout(gtx)
	inset := layout.Inset{Left: unit.Dp(8), Right: unit.Dp(8)}
	return inset.Layout(gtx, widget)
}

type PromptContent struct {
	*material.Theme
	btnYes      *widget.Clickable
	btnNo       *widget.Clickable
	HeaderTxt   string
	ContentText string
}

func NewPromptContent(theme *material.Theme, headerText string, contentText string, btnYes *widget.Clickable, btnNo *widget.Clickable) api.View {
	return &PromptContent{
		Theme:       theme,
		btnYes:      btnYes,
		btnNo:       btnNo,
		HeaderTxt:   headerText,
		ContentText: contentText,
	}
}

func (p *PromptContent) Layout(gtx layout.Context) layout.Dimensions {
	gtx.Constraints.Min.X = gtx.Constraints.Max.X
	inset := layout.UniformInset(unit.Dp(16))
	d := inset.Layout(gtx, func(gtx layout.Context) layout.Dimensions {
		return layout.Flex{Axis: layout.Vertical}.Layout(gtx,
			layout.Rigid(func(gtx layout.Context) layout.Dimensions {
				if p.HeaderTxt == "" {
					return layout.Dimensions{}
				}
				bd := material.Body1(p.Theme, p.HeaderTxt)
				bd.Font.Weight = text.Bold
				bd.Alignment = text.Middle
				return bd.Layout(gtx)
			}),
			layout.Rigid(func(gtx layout.Context) layout.Dimensions {
				return layout.Spacer{Height: unit.Dp(8)}.Layout(gtx)
			}),
			layout.Rigid(func(gtx layout.Context) layout.Dimensions {
				if p.ContentText == "" {
					return layout.Dimensions{}
				}
				bd := material.Body1(p.Theme, p.ContentText)
				bd.Alignment = text.Middle
				return bd.Layout(gtx)
			}),
			layout.Rigid(layout.Spacer{Height: unit.Dp(16)}.Layout),
			layout.Rigid(func(gtx layout.Context) layout.Dimensions {
				return layout.Flex{Spacing: layout.SpaceSides, Alignment: layout.Middle}.Layout(gtx,
					layout.Rigid(func(gtx layout.Context) layout.Dimensions {
						btn := material.Button(p.Theme, p.btnYes, "Yes")
						btn.Background = color.NRGBA(colornames.Red500)
						return btn.Layout(gtx)
					}),
					layout.Rigid(layout.Spacer{Width: unit.Dp(16)}.Layout),
					layout.Rigid(func(gtx layout.Context) layout.Dimensions {
						btn := material.Button(p.Theme, p.btnNo, "No")
						btn.Background = color.NRGBA(colornames.Green500)
						return btn.Layout(gtx)
					}),
				)
			}),
		)
	})
	return d
}
