package view

import (
	"image"

	"gioui.org/layout"
	"gioui.org/op"
	"gioui.org/op/clip"
	"gioui.org/op/paint"
	"gioui.org/widget"
	"gioui.org/widget/material"
	"gioui.org/x/component"
	"github.com/hauntedness/giom/assets"
	"golang.org/x/exp/shiny/materialdesign/icons"
)

type AvatarView struct {
	Size  image.Point
	Image image.Image
	widget.Clickable
	*material.Theme
	Selected      bool
	SelectionMode bool
}

func (v *AvatarView) Layout(gtx layout.Context) layout.Dimensions {
	if v.Size == (image.Point{}) {
		v.Size = image.Point{X: gtx.Dp(48), Y: gtx.Dp(48)}
	}
	gtx.Constraints.Min, gtx.Constraints.Max = v.Size, v.Size
	var imageWidget widget.Image
	if v.Image == nil {
		v.Image = assets.AppIconImage
		imageOps := paint.NewImageOp(v.Image)
		imageWidget = widget.Image{Src: imageOps, Fit: widget.Fill, Position: layout.Center, Scale: 0}
	} else {
		imageOps := paint.NewImageOp(v.Image)
		imageWidget = widget.Image{Src: imageOps, Fit: widget.Fill, Position: layout.Center, Scale: 0}
	}
	stack := layout.Stack{Alignment: layout.SE}
	return stack.Layout(gtx,
		layout.Stacked(func(gtx layout.Context) layout.Dimensions {
			ops := clip.UniformRRect(image.Rectangle{
				Max: image.Point{
					X: gtx.Constraints.Max.X,
					Y: gtx.Constraints.Max.Y,
				},
			}, gtx.Constraints.Max.X/2).Push(gtx.Ops)
			defer ops.Pop()
			return imageWidget.Layout(gtx)
		}),
		layout.Stacked(func(gtx layout.Context) layout.Dimensions {
			if !v.SelectionMode {
				return layout.Dimensions{}
			}
			gtx.Constraints.Max.X = int(float64(v.Size.X) * 0.40)
			gtx.Constraints.Max.Y = int(float64(v.Size.Y) * 0.40)
			gtx.Constraints.Min = gtx.Constraints.Max
			offsetOp := op.Offset(
				image.Pt(gtx.Constraints.Max.X/4, gtx.Constraints.Max.Y/4),
			).Push(gtx.Ops)
			defer offsetOp.Pop()
			color := v.Theme.ContrastBg
			d := component.Rect{Size: gtx.Constraints.Max, Color: color, Radii: gtx.Constraints.Max.X / 2}.Layout(gtx)
			iconColor := v.Theme.ContrastFg
			icon, _ := widget.NewIcon(icons.ActionCheckCircle)
			if !v.Selected {
				return d
			}
			return icon.Layout(gtx, iconColor)
		}),
	)
}
