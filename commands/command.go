package commands

import (
	"bufio"
	"bytes"
	"github.com/charliego3/tools/store"
	"github.com/charliego3/tools/utility"
	"github.com/progrium/macdriver/dispatch"
	"os"
	"os/exec"
	"slices"
	"strings"
	"time"

	"github.com/progrium/macdriver/helper/action"
	"github.com/progrium/macdriver/helper/layout"
	"github.com/progrium/macdriver/helper/widgets"
	"github.com/progrium/macdriver/macos/appkit"
	"github.com/progrium/macdriver/macos/foundation"
	"github.com/progrium/macdriver/objc"
)

var supportedApps = []Executor{
	&Terminal{},
	&Kitty{},
	&Alacritty{},
	&Iterm2{},
}

type Command struct {
	appkit.MenuItem
	menu appkit.Menu
}

func Item() *Command {
	t := &Command{}
	t.MenuItem = appkit.NewMenuItem()
	t.MenuItem.SetTitle("Shell Commands")
	t.MenuItem.SetImage(utility.SymbolImage("terminal.fill", utility.ImageScale(appkit.ImageSymbolScaleSmall)))

	opts := store.Fetch()
	t.menu = appkit.NewMenu()
	var ids []int64
	for id := range opts.Terminals {
		ids = append(ids, id)
	}
	slices.Sort(ids)
	for _, id := range ids {
		m := opts.Terminals[id]
		t.menu.AddItem(t.getShellItem(id, &m))
	}

	font := appkit.Font_BoldSystemFontOfSize(appkit.Font_SystemFontSize())
	item := appkit.NewMenuItem()
	const title = "Add new Command"
	item.SetAttributedTitle(foundation.NewAttributedStringWithStringAttributes(
		title, map[foundation.AttributedStringKey]objc.IObject{"NSFont": font}),
	)
	action.Set(item, t.showWindow(time.Now().UnixNano(), &store.Terminal{}, title, func(id int64, m *store.Terminal) {
		t.menu.InsertItemAtIndex(t.getShellItem(id, m), t.menu.NumberOfItems()-2)
	}))
	t.menu.AddItem(appkit.MenuItem_SeparatorItem())
	t.menu.AddItem(item)
	t.MenuItem.SetSubmenu(t.menu)
	return t
}

func openApp(m *store.Terminal) func(objc.Object) {
	return func(objc.Object) {
		var exe Executor
		for _, app := range supportedApps {
			if !app.Enabled() {
				continue
			}
			if app.BundleIdentifier() == m.App {
				exe = app
				break
			}
		}

		if exe == nil {
			workspace := appkit.Workspace_SharedWorkspace()
			url := workspace.URLForApplicationWithBundleIdentifier(m.App)
			utility.ShowAlert(nil, false, "Run script failed", "This app does not support: "+url.LastPathComponent())
			return
		}

		dispatch.MainQueue().DispatchAsync(func() {
			e := exe.Execute(m)
			if e != nil {
				utility.ShowAlert(nil, false, "Run script failed", e.Error())
			}
		})
	}
}

func (t *Command) getShellItem(id int64, m *store.Terminal) appkit.MenuItem {
	item := appkit.NewMenuItem()
	item.SetTitle(m.Name)
	action.Set(item, openApp(m))
	edit := appkit.NewMenuItem()
	edit.SetTitle("Edit")
	action.Set(edit, t.showWindow(id, m, "Edit Command", func(id int64, em *store.Terminal) {
		item.SetTitle(m.Name)
	}))
	del := appkit.NewMenuItem()
	del.SetTitle("Delete")
	action.Set(del, func(sender objc.Object) {
		utility.ShowAlert(nil, true, "Are you sure you want to delete this Command?",
			"This operation is irreversible and will be deleted soon "+m.Name,
			func(response appkit.ModalResponse) {
				if response == appkit.AlertFirstButtonReturn {
					opts := store.Fetch()
					delete(opts.Terminals, id)
					_ = opts.Save()
					t.menu.RemoveItem(item)
				}
			})
	})
	menu := appkit.NewMenu()
	menu.AddItem(edit)
	menu.AddItem(del)
	item.SetSubmenu(menu)
	return item
}

func (t *Command) getAddCmdItem() appkit.MenuItem {
	item := appkit.NewMenuItem()
	item.SetTitle("New Command")
	action.Set(item, func(_ objc.Object) {
		defaults := foundation.UserDefaults_StandardUserDefaults()
		aaa := defaults.StringForKey("key")

		dialog := appkit.NewAlert()
		dialog.SetInformativeText(aaa)
		dialog.SetMessageText("User Defaults")
		dialog.SetAlertStyle(appkit.AlertStyleInformational)
		dialog.RunModal()
	})
	return item
}

func (t *Command) showWindow(id int64, m *store.Terminal, title string, callback func(int64, *store.Terminal)) func(objc.Object) {
	return func(sender objc.Object) {
		nameField := appkit.NewTextField()
		nameField.SetBezelStyle(appkit.TextFieldRoundedBezel)
		nameField.SetStringValue(m.Name)
		form := widgets.NewFormView()
		shells := t.getShells(m.Shell)
		apps := t.getApps(m.App)
		form.AddRow("Name:", nameField)
		form.AddRow("Shell:", shells)
		form.AddRow("Terminal:", apps)

		matrix := t.getWindowMatrix()
		matrix.SelectCellAtRowColumn(m.Window, 0)
		form.AddRow("Window:", matrix)

		textView := appkit.NewScrollableTextView()
		textView.SetFocusRingType(appkit.FocusRingTypeExterior)
		textView.ContentTextView().SetFont(appkit.Font_FontWithNameSize("Monaco", 13))
		textView.SetWantsLayer(true)
		textView.Layer().SetCornerRadius(3)
		textView.ContentTextView().SetString(strings.Join(m.Cmds, "\r\n"))
		layout.SetHeight(textView, 100)
		form.AddRowWithViews([]appkit.IView{
			appkit.NewLabel("Script:"),
			textView,
		})

		form.SetLabelWidth(55)
		form.SetLabelAlignment(widgets.LabelAlignmentTrailing)
		form.SetTranslatesAutoresizingMaskIntoConstraints(false)
		view := appkit.NewView()
		view.AddSubview(form)

		btn := appkit.NewButtonWithTitle("Save")
		btn.SetTranslatesAutoresizingMaskIntoConstraints(false)
		btn.SetBezelColor(appkit.Color_SystemBlueColor())
		w := utility.NewWindow(title, view, func(w appkit.Window) {
			w.SetStyleMask(w.StyleMask() | appkit.MiniaturizableWindowMask)
		})
		action.Set(btn, func(sender objc.Object) {
			m.Name = nameField.StringValue()
			if m.Name == "" {
				utility.ShowAlert(w, false, "Name is empty!", "Must specify a name to save")
				return
			}

			script := textView.ContentTextView().String()
			if script == "" {
				utility.ShowAlert(w, false, "Script is empty!", "There are no scripts to run here")
				return
			}

			app := foundation.StringFrom(apps.SelectedItem().RepresentedObject().Ptr())
			scanner := bufio.NewScanner(strings.NewReader(script))
			var cmds []string
			for scanner.Scan() {
				cmds = append(cmds, scanner.Text())
			}
			m.App = app.String()
			m.Shell = shells.SelectedItem().Title()
			m.Window = matrix.SelectedRow()
			m.Cmds = cmds
			if id == 0 {
				id = time.Now().UnixNano()
			}
			opts := store.Fetch()
			opts.Terminals[id] = *m
			err := opts.Save()
			if err != nil {
				utility.ShowAlert(w, false, "Failed to save", err.Error())
				return
			}
			w.Close()
			if callback != nil {
				callback(id, m)
			}
		})

		view.AddSubview(btn)
		layout.SetWidth(form, 400)
		utility.LayoutActives(
			form.TopAnchor().ConstraintEqualToAnchorConstant(view.TopAnchor(), 30),
			form.LeadingAnchor().ConstraintGreaterThanOrEqualToAnchorConstant(view.LeadingAnchor(), 10),
			form.TrailingAnchor().ConstraintLessThanOrEqualToAnchorConstant(view.TrailingAnchor(), -10),
			form.BottomAnchor().ConstraintLessThanOrEqualToAnchorConstant(view.BottomAnchor(), -10),
			btn.TopAnchor().ConstraintEqualToAnchorConstant(form.BottomAnchor(), 10),
			btn.BottomAnchor().ConstraintEqualToAnchorConstant(view.BottomAnchor(), -10),
			btn.TrailingAnchor().ConstraintEqualToAnchorConstant(view.TrailingAnchor(), -10),
		)
		w.MakeKeyAndOrderFront(sender)
	}
}

func (t *Command) getShells(selected string) appkit.PopUpButton {
	popup := appkit.NewPopUpButton()
	if selected == "" {
		selected = os.Getenv("SHELL")
	}

	r, _ := exec.Command("cat", "/etc/shells").CombinedOutput()
	reader := bufio.NewScanner(bytes.NewReader(r))
	var shells []string
	for reader.Scan() {
		text := strings.TrimSpace(reader.Text())
		if text == "" || strings.HasPrefix(text, "#") {
			continue
		}
		shells = append(shells, text)
	}

	menu := appkit.NewMenu()
	popup.SetMenu(menu)
	for _, shell := range shells {
		item := appkit.NewMenuItem()
		item.SetTitle(shell)
		menu.AddItem(item)
		if selected == shell {
			popup.SelectItem(item)
		}
	}
	return popup
}

func (t *Command) getWindowMatrix() appkit.Matrix {
	wins := []string{"New instance window", "Reuse the last window"}
	matrix := appkit.NewMatrix()
	matrix.SetMode(appkit.RadioModeMatrix)
	matrix.SetCellClass(appkit.ButtonCellClass)
	matrix.RenewRowsColumns(0, 1)
	matrix.SetControlSize(appkit.ControlSizeLarge)
	matrix.SetCellSize(utility.SizeOf(400, 20))
	for _, win := range wins {
		cell := appkit.NewButtonCell()
		cell.SetButtonType(appkit.ButtonTypeRadio)
		cell.SetTitle(win)
		matrix.AddRowWithCells([]appkit.ICell{cell})
	}
	return matrix
}

func (t *Command) getApps(selectedIdentifier string) appkit.PopUpButton {
	popup := appkit.NewPopUpButton()
	menu := appkit.NewMenu()
	popup.SetMenu(menu)
	workspace := appkit.Workspace_SharedWorkspace()
	for i, app := range supportedApps {
		url := workspace.URLForApplicationWithBundleIdentifier(app.BundleIdentifier())
		if url.IsNil() {
			continue
		}

		icon := workspace.IconForFile(url.Path())
		icon.SetSize(utility.SizeOf(20, 20))
		item := appkit.NewMenuItem()
		item.SetImage(icon)
		item.SetTitle(url.LastPathComponent())
		item.SetTag(i)
		item.SetHidden(!app.Enabled())
		item.SetRepresentedObject(foundation.String_StringWithString(app.BundleIdentifier()))
		menu.AddItem(item)

		if selectedIdentifier == app.BundleIdentifier() {
			popup.SelectItem(item)
		}
	}
	return popup
}
