package commands

import (
	"errors"
	"github.com/charliego3/tools/store"
	"github.com/progrium/macdriver/macos/foundation"
	"github.com/progrium/macdriver/objc"
	"strings"
)

const (
	start = `tell application "Terminal"
    activate
`
	windowWaiting = `
    repeat while contents of window 1 starts with "Executing"
        delay 0.1
    end repeat
`
	end = `
end tell`
	doPrefix = `
    do script "`
	checkWindows = `
    if (count of windows) > 0 then
        set currentTab to first tab of first window
    else
        set currentTab to do script ""
    end if
`
	tabWaiting = `
    repeat
        delay 0.1
        if not busy of currentTab then
            exit repeat
        end if
    end repeat
`
)

type Terminal struct{}

func (a *Terminal) BundleIdentifier() string {
	return "com.apple.Terminal"
}

func (a *Terminal) Enabled() bool {
	return true
}

func (a *Terminal) Execute(m *store.Terminal) error {
	length := len(m.Cmds)
	if length == 0 {
		return errors.New("There are no commands to execute")
	}

	inNewWindow := m.Window == 0
	source := strings.Builder{}
	source.WriteString(start)

	if !inNewWindow {
		source.WriteString(checkWindows)
	}

	for i, cmd := range m.Cmds {
		if len(cmd) == 0 {
			continue
		}

		command := doPrefix + cmd + `"`
		if inNewWindow {
			if i > 0 {
				command += " in front window"
			}
		} else {
			command += " in currentTab"
		}
		if i > 0 {
			if inNewWindow {
			} else {
				command += ""
			}
		}
		source.WriteString(command)
		if length-1 != i {
			if inNewWindow {
				source.WriteString(windowWaiting)
			} else {
				source.WriteString(tabWaiting)
			}
		}
	}
	source.WriteString(end)
	//fmt.Println(source.String())

	errInfo := make(map[string]objc.IObject)
	script := foundation.NewAppleScriptWithSource(source.String())
	script.ExecuteAndReturnError(errInfo)
	//fmt.Printf("%+v", errInfo)
	return nil
}
