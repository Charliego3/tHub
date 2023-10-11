package main

import (
	"bufio"
	"bytes"
	"errors"
	"fmt"
	"os"
	"os/exec"

	"github.com/progrium/macdriver/helper/action"
	"github.com/progrium/macdriver/macos/appkit"
	"github.com/progrium/macdriver/macos/foundation"
	"github.com/progrium/macdriver/objc"
)

const (
	emulatorPathKey = "emulatorPath"
	executableName  = "emulator"
)

func getEmulatorItem() appkit.MenuItem {
	avds := getEmulatorList(executableName)
	item := appkit.NewMenuItem()
	item.SetImage(getSymbolImage("iphone.smartbatterycase.gen2"))
	if len(avds) == 0 {
		item.SetTitle("Setup emulator path")
		target, selector := setEmulatorPath(item)
		item.SetTarget(target)
		item.SetAction(selector)
		item.SetToolTip("Add emulator to the PATH environment variable or manually select the path of emulator")
		return item
	}
	setEmulatorSubItem(item, executableName)
	return item
}

func start(name string) (action.Target, objc.Selector) {
	return action.Wrap(func(_ objc.Object) {
		exec.Command(executableName, "-avd", name).Start()
	})
}

func setEmulatorSubItem(item appkit.MenuItem, path string) {
	avds := getEmulatorList(path)

	menu := appkit.NewMenu()
	for _, avd := range avds {
		subItem := appkit.NewMenuItem()
		subItem.SetTitle(avd)
		target, selector := start(avd)
		subItem.SetTarget(target)
		subItem.SetAction(selector)
	}
	item.SetToolTip("")
	item.SetTitle("Android Emulators")
	item.SetSubmenu(menu)
}

func getEmulatorList(path string) []string {
	cmd := exec.Command(path, "-list-avds")
	cmd.Dir, _ = os.UserHomeDir()
	bs, err := cmd.Output()
	if err == nil {
		var avds []string
		scanner := bufio.NewScanner(bytes.NewReader(bs))
		for scanner.Scan() {
			avd := scanner.Text()
			if len(avd) == 0 {
				continue
			}
			avds = append(avds, avd)
		}
		return avds
	}

	if path == executableName && errors.Is(err, exec.ErrNotFound) {
		fmt.Println(err)
		path = appkit.UserDefaultsController_SharedUserDefaultsController().
			Defaults().StringForKey(emulatorPathKey)
		if len(path) == 0 {
			fmt.Println("not configure emulator path", path)
			return nil
		}
		return getEmulatorList(path)
	}
	fmt.Println(err)
	return nil
}

func setEmulatorPath(item appkit.MenuItem) (action.Target, objc.Selector) {
	return action.Wrap(func(sender objc.Object) {
		panel := appkit.OpenPanel_OpenPanel()
		panel.Center()
		panel.SetCanChooseDirectories(true)
		panel.SetTitle("Choose Android Emulator Path")
		panel.SetTitleVisibility(appkit.WindowTitleVisible)
		panel.BeginWithCompletionHandler(func(result appkit.ModalResponse) {
			if result == appkit.ModalResponseOK {
				path := panel.URL().Path()
				appkit.UserDefaultsController_SharedUserDefaultsController().Defaults().
					SetObjectForKey(foundation.String_StringWithString(path), emulatorPathKey)
				setEmulatorSubItem(item, path)
			}
		})
		panel.OrderFront(nil)
	})
}
