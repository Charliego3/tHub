package main

import (
	"fmt"
	"github.com/ProtonMail/ui"
)

func encrypt(window *ui.Window, f func(text string) string) *EncryptEntry {
	entry := &EncryptEntry{}
	vbox := ui.NewVerticalBox()
	vbox.SetPadded(true)

	line := ui.NewHorizontalBox()
	line.SetPadded(true)
	vbox.Append(line, false)

	form := ui.NewForm()
	form.SetPadded(true)
	entry.Resource = ui.NewEntry()
	entry.Resource.SetReadOnly(true)
	form.Append("Resource", entry.Resource, true)
	line.Append(form, true)

	l2 := ui.NewHorizontalBox()
	l2.SetPadded(true)
	vbox.Append(l2, false)
	entry.Clipboard = ui.NewCheckbox("From Clipboard?")
	entry.Clipboard.SetChecked(true)
	entry.Clipboard.OnToggled(func(checkbox *ui.Checkbox) {
		if checkbox.Checked() {
			entry.Resource.SetReadOnly(true)
		} else {
			entry.Resource.SetReadOnly(false)
		}
	})
	l2.Append(entry.Clipboard, true)

	button := ui.NewButton("Encrypt")
	button.OnClicked(func(button *ui.Button) {
		var text string
		if entry.Clipboard.Checked() {
			var err error
			text, err = getClipboard()
			if err != nil {
				ui.MsgBoxError(window,
					"Encountered an error while getting Clipboard content.",
					"Error details: "+fmt.Sprintf("error: %v\n", err))
				entry.Clear()
				return
			}
			if text == "" {
				ui.MsgBox(window,
					"This is a prompt message.",
					"There is nothing in the clipboard, you can copy the content and try again or manually enter the content that needs to be encrypted.")
				entry.Clear()
				return
			}
			entry.Resource.SetText(text)
		} else {
			text = entry.Resource.Text()
			if text == "" {
				ui.MsgBox(window,
					"This is a prompt message.",
					"Nothing is entered, you can copy the content and try again or manually enter the content that needs to be encrypted.")
				entry.Clear()
				return
			}
		}
		rs := f(text)
		if entry.ResultLine == nil {
			entry.ResultLine = ui.NewHorizontalBox()
			entry.ResultLine.SetPadded(true)
			entry.ResultLine.Append(ui.NewLabel("Result"), false)
			vbox.Append(entry.ResultLine, false)
		}
		if entry.Result != nil {
			entry.ResultLine.Delete(1)
		}
		entry.Result = ui.NewEntry()
		entry.Result.SetReadOnly(true)
		entry.ResultLine.Append(entry.Result, true)
		entry.Result.SetText(rs)
		entry.ResultLine.Show()
		setClipboard(rs)

		if entry.PromptLine == nil {
			entry.PromptLine = ui.NewHorizontalBox()
			entry.PromptLine.SetPadded(true)
			vbox.Append(entry.PromptLine, false)
		}
		if entry.Prompt != nil && !entry.PromptEmpty {
			entry.PromptLine.Delete(0)
		}
		entry.Prompt = ui.NewLabel("The content has been copied to the clipboard and can be used directly.")
		entry.PromptLine.Append(entry.Prompt, true)
		entry.PromptEmpty = false
	})
	line.Append(button, false)

	window.SetChild(vbox)
	return entry
}

type EncryptEntry struct {
	Resource    *ui.Entry
	Clipboard   *ui.Checkbox
	ResultLine  *ui.Box
	Result      *ui.Entry
	Prompt      *ui.Label
	PromptLine  *ui.Box
	PromptEmpty bool
}

func (m *EncryptEntry) Clear() {
	if m.Resource != nil {
		m.Resource.SetText("")
	}
	if m.Result != nil {
		m.Result.SetText("")
	}
	if m.Prompt != nil && m.PromptLine != nil && !m.PromptEmpty {
		m.PromptLine.Delete(0)
		m.PromptEmpty = true
	}
}
