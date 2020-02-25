package main

import (
	"crypto/md5"
	"encoding/hex"
	"github.com/ProtonMail/ui"
	"github.com/shurcooL/trayhost"
)

func md5MenuItem() trayhost.MenuItem {
	return trayhost.MenuItem{
		Title: "MD5 Encrypt",
		Handler: func() {
			md5Window.Show()
			clear(resource, r)
		},
	}
}

func md5OnReady(window *ui.Window) {
	md5Window = window
	encrypt(md5Window, func(text string) string {
		h := md5.New()
		h.Write([]byte(text))
		return hex.EncodeToString(h.Sum(nil))
	})
}

type MD5Entry struct {
	Resource  *ui.Entry
	Clipboard *ui.Checkbox
	Result    *ui.Entry
	Prompt    *ui.Label
}

func (m *MD5Entry) Clear() {
	if m.Resource != nil {
		m.Resource.SetText("")
	}
	if m.Result != nil {
		m.Result.SetText("")
	}
	if m.Prompt != nil {

	}
}
