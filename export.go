package main

import (
	"github.com/ProtonMail/ui"
	"github.com/shurcooL/trayhost"
	"strconv"
)

var exportEntry *ExportEntry

func exportMenu() trayhost.MenuItem {
	return trayhost.MenuItem{
		Title: exportWindow.Title(),
		Handler: func() {
			if exportEntry != nil {
				exportEntry.Clear()
			}
			exportWindow.Show()
		},
	}
}

func exportOnReady(window *ui.Window) {
	exportWindow = window
	exportWindow.OnClosing(func(window *ui.Window) bool {
		window.Hide()
		return false
	})
	export := &ExportEntry{}
	mainBox := ui.NewVerticalBox()
	mainBox.SetPadded(true)

	form := ui.NewForm()
	form.SetPadded(true)
	export.XLSName = ui.NewEntry()
	form.Append("FileName", export.XLSName, false)
	export.Extension = ui.NewCombobox()
	export.Extension.Append(".xlsx")
	export.Extension.Append(".xls")
	export.Extension.SetSelected(0)
	form.Append("Extension", export.Extension, false)
	mainBox.Append(form, false)

	group := ui.NewGroup("SQLEntries")
	mainBox.Append(group, true)

	groups := ui.NewVerticalBox()
	groups.SetPadded(true)
	export.Groups = groups

	groups.Append(entry(1, export), false)
	group.SetChild(groups)

	btnBox := ui.NewHorizontalBox()
	btnBox.SetPadded(true)
	addBtn := ui.NewButton("Add Entry")
	addBtn.OnClicked(func(button *ui.Button) {
		groups.Append(entry(len(export.SQLEntries)+1, export), false)
		//groups.Delete(len(export.SQLEntries) - 1)
		//groups.Append(addBtn, false)
	})
	btnBox.Append(addBtn, false)

	exportBtn := ui.NewButton("Export")
	exportBtn.OnClicked(func(button *ui.Button) {

	})
	btnBox.Append(exportBtn, false)

	mainBox.Append(btnBox, false)
	exportWindow.SetChild(mainBox)
	exportEntry = export
}

func entry(entrySize int, export *ExportEntry) *ui.Group {
	entryGroup := ui.NewGroup("Entry: " + strconv.Itoa(entrySize))
	sqlEntry := &SQLEntry{}
	entryBox := ui.NewVerticalBox()
	entryBox.SetPadded(true)
	form := ui.NewForm()
	form.SetPadded(true)
	input := ui.NewEntry()
	sqlEntry.URL = input
	form.Append("URL", input, false)
	input = ui.NewEntry()
	sqlEntry.SQL = input
	form.Append("SQL", input, false)
	input = ui.NewEntry()
	sqlEntry.Args = input
	form.Append("Args", input, false)
	input = ui.NewEntry()
	sqlEntry.Title = input
	form.Append("Titles", input, false)
	input = ui.NewEntry()
	sqlEntry.SheetName = input
	form.Append("Sheet", input, false)
	var delBtn *ui.Button
	if entrySize > 1 {
		delBtn = ui.NewButton("Delete")
		delBtn.OnClicked(func(button *ui.Button) {
			if len(export.SQLEntries) > 1 {
				export.Groups.Delete(entrySize - 1)
			}
		})
	}
	entryBox.Append(form, false)
	if delBtn != nil {
		entryBox.Append(delBtn, false)
	}
	entryGroup.SetMargined(true)
	entryGroup.SetChild(entryBox)
	sqlEntry.Group = entryGroup
	export.SQLEntries = append(export.SQLEntries, sqlEntry)
	return entryGroup
}

type ExportEntry struct {
	XLSName    *ui.Entry
	SQLEntries []*SQLEntry
	Groups     *ui.Box
	Extension  *ui.Combobox
}

type SQLEntry struct {
	Group     *ui.Group
	URL       *ui.Entry
	SQL       *ui.Entry
	Args      *ui.Entry
	Title     *ui.Entry
	SheetName *ui.Entry
}

func (e *ExportEntry) Clear() {
	if e.Extension != nil {
		e.Extension.SetSelected(0)
	}
	if e.XLSName != nil {
		e.XLSName.SetText("")
	}
	if len(e.SQLEntries) > 1 {
		sqlEntry := e.SQLEntries[0]
		sqlEntry.SheetName.SetText("")
		sqlEntry.Title.SetText("")
		sqlEntry.Args.SetText("")
		sqlEntry.SQL.SetText("")
		sqlEntry.URL.SetText("")
	}
	if e.Groups != nil && len(e.SQLEntries) > 1 {
		for i := 1; i < len(e.SQLEntries); {
			e.Groups.Delete(i)
			temp := make([]*SQLEntry, 0)
			temp = append(temp, e.SQLEntries[:1]...)
			if len(e.SQLEntries) >= 2 {
				temp = append(temp, e.SQLEntries[2:]...)
			}
			e.SQLEntries = temp
		}
	}
}
