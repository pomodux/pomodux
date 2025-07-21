package main

import (
	"github.com/rivo/tview"
)

func main() {
	app := tview.NewApplication()
	modal := tview.NewModal().
		SetText("If you see this, tview is working!").
		AddButtons([]string{"OK"}).
		SetDoneFunc(func(buttonIndex int, buttonLabel string) {
			app.Stop()
		})
	_ = app.SetRoot(modal, false).Run()
}
