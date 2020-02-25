package main

import (
	"github.com/ProtonMail/ui"
	"github.com/shurcooL/trayhost"
)

func exportMenu() trayhost.MenuItem {
	return trayhost.MenuItem{
		Title: "Export Data From MySQL",
		Handler: func() {

		},
	}
}

func exportOnReady(window *ui.Window) {

}
