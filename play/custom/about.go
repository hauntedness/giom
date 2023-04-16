package custom

import (
	"gioui.org/layout"
	"gioui.org/unit"
	"gioui.org/widget"
	"gioui.org/widget/material"
	"github.com/hauntedness/giom/ui/api"
	"github.com/hauntedness/giom/ui/view"
	"golang.org/x/exp/shiny/materialdesign/icons"
)

type AboutPage struct {
	Theme            *material.Theme
	title            string
	buttonNavigation widget.Clickable
	navigationIcon   *widget.Icon
}

func NewAboutPage(theme *material.Theme) *AboutPage {
	navIcon, _ := widget.NewIcon(icons.NavigationArrowBack)
	return &AboutPage{
		Theme:          theme,
		title:          "About",
		navigationIcon: navIcon,
	}
}

func (p *AboutPage) Layout(gtx layout.Context) layout.Dimensions {
	flex := layout.Flex{Axis: layout.Vertical, Spacing: layout.SpaceEnd, Alignment: layout.Start}
	greetings := view.Greetings{}
	d := flex.Layout(gtx,
		layout.Rigid(p.DrawAppBar),
		layout.Rigid(greetings.Layout),
	)
	return d
}

func (p *AboutPage) DrawAppBar(gtx layout.Context) layout.Dimensions {
	gtx.Constraints.Max.Y = gtx.Dp(56)
	th := p.Theme

	return view.DrawAppBarLayout(gtx, th, func(gtx layout.Context) layout.Dimensions {
		return layout.Flex{Alignment: layout.Middle, Spacing: layout.SpaceBetween}.Layout(gtx,
			layout.Rigid(func(gtx layout.Context) layout.Dimensions {
				return layout.Flex{Alignment: layout.Middle}.Layout(gtx,
					layout.Rigid(func(gtx layout.Context) layout.Dimensions {
						navigationIcon := p.navigationIcon
						button := material.IconButton(th, &p.buttonNavigation, navigationIcon, "Nav Icon Button")
						button.Size = unit.Dp(40)
						button.Background = th.Palette.ContrastBg
						button.Color = th.Palette.ContrastFg
						button.Inset = layout.UniformInset(unit.Dp(8))
						return button.Layout(gtx)
					}),
					layout.Rigid(func(gtx layout.Context) layout.Dimensions {
						return layout.Inset{Left: unit.Dp(16)}.Layout(gtx, func(gtx layout.Context) layout.Dimensions {
							titleText := p.title
							title := material.Body1(th, titleText)
							title.Color = th.Palette.ContrastFg
							title.TextSize = unit.Sp(18)
							return title.Layout(gtx)
						})
					}),
				)
			}),
		)
	})
}

func (p *AboutPage) URL() api.URL {
	return api.AboutPageURL
}
