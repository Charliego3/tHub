package commands

import (
	"github.com/charliego3/tools/store"
)

type Kitty struct{}

func (a *Kitty) BundleIdentifier() string {
	return "net.kovidgoyal.kitty"
}

func (a *Kitty) Enabled() bool {
	return true
}

func (a *Kitty) Execute(m *store.Terminal) error {
	wrapExecute(m, func() {

	})
	return nil
}
