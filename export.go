package main

import (
	"fmt"
	"github.com/ProtonMail/ui"
	"github.com/shurcooL/trayhost"
	"strconv"
	"strings"
)

var (
	exportEntry *ExportEntry
	extensions  = []string{".xlsx", ".xls"}
	prompts     = []string{
		`1. URL: username:password@tcp(ip:port)/db?Charset=utf8`,
		`2. SQL: select * from user where user_id = ? and name like ?`,
		`3. Args: 666,tools (If the parameter contains[,] when, use [\.] to avoid this)`,
		`4. Titles: ID,姓名,年龄... (This is excel sheet column title)`,
		`5. Sheet: 用户统计 (This is excel sheet name)`,
		`Tips: When multiple Sheets use the same link, other link items need not be filled in`,
	}
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

	// SavePath
	defaultDownload := downloadPath()
	savePathBox := ui.NewHorizontalBox()
	savePathBox.SetPadded(true)
	savePath := ui.NewEntry()
	savePath.SetReadOnly(true)
	savePath.SetText(defaultDownload)
	selectBtn := ui.NewButton("Choose")
	selectBtn.OnClicked(func(button *ui.Button) {
		filename := ui.SaveFile(exportWindow)
		if filename == "" {
			filename = defaultDownload
		}
		if strings.HasSuffix(filename, "/Untitled") {
			filename = filename[:strings.LastIndex(filename, "/")]
		}
		savePath.SetText(filename)
	})
	savePathBox.Append(selectBtn, false)
	savePathBox.Append(savePath, true)
	form.Append("Download", savePathBox, false)

	exportEntry.Extension = ui.NewCombobox()
	exportEntry.Extension.Append(extensions[0])
	exportEntry.Extension.Append(extensions[1])
	exportEntry.Extension.SetSelected(0)
	form.Append("Extension", exportEntry.Extension, false)
	mainBox.Append(form, false)

	// Radio Buttons for Same connection URL from checkbox impl
	yes := ui.NewCheckbox("Yes")
	no := ui.NewCheckbox("No")
	yes.SetChecked(true)
	yes.OnToggled(func(checkbox *ui.Checkbox) {
		// checked yes
		if checkbox.Checked() {
			no.SetChecked(false)
			exportEntry.UseOneURL = true
		} else {
			// remove TabSheet's URL form line
			checkbox.SetChecked(true)
		}
	})
	no.OnToggled(func(checkbox *ui.Checkbox) {
		if checkbox.Checked() {
			yes.SetChecked(false)
			exportEntry.UseOneURL = false
		} else {
			// add TabSheet's URL form line
			checkbox.SetChecked(true)
		}
	})
	radioBox := ui.NewHorizontalBox()
	radioBox.SetPadded(true)
	radioBox.Append(ui.NewLabel("Use the same connection URL for multi sheet?"), false)
	radioBox.Append(yes, false)
	radioBox.Append(no, false)
	mainBox.Append(radioBox, false)

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

	// Prompt Form format
	separator := ui.NewHorizontalSeparator()
	mainBox.Append(separator, false)
	prompt(mainBox)

	exportWindow.SetChild(mainBox)
}

func prompt(mainBox *ui.Box) {
	for index, p := range prompts {
		if index == 0 {
			box := ui.NewHorizontalBox()
			box.SetPadded(true)
			label := ui.NewLabel(p)
			button := ui.NewButton("BuildURL")
			button.OnClicked(func(button *ui.Button) {

			})
			box.Append(label, false)
			box.Append(button, false)
			mainBox.Append(box, false)
		} else {
			mainBox.Append(ui.NewLabel(p), false)
		}
	}
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
	var input *ui.Entry
	if !exportEntry.UseOneURL || len(exportEntry.SQLEntries) == 0 {
		input = ui.NewEntry()
		entry.URL = input
		form.Append("URL", input, false)
	}
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
	SavePath   *ui.Entry
	SQLEntries []*SQLEntry
	Extension  *ui.Combobox
	TabEntries map[int]*ui.Grid
	Tab        *ui.Tab
	DeletedTab int
	UseOneURL  bool
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
	if e.SavePath != nil {
		e.SavePath.SetText(downloadPath())
	}
	if e.Extension != nil {
		e.Extension.SetSelected(0)
	}
	if e.XLSName != nil {
		e.XLSName.SetText("")
	}
	// TODO exportEntry.Tab clear
}
