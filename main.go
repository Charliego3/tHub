package main

import (
	"github.com/charliego3/tools/commands"
	"github.com/charliego3/tools/utility"
	"github.com/progrium/macdriver/helper/action"
	"github.com/progrium/macdriver/macos/appkit"
	"github.com/progrium/macdriver/macos/foundation"
	"github.com/progrium/macdriver/objc"
	"runtime"
)

var MenuIcon string

func main() {
	if len(MenuIcon) == 0 {
		MenuIcon = "square.dashed.inset.filled"
	}

	runtime.LockOSThread()
	defer runtime.UnlockOSThread()
	app := appkit.Application_SharedApplication()
	delegate := &appkit.ApplicationDelegate{}
	delegate.SetApplicationDidFinishLaunching(func(notification foundation.Notification) {
		launched(app, delegate)
	})
	delegate.SetApplicationWillFinishLaunching(func(foundation.Notification) {
		setMainMenu(app)
	})
	delegate.SetApplicationShouldTerminateAfterLastWindowClosed(func(appkit.Application) bool {
		return true
	})
	app.SetDelegate(delegate)
	app.Run()
}

func launched(app appkit.Application, delegate *appkit.ApplicationDelegate) {
	delegate.SetApplicationShouldTerminateAfterLastWindowClosed(func(appkit.Application) bool {
		return false
	})

	bar := appkit.StatusBar_SystemStatusBar().StatusItemWithLength(-1)
	objc.Retain(&bar)
	bar.Button().SetImage(utility.SymbolImage(
		MenuIcon,
		appkit.ImageSymbolConfiguration_ConfigurationWithScale(appkit.ImageSymbolScaleLarge),
	))

	menu := appkit.NewMenu()
	menu.AddItem(getGeneratePasswordItem(bar))
	menu.AddItem(getEmulatorItem())
	menu.AddItem(commands.Item())
	menu.AddItem(appkit.MenuItem_SeparatorItem())
	menu.AddItem(getAutoStart())
	menu.AddItem(getPreferencesItem())
	menu.AddItem(getWebsite())
	menu.AddItem(getQuit())
	bar.SetMenu(menu)

	app.SetActivationPolicy(appkit.ApplicationActivationPolicyRegular)
	app.ActivateIgnoringOtherApps(true)
}

func setMainMenu(app appkit.Application) {
	menu := appkit.NewMenuWithTitle("main")
	app.SetMainMenu(menu)

	item := appkit.NewMenuItemWithSelector("", "", objc.Selector{})
	edit := appkit.NewMenuWithTitle("Edit")
	edit.AddItem(appkit.NewMenuItemWithSelector("Select All", "a", objc.Sel("selectAll:")))
	edit.AddItem(appkit.MenuItem_SeparatorItem())
	edit.AddItem(appkit.NewMenuItemWithSelector("Copy", "c", objc.Sel("copy:")))
	edit.AddItem(appkit.NewMenuItemWithSelector("Paste", "v", objc.Sel("paste:")))
	edit.AddItem(appkit.NewMenuItemWithSelector("Cut", "x", objc.Sel("cut:")))
	edit.AddItem(appkit.NewMenuItemWithSelector("Undo", "z", objc.Sel("undo:")))
	edit.AddItem(appkit.NewMenuItemWithSelector("Redo", "Z", objc.Sel("redo:")))
	item.SetSubmenu(edit)
	menu.AddItem(item)
}

func getWebsite() appkit.MenuItem {
	item := appkit.NewMenuItem()
	item.SetTitle("Visit Website")
	item.SetImage(utility.SymbolImage("safari"))
	action.Set(item, func(sender objc.Object) {
		url := foundation.URL_URLWithString("https://github.com/Charliego3/tools")
		workspace := appkit.Workspace_SharedWorkspace()
		workspace.OpenURL(url)
	})
	return item
}

func getQuit() appkit.MenuItem {
	quit := appkit.NewMenuItem()
	quit.SetTitle("Quit")
	quit.SetImage(utility.SymbolImage("power"))
	quit.SetAction(objc.Sel("terminate:"))
	quit.SetKeyEquivalentModifierMask(appkit.EventModifierFlagCommand)
	quit.SetAllowsAutomaticKeyEquivalentMirroring(true)
	quit.SetKeyEquivalent("q")
	return quit
}

func getAutoStart() appkit.MenuItem {
	item := appkit.NewMenuItem()
	setImage := func(autoStartup bool) {
		symbol := "autostartstop"
		if !autoStartup {
			symbol = "autostartstop.slash"
		}
		item.SetImage(utility.SymbolImage(symbol))
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
