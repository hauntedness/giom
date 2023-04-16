package ui

import (
	"image"
	"sync"
	"time"

	"gioui.org/app"
	"gioui.org/io/system"
	"gioui.org/layout"
	"gioui.org/op"
	"gioui.org/op/clip"
	"gioui.org/op/paint"
	"gioui.org/unit"
	"gioui.org/widget/material"
	"gioui.org/x/component"
	"gioui.org/x/notify"
	"github.com/hauntedness/giom/assets/fonts"
	"github.com/hauntedness/giom/internal/log"
	"github.com/hauntedness/giom/service"
	"github.com/hauntedness/giom/ui/api"
	"github.com/hauntedness/giom/ui/page/about"
	"github.com/hauntedness/giom/ui/page/accounts"
	"github.com/hauntedness/giom/ui/page/chat"
	"github.com/hauntedness/giom/ui/page/contacts"
	"github.com/hauntedness/giom/ui/page/help"
	"github.com/hauntedness/giom/ui/page/notifications"
	"github.com/hauntedness/giom/ui/page/settings"
	"github.com/hauntedness/giom/ui/page/theme"
	"github.com/hauntedness/giom/ui/view"
)

// AppManager Always call NewAppManager function to create AppManager instance
type AppManager struct {
	window *app.Window
	view.Greetings
	chatSideBar     api.Page
	settingsSideBar api.Page
	theme           *material.Theme
	service         service.Service
	Constraints     layout.Constraints
	Metric          unit.Metric
	notifier        notify.Notifier
	system.Insets
	// isStageRunning, true value indicates app is running in foreground,
	// false indicates running in background
	isStageRunning   bool
	modal            api.Modal
	pagesStack       []api.Page
	prePushedView    api.Page
	poppedUpView     api.Page
	pageAnimation    component.VisibilityAnimation
	afterPageChange  func()
	snackbar         api.Snackbar
	initialized      bool
	initializedMutex sync.RWMutex
}

func (m *AppManager) Theme() *material.Theme {
	return m.theme
}

func (m *AppManager) Service() service.Service {
	return m.service
}

func (m *AppManager) SystemInsets() system.Insets {
	return m.Insets
}

func (m *AppManager) Window() *app.Window {
	return m.window
}

func (m *AppManager) Notifier() notify.Notifier {
	return m.notifier
}

func (m *AppManager) Snackbar() api.Snackbar {
	return m.snackbar
}

func (m *AppManager) Initialized() bool {
	m.initializedMutex.RLock()
	initialized := m.initialized
	m.initializedMutex.RUnlock()
	return initialized && m.Service().Initialized()
}

func (m *AppManager) setInitialized(initialized bool) {
	m.initializedMutex.Lock()
	m.initialized = initialized
	m.initializedMutex.Unlock()
}

func (m *AppManager) init() {
	m.theme = fonts.NewTheme()
	chatPage := chat.New(m)
	settingsPage := settings.New(m)
	m.chatSideBar = chatPage
	m.settingsSideBar = settingsPage
	var err error
	m.notifier, err = notify.NewNotifier()
	if err != nil {
		log.Errors(err)
	}
	m.SetModal(view.NewModalStack())
	m.snackbar = view.NewSnackBar(m)
	m.pageAnimation = component.VisibilityAnimation{
		Duration: time.Millisecond * 250,
		State:    component.Invisible,
		Started:  time.Time{},
	}
	m.setInitialized(true)
}

func (m *AppManager) Layout(gtx layout.Context) layout.Dimensions {
	d := layout.Flex{Axis: layout.Vertical}.Layout(gtx,
		layout.Rigid(func(gtx layout.Context) layout.Dimensions {
			return component.Rect{
				Color: m.Theme().ContrastBg,
				Size:  image.Point{X: gtx.Constraints.Max.X, Y: gtx.Dp(m.Insets.Top)},
				Radii: 0,
			}.Layout(gtx)
		}),
		layout.Flexed(1, func(gtx layout.Context) layout.Dimensions {
			size := image.Point{X: gtx.Constraints.Max.X, Y: gtx.Constraints.Max.Y - gtx.Dp(m.Insets.Bottom)}
			bounds := image.Rectangle{Max: size}
			paint.FillShape(gtx.Ops, m.Theme().Bg, clip.UniformRRect(bounds, 0).Op(gtx.Ops))
			d := m.drawPage(gtx)
			m.Snackbar().Layout(gtx)
			m.Modal().Layout(gtx)
			return d
		}),
		layout.Rigid(func(gtx layout.Context) layout.Dimensions {
			return component.Rect{
				Color: m.Theme().ContrastBg,
				Size:  image.Point{X: gtx.Constraints.Max.X, Y: gtx.Dp(m.Insets.Bottom)},
				Radii: 0,
			}.Layout(gtx)
		}),
	)
	return d
}

func (m *AppManager) CurrentPage() api.Page {
	stackSize := len(m.pagesStack)
	if stackSize > 0 {
		return m.pagesStack[stackSize-1]
	}
	m.pagesStack = []api.Page{chat.New(m)}
	return m.chatSideBar
}

func (m *AppManager) GetWindowWidthInDp() int {
	width := int(float32(m.Constraints.Max.X) / m.Metric.PxPerDp)
	return width
}

func (m *AppManager) GetWindowWidthInPx() int {
	return m.Constraints.Max.X
}

func (m *AppManager) GetWindowHeightInDp() int {
	width := int(float32(m.Constraints.Max.Y) / m.Metric.PxPerDp)
	return width
}

func (m *AppManager) GetWindowHeightInPx() int {
	return m.Constraints.Max.Y
}

func (m *AppManager) IsStageRunning() bool {
	return m.isStageRunning
}

func (m *AppManager) ShouldDrawSidebar() bool {
	minWidth := 800 // 800 is value in Dp
	winWidth := m.GetWindowWidthInDp()
	return winWidth >= minWidth
}

func (m *AppManager) Modal() api.Modal {
	return m.modal
}

func (m *AppManager) SetModal(modal api.Modal) {
	m.modal = modal
}

func (m *AppManager) PageFromUrl(url api.URL) api.Page {
	switch url {
	case api.SettingsPageURL:
		return m.settingsSideBar
	case api.ChatPageURL:
		fallthrough
	default:
		return m.chatSideBar
	}
}

func (m *AppManager) drawPage(gtx layout.Context) layout.Dimensions {
	maxDim := gtx.Constraints.Max

	d := layout.Flex{}.Layout(gtx,
		layout.Rigid(func(gtx layout.Context) (d layout.Dimensions) {
			if !m.ShouldDrawSidebar() {
				return d
			}
			gtx.Constraints.Max.X = int(float32(maxDim.X) * 0.40)
			gtx.Constraints.Min = gtx.Constraints.Max
			pageUrl := m.CurrentPage().URL()
			switch pageUrl {
			case api.SettingsPageURL, api.AccountsPageURL, api.ContactsPageURL,
				api.ThemePageURL, api.HelpPageURL, api.AboutPageURL, api.NotificationsPageURL:
				d = m.settingsSideBar.Layout(gtx)
			case api.ChatPageURL, api.ChatRoomPageURL:
				fallthrough
			default:
				d = m.chatSideBar.Layout(gtx)
			}
			return d
		}),
		layout.Rigid(func(gtx layout.Context) layout.Dimensions {
			if m.ShouldDrawSidebar() {
				gtx.Constraints.Max.X = int(float32(maxDim.X) * 0.60)
			}
			gtx.Constraints.Min = gtx.Constraints.Max
			maxDim := gtx.Constraints.Max
			gtx.Constraints.Min = maxDim
			currentPage := m.CurrentPage()
			prePushedPage := m.prePushedView
			poppedUpPage := m.poppedUpView
			switch {
			case prePushedPage != nil:
				progress := m.pageAnimation.Revealed(gtx)
				if m.pageAnimation.Animating() || m.pageAnimation.State == component.Invisible {
					m.pageAnimation.Appear(gtx.Now)
				}
				if progress == 1 {
					m.prePushedView = nil
					m.pageAnimation.State = component.Invisible
				}

				cl := clip.Rect{Max: maxDim}.Push(gtx.Ops)
				layout.Stack{Alignment: layout.NW}.Layout(gtx,
					layout.Stacked(func(gtx layout.Context) layout.Dimensions {
						gtx.Constraints.Min = maxDim
						offsetOperation := op.Offset(image.Point{
							X: -int(float32(maxDim.X) * progress),
							Y: 0,
						}).Push(gtx.Ops)
						defer offsetOperation.Pop()
						if (prePushedPage.URL() == api.SettingsPageURL ||
							prePushedPage.URL() == api.ChatPageURL) && m.ShouldDrawSidebar() {
							m.Greetings.Layout(gtx)
							return layout.Dimensions{Size: maxDim}
						}
						prePushedPage.Layout(gtx)
						return layout.Dimensions{Size: maxDim}
					}),
					layout.Stacked(func(gtx layout.Context) layout.Dimensions {
						gtx.Constraints.Min = maxDim
						offsetOperation := op.Offset(image.Point{
							X: int(float32(maxDim.X) * (1 - progress)),
							Y: 0,
						}).Push(gtx.Ops)
						defer offsetOperation.Pop()
						if (currentPage.URL() == api.SettingsPageURL ||
							currentPage.URL() == api.ChatPageURL) && m.ShouldDrawSidebar() {
							m.Greetings.Layout(gtx)
							return layout.Dimensions{Size: maxDim}
						}
						currentPage.Layout(gtx)
						return layout.Dimensions{Size: maxDim}
					}),
				)
				cl.Pop()
				return layout.Dimensions{Size: maxDim}
			case poppedUpPage != nil:
				progress := m.pageAnimation.Revealed(gtx)
				if m.pageAnimation.Animating() || m.pageAnimation.State == component.Invisible {
					m.pageAnimation.Appear(gtx.Now)
				}
				if progress == 1 {
					m.poppedUpView = nil
					m.pageAnimation.State = component.Invisible
				}

				cl := clip.Rect{Max: maxDim}.Push(gtx.Ops)
				layout.Stack{Alignment: layout.NW}.Layout(gtx,
					layout.Stacked(func(gtx layout.Context) layout.Dimensions {
						offsetOperation := op.Offset(image.Point{
							X: -int(float32(maxDim.X) * (1 - progress)),
							Y: 0,
						}).Push(gtx.Ops)
						defer offsetOperation.Pop()
						gtx.Constraints.Min = maxDim
						if (currentPage.URL() == api.SettingsPageURL ||
							currentPage.URL() == api.ChatPageURL) && m.ShouldDrawSidebar() {
							m.Greetings.Layout(gtx)
							return layout.Dimensions{Size: maxDim}
						}
						currentPage.Layout(gtx)
						return layout.Dimensions{Size: maxDim}
					}),
					layout.Stacked(func(gtx layout.Context) layout.Dimensions {
						offsetOperation := op.Offset(image.Point{
							X: int(float32(maxDim.X) * progress),
							Y: 0,
						}).Push(gtx.Ops)
						defer offsetOperation.Pop()
						gtx.Constraints.Min = maxDim
						if (poppedUpPage.URL() == api.SettingsPageURL ||
							poppedUpPage.URL() == api.ChatPageURL) && m.ShouldDrawSidebar() {
							m.Greetings.Layout(gtx)
							return layout.Dimensions{Size: maxDim}
						}
						poppedUpPage.Layout(gtx)
						return layout.Dimensions{Size: maxDim}
					}),
				)
				cl.Pop()
				return layout.Dimensions{Size: maxDim}
			default:

			}
			if (currentPage.URL() == api.SettingsPageURL ||
				currentPage.URL() == api.ChatPageURL) && m.ShouldDrawSidebar() {
				m.Greetings.Layout(gtx)
				return layout.Dimensions{Size: maxDim}
			}
			return currentPage.Layout(gtx)
		}),
	)
	state := m.pageAnimation.State
	progress := m.pageAnimation.Revealed(gtx)
	shouldCall := state == component.Invisible && progress == 0 && m.afterPageChange != nil &&
		!m.pageAnimation.Animating()
	if shouldCall {
		m.afterPageChange()
		m.afterPageChange = nil
	}
	return d
}

func (m *AppManager) NavigateToPage(page api.Page, AfterNavCallback func()) {
	pageURL := page.URL()
	m.afterPageChange = AfterNavCallback
	if pageURL == m.CurrentPage().URL() {
		if m.afterPageChange != nil {
			m.afterPageChange()
		}
		m.afterPageChange = nil
		return
	}
	switch pageURL {
	case api.SettingsPageURL:
		m.pagesStack = []api.Page{m.settingsSideBar}
	case api.ChatPageURL:
		m.pagesStack = []api.Page{m.chatSideBar}
	default:
		m.prePushedView = m.CurrentPage()
		m.pagesStack = append(m.pagesStack, page)
	}
	m.pageAnimation.Disappear(time.Now())
}

func (m *AppManager) NavigateToUrl(pageURL api.URL, AfterNavCallback func()) {
	m.afterPageChange = AfterNavCallback
	if pageURL == m.CurrentPage().URL() {
		if m.afterPageChange != nil {
			m.afterPageChange()
		}
		m.afterPageChange = nil
		return
	}
	var page api.Page
	switch pageURL {
	case api.SettingsPageURL:
		m.pagesStack = []api.Page{m.settingsSideBar}
	case api.ChatPageURL:
		m.pagesStack = []api.Page{m.chatSideBar}
	case api.AccountsPageURL:
		page = accounts.New(m)
	case api.ContactsPageURL:
		page = contacts.New(m)
	case api.ThemePageURL:
		page = theme.New(m)
	case api.NotificationsPageURL:
		page = notifications.New(m)
	case api.HelpPageURL:
		page = help.New(m)
	case api.AboutPageURL:
		page = about.New(m)
	}
	if page != nil {
		m.prePushedView = m.CurrentPage()
		m.pagesStack = append(m.pagesStack, page)
	}
	m.pageAnimation.Disappear(time.Now())
}

func (m *AppManager) PopUp() {
	stackLength := len(m.pagesStack)
	if stackLength > 1 {
		m.poppedUpView = m.pagesStack[stackLength-1]
		m.pagesStack = m.pagesStack[0 : stackLength-1]
	}
	m.pageAnimation.Disappear(time.Now())
}
