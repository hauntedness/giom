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
	"github.com/hauntedness/giom/service"
	"github.com/hauntedness/giom/ui/api"
	"golang.org/x/exp/shiny/materialdesign/icons"
)

type NoContactView struct {
	api.Manager
	buttonAddContact *IconButton
	*material.Theme
	*widget.Icon
	ContactFormView api.View
	*ModalContent
}

func NewNoContact(manager api.Manager, onSuccess func(contactAddr string), btnText string) *NoContactView {
	btnIcon, _ := widget.NewIcon(icons.CommunicationContacts)
	if btnText == "" {
		btnText = "Add Contact"
	}
	nc := NoContactView{
		Manager: manager,
		Theme:   manager.Theme(),
		buttonAddContact: &IconButton{
			Theme: manager.Theme(),
			Icon:  btnIcon,
			Text:  btnText,
		},
	}
	nc.ContactFormView = NewContactForm(manager, service.Contact{}, onSuccess)
	nc.ModalContent = NewModalContent(func() { nc.Modal().Dismiss(nil) })
	return &nc
}

func (nc *NoContactView) Layout(gtx layout.Context) layout.Dimensions {
	flex := layout.Flex{Axis: layout.Vertical, Spacing: layout.SpaceSides, Alignment: layout.Middle}
	gtx.Constraints.Min.Y = gtx.Constraints.Max.Y
	if nc.buttonAddContact.Button.Clicked() {
		nc.Modal().Show(nc.drawModalContent, nil, component.VisibilityAnimation{
			Duration: time.Millisecond * 250,
			State:    component.Invisible,
			Started:  time.Time{},
		})
	}
	d := flex.Layout(gtx,
		layout.Rigid(func(gtx layout.Context) layout.Dimensions {
			return DrawProtonetImageCenter(gtx, nc.Theme)
		}),
		layout.Rigid(layout.Spacer{Height: unit.Dp(16)}.Layout),
		layout.Rigid(func(gtx layout.Context) layout.Dimensions {
			return layout.Center.Layout(gtx, func(gtx layout.Context) layout.Dimensions {
				body := material.Body1(nc.Theme, "No Contact(s) Found")
				body.Alignment = text.Middle
				body.Font.Weight = text.Black
				body.Color = color.NRGBA{R: 102, G: 117, B: 127, A: 255}
				return body.Layout(gtx)
			})
		}),
		layout.Rigid(layout.Spacer{Height: unit.Dp(16)}.Layout),
		layout.Rigid(func(gtx layout.Context) layout.Dimensions {
			return layout.Flex{Spacing: layout.SpaceSides}.Layout(gtx, layout.Rigid(func(gtx layout.Context) layout.Dimensions {
				gtx.Constraints.Max.X = gtx.Dp(250)
				return nc.buttonAddContact.Layout(gtx)
			}))
		}),
		layout.Rigid(layout.Spacer{Height: unit.Dp(16)}.Layout),
	)
	return d
}

func (nc *NoContactView) drawModalContent(gtx layout.Context) layout.Dimensions {
	gtx.Constraints.Max.X = int(float32(gtx.Constraints.Max.X) * 0.85)
	gtx.Constraints.Max.Y = int(float32(gtx.Constraints.Max.Y) * 0.85)
	return nc.ModalContent.DrawContent(gtx, nc.Theme, nc.ContactFormView.Layout)
}
