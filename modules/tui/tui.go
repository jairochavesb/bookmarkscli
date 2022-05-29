package tui

import (
	"bookmarksV2/modules/db"
	"fmt"
	"os"
	"strings"

	ui "github.com/jairochavesb/clui"
	termbox "github.com/nsf/termbox-go"

	"golang.org/x/term"
)

var txtName, txtPath, txtTags, txtSearch *ui.EditField
var listboxResults *ui.ListBox

var lblStatus *ui.Label

var searchCol string

var id string = ""

func MainLoop() {
	ui.InitLibrary()
	defer ui.DeinitLibrary()

	_, err := os.Stat("themes/simpleDark.theme")
	if err == nil {
		ui.SetThemePath("themes/")
		ui.SetCurrentTheme("simpleDark")
	}

	initUI()

	ui.MainLoop()
}

func initUI() {
	width, height, _ := term.GetSize(0)
	frameWidth := width - 10
	widgetWidth := frameWidth - 5

	mainWindow := ui.AddWindow(0, 0, width, height, "")
	mainWindow.SetBorder(0)
	mainWindow.SetTitleButtons(0)
	mainWindow.SetPack(ui.Vertical)

	// FORM WITH BOOKMARK WIDGETS (PATH, NAME, TAGS)
	frmShowInsertData := ui.CreateFrame(mainWindow, frameWidth, 5, ui.BorderThin, ui.Fixed)
	frmShowInsertData.SetPack(ui.Vertical)
	frmShowInsertData.SetTitle("[BOOKMARK]")
	frmShowInsertData.SetGaps(0, 1)

	frmWidgetsName := ui.CreateFrame(frmShowInsertData, 1, 1, ui.BorderNone, ui.Fixed)
	frmWidgetsName.SetPack(ui.Horizontal)
	_ = ui.CreateLabel(frmWidgetsName, 7, 1, "  Name ", ui.Fixed)
	txtName = ui.CreateEditField(frmWidgetsName, widgetWidth, "", ui.Fixed)

	frmWidgetsPath := ui.CreateFrame(frmShowInsertData, 1, 1, ui.BorderNone, ui.Fixed)
	frmWidgetsPath.SetPack(ui.Horizontal)
	_ = ui.CreateLabel(frmWidgetsPath, 7, 1, "  Path ", ui.Fixed)
	txtPath = ui.CreateEditField(frmWidgetsPath, widgetWidth, "", ui.Fixed)

	frmWidgetsTags := ui.CreateFrame(frmShowInsertData, 1, 1, ui.BorderNone, ui.Fixed)
	frmWidgetsTags.SetPack(ui.Horizontal)
	_ = ui.CreateLabel(frmWidgetsTags, 7, 1, "  Tags ", ui.Fixed)
	txtTags = ui.CreateEditField(frmWidgetsTags, widgetWidth, "", ui.AutoSize)

	_ = ui.CreateLabel(frmShowInsertData, widgetWidth, 1, "       <t:cyan>SAVE=<f:>Enter <t:cyan>CLEAR FIELD=<t:>Ctrl+R <t:cyan>CLEAR ALL=<t:>Ctrl+A", ui.Fixed)

	// FRAME WITH WIDGETS TO SEARCH DATA
	frmMainSearchData := ui.CreateFrame(mainWindow, frameWidth, 1, ui.BorderThin, ui.Fixed)
	frmMainSearchData.SetPack(ui.Vertical)
	frmMainSearchData.SetTitle("[SEARCH]")
	frmMainSearchData.SetGaps(0, 1)

	frmRadios := ui.CreateFrame(frmMainSearchData, frameWidth, 1, ui.BorderNone, ui.Fixed)
	frmRadios.SetPack(ui.Horizontal)

	_ = ui.CreateLabel(frmRadios, 8, 1, "", ui.Fixed)
	radioName := ui.CreateRadio(frmRadios, 10, "Name", ui.Fixed)
	radioName.SetSelected(true)
	searchCol = "name"

	radioPath := ui.CreateRadio(frmRadios, 10, "Path", ui.Fixed)
	radioTag := ui.CreateRadio(frmRadios, 10, "Tag", ui.Fixed)
	radioGroup := ui.CreateRadioGroup()
	radioGroup.AddItem(radioName)
	radioGroup.AddItem(radioPath)
	radioGroup.AddItem(radioTag)

	frmSearchData := ui.CreateFrame(frmMainSearchData, frameWidth, 1, ui.BorderNone, ui.Fixed)
	frmSearchData.SetPack(ui.Horizontal)
	_ = ui.CreateLabel(frmSearchData, 8, 1, "KEYWORD ", ui.Fixed)
	txtSearch = ui.CreateEditField(frmSearchData, widgetWidth, "", ui.AutoSize)

	_ = ui.CreateLabel(frmMainSearchData, 1, 1, "        <t:cyan>SEARCH=<t:>Enter", ui.Fixed)

	// LISTBOX WITH THE SEARCH RESULTS.
	frmResults := ui.CreateFrame(mainWindow, widgetWidth, height/2, ui.BorderThin, ui.Fixed)
	frmResults.SetPack(ui.Vertical)
	frmResults.SetTitle("[RESULTS]")
	frmResults.SetGaps(0, 1)
	listboxResults = ui.CreateListBox(frmResults, -1, (height/2)-6, ui.AutoSize)
	_ = ui.CreateLabel(frmResults, widgetWidth-1, 1, "<t:cyan>OPEN=<f:>Ctrl+O <t:cyan>EDIT=<t:>Ctrl+E", ui.Fixed)

	// Status label
	frmStatus := ui.CreateFrame(mainWindow, widgetWidth, 1, ui.BorderThin, ui.AutoSize)
	frmStatus.SetTitle("[STATUS]")
	lblStatus = ui.CreateLabel(frmStatus, 100, 1, "Hello user!", ui.AutoSize)

	ui.ActivateControl(mainWindow, txtPath)

	txtName.OnKeyPress(func(key termbox.Key, r rune) bool {
		if key == termbox.KeyEnter {
			save()
		} else if key == termbox.KeyCtrlA {
			clearAll()
		}
		return false
	})

	txtPath.OnKeyPress(func(key termbox.Key, r rune) bool {
		if key == termbox.KeyEnter {
			save()
		} else if key == termbox.KeyCtrlA {
			clearAll()
		}
		return false
	})

	txtTags.OnKeyPress(func(key termbox.Key, r rune) bool {
		if key == termbox.KeyEnter {
			save()
		} else if key == termbox.KeyCtrlA {
			clearAll()
		}
		return false
	})

	txtSearch.OnKeyPress(func(key termbox.Key, r rune) bool {
		if key == termbox.KeyEnter {
			search()
		}
		return false
	})

	listboxResults.OnKeyPress(func(key termbox.Key) bool {
		if key == termbox.KeyEnter {
			open()
		} else if key == termbox.KeyCtrlE {
			edit()
		}
		return false
	})

	radioName.OnChange(func(b bool) {
		if b {
			searchCol = "name"
		}
	})
	radioPath.OnChange(func(b bool) {
		if b {
			searchCol = "path"
		}
	})
	radioTag.OnChange(func(b bool) {
		if b {
			searchCol = "tags"
		}
	})
}

func edit() {
	l := listboxResults.SelectedItemText()
	id = strings.Split(l, "◄►")[0]
	txtName.SetTitle(strings.Split(l, "◄►")[1])
	txtPath.SetTitle(strings.Split(l, "◄►")[2])
	txtTags.SetTitle(strings.Split(l, "◄►")[3])
}

func open() {
	fmt.Println("open")
}

func save() {
	name := txtName.Title()
	path := txtPath.Title()
	tags := txtTags.Title()

	if tags == "" || name == "" || path == "" {
		updateStatusLabel("Please fill all the bookmark input fields! >:(", 1)
	} else {
		if id == "" {
			db.InsertData(name, path, tags)
			updateStatusLabel("Bookmark saved successfully :)", 0)
		} else {
			db.UpdateData(id, name, path, tags)
			updateStatusLabel("Bookmark updated successfully :)", 0)
		}

		id = ""
	}

	clearAll()
}

func search() {
	listboxResults.Clear()

	results := db.SearchData(txtSearch.Title(), searchCol)

	for _, v := range results {
		listboxResults.AddItem(v)
	}

	updateStatusLabel("Search finish successfully :)", 0)
}

func clearAll() {
	txtName.SetTitle("")
	txtTags.SetTitle("")
	txtPath.SetTitle("")
	id = ""
}

func updateStatusLabel(m string, t int) {
	var s string
	if t == 0 {
		s = "<b:cyan>" + m
	} else {
		s = "<b:red>" + m
	}

	lblStatus.SetTitle(s)
}
