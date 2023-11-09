package commands

import "github.com/charliego3/tools/store"

type Iterm2 struct{}

func (a *Iterm2) BundleIdentifier() string {
	return "com.googlecode.iterm2"
}

func (a *Iterm2) Enabled() bool {
	return false
}

func (a *Iterm2) Execute(m *store.Terminal) error {
	return nil
}
