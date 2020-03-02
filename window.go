package main

import (
	"github.com/ProtonMail/ui"
	"github.com/shurcooL/trayhost"
	"os"
	"runtime"
	"strconv"
	"time"
)

func createWindow(title string, width, height int) *ui.Window {
	window := ui.NewWindow(title, width, height, false)
	window.Center()
	window.SetMargined(true)
	return window
}

func notification(title, body string, handler func()) {
	notification := trayhost.Notification{
		Title:   title,
		Body:    body,
		Timeout: 3 * time.Second,
		Handler: handler,
	}
	if cc, err := trayhost.GetClipboardContent(); err == nil && cc.Image.Kind != "" {
		// Use image from clipboard as notification image.
		notification.Image = cc.Image
	}
	notification.Display()
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

func downloadPath() string {
	dir, err := os.UserHomeDir()
	if err != nil {
		switch runtime.GOOS {
		case "windows":
			dir = "C:/"
		default:
			dir = "/"
		}
	}
	switch runtime.GOOS {
	case "windows":
		dir += ""
	default:
		dir += "/Downloads"
	}
	return dir
}

type Window interface {
	Clear()
}

type ExportWindow struct {
	*ui.Window
	Showing bool
}

func (e *ExportWindow) Hide() {
	e.Window.Hide()
	e.Showing = false
}

func (e *ExportWindow) Show() {
	if !exportWindow.Showing {
		if exportEntry != nil {
			exportEntry.Clear()
			e.hideProgress()
		}
		exportWindow.Window.Show()
		exportWindow.Showing = true
		for i, p := range prompts {
			exportEntry.PromptLabels[i].SetText(p)
		}
	}
}

func (e *ExportWindow) showProgress() {
	exportPadding.Hide()
	progressBar.SetValue(0)
	progressBar.Show()
	progressVal.SetText(" 0% ")
	progressVal.Show()
}

func (e *ExportWindow) hideProgress() {
	exportPadding.Show()
	progressBar.Hide()
	progressVal.Hide()
}

func (e *ExportWindow) progressFinish() {
	progressBar.SetValue(100)
	progressVal.SetText("100%")
}

func (e *ExportWindow) setProgress(progress int) {
	progressBar.SetValue(progress)
	p := strconv.Itoa(progress)
	length := len([]rune(p))
	if length == 1 {
		p = " " + p + " %"
	} else if length == 2 {
		p = " " + p + "%"
	} else {
		p += "%"
	}
	progressVal.SetText(p)
	progressBar.Handle()
	progressVal.Handle()
}
