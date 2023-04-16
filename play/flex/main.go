package main

import (
	"image"
	"os"

	"gioui.org/app"
	"gioui.org/font/gofont"
	"gioui.org/io/system"
	"gioui.org/layout"
	"gioui.org/op"
	"gioui.org/unit"
	"gioui.org/widget"
	"gioui.org/widget/material"
	"gioui.org/x/component"
	"github.com/hauntedness/giom/internal/log"
	"github.com/hauntedness/giom/play/common"
	"golang.org/x/exp/shiny/materialdesign/icons"
)

func main() {
	go func() {
		window := app.NewWindow(app.Title("Protonet"))
		theme := material.NewTheme(gofont.Collection())
		window.Option(func(m unit.Metric, c *app.Config) {
			c.MaxSize = image.Point{
				X: m.Dp(800),
				Y: m.Dp(640),
			}
		})
		var ops op.Ops
		for e := range window.Events() {
			switch e := e.(type) {
			case system.DestroyEvent:
				log.Errors(e.Err)
			case system.FrameEvent:
				gtx := layout.NewContext(&ops, e)
				children := []layout.Widget{
					func(gtx layout.Context) layout.Dimensions {
						gtx.Constraints.Max = image.Pt(100, 100)
						button := material.Button(theme, &widget.Clickable{}, "hello")
						button.Inset = layout.UniformInset(unit.Dp(10))
						// button.Background = common.ColorBackground
						return button.Layout(gtx)
					},
					func(gtx layout.Context) layout.Dimensions {
						tipIcon := component.TipIconButton(
							theme, &component.TipArea{}, &widget.Clickable{}, "hahaha", common.MustIcon(icons.ActionHome),
						)
						tipIcon.Background = common.ColorBackground
						return tipIcon.Layout(gtx)
					},
					func(gtx layout.Context) layout.Dimensions {
						gtx.Constraints.Max = image.Pt(100, 100)
						button := component.SimpleIconButton(
							common.ColorBackground, common.ColorBlue, &widget.Clickable{}, common.MustIcon(icons.AVAVTimer),
						)
						return button.Layout(gtx)
					},
					func(gtx layout.Context) layout.Dimensions {
						ds := component.Discloser(theme, &component.DiscloserState{})
						return ds.Layout(gtx,
							material.H3(theme, "hahaha").Layout,
							material.H3(theme, "hahaha").Layout,
							material.H3(theme, "hahaha").Layout,
						)
					},
				}
				layout.Flex{
					Axis:      0,
					Spacing:   layout.SpaceAround,
					Alignment: 0,
					WeightSum: 0,
				}.Layout(gtx,
					func(ws []layout.Widget) []layout.FlexChild {
						ret := make([]layout.FlexChild, len(ws))
						for i := range ws {
							ret[i] = layout.Flexed(1.0, ws[i])
						}
						return ret
					}(children)...,
				)
				e.Frame(gtx.Ops)
			}
		}
		os.Exit(0)
	}()
	app.Main()
}
