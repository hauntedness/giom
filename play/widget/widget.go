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

type state struct {
	checked         *widget.Bool
	buttonClickable widget.Clickable
}

// App Bar looks not so powerfull at present
type App struct {
	theme *material.Theme
	anim  *component.VisibilityAnimation
	item1 *int
	state state
}

func (app *App) Layout(gtx layout.Context) {
	layout.Flex{}.Layout(gtx, layout.Flexed(1, func(gtx layout.Context) layout.Dimensions {
		flex := layout.Flex{
			Axis:      layout.Vertical,
			Spacing:   0,
			Alignment: 0,
			WeightSum: 0,
		}
		dim := flex.Layout(gtx,
			layout.Flexed(1, func(gtx layout.Context) layout.Dimensions {
				checkBox := material.CheckBox(app.theme, app.state.checked, "check me")
				checkBox.Color = common.ColorBlue
				if checkBox.CheckBox.Value && checkBox.CheckBox.Changed() {
					log.Info("check me checked")
				}
				if gtx.Constraints.Max.Y > 100 {
					gtx.Constraints.Max.Y = 100
					gtx.Constraints.Min.Y = 0
				}
				return checkBox.Layout(gtx)
			}),
			layout.Flexed(1, func(gtx layout.Context) layout.Dimensions {
				dim := component.Rect{
					Color: common.ColorBlue,
					Size: image.Point{
						X: 100,
						Y: 100,
					},
					Radii: 30,
				}.Layout(gtx)
				return dim
			}),
			layout.Flexed(1, func(gtx layout.Context) layout.Dimensions {
				return material.IconButton(
					app.theme, &app.state.buttonClickable, common.MustIcon(icons.AVAVTimer), "hello",
				).Layout(gtx)
			}),
		)
		return dim
	}))
}

func main() {
	th := material.NewTheme(gofont.Collection())
	item1 := 1
	a := App{
		theme: th,
		anim:  &component.VisibilityAnimation{},
		item1: &item1,
		state: state{
			checked: &widget.Bool{},
		},
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
		if err := Loop(window, &a); err != nil {
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
			app.Layout(gtx)
			e.Frame(gtx.Ops)
		}
	}
	return nil
}

func Add(a int, b int) int {
	return a + b
}
