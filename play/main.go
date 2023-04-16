package main

import (
	"image"
	"os"

	"gioui.org/app"
	"gioui.org/font/gofont"
	"gioui.org/io/system"
	"gioui.org/layout"
	"gioui.org/op"
	"gioui.org/text"
	"gioui.org/unit"
	"gioui.org/widget"
	"gioui.org/widget/material"
	"gioui.org/x/component"
	"github.com/hauntedness/giom/play/custom"
	"golang.org/x/exp/shiny/materialdesign/icons"
)

type App struct {
	theme       *material.Theme
	button      *custom.ButtonVisual
	currentPage custom.AboutPage
}

func main() {
	th := material.NewTheme(gofont.Collection())
	giomApp := App{
		theme:       th,
		button:      &custom.ButtonVisual{},
		currentPage: *custom.NewAboutPage(th),
	}
	go func() {
		window := app.NewWindow(app.Title("Protonet"))
		window.Option(
			func(m unit.Metric, c *app.Config) {
				c.MaxSize = image.Point{
					X: m.Dp(320),
					Y: m.Dp(640),
				}
			},
		)
		if err := Loop(window, &giomApp); err != nil {
			panic(err)
		}
		os.Exit(0)
	}()
	app.Main()
}

func Loop(w *app.Window, app *App) error {
	// theme := material.NewTheme(gofont.Collection())
	// backClickTag is meant for tracking user's backClick action, specially on mobile
	var ops op.Ops
	for e := range w.Events() {
		switch e := e.(type) {
		case system.DestroyEvent:
			return e.Err
		case system.FrameEvent:
			gtx := layout.NewContext(&ops, e)
			draw(gtx, app)
			e.Frame(gtx.Ops)
		}
	}
	return nil
}

func draw(gtx layout.Context, app *App) {
	layout.Flex{
		Axis:      layout.Vertical,
		Spacing:   0,
		Alignment: layout.Middle,
		WeightSum: 0,
	}.Layout(
		gtx,
		// top
		layout.Rigid(func(gtx layout.Context) layout.Dimensions {
			style := material.H6(app.theme, "hahaha")
			style.Alignment = text.Middle
			return style.Layout(gtx)
		}),
		// body
		layout.Flexed(1, func(gtx layout.Context) layout.Dimensions {
			style := material.Body1(app.theme, "hahaha hahaha")
			style.Alignment = text.Start
			return style.Layout(gtx)
		}),
		// bottom
		layout.Rigid(func(gtx layout.Context) layout.Dimensions {
			flex := layout.Flex{
				Axis:      0,
				Spacing:   layout.SpaceSides,
				Alignment: layout.Middle,
				WeightSum: 0,
			}
			return flex.Layout(
				gtx,
				// home menu
				layout.Flexed(0.1, func(gtx layout.Context) layout.Dimensions {
					ml := component.NewModal()
					bar := component.NewAppBar(ml)
					bar.ContextualTitle = "ContextualTitle"
					bar.NavigationButton = widget.Clickable{}
					var err error
					bar.NavigationIcon, err = widget.NewIcon(icons.ActionHome)
					if err != nil {
						panic(err)
					}
					bar.Anchor = component.Bottom
					return bar.Layout(gtx, app.theme, "hahaha", "ActionHome")
				}),
				// search menu
				layout.Flexed(0.1, func(gtx layout.Context) layout.Dimensions {
					ml := component.NewModal()
					bar := component.NewAppBar(ml)
					bar.ContextualTitle = "ContextualTitle"
					bar.NavigationButton = widget.Clickable{}
					bar.ContextualTitle = "ContextTitle"
					var err error
					bar.NavigationIcon, err = widget.NewIcon(icons.ActionSearch)
					if err != nil {
						panic(err)
					}
					return bar.Layout(gtx, app.theme, "hahaha", "ActionSearch")
				}),
				layout.Flexed(0.1, func(gtx layout.Context) layout.Dimensions {
					ml := component.NewModal()
					bar := component.NewAppBar(ml)
					bar.ContextualTitle = "ContextualTitle"
					bar.NavigationButton = widget.Clickable{}
					var err error
					bar.NavigationIcon, err = widget.NewIcon(icons.ActionAccountBalance)
					if err != nil {
						panic(err)
					}
					return bar.Layout(gtx, app.theme, "hahaha", "ActionAccountBalance")
				}),
				layout.Flexed(0.1, func(gtx layout.Context) layout.Dimensions {
					ml := component.NewModal()
					bar := component.NewAppBar(ml)
					bar.ContextualTitle = "ContextualTitle"
					bar.NavigationButton = widget.Clickable{}
					var err error
					bar.NavigationIcon, err = widget.NewIcon(icons.ActionPermMedia)
					if err != nil {
						panic(err)
					}
					return bar.Layout(gtx, app.theme, "navDesc", "overflowDesc")
				}),
			)
		}),
	)
}
