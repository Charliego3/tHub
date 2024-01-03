package commands

import (
	"errors"
	"github.com/charliego3/thub/store"
	"github.com/charliego3/thub/utility"
	"github.com/progrium/macdriver/dispatch"
	"github.com/progrium/macdriver/macos/appkit"
	"github.com/progrium/macdriver/macos/foundation"
	"os"
	"os/exec"
	"time"
)

type Executor interface {
	BundleIdentifier() string
	Enabled() bool
	Execute(m *store.Terminal) error
}

func activate(m *store.Terminal, callback func()) {
	cfg := appkit.NewWorkspaceOpenConfiguration()
	cfg.SetPromptsUserIfNeeded(true)
	cfg.SetActivates(true)
	if m.Window == 0 {
		cfg.SetCreatesNewApplicationInstance(true)
	} else {
		cfg.SetAllowsRunningApplicationSubstitution(true)
	}

	workspace := appkit.Workspace_SharedWorkspace()
	workspace.OpenApplicationAtURLConfigurationCompletionHandler(
		workspace.URLForApplicationWithBundleIdentifier(m.App),
		cfg, func(app appkit.RunningApplication, err foundation.Error) {
			if !err.IsNil() {
				utility.ShowAlert(nil, false, "Run script failed", err.LocalizedDescription())
				return
			}
			if callback == nil {
				return
			}

			for !app.IsFinishedLaunching() {
				time.Sleep(time.Millisecond * 100)
			}
			callback()
		},
	)
}

func doCmd(command string, tips ...string) []byte {
	output, err := exec.Command("bash", "-c", command).Output()
	if err != nil {
		title := "Failed to run Script"
		if len(tips) > 0 {
			title = tips[0]
		}
		desc := err.Error()
		var e *exec.ExitError
		if errors.As(err, &e) {
			desc = string(e.Stderr)
		}
		dispatch.MainQueue().DispatchAsync(func() {
			utility.ShowAlert(nil, false, title, desc)
		})
	}
	return output
}

func toShell(prefix, to string) bool {
	shell := os.Getenv("SHELL")
	if len(shell) > 0 && shell == to {
		return false
	}
	doCmd(prefix + to + "\r\n")
	return true
}
