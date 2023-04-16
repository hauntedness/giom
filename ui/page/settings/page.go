package settings

import (
	"bytes"
	"fmt"
	"image"
	"image/color"
	_ "image/gif"
	_ "image/jpeg"
	_ "image/png"
	"time"

	"gioui.org/layout"
	"gioui.org/op/clip"
	"gioui.org/op/paint"
	"gioui.org/unit"
	"gioui.org/widget"
	"gioui.org/widget/material"
	"gioui.org/x/component"
	"github.com/hauntedness/giom/internal/log"
	"github.com/hauntedness/giom/ui/api"
	"github.com/hauntedness/giom/ui/view"
	"golang.org/x/exp/shiny/materialdesign/colornames"
	"golang.org/x/exp/shiny/materialdesign/icons"
)

type page struct {
	layout.List
	api.Manager
	buttonNavIcon      widget.Clickable
	btnAddAccount      widget.Clickable
	btnShowAccounts    widget.Clickable
	menuIcon           *widget.Icon
	items              []*pageItem
	AccountForm        api.View
	AccountsView       api.View
	menuVisibilityAnim component.VisibilityAnimation
	*view.ModalContent
}

func New(manager api.Manager) api.Page {
	menuIcon, _ := widget.NewIcon(icons.ContentAddCircle)
	accountsIcon, _ := widget.NewIcon(icons.SocialGroup)
	contactsIcon, _ := widget.NewIcon(icons.CommunicationContacts)
	themeIcon, _ := widget.NewIcon(icons.ImagePalette)
	notificationsIcon, _ := widget.NewIcon(icons.SocialNotifications)
	helpIcon, _ := widget.NewIcon(icons.ActionHelp)
	aboutIcon, _ := widget.NewIcon(icons.ActionInfo)
	p := page{
		Manager:  manager,
		List:     layout.List{Axis: layout.Vertical},
		menuIcon: menuIcon,
		menuVisibilityAnim: component.VisibilityAnimation{
			Duration: time.Millisecond * 250,
			State:    component.Invisible,
			Started:  time.Time{},
		},
		items: []*pageItem{
			{
				Manager: manager,
				Theme:   manager.Theme(),
				Title:   "Accounts",
				Icon:    accountsIcon,
				url:     api.AccountsPageURL,
			},
			{
				Manager: manager,
				Theme:   manager.Theme(),
				Title:   "Contacts",
				Icon:    contactsIcon,
				url:     api.ContactsPageURL,
			},
			{
				Manager: manager,
				Theme:   manager.Theme(),
				Title:   "Theme",
				Icon:    themeIcon,
				url:     api.ThemePageURL,
			},
			{
				Manager: manager,
				Theme:   manager.Theme(),
				Title:   "Notifications",
				Icon:    notificationsIcon,
				url:     api.NotificationsPageURL,
			},
			{
				Manager: manager,
				Theme:   manager.Theme(),
				Title:   "Help",
				Icon:    helpIcon,
				url:     api.HelpPageURL,
			},
			{
				Manager: manager,
				Theme:   manager.Theme(),
				Title:   "About",
				Icon:    aboutIcon,
				url:     api.AboutPageURL,
			},
		},
	}
	p.AccountForm = view.NewAccountFormView(manager, p.onAddAccountSuccess)
	p.AccountsView = view.NewAccountsView(manager, p.onAccountChange)
	p.ModalContent = view.NewModalContent(func() { p.Modal().Dismiss(nil) })
	return &p
}

func (p *page) Layout(gtx layout.Context) (d layout.Dimensions) {
	if p.items == nil {
		p.items = make([]*pageItem, 0)
	}
	if p.btnAddAccount.Clicked() {
		p.AccountForm = view.NewAccountFormView(p.Manager, p.onAddAccountSuccess)
		p.menuVisibilityAnim.Disappear(gtx.Now)
		p.Modal().Show(p.drawAddAccountModal, nil, component.VisibilityAnimation{
			Duration: time.Millisecond * 250,
			State:    component.Invisible,
			Started:  time.Time{},
		})
	}

	if p.btnShowAccounts.Clicked() {
		p.AccountsView = view.NewAccountsView(p.Manager, p.onAccountChange)
		p.menuVisibilityAnim.Disappear(gtx.Now)
		p.Modal().Show(p.drawShowAccountsModal, nil, component.VisibilityAnimation{
			Duration: time.Millisecond * 250,
			State:    component.Invisible,
			Started:  time.Time{},
		})
	}

	flex := layout.Flex{
		Axis:      layout.Vertical,
		Spacing:   layout.SpaceEnd,
		Alignment: layout.Start,
	}
	d = flex.Layout(gtx,
		layout.Rigid(p.DrawAppBar),
		layout.Rigid(p.drawItems),
	)
	return d
}

func (p *page) DrawAppBar(gtx layout.Context) layout.Dimensions {
	if p.buttonNavIcon.Clicked() {
		p.Manager.NavigateToUrl(api.ChatPageURL, nil)
	}

	return view.DrawAppBarLayout(gtx, p.Manager.Theme(), func(gtx layout.Context) layout.Dimensions {
		return layout.Flex{Alignment: layout.Middle, Spacing: layout.SpaceBetween}.Layout(gtx,
			layout.Rigid(func(gtx layout.Context) layout.Dimensions {
				return layout.Flex{Alignment: layout.Middle}.Layout(gtx,
					layout.Rigid(func(gtx layout.Context) layout.Dimensions {
						return material.ButtonLayoutStyle{
							Background:   p.Manager.Theme().ContrastBg,
							Button:       &p.buttonNavIcon,
							CornerRadius: unit.Dp(56 / 2),
						}.Layout(gtx,
							func(gtx layout.Context) layout.Dimensions {
								return view.DrawAppImageForNav(gtx, p.Manager.Theme())
							},
						)
					}),
					layout.Rigid(func(gtx layout.Context) layout.Dimensions {
						gtx.Constraints.Max.X = gtx.Constraints.Max.X - gtx.Dp(56)
						return layout.Inset{Left: unit.Dp(12)}.Layout(gtx, func(gtx layout.Context) layout.Dimensions {
							titleText := "Settings"
							label := material.Label(p.Manager.Theme(), unit.Sp(18), titleText)
							label.Color = p.Manager.Theme().Palette.ContrastFg
							return component.TruncatingLabelStyle(label).Layout(gtx)
						})
					}),
				)
			}),
			layout.Rigid(func(gtx layout.Context) layout.Dimensions {
				var img image.Image
				var err error
				a := p.Service().Account()
				if a.PublicKey != "" && len(a.PublicImage) != 0 {
					img, _, err = image.Decode(bytes.NewReader(a.PublicImage))
					if err != nil {
						log.Errors(err)
					}
				}
				if img != nil {
					return p.btnShowAccounts.Layout(gtx, func(gtx layout.Context) layout.Dimensions {
						radii := gtx.Dp(20)
						gtx.Constraints.Max.X, gtx.Constraints.Max.Y = radii*2, radii*2
						bounds := image.Rect(0, 0, radii*2, radii*2)
						clipOp := clip.UniformRRect(bounds, radii).Push(gtx.Ops)
						imgOps := paint.NewImageOp(img)
						imgWidget := widget.Image{Src: imgOps, Fit: widget.Contain, Position: layout.Center, Scale: 0}
						d := imgWidget.Layout(gtx)
						clipOp.Pop()
						return d
					})
				}
				button := material.IconButton(p.Manager.Theme(), &p.btnAddAccount, p.menuIcon, "Context Menu")
				button.Size = unit.Dp(40)
				button.Background = p.Manager.Theme().Palette.ContrastBg
				button.Color = p.Manager.Theme().Palette.ContrastFg
				button.Inset = layout.UniformInset(unit.Dp(8))
				d := button.Layout(gtx)
				return d
			}),
		)
	})
}

func (p *page) drawItems(gtx layout.Context) layout.Dimensions {
	return p.List.Layout(gtx, len(p.items), func(gtx layout.Context, index int) (d layout.Dimensions) {
		inset := layout.Inset{Top: unit.Dp(0), Bottom: unit.Dp(0)}
		return inset.Layout(gtx, func(gtx layout.Context) layout.Dimensions {
			return layout.Flex{Axis: layout.Vertical}.Layout(gtx,
				layout.Rigid(p.items[index].Layout),
				layout.Rigid(func(gtx layout.Context) layout.Dimensions {
					size := image.Pt(gtx.Constraints.Max.X, gtx.Dp(1))
					bounds := image.Rectangle{Max: size}
					bgColor := color.NRGBA(colornames.Grey500)
					bgColor.A = 75
					paint.FillShape(gtx.Ops, bgColor, clip.UniformRRect(bounds, 0).Op(gtx.Ops))
					return layout.Dimensions{Size: image.Pt(size.X, size.Y)}
				}),
			)
		})
	})
}

func (p *page) onAddAccountSuccess() {
	p.Modal().Dismiss(func() {
		p.NavigateToUrl(api.ChatPageURL, nil)
	})
}

func (p *page) drawAddAccountModal(gtx layout.Context) layout.Dimensions {
	gtx.Constraints.Max.X = int(float32(gtx.Constraints.Max.X) * 0.85)
	gtx.Constraints.Max.Y = int(float32(gtx.Constraints.Max.Y) * 0.85)
	return p.ModalContent.DrawContent(gtx, p.Theme(), p.AccountForm.Layout)
}

func (p *page) drawShowAccountsModal(gtx layout.Context) layout.Dimensions {
	gtx.Constraints.Max.X = int(float32(gtx.Constraints.Max.X) * 0.85)
	gtx.Constraints.Max.Y = int(float32(gtx.Constraints.Max.Y) * 0.85)
	return p.ModalContent.DrawContent(gtx, p.Theme(), p.AccountsView.Layout)
}

func (p *page) onAccountChange() {
	p.Modal().Dismiss(p.afterAccountsModalDismissed)
}

func (p *page) afterAccountsModalDismissed() {
	p.NavigateToUrl(api.ChatPageURL, func() {
		a := p.Service().Account()
		txt := fmt.Sprintf("Switched to %s account", a.PublicKey)
		p.Snackbar().Show(txt, nil, color.NRGBA{}, "")
	})
}

func (p *page) URL() api.URL {
	return api.SettingsPageURL
}
