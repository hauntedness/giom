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
	"gioui.org/widget/material"
	"gioui.org/x/component"
)

// App Bar looks not so powerfull at present
type App struct {
	theme *material.Theme
	anim  *component.VisibilityAnimation
	item1 *int
}

func main() {
	th := material.NewTheme(gofont.Collection())
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
		resize := component.Resize{
			Axis:  layout.Vertical,
			Ratio: 0,
		}
		dim := resize.Layout(gtx,
			material.Body1(app.theme, "body1111111111111111111111111111").Layout,
			material.Body1(app.theme, "body2222222222222222222222222222").Layout,
			component.Divider(app.theme).Layout,
		)
		return dim
	}))
}
