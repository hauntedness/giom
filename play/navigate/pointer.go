package main

import (
	"gioui.org/io/pointer"
	"gioui.org/layout"
	"github.com/hauntedness/giom/internal/log"
)

type ClickState struct {
	pressed bool
}

func (b *ClickState) Layout(gtx layout.Context) layout.Dimensions {
	// here we loop through all the events associated with this button.
	for _, e := range gtx.Events(b) {
		if e, ok := e.(pointer.Event); ok {
			switch e.Type {
			case pointer.Press:
				if e.Buttons.Contain(pointer.ButtonPrimary) {
					log.Info("primary")
				} else if e.Buttons.Contain(pointer.ButtonSecondary) {
					log.Info("secondary")
				} else if e.Buttons.Contain(pointer.ButtonTertiary) {
					log.Info("tertiary")
				}
				b.pressed = true
			case pointer.Release:
				b.pressed = false
			}
		}
	}
	return layout.Dimensions{
		Size:     gtx.Constraints.Max,
		Baseline: 0,
	}
}
