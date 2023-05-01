package main

import (
	"embed"
	"os"
	"strings"
	"time"

	"gioui.org/app"
	"gioui.org/io/system"
	"gioui.org/layout"
	"gioui.org/op"
	"gioui.org/unit"
	"gioui.org/widget"
	"gioui.org/widget/material"
	"gioui.org/x/component"
	"github.com/hauntedness/giom/assets/fonts"
	"github.com/hauntedness/giom/internal/log"
	"github.com/hauntedness/giom/play/common"
	"github.com/hauntedness/giom/play/navigate/page"
	"github.com/hauntedness/giom/ui/icon"
)

// App Bar looks not so powerfull at present
type App struct {
	theme               *material.Theme
	accountClickable    widget.Clickable
	restaurantClickable widget.Clickable
	settingsClickable   widget.Clickable
	heartClickable      widget.Clickable
	h1Anim              component.VisibilityAnimation
	h1Text              string
	h2Text              string
	body1Text           string
	body1State          widget.Selectable
	modalState          component.ModalState
}

var _ embed.FS

//go:embed page/testdata/逍遥小贵婿.txt
var data string
var p page.Page

func init() {
	p = page.From(data, 100, 10)
	total := p.Total
	go func() {
		for t := range time.Tick(time.Second * 10) {
			unixSecond := t.Unix()
			offset := int(unixSecond) % total
			p = page.From(data, offset, 10)
			log.Infos("offset", offset)
		}
	}()
}

func (app *App) Layout(gtx layout.Context) layout.Dimensions {
	for _, e := range gtx.Events(app) {
		log.Infos("enter app event", e)
	}
	return layout.Flex{Axis: layout.Vertical}.Layout(
		gtx,
		// header
		layout.Flexed(1, func(gtx layout.Context) layout.Dimensions {
			nd := component.NewNav(app.h1Text, app.h2Text)
			return nd.Layout(gtx, app.theme, &app.h1Anim)
		}),
		// todo
		layout.Flexed(1, func(gtx layout.Context) layout.Dimensions {
			return layout.Center.Layout(gtx, _kline.Layout)
		}),
		// body
		layout.Flexed(8, func(gtx layout.Context) layout.Dimensions {
			txt := strings.Join(p.Lines, "\n")
			body1 := material.Body1(app.theme, txt)
			body1.State = &app.body1State
			if body1.State.SelectionLen() > 0 {
				if body1.State.Focused() {
					log.Infos("select", body1.State.SelectedText())
					ms := component.Modal(app.theme, &app.modalState)
					if ms.Clicked() {
						log.Info("modal clicked")
					}
					ms.Appear(time.Now())
					ms.Layout(gtx)
				}
			}
			return body1.Layout(gtx)
		}),
		// 底栏上面的分割线
		layout.Rigid(
			func(gtx layout.Context) layout.Dimensions {
				div := common.Div{
					NoDirection: true,
					Border:      widget.Border{Color: common.ColorGreen, Width: unit.Dp(0)},
					Inset:       layout.Inset{Top: unit.Dp(10), Bottom: unit.Dp(10), Left: unit.Dp(1), Right: unit.Dp(1)},
				}
				dim := div.Layout(gtx, func(gtx layout.Context) layout.Dimensions {
					divider := component.Divider(app.theme)
					divider.Fill = common.ColorRed
					divider.Thickness = unit.Dp(1)
					return divider.Layout(gtx)
				})
				return dim
			},
		),
		// 底部导航
		layout.Rigid(
			func(gtx layout.Context) layout.Dimensions {
				// 添加border 方便调试
				var dim layout.Dimensions = common.Div{
					Inset: layout.UniformInset(unit.Dp(8)),
					Border: widget.Border{
						Color:        common.ColorBlue,
						CornerRadius: 0,
						Width:        1,
					},
					Direction: layout.Center,
				}.Layout(
					gtx,
					func(gtx layout.Context) layout.Dimensions {
						return layout.Flex{
							Axis:      layout.Horizontal,
							Spacing:   layout.SpaceAround,
							Alignment: layout.Middle,
							WeightSum: 0,
						}.Layout(
							gtx,
							layout.Flexed(
								1,
								func(gtx layout.Context) layout.Dimensions {
									return layout.Center.Layout(
										gtx,
										func(gtx layout.Context) layout.Dimensions {
											return material.IconButton(
												app.theme, &app.accountClickable, icon.AccountBalanceIcon, "account balance",
											).Layout(gtx)
										},
									)
								},
							),
							layout.Flexed(
								1,
								func(gtx layout.Context) layout.Dimensions {
									return layout.Center.Layout(
										gtx,
										func(gtx layout.Context) layout.Dimensions {
											return material.IconButton(
												app.theme, &app.restaurantClickable, icon.RestaurantMenuIcon, "restaurant",
											).Layout(gtx)
										},
									)
								},
							),
							layout.Flexed(
								1,
								func(gtx layout.Context) layout.Dimensions {
									return layout.Center.Layout(gtx,
										func(gtx layout.Context) layout.Dimensions {
											return material.IconButton(
												app.theme, &app.settingsClickable, icon.SettingsIcon, "settings",
											).Layout(gtx)
										},
									)
								},
							),
							layout.Flexed(
								1,
								func(gtx layout.Context) layout.Dimensions {
									return layout.Center.Layout(gtx,
										func(gtx layout.Context) layout.Dimensions {
											return material.IconButton(
												app.theme, &app.heartClickable, icon.HeartIcon, "heart",
											).Layout(gtx)
										},
									)
								},
							),
						)
					})
				return dim
			},
		),
	)
}

func main() {
	th := fonts.NewTheme()
	giomApp := &App{
		theme:     th,
		h1Text:    time.Now().Format(time.DateTime),
		body1Text: "",
	}

	go func() {
		window := app.NewWindow(app.Title("giom"))
		if err := giomApp.Loop(window); err != nil {
			panic(err)
		}
		os.Exit(0)
	}()
	app.Main()
}

func (app *App) Loop(w *app.Window) error {
	go func() {
		for range time.NewTicker(time.Second * 3).C {
			app.h1Text = Fetch(time.Now())
			app.body1Text = Fetch(time.Now()) + time.Now().Format(time.RFC3339)
			// fees, err := starknet.EstimateFees()
			// if err != nil {
			// 	app.body1Text = err.Error()
			// } else {
			// 	app.body1Text = fmt.Sprintf("%v", fees)
			// }
			w.Invalidate()
		}
	}()
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

func Fetch(t time.Time) string {
	return "一些中文 和 some english"
}
