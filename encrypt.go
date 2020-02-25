package main

import (
	"fmt"
	"github.com/ProtonMail/ui"
)

var (
	lr *ui.Box
	r  *ui.Entry

	lp *ui.Box
	p  *ui.Label
	pr bool

	resource *ui.Entry
)

func encrypt(window *ui.Window, f func(text string) string) {
	vbox := ui.NewVerticalBox()
	vbox.SetPadded(true)

	line := ui.NewHorizontalBox()
	line.SetPadded(true)
	vbox.Append(line, false)

	form := ui.NewForm()
	form.SetPadded(true)
	resource = ui.NewEntry()
	resource.SetReadOnly(true)
	form.Append("Resource", resource, true)
	line.Append(form, true)

	l2 := ui.NewHorizontalBox()
	l2.SetPadded(true)
	vbox.Append(l2, false)
	checkbox := ui.NewCheckbox("From Clipboard?")
	checkbox.SetChecked(true)
	checkbox.OnToggled(func(checkbox *ui.Checkbox) {
		if checkbox.Checked() {
			resource.SetReadOnly(true)
		} else {
			resource.SetReadOnly(false)
		}
	})
	l2.Append(checkbox, true)

	button := ui.NewButton("Encrypt")
	button.OnClicked(func(button *ui.Button) {
		var text string
		if checkbox.Checked() {
			var err error
			text, err = getClipboard()
			if err != nil {
				ui.MsgBoxError(window,
					"Encountered an error while getting Clipboard content.",
					"Error details: "+fmt.Sprintf("error: %v\n", err))
				if !pr {
					lp.Delete(0)
					pr = true
				}
				clear(resource, r)
				return
			}
			if text == "" {
				ui.MsgBox(window,
					"This is a prompt message.",
					"There is nothing in the clipboard, you can copy the content and try again or manually enter the content that needs to be encrypted.")
				if !pr {
					lp.Delete(0)
					pr = true
				}
				clear(resource, r)
				return
			}
			resource.SetText(text)
		} else {
			text = resource.Text()
			if text == "" {
				ui.MsgBox(window,
					"This is a prompt message.",
					"Nothing is entered, you can copy the content and try again or manually enter the content that needs to be encrypted.")
				if !pr {
					lp.Delete(0)
					pr = true
				}
				clear(resource, r)
				return
			}
		}
		rs := f(text)
		if lr == nil {
			lr = ui.NewHorizontalBox()
			lr.SetPadded(true)
			lr.Append(ui.NewLabel("Result"), false)
			vbox.Append(lr, false)
		}
		if r != nil {
			lr.Delete(1)
		}
		r = ui.NewEntry()
		r.SetReadOnly(true)
		lr.Append(r, true)
		r.SetText(rs)
		setClipboard(rs)

		if lp == nil {
			lp = ui.NewHorizontalBox()
			lp.SetPadded(true)
			vbox.Append(lp, false)
		}
		if p != nil && !pr {
			lp.Delete(0)
		}
		p = ui.NewLabel("The content has been copied to the clipboard and can be used directly.")
		lp.Append(p, true)
		pr = false
	})
	line.Append(button, false)

	window.SetChild(vbox)
}

func clear(resource, r *ui.Entry) {
	if resource != nil {
		resource.SetText("")
	}
	if r != nil {
		r.SetText("")
	}
}

type Encrypt interface {
	Clear()
}
