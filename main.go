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

	bar := appkit.StatusBar_SystemStatusBar().StatusItemWithLength(-1)
	objc.Retain(&bar)
	bar.Button().SetImage(getSymbolImage(
		"command.square.fill",
		appkit.ImageSymbolConfiguration_ConfigurationWithScale(appkit.ImageSymbolScaleLarge),
	))

	menu := appkit.NewMenu()
	menu.AddItem(getGeneratePasswordItem(bar))
	menu.AddItem(getEmulatorItem())
	menu.AddItem(getTerminalItem())
	menu.AddItem(appkit.MenuItem_SeparatorItem())
	menu.AddItem(getAutoStart())
	menu.AddItem(getPreferencesItem())
	menu.AddItem(getQuit())
	bar.SetMenu(menu)

	app.SetActivationPolicy(appkit.ApplicationActivationPolicyRegular)
	app.ActivateIgnoringOtherApps(true)
}

func getQuit() appkit.MenuItem {
	quit := appkit.NewMenuItem()
	quit.SetTitle("Quit")
	quit.SetImage(getSymbolImage("power"))
	quit.SetAction(objc.Sel("terminate:"))
	quit.SetKeyEquivalentModifierMask(appkit.EventModifierFlagCommand)
	quit.SetAllowsAutomaticKeyEquivalentMirroring(true)
	quit.SetKeyEquivalent("q")
	return quit
}

func getAutoStart() appkit.MenuItem {
	item := appkit.NewMenuItem()
	setImage := func(autoStartup bool) {
		fmt.Println(autoStartup)
		symbol := "autostartstop"
		if !autoStartup {
			symbol = "autostartstop.slash"
		}
		item.SetImage(getSymbolImage(symbol))
	}

	const key = "launchAtStartup"
	defaults := appkit.UserDefaultsController_SharedUserDefaultsController().Defaults()
	setImage(defaults.BoolForKey(key))
	item.SetTitle("Launch at startup")
	target, selector := action.Wrap(func(objc.Object) {
		autoStartup := defaults.BoolForKey(key)
		if autoStartup {
			defaults.RemoveObjectForKey(key)
			appkit.Application_SharedApplication().DisableRelaunchOnLogin()
		} else {
			defaults.SetBoolForKey(true, key)
			appkit.Application_SharedApplication().EnableRelaunchOnLogin()
		}
		setImage(!autoStartup)
	})
	item.SetTarget(target)
	item.SetAction(selector)
	return item
}
