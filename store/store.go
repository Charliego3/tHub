package store

import (
	"encoding/json"
	"errors"
	"github.com/progrium/macdriver/dispatch"
	"github.com/progrium/macdriver/macos/appkit"
	"io/fs"
	"os"
	"path/filepath"
	"sync"
)

type Option struct {
	Terminals map[int64]Terminal
}

type Terminal struct {
	Name   string
	App    string
	Shell  string
	Window int
	Cmds   []string
}

var (
	once     sync.Once
	instance *Option
	path     string
)

func Fetch() *Option {
	once.Do(func() {
		path = os.Getenv("$XDG_CONFIG_HOME")
		if path == "" {
			home, _ := os.UserHomeDir()
			path = filepath.Join(home, ".config")
		}

		path = filepath.Join(path, "thub", "settings.json")
		_, err := os.Stat(path)
		instance = &Option{Terminals: make(map[int64]Terminal)}
		if errors.Is(err, fs.ErrNotExist) {
			showError(instance.Save())
			return
		}

		bs, err := os.ReadFile(path)
		if err != nil {
			showError(err)
			return
		}
		showError(json.Unmarshal(bs, instance))
	})
	return instance
}

func showError(err error) {
	if err == nil {
		return
	}
	dispatch.MainQueue().DispatchAsync(func() {
		dialog := appkit.NewAlert()
		dialog.SetAlertStyle(appkit.AlertStyleCritical)
		dialog.SetMessageText("Failed read settings")
		dialog.SetInformativeText(err.Error())
		dialog.RunModal()
	})
	_ = os.Remove(path)
}

func (o *Option) Save() error {
	bs, err := json.MarshalIndent(instance, "", "    ")
	if err != nil {
		return err
	}

	dir := filepath.Dir(path)
	err = os.MkdirAll(dir, os.ModePerm)
	if err != nil {
		return err
	}

	return os.WriteFile(path, bs, os.ModePerm)
}
