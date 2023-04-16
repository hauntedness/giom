package main

import (
	"os"

	"gioui.org/app"
	"github.com/hauntedness/giom/ui"
)

func main() {
	go func() {
		window := app.NewWindow(app.Title("Protonet"))
		if err := ui.Loop(window); err != nil {
			panic(err)
		}
		os.Exit(0)
	}()
	app.Main()
}
