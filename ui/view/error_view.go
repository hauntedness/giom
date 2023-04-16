package view

import (
	"gioui.org/layout"
	"gioui.org/widget/material"
	"github.com/hauntedness/giom/assets/fonts"
	"github.com/hauntedness/giom/ui/api"
)

type ErrorView struct {
	*material.Theme
	api.Manager
	Error string
}

func (i *ErrorView) Layout(gtx layout.Context) (d layout.Dimensions) {
	if i.Theme == nil {
		i.Theme = fonts.NewTheme()
	}
	if i.Error != "" {
		return layout.Flex{Axis: layout.Vertical}.Layout(gtx,
			layout.Rigid(func(gtx layout.Context) layout.Dimensions {
				return material.Body1(i.Theme, i.Error).Layout(gtx)
			}),
		)
	}
	return d
}
