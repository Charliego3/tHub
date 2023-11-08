package commands

import (
	"github.com/charliego3/tools/store"
	"github.com/charliego3/tools/utility"
	"github.com/progrium/macdriver/macos/appkit"
	"github.com/progrium/macdriver/macos/foundation"
)

type Executor interface {
	BundleIdentifier() string
	Enabled() bool
	Execute(m *store.Terminal) error
}

func wrapExecute(m *store.Terminal, callback func()) {
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
		cfg, func(_ appkit.RunningApplication, err foundation.Error) {
			if !err.IsNil() {
				utility.ModalAlert(nil, false, "Run script failed", err.LocalizedDescription())
				return
			}
			if callback == nil {
				return
			}
			callback()
		},
	)
}
