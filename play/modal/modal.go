package main

import (
	"image"
	"os"
	"strconv"

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
)

// App Bar looks not so powerfull at present
type App struct {
	theme      *material.Theme
	anim       component.VisibilityAnimation
	modalState component.ModalState
	clickables []widget.Clickable
}

func main() {
	th := material.NewTheme(gofont.Collection())
	giomApp := App{
		theme: th,
		clickables: []widget.Clickable{
			{}, {}, {},
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
			layout.Flexed(1, func(gtx layout.Context) layout.Dimensions {
				sheet := component.NewSheet()
				dim := sheet.Layout(gtx, app.theme, &app.anim, func(gtx layout.Context) layout.Dimensions {
					ls := material.Body2(app.theme, "body1")
					ls.Color = common.ColorRed
					return ls.Layout(gtx)
				})
				return dim
			}),
			layout.Flexed(1, func(gtx layout.Context) layout.Dimensions {
				ms := component.Modal(app.theme, &app.modalState)
				ms.Scrim.Color = common.ColorBackground
				ms.Show(gtx.Now, func(gtx layout.Context) layout.Dimensions {
					list := layout.List{
						Axis:        layout.Vertical,
						ScrollToEnd: false,
						Alignment:   0,
						Position:    layout.Position{},
					}
					return list.Layout(gtx, 3, func(gtx layout.Context, index int) layout.Dimensions {
						btn := material.Button(app.theme, &app.clickables[index], "button"+strconv.Itoa(index))
						if (&app.clickables[index]).Clicked() {
							log.Info("button clicked", "button", btn.Text, "index", index)
						}
						return btn.Layout(gtx)
					})
				})
				return ms.Layout(gtx)
			}),
		)
		return dim
	}))
}
