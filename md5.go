package main

import (
	"crypto/md5"
	"encoding/hex"
	"github.com/ProtonMail/ui"
	"github.com/shurcooL/trayhost"
)

var md5Entry *EncryptEntry

func md5MenuItem() trayhost.MenuItem {
	return trayhost.MenuItem{
		Title: md5Window.Title(),
		Handler: func() {
			md5Window.Show()
			if md5Entry != nil {
				md5Entry.Clear()
			}
		},
	}
}

func md5OnReady(window *ui.Window) {
	md5Window = window
	md5Entry = encrypt(md5Window, func(text string) string {
		h := md5.New()
		h.Write([]byte(text))
		return hex.EncodeToString(h.Sum(nil))
	})
	md5Window.OnClosing(func(window *ui.Window) bool {
		md5Entry.ResultLine.Hide()
		window.Hide()
		return false
	})
}
