package main

import (
	"encoding/base64"
	"github.com/ProtonMail/ui"
	"github.com/shurcooL/trayhost"
)

func base64MenuItem() trayhost.MenuItem {
	return trayhost.MenuItem{
		Title: "Base64 Encrypt",
		Handler: func() {
			base64Window.Show()
		},
	}
}

func base64OnReady(window *ui.Window) {
	base64Window = window
	encrypt(base64Window, func(text string) string {
		return base64.RawStdEncoding.EncodeToString([]byte(text))
	})
}
