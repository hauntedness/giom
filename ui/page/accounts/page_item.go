package accounts

import (
	"image/color"
	"strings"
	"time"

	"gioui.org/layout"
	"gioui.org/op"
	"gioui.org/op/clip"
	"gioui.org/text"
	"gioui.org/unit"
	"gioui.org/widget"
	"gioui.org/widget/material"
	"gioui.org/x/component"
	"github.com/hauntedness/giom/service"
	"github.com/hauntedness/giom/ui/api"
	"github.com/hauntedness/giom/ui/view"
	"golang.org/x/exp/shiny/materialdesign/colornames"
	"golang.org/x/exp/shiny/materialdesign/icons"
)

type pageItem struct {
	*material.Theme
	widget.Clickable
	buttonIconMore        widget.Clickable
	btnSetCurrentIdentity widget.Clickable
	btnAccountDetails     widget.Clickable
	btnMenuContent        widget.Clickable
	listAccountDetails    layout.List
	buttonIconMoreDim     layout.Dimensions
	api.Manager
	service.Account
	PressedStamp int64
	view.AvatarView
	iconMore           *widget.Icon
	menuVisibilityAnim component.VisibilityAnimation
	AccountDetails     *view.AccountDetails
	ModalContent       *view.ModalContent
}

func (i *pageItem) Layout(gtx layout.Context) layout.Dimensions {
	if i.Theme == nil {
		i.Theme = i.Manager.Theme()
	}
	return i.layoutContent(gtx)
}

func (i *pageItem) layoutContent(gtx layout.Context) layout.Dimensions {
	if i.menuVisibilityAnim == (component.VisibilityAnimation{}) {
		i.menuVisibilityAnim = component.VisibilityAnimation{
			Duration: time.Millisecond * 250,
			State:    component.Invisible,
			Started:  time.Time{},
		}
	}

	if i.buttonIconMore.Clicked() {
		i.Clickable.Clicked()
		i.menuVisibilityAnim.Appear(gtx.Now)
	}

	if i.btnSetCurrentIdentity.Clicked() {
		i.Manager.Service().SetAsCurrentAccount(i.Account)
		i.menuVisibilityAnim.Disappear(gtx.Now)
	}
	if i.btnAccountDetails.Clicked() {
		i.menuVisibilityAnim.Disappear(gtx.Now)
		if i.AccountDetails == nil {
			i.AccountDetails = view.NewAccountDetails(i.Manager, i.Account)
		}
		i.Modal().Show(i.drawAccountDetailsModal, nil, component.VisibilityAnimation{
			Duration: time.Millisecond * 250,
			State:    component.Invisible,
			Started:  time.Time{},
		})
	}

	if i.Clickable.Clicked() {
		if !i.menuVisibilityAnim.Visible() {
			if i.SelectionMode {
				i.Selected = !i.Selected
			}
			if i.PressedStamp != 0 {
				diff := time.Now().UnixMilli() - i.PressedStamp
				if diff < 350 {
					i.SelectionMode = !i.SelectionMode
					i.Selected = !i.Selected
				}
			}
		}
		if !i.btnMenuContent.Pressed() {
			i.menuVisibilityAnim.Disappear(gtx.Now)
		}
		i.PressedStamp = time.Now().UnixMilli()
	}

	btnStyle := material.ButtonLayoutStyle{Background: i.Theme.ContrastBg, Button: &i.Clickable}

	if i.Selected || i.Clickable.Hovered() {
		btnStyle.Background.A = 50
	} else {
		btnStyle.Background.A = 10
	}
	if i.AvatarView.Theme == nil {
		i.AvatarView.Theme = i.Theme
	}

	gtx.Constraints.Min.X = gtx.Constraints.Max.X
	d := btnStyle.Layout(gtx, func(gtx layout.Context) layout.Dimensions {
		inset := layout.Inset{Top: unit.Dp(16), Bottom: unit.Dp(16), Left: unit.Dp(8), Right: unit.Dp(8)}
		d := inset.Layout(gtx, func(gtx layout.Context) layout.Dimensions {
			flex := layout.Flex{Spacing: layout.SpaceEnd, Alignment: layout.Middle}
			d := flex.Layout(gtx,
				layout.Rigid(func(gtx layout.Context) layout.Dimensions {
					flex := layout.Flex{Spacing: layout.SpaceSides, Alignment: layout.Start, Axis: layout.Vertical}
					d := flex.Layout(gtx, layout.Rigid(i.AvatarView.Layout))
					return d
				}),
				layout.Rigid(layout.Spacer{Width: unit.Dp(8)}.Layout),
				layout.Rigid(func(gtx layout.Context) layout.Dimensions {
					gtx.Constraints.Max.X = gtx.Constraints.Max.X - gtx.Dp(80)
					gtx.Constraints.Min.X = gtx.Constraints.Max.X
					return i.listAccountDetails.Layout(gtx, 1, func(gtx layout.Context, _ int) layout.Dimensions {
						flex := layout.Flex{Spacing: layout.SpaceSides, Alignment: layout.Start, Axis: layout.Vertical}
						inset := layout.Inset{Right: unit.Dp(16), Left: unit.Dp(16)}
						d := inset.Layout(gtx, func(gtx layout.Context) layout.Dimensions {
							d := flex.Layout(gtx,
								layout.Rigid(func(gtx layout.Context) layout.Dimensions {
									b := material.Body1(i.Theme, i.Account.PublicKey)
									b.Font.Weight = text.Bold
									return b.Layout(gtx)
								}),
								layout.Rigid(func(gtx layout.Context) layout.Dimensions {
									b := material.Body1(i.Theme, strings.Trim(i.Account.PublicKey, "\n"))
									b.Color = color.NRGBA(colornames.Grey600)
									return b.Layout(gtx)
								}),
							)
							return d
						})
						return d
					})
				}),
				layout.Rigid(func(gtx layout.Context) layout.Dimensions {
					flex := layout.Flex{Axis: layout.Vertical, Spacing: layout.SpaceBetween, Alignment: layout.Middle}
					return flex.Layout(gtx,
						layout.Rigid(func(gtx layout.Context) layout.Dimensions {
							if i.iconMore == nil {
								i.iconMore, _ = widget.NewIcon(icons.NavigationMoreVert)
							}
							button := material.IconButton(i.Theme, &i.buttonIconMore, i.iconMore, "Vertical Icon For Options")
							button.Size = unit.Dp(24)
							button.Background = color.NRGBA{}
							button.Color = i.Theme.ContrastBg
							button.Inset = layout.UniformInset(unit.Dp(16))
							i.buttonIconMoreDim = button.Layout(gtx)
							return i.buttonIconMoreDim
						}),
						layout.Rigid(func(gtx layout.Context) layout.Dimensions {
							a := i.Manager.Service().Account()
							if a.PublicKey == i.Account.PublicKey {
								icon, _ := widget.NewIcon(icons.ActionCheckCircle)
								return icon.Layout(gtx, i.Theme.ContrastBg)
							}
							return layout.Dimensions{}
						}),
					)
				}),
			)
			return d
		})
		return d
	})

	gtx.Constraints.Max.Y = d.Size.Y
	i.drawMenuLayout(gtx)
	return d
}

func (i *pageItem) drawMenuLayout(gtx layout.Context) layout.Dimensions {
	layout.Stack{Alignment: layout.NE}.Layout(gtx,
		layout.Stacked(func(gtx layout.Context) layout.Dimensions {
			progress := i.menuVisibilityAnim.Revealed(gtx)
			gtx.Constraints.Max.X = int(float32(gtx.Constraints.Max.X) * progress)
			gtx.Constraints.Max.Y = int(float32(gtx.Constraints.Max.Y) * progress)
			return component.Rect{Size: gtx.Constraints.Max, Color: color.NRGBA{A: 200}}.Layout(gtx)
		}),
		layout.Stacked(func(gtx layout.Context) layout.Dimensions {
			progress := i.menuVisibilityAnim.Revealed(gtx)
			macro := op.Record(gtx.Ops)
			d := i.btnMenuContent.Layout(gtx, i.drawMenuItems)
			call := macro.Stop()
			d.Size.X = int(float32(d.Size.X) * progress)
			d.Size.Y = int(float32(d.Size.Y) * progress)
			component.Rect{Size: d.Size, Color: color.NRGBA(colornames.White)}.Layout(gtx)
			clipOp := clip.Rect{Max: d.Size}.Push(gtx.Ops)
			call.Add(gtx.Ops)
			clipOp.Pop()
			return d
		}),
	)
	return layout.Dimensions{}
}

func (i *pageItem) drawMenuItems(gtx layout.Context) layout.Dimensions {
	inset := layout.UniformInset(unit.Dp(12))
	gtx.Constraints.Max.X = int(float32(gtx.Constraints.Max.X) / 1.5)
	gtx.Constraints.Min.X = gtx.Constraints.Max.X
	return layout.Flex{Axis: layout.Vertical}.Layout(gtx,
		layout.Rigid(func(gtx layout.Context) layout.Dimensions {
			btnStyle := material.ButtonLayoutStyle{Button: &i.btnSetCurrentIdentity}
			btnStyle.Background = color.NRGBA(colornames.White)
			return btnStyle.Layout(gtx,
				func(gtx layout.Context) layout.Dimensions {
					inset := inset
					return inset.Layout(gtx, func(gtx layout.Context) layout.Dimensions {
						return layout.Flex{Spacing: layout.SpaceEnd}.Layout(gtx,
							layout.Flexed(1, func(gtx layout.Context) layout.Dimensions {
								bd := material.Body1(i.Theme, "Set As Current Account")
								bd.Color = color.NRGBA(colornames.Black)
								bd.Alignment = text.Start
								return bd.Layout(gtx)
							}),
						)
					})
				},
			)
		}),
		layout.Rigid(func(gtx layout.Context) layout.Dimensions {
			btnStyle := material.ButtonLayoutStyle{Button: &i.btnAccountDetails}
			btnStyle.Background = color.NRGBA(colornames.White)
			return btnStyle.Layout(gtx,
				func(gtx layout.Context) layout.Dimensions {
					inset := inset
					return inset.Layout(gtx, func(gtx layout.Context) layout.Dimensions {
						return layout.Flex{Spacing: layout.SpaceEnd}.Layout(gtx,
							layout.Flexed(1, func(gtx layout.Context) layout.Dimensions {
								bd := material.Body1(i.Theme, "Account Details")
								bd.Color = color.NRGBA(colornames.Black)
								bd.Alignment = text.Start
								return bd.Layout(gtx)
							}),
						)
					})
				},
			)
		}),
	)
}

func (i *pageItem) drawAccountDetailsModal(gtx layout.Context) layout.Dimensions {
	gtx.Constraints.Max.X = int(float32(gtx.Constraints.Max.X) * 0.85)
	gtx.Constraints.Max.Y = int(float32(gtx.Constraints.Max.Y) * 0.85)
	return i.ModalContent.DrawContent(gtx, i.Theme, i.AccountDetails.Layout)
}
