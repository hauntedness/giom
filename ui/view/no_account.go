package view

import (
	"image/color"
	"time"

	"gioui.org/layout"
	"gioui.org/text"
	"gioui.org/unit"
	"gioui.org/widget"
	"gioui.org/widget/material"
	"gioui.org/x/component"
	"github.com/hauntedness/giom/assets/fonts"
	"github.com/hauntedness/giom/ui/api"
	"golang.org/x/exp/shiny/materialdesign/colornames"
	"golang.org/x/exp/shiny/materialdesign/icons"
)

type NoAccountView struct {
	api.Manager
	buttonNewAccount *IconButton
	*material.Theme
	*widget.Icon
	inActiveTh      *material.Theme
	iconCreateNewID *widget.Icon
	AccountFormView api.View
	*ModalContent
}

func NewNoAccount(manager api.Manager) *NoAccountView {
	acc := NoAccountView{Manager: manager, Theme: manager.Theme()}
	acc.AccountFormView = NewAccountFormView(manager, acc.onSuccess)
	acc.ModalContent = NewModalContent(func() {
		acc.Modal().Dismiss(nil)
	})
	return &acc
}

func (na *NoAccountView) Layout(gtx layout.Context) layout.Dimensions {
	if na.Theme == nil {
		na.Theme = fonts.NewTheme()
	}
	if na.Icon == nil {
		na.Icon, _ = widget.NewIcon(icons.ActionAccountCircle)
	}
	if na.inActiveTh == nil {
		inActiveTh := *fonts.NewTheme()
		inActiveTh.ContrastBg = color.NRGBA(colornames.Grey500)
		na.inActiveTh = &inActiveTh
	}
	if na.iconCreateNewID == nil {
		na.iconCreateNewID, _ = widget.NewIcon(icons.ContentCreate)
	}
	if na.buttonNewAccount == nil {
		na.buttonNewAccount = &IconButton{
			Theme: na.Theme,
			Icon:  na.Icon,
			Text:  "Add/Edit Account",
		}
	}

	flex := layout.Flex{Axis: layout.Vertical, Spacing: layout.SpaceSides, Alignment: layout.Middle}
	gtx.Constraints.Min.Y = gtx.Constraints.Max.Y
	if na.buttonNewAccount.Button.Clicked() {
		na.Modal().Show(na.drawModalContent, nil, component.VisibilityAnimation{
			Duration: time.Millisecond * 250,
			State:    component.Invisible,
			Started:  time.Time{},
		})
	}
	d := flex.Layout(gtx,
		layout.Rigid(func(gtx layout.Context) layout.Dimensions {
			return DrawProtonetImageCenter(gtx, na.Theme)
		}),
		layout.Rigid(layout.Spacer{Height: unit.Dp(16)}.Layout),
		layout.Rigid(func(gtx layout.Context) layout.Dimensions {
			return layout.Center.Layout(gtx, func(gtx layout.Context) layout.Dimensions {
				bdy := material.Body1(na.Theme, "No Account(s) Created")
				bdy.Alignment = text.Middle
				bdy.Font.Weight = text.Black
				bdy.Color = color.NRGBA{R: 102, G: 117, B: 127, A: 255}
				return bdy.Layout(gtx)
			})
		}),
		layout.Rigid(layout.Spacer{Height: unit.Dp(16)}.Layout),
		layout.Rigid(func(gtx layout.Context) layout.Dimensions {
			return layout.Flex{Spacing: layout.SpaceSides}.Layout(gtx, layout.Rigid(func(gtx layout.Context) layout.Dimensions {
				gtx.Constraints.Max.X = gtx.Dp(250)
				return na.buttonNewAccount.Layout(gtx)
			}))
		}),
		layout.Rigid(layout.Spacer{Height: unit.Dp(16)}.Layout),
	)
	return d
}

func (na *NoAccountView) drawModalContent(gtx layout.Context) layout.Dimensions {
	gtx.Constraints.Max.X = int(float32(gtx.Constraints.Max.X) * 0.85)
	gtx.Constraints.Max.Y = int(float32(gtx.Constraints.Max.Y) * 0.85)
	return na.ModalContent.DrawContent(gtx, na.Theme, na.AccountFormView.Layout)
}

func (na *NoAccountView) onSuccess() {
	na.Modal().Dismiss(func() {
		na.Manager.NavigateToUrl(api.ChatPageURL, nil)
	})
}
