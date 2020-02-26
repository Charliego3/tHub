package main

import (
	"encoding/base64"
	"github.com/ProtonMail/ui"
	"github.com/shurcooL/trayhost"
)

var base64Entry *EncryptEntry

func base64MenuItem() trayhost.MenuItem {
	return trayhost.MenuItem{
		Title: base64Window.Title(),
		Handler: func() {
			base64Window.Show()
			base64Entry.Clear()
		},
	}
}

func base64OnReady(window *ui.Window) {
	base64Window = window
	base64Entry = encrypt(base64Window, func(text string) string {
		return base64.RawStdEncoding.EncodeToString([]byte(text))
	})
	base64Window.OnClosing(func(window *ui.Window) bool {
		if base64Entry.ResultLine != nil {
			base64Entry.ResultLine.Hide()
			window.Handle()
			window.SetContentSize(480, 84)
		}
		window.Hide()
		return false
	})
}
