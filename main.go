package main

import (
	"github.com/ProtonMail/ui"
	"github.com/shurcooL/trayhost"
	"os"
)

var (
	md5Window    *ui.Window
	base64Window *ui.Window
	exportWindow *ExportWindow
)

func main() {
	err := ui.Main(func() {
		md5OnReady(createWindow("MD5 Encrypt", 480, 84))
		base64OnReady(createWindow("Base64 Encrypt", 480, 84))
		exportOnReady(createWindow("Export Excel From MySQL", 608, 115))
		menus()
	})
	if err != nil {
		os.Exit(1)
	}
}

func menus() {
	menuItems := []trayhost.MenuItem{
		md5MenuItem(),
		base64MenuItem(),
		exportMenu(),
		trayhost.SeparatorMenuItem(),
		quitMenu(),
	}
	trayhost.Initialize("Tools", iconData, menuItems)
	trayhost.EnterLoop()
}

func quitMenu() trayhost.MenuItem {
	return trayhost.MenuItem{
		Title:   "Quit",
		Handler: trayhost.Exit,
	}
}

/*
{
	// Displaying notifications requires a proper app bundle and won't work without one.
	// See https://godoc.org/github.com/shurcooL/trayhost#hdr-Notes.
	Title: "Display Notification",
	Handler: func() {
		notification := trayhost.Notification{
			Title:   "Example Notification",
			Body:    "Notification body text is here.",
			Timeout: 3 * time.Second,
			Handler: func() {
				fmt.Println("do stuff when notification is clicked")
			},
		}
		if cc, err := trayhost.GetClipboardContent(); err == nil && cc.Image.Kind != "" {
			// Use image from clipboard as notification image.
			notification.Image = cc.Image
		}
		notification.Display()
	},
},
*/
