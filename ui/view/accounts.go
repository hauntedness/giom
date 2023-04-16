package view

import (
	"bytes"
	"image"
	"image/color"

	"gioui.org/layout"
	"gioui.org/op/clip"
	"gioui.org/op/paint"
	"gioui.org/text"
	"gioui.org/unit"
	"gioui.org/widget"
	"gioui.org/widget/material"
	"github.com/hauntedness/giom/assets"
	"github.com/hauntedness/giom/internal/log"
	"github.com/hauntedness/giom/ui/api"
	"golang.org/x/exp/shiny/materialdesign/colornames"
)

type accountsView struct {
	layout.List
	api.Manager
	theme                 *material.Theme
	title                 string
	accountsItems         []*accountsItem
	currentAccountLayout  layout.List
	enum                  widget.Enum
	accountChangeCallback func()
}

func NewAccountsView(manager api.Manager, accountChangeCallback func()) api.View {
	errorTh := *manager.Theme()
	errorTh.ContrastBg = color.NRGBA(colornames.Red500)
	p := accountsView{
		Manager:               manager,
		theme:                 manager.Theme(),
		title:                 "Accounts",
		List:                  layout.List{Axis: layout.Vertical},
		accountsItems:         []*accountsItem{},
		accountChangeCallback: accountChangeCallback,
	}
	return &p
}

func (p *accountsView) Layout(gtx layout.Context) layout.Dimensions {
	a := p.Service().Account()
	p.enum.Value = a.PublicKey
	flex := layout.Flex{
		Axis:      layout.Vertical,
		Spacing:   layout.SpaceEnd,
		Alignment: layout.Start,
	}

	d := flex.Layout(gtx,
		layout.Rigid(p.drawIdentitiesItems),
	)

	return d
}

func (p *accountsView) drawIdentitiesItems(gtx layout.Context) layout.Dimensions {
	if p.isProcessingRequired() {
		p.accountsItems = make([]*accountsItem, 0, len(p.Service().Accounts()))
		for _, userID := range <-p.Service().Accounts() {
			p.accountsItems = append(p.accountsItems, &accountsItem{
				Theme:   p.theme,
				Manager: p.Manager,
				Account: userID,
				Enum:    &p.enum,
			})
		}
	}
	return layout.Flex{Axis: layout.Vertical}.Layout(gtx,
		layout.Rigid(func(gtx layout.Context) layout.Dimensions {
			inset := layout.UniformInset(unit.Dp(16))
			return inset.Layout(gtx, func(gtx layout.Context) layout.Dimensions {
				flex := layout.Flex{Alignment: layout.Middle}
				a := p.Service().Account()
				d := flex.Layout(gtx,
					layout.Rigid(func(gtx layout.Context) layout.Dimensions {
						var img image.Image
						var err error
						img, _, err = image.Decode(bytes.NewReader(a.PublicImage))
						if err != nil {
							log.Errors(err)
						}
						if img == nil {
							img = assets.AppIconImage
						}
						radii := gtx.Dp(24)
						gtx.Constraints.Max.X, gtx.Constraints.Max.Y = radii*2, radii*2
						bounds := image.Rect(0, 0, radii*2, radii*2)
						clipOp := clip.UniformRRect(bounds, radii).Push(gtx.Ops)
						imgOps := paint.NewImageOp(img)
						imgWidget := widget.Image{Src: imgOps, Fit: widget.Contain, Position: layout.Center, Scale: 0}
						d := imgWidget.Layout(gtx)
						clipOp.Pop()
						return d
					}),
					layout.Rigid(func(gtx layout.Context) layout.Dimensions {
						gtx.Constraints.Min.X = gtx.Constraints.Max.X
						return p.currentAccountLayout.Layout(gtx, 1, func(gtx layout.Context, index int) layout.Dimensions {
							flex := layout.Flex{Spacing: layout.SpaceSides, Alignment: layout.Start, Axis: layout.Vertical}
							inset := layout.Inset{Right: unit.Dp(8), Left: unit.Dp(8)}
							d := inset.Layout(gtx, func(gtx layout.Context) layout.Dimensions {
								d := flex.Layout(gtx,
									layout.Rigid(func(gtx layout.Context) layout.Dimensions {
										b := material.Body1(p.Theme(), a.PublicKey)
										b.Font.Weight = text.Bold
										return b.Layout(gtx)
									}),
									//layout.Rigid(func(gtx layout.Context) layout.Dimensions {
									//	b := material.Body1(p.Theme(), strings.Trim(string(p.currentAccount.Contents), "\n"))
									//	b.Color = color.NRGBA(colornames.Grey600)
									//	return b.Layout(gtx)
									//}),
								)
								return d
							})
							return d
						})
					}),
				)
				return d
			})
		}),
		layout.Rigid(func(gtx layout.Context) layout.Dimensions {
			return p.List.Layout(gtx, len(p.accountsItems), func(gtx layout.Context, index int) (d layout.Dimensions) {
				accountItem := p.accountsItems[index]
				if accountItem.Clickable.Pressed() {
					p.Manager.Service().SetAsCurrentAccount(accountItem.Account)
					if p.accountChangeCallback != nil {
						p.accountChangeCallback()
					}
				}
				return p.accountsItems[index].Layout(gtx)
			})
		}),
	)
}

// isProcessingRequired
func (p *accountsView) isProcessingRequired() bool {
	isRequired := len(<-p.Service().Accounts()) != len(p.accountsItems)
	return isRequired
}
