package ui

import (
	"gioui.org/app"

	"github.com/hauntedness/giom/ui/pages"
	"github.com/hauntedness/giom/ui/pages/about"
	"github.com/hauntedness/giom/ui/pages/appbar"
	"github.com/hauntedness/giom/ui/pages/discloser"
	"github.com/hauntedness/giom/ui/pages/menu"
	"github.com/hauntedness/giom/ui/pages/navdrawer"
	"github.com/hauntedness/giom/ui/pages/textfield"

	"gioui.org/font/gofont"
	"gioui.org/io/system"
	"gioui.org/layout"
	"gioui.org/op"
	"gioui.org/widget/material"
)

func Loop(w *app.Window) error {
	th := material.NewTheme(gofont.Collection())
	var ops op.Ops

	router := pages.NewRouter()
	router.Register(0, appbar.New(&router))
	router.Register(1, navdrawer.New(&router))
	router.Register(2, textfield.New(&router))
	router.Register(3, menu.New(&router))
	router.Register(4, discloser.New(&router))
	router.Register(5, about.New(&router))

	for {
		select {
		case e := <-w.Events():
			switch e := e.(type) {
			case system.DestroyEvent:
				return e.Err
			case system.FrameEvent:
				gtx := layout.NewContext(&ops, e)
				router.Layout(gtx, th)
				e.Frame(gtx.Ops)
			}
		}
	}
}
