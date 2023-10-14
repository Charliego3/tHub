package main

import (
	"github.com/progrium/macdriver/helper/action"
	"github.com/progrium/macdriver/macos/appkit"
	"github.com/progrium/macdriver/objc"
)

type terminal struct {
	appkit.MenuItem
}

func getTerminalItem() *terminal {
	t := &terminal{}
	t.MenuItem = appkit.NewMenuItem()
	t.MenuItem.SetTitle("Shell Commands")
	t.MenuItem.SetImage(getSymbolImage("terminal.fill", getImageScale(appkit.ImageSymbolScaleSmall)))
	t.MenuItem.SetKeyEquivalentModifierMask(appkit.EventModifierFlagCommand)
	t.MenuItem.SetKeyEquivalent("t")
	target, selector := action.Wrap(t.showAddWindow)
	t.MenuItem.SetTarget(target)
	t.MenuItem.SetAction(selector)

	menu := appkit.NewMenu()
	menu.AddItem(t.getAddCmdItem())
	t.MenuItem.SetSubmenu(menu)
	return t
}

func (t *terminal) getAddCmdItem() appkit.MenuItem {
	item := appkit.NewMenuItem()
	item.SetTitle("Command Settings")
	target, selector := action.Wrap(t.showAddWindow)
	item.SetTarget(target)
	item.SetAction(selector)
	return item
}

func (t *terminal) showAddWindow(_ objc.Object) {
	view := appkit.NewView()
	w := NewWindow("Command Settings", view, func(w appkit.Window) {
		w.SetContentSize(sizeOf(400, 200))
		w.SetLevel(appkit.NormalWindowLevel)
		w.SetStyleMask(w.StyleMask() | appkit.MiniaturizableWindowMask)
	})
	w.SetToolbarStyle(appkit.WindowToolbarStylePreference)
	w.MakeKeyAndOrderFront(nil)
}
