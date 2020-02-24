package main

import (
	"github.com/andlabs/ui"
	"github.com/shurcooL/trayhost"
	"github.com/ying32/govcl/vcl"
	"github.com/ying32/govcl/vcl/types"
)

func base64MenuItem() trayhost.MenuItem {
	return trayhost.MenuItem{
		Title:   "Base64 Encrypt",
		Handler: func() {
			_ = ui.Main(func() {
				window := ui.NewWindow("MongoDB Client", 600, 600, false)
				window.OnClosing(func(window *ui.Window) bool {
					window.Hide()
					return true
				})
				ui.OnShouldQuit(func() bool {
					window.Destroy()
					return true
				})

				vbox := ui.NewVerticalBox()
				vbox.SetPadded(true)

				hbox := ui.NewHorizontalBox()
				hbox.SetPadded(true)
				vbox.Append(hbox, false)

				hbox.Append(ui.NewButton("Button"), false)
				hbox.Append(ui.NewCheckbox("Checkbox"), false)

				vbox.Append(ui.NewLabel("This is a label. Right now, labels can only span one line."), false)

				vbox.Append(ui.NewHorizontalSeparator(), false)
				vbox.Append(newConnectionButton(), false)

				window.SetChild(vbox)
				//window.SetChild(mainBox(window))
				window.Show()
			})
		},
	}
}

type Base64Window struct {
	*vcl.TForm
	EncryptBtn *vcl.TButton
}

func (f *Base64Window) OnFormCreate(sender vcl.IObject) {
	f.SetCaption("Base64 Encrypt")
	f.ScreenCenter()
	f.SetWidth(300)
	f.SetHeight(100)
	f.SetOnClose(func(sender vcl.IObject, action *types.TCloseAction) {
		f.Hide()
	})
	f.EncryptBtn = vcl.NewButton(f)
	f.EncryptBtn.SetParent(f)
	f.EncryptBtn.SetTop(10)
	f.EncryptBtn.SetLeft(10)
	f.EncryptBtn.SetCaption("Encrypt")
	f.EncryptBtn.SetOnClick(f.OnButtonClick)
}

func (f *Base64Window) OnButtonClick(sender vcl.IObject) {
	vcl.ShowMessage("Clicked Button")
}