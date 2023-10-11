package main

import (
	"fmt"

	"github.com/progrium/macdriver/helper/action"
	"github.com/progrium/macdriver/macos"
	"github.com/progrium/macdriver/macos/appkit"
	"github.com/progrium/macdriver/objc"
)

func main() {
	macos.RunApp(launched)
}

func launched(app appkit.Application, delegate *appkit.ApplicationDelegate) {
	delegate.SetApplicationShouldTerminateAfterLastWindowClosed(func(appkit.Application) bool {
		return false
	})

	item := appkit.StatusBar_SystemStatusBar().StatusItemWithLength(-1)
	objc.Retain(&item)
	item.Button().SetImage(getSymbolImage(
		"command.square.fill",
		appkit.ImageSymbolConfiguration_ConfigurationWithScale(appkit.ImageSymbolScaleLarge),
	))

	quit := appkit.NewMenuItem()
	quit.SetTitle("Quit")
	quit.SetAction(objc.Sel("terminate:"))
	quit.SetKeyEquivalentModifierMask(appkit.EventModifierFlagCommand)
	quit.SetAllowsAutomaticKeyEquivalentMirroring(true)
	quit.SetKeyEquivalent("q")

	menu := appkit.NewMenu()
	menu.AddItem(getGeneratePasswordItem(item))
	menu.AddItem(getEmulatorItem())
	menu.AddItem(getTerminalItem())
	menu.AddItem(appkit.MenuItem_SeparatorItem())
	menu.AddItem(getAutoStart())
	menu.AddItem(quit)
	item.SetMenu(menu)

	app.SetActivationPolicy(appkit.ApplicationActivationPolicyRegular)
	app.ActivateIgnoringOtherApps(true)
	fmt.Println("started")
}

func getAutoStart() appkit.MenuItem {
	item := appkit.NewMenuItem()
	item.SetImage(getSymbolImage("autostartstop"))
	item.SetTitle("Launched at login")
	target, selector := action.Wrap(func(objc.Object) {})
	item.SetEnabled(true)
	item.SetTarget(target)
	item.SetAction(selector)
	return item
}
