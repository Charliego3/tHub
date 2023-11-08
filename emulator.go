package main

import (
	"bufio"
	"bytes"
	"errors"
	"fmt"
	"github.com/charliego3/tools/utility"
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
	item.SetImage(utility.SymbolImage("iphone.smartbatterycase.gen2"))
	if len(avds) == 0 {
		item.SetTitle("Setup emulator path")
		target, selector := setEmulatorPath(item)
		item.SetTarget(target)
		item.SetAction(selector)
		item.SetToolTip("Add emulator to the PATH environment variable or manually select the path of emulator")
		return item
	}
	setEmulatorSubItem(item, avds)
	return item
}

func start(item appkit.MenuItem, name string) (action.Target, objc.Selector) {
	return action.Wrap(func(_ objc.Object) {
		cmd := exec.Command(executableName, "-avd", name)
		err := cmd.Start()
		if err != nil {
			dialog := appkit.NewAlert()
			dialog.SetAlertStyle(appkit.AlertStyleWarning)
			dialog.SetIcon(appkit.Image_ImageNamed(appkit.ImageNameStatusUnavailable))
			dialog.SetInformativeText("Start fialed with " + name)
			dialog.SetMessageText(err.Error())
			dialog.AddButtonWithTitle("OK")
			dialog.RunModal()
			return
		}
		go func() {
			_ = cmd.Wait()
			item.SetState(appkit.ControlStateValueOff)
		}()
		item.SetState(appkit.ControlStateValueOn)
	})
}

func setEmulatorSubItem(item appkit.MenuItem, avds []string) {
	menu := appkit.NewMenu()
	for _, avd := range avds {
		subItem := appkit.NewMenuItem()
		subItem.SetTitle(avd)
		target, selector := start(subItem, avd)
		subItem.SetTarget(target)
		subItem.SetAction(selector)
		subItem.SetOnStateImage(utility.SymbolImage("circle.inset.filled",
			appkit.ImageSymbolConfiguration_ConfigurationWithScale(appkit.ImageSymbolScaleSmall),
			appkit.ImageSymbolConfiguration_ConfigurationWithHierarchicalColor(appkit.Color_SystemGreenColor()),
		))
		menu.AddItem(subItem)
	}
	item.SetToolTip("")
	item.SetTitle("Android Emulators")
	item.SetSubmenu(menu)
}

func getEmulatorList(path string) []string {
	cmd := exec.Command(path, "-list-avds")
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
				setEmulatorSubItem(item, getEmulatorList(path))
			}
		})
		panel.OrderFront(nil)
	})
}
