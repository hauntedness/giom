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
	"github.com/hauntedness/giom/internal/log"
)

func main() {
	go func() {
		window := app.NewWindow(app.Title("Protonet"))
		app_ := &application{
			theme:     material.NewTheme(gofont.Collection()),
			Clickable: &widget.Clickable{},
		}
		_ = app_

		window.Option(func(m unit.Metric, c *app.Config) {
			c.MaxSize = image.Point{
				X: m.Dp(400),
				Y: m.Dp(600),
			}
		})
		var ops op.Ops
		for e := range window.Events() {
			switch e := e.(type) {
			case system.DestroyEvent:
				log.Errors(e.Err)
			case system.FrameEvent:
				gtx := layout.NewContext(&ops, e)
				app_.Layout(gtx)
				e.Frame(gtx.Ops)
			}
		}
		os.Exit(0)
	}()
	app.Main()
}

type application struct {
	theme     *material.Theme
	Clickable *widget.Clickable
	state     bool
}

func (app *application) Layout(gtx layout.Context) layout.Dimensions {
	button := material.Button(
		app.theme, app.Clickable, "submit",
	)
	app.state = !app.state
	dim := button.Layout(gtx)
	dim.Size.Y = 50
	for i, c := range button.Button.Clicks() {
		log.Info("button clicked", "count", i, "modifiers", c.Modifiers.String())
	}
	return dim
}
