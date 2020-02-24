package main

import (
	"github.com/shurcooL/trayhost"
	"github.com/ying32/govcl/vcl"
	"github.com/ying32/govcl/vcl/types"
	"github.com/andlabs/ui"
)

func md5MenuItem() trayhost.MenuItem {
	return trayhost.MenuItem{
		Title: "MD5 Encrypt",
		Handler: func() {
			_ = ui.Main(func() {
				window := ui.NewWindow("MD5 Encrypt", 600, 600, false)
				window.OnClosing(func(window *ui.Window) bool {
					window.Hide()
					return true
				})
				ui.OnShouldQuit(func() bool {
					window.Destroy()
					return true
				})
				window.SetMargined(true)
				window.Handle()

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

func newConnectionButton() *ui.Button {
	button := ui.NewButton("New Connection😊")
	button.OnClicked(func(button *ui.Button) {
		window := ui.NewWindow("New Connection", 500, 500, true)
		window.OnClosing(func(window *ui.Window) bool {
			window.Destroy()
			return false
		})
		window.SetMargined(true)

		box := ui.NewVerticalBox()
		window.SetChild(box)

		tab := ui.NewTab()
		tab.Append("Basic Controls", newConnectionForm())
		tab.SetMargined(0, true)
		tab.Append("Numbers and Lists", ui.NewLabel("Numbers"))
		tab.SetMargined(1, true)
		tab.Append("Data Choosers", ui.NewLabel("Data Choosers"))
		tab.SetMargined(2, true)

		box.Append(tab, true)

		grid := ui.NewGrid()
		grid.SetPadded(true)
		testBtn := ui.NewButton("Test Connection")
		grid.Append(testBtn, 0, 0, 1, 1, false, ui.AlignStart, false, ui.AlignEnd)
		box.Append(grid, true)

		window.Show()
	})
	return button
}
func newConnectionForm() ui.Control {
	entryForm := ui.NewForm()
	entryForm.SetPadded(true)

	entryForm.Append("Entry", ui.NewEntry(), false)
	entryForm.Append("Password Entry", ui.NewPasswordEntry(), false)
	entryForm.Append("Search Entry", ui.NewSearchEntry(), false)
	return entryForm
}

type MD5Window struct {
	*vcl.TForm
	EncryptBtn *vcl.TButton
}

func (f *MD5Window) OnFormCreate(sender vcl.IObject) {
	f.SetCaption("MD5 Encrypt")
	f.ScreenCenter()
	f.SetWidth(100)
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

func (f *MD5Window) OnButtonClick(sender vcl.IObject) {
	vcl.ShowMessage("Clicked Button")
}