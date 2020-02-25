package main

import (
	"github.com/ProtonMail/ui"
	"github.com/shurcooL/trayhost"
)

func createWindow(title string, width, height int) *ui.Window {
	window := ui.NewWindow(title, width, height, true)
	window.Center()
	window.SetMargined(true)
	return window
}

func getClipboard() (string, error) {
	cc, err := trayhost.GetClipboardContent()
	if err != nil {
		return "", err
	}
	return cc.Text, nil
}

func setClipboard(text string) {
	trayhost.SetClipboardText(text)
}

type Window interface {
	Clear()
}
