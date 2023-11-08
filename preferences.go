package main

import (
	"github.com/charliego3/tools/utility"
	"github.com/progrium/macdriver/helper/action"
	"github.com/progrium/macdriver/macos/appkit"
	"github.com/progrium/macdriver/objc"
)

type PreferencesMenuItem struct {
	appkit.MenuItem
	showing bool
}

func getPreferencesItem() *PreferencesMenuItem {
	p := &PreferencesMenuItem{}
	p.MenuItem = appkit.NewMenuItem()
	p.MenuItem.SetTitle("Preferences")
	p.MenuItem.SetImage(utility.SymbolImage("gear"))
	p.MenuItem.SetKeyEquivalentModifierMask(appkit.EventModifierFlagCommand)
	p.MenuItem.SetKeyEquivalent(",")
	target, selector := action.Wrap(p.showWindow)
	p.SetTarget(target)
	p.SetAction(selector)
	return p
}

func (p *PreferencesMenuItem) showWindow(sender objc.Object) {
	if p.showing {
		return
	}

}
