package ui

import (
	"image"
	"os/exec"
	"strings"
	"time"

	"gioui.org/app"
	"gioui.org/io/key"
	"gioui.org/io/pointer"
	"gioui.org/io/system"
	"gioui.org/layout"
	"gioui.org/op"
	"gioui.org/op/clip"
	"gioui.org/op/paint"
	"github.com/hauntedness/giom/internal/log"
	"github.com/hauntedness/giom/service"
	"github.com/hauntedness/giom/ui/api"
)

// FixTimezone https://github.com/golang/go/issues/20455
func FixTimezone() {
	out, err := exec.Command("/system/bin/getprop", "persist.sys.timezone").Output()
	if err != nil {
		return
	}
	z, err := time.LoadLocation(strings.TrimSpace(string(out)))
	if err != nil {
		return
	}
	time.Local = z
}

var appManager = AppManager{service: service.GetServiceInstance()}

func init() {
	FixTimezone()
	go appManager.init()
}

func Loop(w *app.Window) error {
	var ops op.Ops
	appManager.window = w

	// backClickTag is meant for tracking user's backClick action, specially on mobile
	var backClickTag struct{}

	subscription := appManager.Service().Subscribe()

	for {
		select {
		case e := <-w.Events():
			switch e := e.(type) {
			case system.DestroyEvent:
				log.Error("system.DestroyEvent called", e.Err)
				return e.Err
			case system.FrameEvent:
				appManager.Insets = e.Insets
				e.Insets = system.Insets{}
				gtx := layout.NewContext(&ops, e)
				for _, event := range gtx.Events(&backClickTag) {
					switch e := event.(type) {
					case key.Event:
						switch e.Name {
						case key.NameBack:
							if len(appManager.pagesStack) > 1 {
								appManager.PopUp()
							}
						}
					}
				}
				// Listen to back command only when appManager.pagesStack is greater than 1,
				//  so we can pop up page else we want the android's default behavior
				if len(appManager.pagesStack) > 1 {
					key.InputOp{Tag: &backClickTag, Keys: key.NameBack}.Add(gtx.Ops)
				}
				appManager.Constraints = gtx.Constraints
				appManager.Metric = gtx.Metric
				// Create a clip area the size of the window.
				areaStack := clip.Rect(image.Rectangle{Max: gtx.Constraints.Max}).Push(gtx.Ops)
				// In desktop layout, sidebar exists and needs to listen to entire window's pointer event
				// hence added here. It avoids conflict with page that contains sidebar
				for _, elem := range []any{appManager.CurrentPage(), appManager.settingsSideBar, appManager.chatSideBar} {
					pointer.InputOp{
						Types: pointer.Enter | pointer.Leave | pointer.Drag | pointer.Press | pointer.Release | pointer.Scroll | pointer.Move,
						Tag:   elem,
					}.Add(gtx.Ops)
				}
				layout.Flex{Axis: layout.Vertical}.Layout(
					gtx,
					layout.Flexed(1, func(gtx layout.Context) layout.Dimensions {
						size := image.Point{X: gtx.Constraints.Max.X, Y: gtx.Constraints.Max.Y}
						bounds := image.Rectangle{Max: size}
						paint.FillShape(gtx.Ops, appManager.Theme().Bg, clip.UniformRRect(bounds, 0).Op(gtx.Ops))
						return appManager.Layout(gtx)
					}),
				)
				areaStack.Pop()
				e.Frame(gtx.Ops)
			case system.StageEvent:
				if e.Stage == system.StagePaused {
					log.Info("window is running in background")
					appManager.isStageRunning = false
				} else if e.Stage == system.StageRunning {
					log.Info("window is running in foreground")
					appManager.isStageRunning = true
				}
			}
		case event := <-subscription.Events():
			var chatBarFound, settingsBarFound bool
			for _, eachPage := range appManager.pagesStack {
				if l, ok := eachPage.(api.DatabaseListener); ok {
					l.OnDatabaseChange(event)
				}
				chatBarFound = eachPage == appManager.chatSideBar
				settingsBarFound = eachPage == appManager.settingsSideBar
			}
			if l, ok := appManager.chatSideBar.(api.DatabaseListener); ok && !chatBarFound {
				l.OnDatabaseChange(event)
			}
			if l, ok := appManager.settingsSideBar.(api.DatabaseListener); ok && !settingsBarFound {
				l.OnDatabaseChange(event)
			}
		}
	}
}