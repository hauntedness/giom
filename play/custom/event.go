package main

import (
	"image"
	"image/color"

	"gioui.org/font/gofont"
	"gioui.org/io/event"
	"gioui.org/io/pointer"
	"gioui.org/layout"
	"gioui.org/op"
	"gioui.org/op/clip"
	"gioui.org/op/paint"
	"gioui.org/widget/material"
)

var Pressed = false

func DrawButton(ops *op.Ops, gtx layout.Context, tag event.Tag) {
	// Process events that arrived between the last frame and this one.
	for _, event := range gtx.Events(tag) {
		if x, ok := event.(pointer.Event); ok {
			switch x.Type {
			case pointer.Press:
				Pressed = true
			case pointer.Release:
				Pressed = false
			}
		}
	}

	// Confine the area of interest to a 100x100 rectangle.
	defer clip.Rect{Max: image.Pt(100, 100)}.Push(ops).Pop()

	// Declare the tag.
	pointer.InputOp{
		Tag:   tag,
		Types: pointer.Press | pointer.Release,
	}.Add(ops)

	var c color.NRGBA
	if Pressed {
		c = color.NRGBA{R: 0xFF, A: 0xFF}
	} else {
		c = color.NRGBA{G: 0xFF, A: 0xFF}
	}
	paint.ColorOp{Color: c}.Add(ops)
	paint.PaintOp{}.Add(ops)
	material.Body1(material.NewTheme(gofont.Collection()), "hahaha").Layout(gtx)
}
