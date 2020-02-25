package main

import (
	"fmt"
	"github.com/ProtonMail/ui"
	"github.com/shurcooL/trayhost"
	"strconv"
)

var (
	exportEntry *ExportEntry
	extensions  = []string{".xlsx", ".xls"}
)

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
	exportEntry = &ExportEntry{
		TabEntries: make(map[int]*ui.Grid),
	}
	mainBox := ui.NewVerticalBox()
	mainBox.SetPadded(true)

	form := ui.NewForm()
	form.SetPadded(true)
	exportEntry.XLSName = ui.NewEntry()
	form.Append("FileName", exportEntry.XLSName, false)
	exportEntry.Extension = ui.NewCombobox()
	exportEntry.Extension.Append(extensions[0])
	exportEntry.Extension.Append(extensions[1])
	exportEntry.Extension.SetSelected(0)
	form.Append("Extension", exportEntry.Extension, false)
	mainBox.Append(form, false)

	exportEntry.Tab = ui.NewTab()
	addNewTab()
	exportEntry.Tab.SetMargined(0, true)
	mainBox.Append(exportEntry.Tab, false)

	exportBtnLine := ui.NewGrid()
	exportBtnLine.SetPadded(true)
	exportBtn := ui.NewButton("Export")
	exportBtn.OnClicked(onExportBtnClicked)
	exportBtnLine.Append(exportBtn, 0, 0, 1, 1, false, ui.AlignEnd, false, ui.AlignFill)
	mainBox.Append(exportBtnLine, false)
	exportWindow.SetChild(mainBox)
}

func onExportBtnClicked(button *ui.Button) {
	button.Disable()
	defer func() {
		if err := recover(); err != nil {
			ui.MsgBoxError(exportWindow,
				"Error generating Excel document.",
				"Error details: "+fmt.Sprintf("error: %v\n", err))
		}
		button.Enable()
	}()
	xlsName := exportEntry.XLSName.Text()
	extension := extensions[exportEntry.Extension.Selected()]
	for _, entry := range exportEntry.SQLEntries {
		fmt.Printf("XLSName: %s, Extension: %s, URL: %s, SQL: %s, Args: %+v, Titles: %+v, SheetName: %s\n",
			xlsName, extension, entry.URL.Text(), entry.SQL.Text(), entry.Args.Text(), entry.Titles.Text(), entry.SheetName.Text())
	}
}

// TODO fix add and delete tab bug
func onAddBtnClicked(button *ui.Button) {
	// Add new TabSheet to Tab
	index := len(exportEntry.TabEntries)
	addNewTab()
	exportEntry.Tab.SetMargined(index, true)
	// AddEntry Button replace to DeleteButton
	tl := index + exportEntry.DeletedTab
	println("Change delete:", tl)
	btnGrid := exportEntry.TabEntries[tl]
	btnGrid.Delete(0)
	delBtn := ui.NewButton("Delete")
	delBtn.OnClicked(func(button *ui.Button) {
		delete(exportEntry.TabEntries, index)
		println("Grid Index:", index, "NumPages:", exportEntry.Tab.NumPages())
		deli := index - exportEntry.DeletedTab - 1
		if deli < 0 {
			deli = 0
		}
		var temp []*SQLEntry
		temp = append(temp, exportEntry.SQLEntries[:index-1]...)
		temp = append(temp, exportEntry.SQLEntries[index:]...)
		exportEntry.SQLEntries = temp
		exportEntry.Tab.Delete(deli)
		exportEntry.DeletedTab += 1
	})
	btnGrid.Append(delBtn, 0, 0, 1, 1, false, ui.AlignEnd, false, ui.AlignFill)
}

func addNewTab() {
	exportEntry.Tab.Append("Sheet-"+strconv.Itoa(len(exportEntry.SQLEntries)+1), newTabEntry())
}

func newTabEntry() *ui.Box {
	entryBox := ui.NewVerticalBox()
	entryBox.SetPadded(true)
	entry := &SQLEntry{}
	form := ui.NewForm()
	form.SetPadded(true)
	input := ui.NewEntry()
	entry.URL = input
	form.Append("URL", input, false)
	input = ui.NewEntry()
	entry.SQL = input
	form.Append("SQL", input, false)
	input = ui.NewEntry()
	entry.Args = input
	form.Append("Args", input, false)
	input = ui.NewEntry()
	entry.Titles = input
	form.Append("Titles", input, false)
	input = ui.NewEntry()
	entry.SheetName = input
	form.Append("Sheet", input, false)
	entryBox.Append(form, false)
	addBtnLine := ui.NewGrid()
	addBtnLine.SetPadded(true)
	addBtn := ui.NewButton("Add Entry")
	addBtn.OnClicked(onAddBtnClicked)
	addBtnLine.Append(addBtn, 0, 0, 1, 1, false, ui.AlignEnd, false, ui.AlignFill)
	entryBox.Append(addBtnLine, false)
	exportEntry.SQLEntries = append(exportEntry.SQLEntries, entry)
	exportEntry.TabEntries[len(exportEntry.TabEntries)+1] = addBtnLine
	return entryBox
}

type ExportEntry struct {
	XLSName    *ui.Entry
	SQLEntries []*SQLEntry
	Extension  *ui.Combobox
	TabEntries map[int]*ui.Grid
	Tab        *ui.Tab
	DeletedTab int
}

type SQLEntry struct {
	URL       *ui.Entry
	SQL       *ui.Entry
	Args      *ui.Entry
	Titles    *ui.Entry
	SheetName *ui.Entry
}

func (e *ExportEntry) Clear() {
	e.SQLEntries = e.SQLEntries[:1]
	e.SQLEntries[0].URL.SetText("")
	e.SQLEntries[0].SQL.SetText("")
	e.SQLEntries[0].Args.SetText("")
	e.SQLEntries[0].Titles.SetText("")
	e.SQLEntries[0].SheetName.SetText("")
	e.TabEntries = make(map[int]*ui.Grid)
	e.DeletedTab = 0
	if e.Extension != nil {
		e.Extension.SetSelected(0)
	}
	if e.XLSName != nil {
		e.XLSName.SetText("")
	}
	// TODO exportEntry.Tab clear
}
