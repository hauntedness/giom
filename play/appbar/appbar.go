package main

import (
	"image"
	"image/color"
	"os"
	"time"

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

// App Bar looks not so powerfull at present
type App struct {
	theme *material.Theme
	modal *component.ModalLayer
}

func main() {
	th := material.NewTheme(gofont.Collection())
	giomApp := App{
		theme: th,
		modal: component.NewModal(),
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
	layout.Flex{}.Layout(
		gtx,
		layout.Rigid(func(gtx layout.Context) layout.Dimensions {
			appBar := component.AppBar{
				NavigationButton: widget.Clickable{},
				NavigationIcon:   common.MustIcon(icons.AVArtTrack),
				Title:            "Title",
				ContextualTitle:  "ContextualTitle",
				ModalLayer:       app.modal,
				Anchor:           0,
			}
			appBar.Widget = func(gtx layout.Context, th *material.Theme, anim *component.VisibilityAnimation) layout.Dimensions {
				ls := material.Body1(th, "hello this is body")
				anim.Appear(time.Now())
				return ls.Layout(gtx)
			}
			appBar.SetActions(
				[]component.AppBarAction{
					{
						OverflowAction: component.OverflowAction{},
						Layout: func(gtx layout.Context, bg color.NRGBA, fg color.NRGBA) layout.Dimensions {
							return material.Body1(app.theme, "body").Layout(gtx)
						},
					},
				}, []component.OverflowAction{
					{
						Name: "over flow",
						Tag:  nil,
					},
				},
			)
			if appBar.Pressed() {
				// doesn't work
				log.Info("appbar pressed")
			}
			return appBar.Layout(gtx, app.theme, "navDesc", "overflowDesc")
		}),
	)
}
