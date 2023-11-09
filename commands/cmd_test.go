package commands

import (
	"fmt"
	"github.com/charliego3/tools/store"
	"testing"
)

func TestKittyConf(t *testing.T) {
	k := &Kitty{}
	remote, bytes := k.remoteControlEnabled(getKittyConf())
	fmt.Println(remote)
	fmt.Println(string(bytes))
}

func TestKitty_Execute(t *testing.T) {
	k := &Kitty{}
	m := &store.Terminal{
		Name:   "Test Kitty",
		App:    k.BundleIdentifier(),
		Shell:  "/usr/local/bin/fish",
		Window: 1,
		Cmds: []string{
			"pwd",
			"sleep 5",
			"ll",
		},
	}
	_ = k.Execute(m)
}
