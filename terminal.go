package main

import (
	"github.com/progrium/macdriver/helper/action"
	"github.com/progrium/macdriver/macos/appkit"
	"github.com/progrium/macdriver/objc"
)

type terminal struct {
	appkit.MenuItem
	hasOpened bool
}

func getTerminalItem() *terminal {
	t := &terminal{}
	t.MenuItem = appkit.NewMenuItem()
	t.MenuItem.SetTitle("Terminal Commands")
	t.MenuItem.SetImage(t.getImage())
	t.MenuItem.SetKeyEquivalentModifierMask(appkit.EventModifierFlagCommand)
	t.MenuItem.SetKeyEquivalent("t")
	target, selector := action.Wrap(t.addCommand)
	t.MenuItem.SetTarget(target)
	t.MenuItem.SetAction(selector)
    t.MenuItem.Image().ImageWithSymbolConfiguration(appkit.ImageSymbolConfiguration_ConfigurationWithScale(appkit.ImageSymbolScaleLarge))
	return t
}

func (t terminal) getImage() appkit.Image {
	name := "terminal"
	if t.hasOpened {
		name += ".fill"
	}
	return getSymbolImage(name, getImageScale(appkit.ImageSymbolScaleSmall))
}

func (t *terminal) addCommand(_ objc.Object) {
	t.hasOpened = !t.hasOpened
	t.MenuItem.SetImage(t.getImage())
}
