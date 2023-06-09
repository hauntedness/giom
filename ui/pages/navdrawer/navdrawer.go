package navdrawer

import (
	"gioui.org/layout"
	"gioui.org/widget"
	"gioui.org/widget/material"
	"gioui.org/x/component"
	"github.com/hauntedness/giom/ui/applayout"
	"github.com/hauntedness/giom/ui/icon"
	"github.com/hauntedness/giom/ui/pages"
)

// Page holds the state for a page demonstrating the features of
// the NavDrawer component.
type Page struct {
	nonModalDrawer widget.Bool
	widget.List
	*pages.Router
}

// New constructs a Page with the provided router.
func New(router *pages.Router) *Page {
	return &Page{
		Router: router,
	}
}

var _ pages.Page = &Page{}

func (p *Page) Actions() []component.AppBarAction {
	return []component.AppBarAction{}
}

func (p *Page) Overflow() []component.OverflowAction {
	return []component.OverflowAction{}
}

func (p *Page) NavItem() component.NavItem {
	return component.NavItem{
		Name: "Nav Drawer Features",
		Icon: icon.SettingsIcon,
	}
}

func (p *Page) Layout(gtx layout.Context, th *material.Theme) layout.Dimensions {
	p.List.Axis = layout.Vertical
	return material.List(th, &p.List).Layout(gtx, 1, func(gtx layout.Context, _ int) layout.Dimensions {
		return layout.Flex{
			Alignment: layout.Middle,
			Axis:      layout.Vertical,
		}.Layout(gtx,
			layout.Rigid(func(gtx layout.Context) layout.Dimensions {
				return applayout.DefaultInset.Layout(gtx, material.Body1(th, `The nav drawer widget provides a consistent interface element for navigation.

The controls below allow you to see the various features available in our Navigation Drawer implementation.`).Layout)
			}),
			layout.Rigid(func(gtx layout.Context) layout.Dimensions {
				return applayout.DetailRow{}.Layout(gtx,
					material.Body1(th, "Use non-modal drawer").Layout,
					func(gtx layout.Context) layout.Dimensions {
						if p.nonModalDrawer.Changed() {
							p.Router.NonModalDrawer = p.nonModalDrawer.Value
							if p.nonModalDrawer.Value {
								p.Router.NavAnim.Appear(gtx.Now)
							} else {
								p.Router.NavAnim.Disappear(gtx.Now)
							}
						}
						return material.Switch(th, &p.nonModalDrawer, "Use Non-Modal Navigation Drawer").Layout(gtx)
					})
			}),
			layout.Rigid(func(gtx layout.Context) layout.Dimensions {
				return applayout.DetailRow{}.Layout(gtx,
					material.Body1(th, "Drag to Close").Layout,
					material.Body2(th, "You can close the modal nav drawer by dragging it to the left.").Layout)
			}),
			layout.Rigid(func(gtx layout.Context) layout.Dimensions {
				return applayout.DetailRow{}.Layout(gtx,
					material.Body1(th, "Touch Scrim to Close").Layout,
					material.Body2(th, "You can close the modal nav drawer touching anywhere in the translucent scrim to the right.").Layout)
			}),
			layout.Rigid(func(gtx layout.Context) layout.Dimensions {
				return applayout.DetailRow{}.Layout(gtx,
					material.Body1(th, "Bottom content anchoring").Layout,
					material.Body2(th, "If you toggle support for the bottom app bar in the App Bar settings, nav drawer content will anchor to the bottom of the drawer area instead of the top.").Layout)
			}),
		)
	})
}
