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
	"github.com/hauntedness/giom/play/common"
	"golang.org/x/exp/shiny/materialdesign/icons"
)

// App Bar looks not so powerfull at present
type App struct {
	theme *material.Theme
	anim  *component.VisibilityAnimation
	item1 *int
}

func main() {
	th := material.NewTheme(gofont.Collection())
	th.TextSize = unit.Sp(10)
	item1 := 1
	giomApp := App{
		theme: th,
		anim:  &component.VisibilityAnimation{},
		item1: &item1,
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
	layout.Flex{}.Layout(gtx, layout.Flexed(1, func(gtx layout.Context) layout.Dimensions {
		dim := layout.Flex{
			Axis:      layout.Vertical,
			Spacing:   0,
			Alignment: 0,
			WeightSum: 0,
		}.Layout(gtx,
			layout.Rigid(func(gtx layout.Context) layout.Dimensions {
				textField := component.TextField{
					Editor:    widget.Editor{},
					Alignment: 0,
					Helper:    "Helper",
					CharLimit: 0,
					Prefix:    nil,
					Suffix: func(gtx layout.Context) layout.Dimensions {
						return material.IconButton(app.theme, &widget.Clickable{}, common.MustIcon(icons.ActionSearch), "desc").Layout(gtx)
					},
				}
				textField.Editor.Insert("search")
				return textField.Layout(gtx, app.theme, "some hint")
			}),
		)
		return dim
	}))
}
