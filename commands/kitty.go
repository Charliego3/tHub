package commands

import (
	"bufio"
	"bytes"
	"encoding/json"
	"os"
	"path/filepath"
	"strconv"
	"strings"
	"time"

	"github.com/charliego3/thub/store"
	"github.com/charliego3/thub/utility"
	"github.com/progrium/macdriver/macos/appkit"
	"github.com/progrium/macdriver/macos/foundation"
)

type Kitty struct {
	listening string
}

type window struct {
	Tabs []struct {
		Wins []struct {
			ID int `json:"id"`
		} `json:"windows"`
	} `json:"tabs"`
}

func (k *Kitty) BundleIdentifier() string {
	return "net.kovidgoyal.kitty"
}

func (k *Kitty) Enabled() bool {
	return true
}

func (k *Kitty) Execute(m *store.Terminal) error {
	if !k.enableRemoteControl() {
		return nil
	}

	if strings.HasPrefix(k.listening, "unix") {
		fd := k.getListeningFd()
		if len(fd) == 0 {
			m.Window = 0
			activate(m, func() {
				fd = k.getListeningFd()
				for len(fd) == 0 {
					time.Sleep(time.Millisecond * 100)
					fd = k.getListeningFd()
				}
				id := k.getLastId(fd)
				k.doScript(id, fd, m)
			})
			return nil
		}

		var id string
		if m.Window == 0 {
			id = strings.TrimSuffix(string(doCmd(k.getCmdPrefix(fd)+" launch --type os-window")), "\n")
		} else {
			id = k.getLastId(fd)
		}
		if id == "" {
			return nil
		}
		doCmd(k.getCmdPrefix(fd) + " focus-window --match 'id:" + id + "'")
		k.doScript(id, fd, m)
	} else {
		utility.ShowAlert(nil, false, "unsupport operation", "kitty listening type is unsupport")
	}
	return nil
}

func (k *Kitty) doScript(id, fd string, m *store.Terminal) {
	prefix := k.getCmdPrefix(fd) + " send-text --match 'id:" + id + "' "
	exit := toShell(prefix, m.Shell)
	for _, cmd := range m.Cmds {
		if cmd == "" {
			continue
		}
		doCmd(prefix + cmd + "\r\n")
	}
	if exit {
		doCmd(prefix + "exit\r\n")
	}
}

func (k *Kitty) getLastId(fd string) string {
	activated := doCmd(k.getCmdPrefix(fd) + " ls")
	var wins []window
	err := json.Unmarshal(activated, &wins)
	if err != nil {
		return ""
	}
	tabs := wins[len(wins)-1].Tabs
	tabWins := tabs[len(tabs)-1].Wins
	return strconv.Itoa(tabWins[len(tabWins)-1].ID)
}

func (k *Kitty) getCmdPrefix(fd string) string {
	return "kitty @ --to " + fd
}

func (k *Kitty) getListeningFd() string {
	path := strings.TrimPrefix(k.listening, "unix:")
	path = filepath.Clean(path)
	dir := filepath.Dir(path)
	var fd string
	dirs, err := os.ReadDir(dir)
	if err != nil {
		utility.ShowAlert(nil, false, "Kitty listening Socket cannot be found", err.Error())
		return ""
	}

	for _, e := range dirs {
		if e.IsDir() {
			continue
		}

		if strings.HasPrefix(e.Name(), filepath.Base(path)) {
			fd = filepath.Join(dir, e.Name())
			break
		}
	}
	if fd == "" {
		return ""
	}
	return "unix:" + fd
}

func (k *Kitty) enableRemoteControl() bool {
	path := getKittyConf()
	enableRemote, conf := k.remoteControlEnabled(path)
	if !enableRemote {
		var resp appkit.ModalResponse
		utility.ShowAlertWithCfg(nil, true,
			"Kitty turns off Remote Control",
			"Do you need to turn on Kitty’s Remote Control function? It can only be used normally after it is turned on.",
			func(alert appkit.Alert) {
				delegate := &appkit.AlertDelegate{}
				delegate.SetAlertShowHelp(func(_ appkit.Alert) bool {
					helpURL := "https://sw.kovidgoyal.net/kitty/remote-control/#control-kitty-from-scripts"
					workspace := appkit.Workspace_SharedWorkspace()
					workspace.OpenURL(foundation.URL_URLWithString(helpURL))
					return true
				})
				alert.SetShowsHelp(true)
				alert.SetDelegate(delegate)
			},
			func(response appkit.ModalResponse) {
				resp = response
			})
		if resp != appkit.AlertFirstButtonReturn {
			return false
		}

		err := os.WriteFile(path, conf, os.ModePerm)
		if err != nil {
			utility.ShowAlert(nil, false, "Failed to enable Remote Control", err.Error())
			return false
		}
	}
	return true
}

func (k *Kitty) remoteControlEnabled(path string) (bool, []byte) {
	f, err := os.OpenFile(path, os.O_RDWR, os.ModePerm)
	if err != nil {
		return false, nil
	}
	defer f.Close()
	scanner := bufio.NewScanner(f)
	var enableRemote bool
	enableListenOn := true
	buf := bytes.Buffer{}
	for scanner.Scan() {
		line := scanner.Text()
		comment := strings.HasPrefix(line, "#")
		txt := strings.TrimLeft(line, "#")
		txt = strings.TrimSpace(txt)
		if !enableRemote && strings.HasPrefix(txt, "allow_remote_control") {
			values := strings.SplitN(txt, " ", 2)
			if !comment && strings.EqualFold(strings.TrimSpace(values[1]), "yes") {
				enableRemote = true
				buf.WriteString(line)
			} else {
				buf.WriteString(values[0])
				buf.WriteByte(' ')
				buf.WriteString("yes")
			}
		} else if strings.HasPrefix(txt, "listen_on") {
			values := strings.SplitN(txt, " ", 2)
			if comment || strings.EqualFold(strings.TrimSpace(values[1]), "none") {
				enableListenOn = false
				k.listening = "unix:/tmp/kitty"
				buf.WriteString(values[0])
				buf.WriteByte(' ')
				buf.WriteString(k.listening)
			} else {
				buf.WriteString(line)
				k.listening = values[1]
			}
		} else {
			buf.WriteString(line)
		}
		buf.WriteString("\n")
	}
	return enableRemote && enableListenOn, buf.Bytes()
}

func getKittyConf() string {
	dir := os.Getenv("XDG_CONFIG_HOME")
	if len(dir) > 0 {
		return filepath.Join(dir, "kitty.conf")
	}

	home, _ := os.UserHomeDir()
	path := filepath.Join(home, ".config/kitty/kitty.conf")
	_, err := os.Stat(path)
	if err == nil {
		return path
	}

	dir = os.Getenv("XDG_CONFIG_DIRS")
	if len(dir) > 0 {
		return filepath.Join(dir, "kitty.conf")
	}
	return ""
}
