package commands

import "github.com/charliego3/thub/store"

type Alacritty struct{}

func (a *Alacritty) BundleIdentifier() string {
	return "org.alacritty"
}

func (a *Alacritty) Enabled() bool {
	return false
}

func (a *Alacritty) Execute(m *store.Terminal) error {
	return nil
}
