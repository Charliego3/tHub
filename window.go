package main

import (
	"github.com/ProtonMail/ui"
	"github.com/shurcooL/trayhost"
	"os"
	"runtime"
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
