package tui

import (
	"bookmarksV2/modules/config"
	"bookmarksV2/modules/db"
	"log"
	"os"
	"os/exec"
	"strings"

	ui "github.com/jairochavesb/clui"
	termbox "github.com/nsf/termbox-go"

	"golang.org/x/term"
)

var txtName, txtUrl, txtTags, txtSearch *ui.EditField
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
	fd := int(os.Stdout.Fd())
	width, height, _ := term.GetSize(fd)
	frameWidth := width - 10
	widgetWidth := frameWidth - 5

	mainWindow := ui.AddWindow(0, 0, width, height, "")
	mainWindow.SetBorder(0)
	mainWindow.SetTitleButtons(0)
	mainWindow.SetPack(ui.Vertical)

	// FORM WITH BOOKMARK WIDGETS (Url, NAME, TAGS)
	frmShowInsertData := ui.CreateFrame(mainWindow, frameWidth, 5, ui.BorderThin, ui.Fixed)
	frmShowInsertData.SetPack(ui.Vertical)
	frmShowInsertData.SetTitle("[BOOKMARK]")
	frmShowInsertData.SetGaps(0, 1)

	frmWidgetsUrl := ui.CreateFrame(frmShowInsertData, 1, 1, ui.BorderNone, ui.Fixed)
	frmWidgetsUrl.SetPack(ui.Horizontal)
	_ = ui.CreateLabel(frmWidgetsUrl, 7, 1, "  Url ", ui.Fixed)
	txtUrl = ui.CreateEditField(frmWidgetsUrl, widgetWidth, "", ui.Fixed)

	frmWidgetsName := ui.CreateFrame(frmShowInsertData, 1, 1, ui.BorderNone, ui.Fixed)
	frmWidgetsName.SetPack(ui.Horizontal)
	_ = ui.CreateLabel(frmWidgetsName, 7, 1, "  Name ", ui.Fixed)
	txtName = ui.CreateEditField(frmWidgetsName, widgetWidth, "", ui.Fixed)

	frmWidgetsTags := ui.CreateFrame(frmShowInsertData, 1, 1, ui.BorderNone, ui.Fixed)
	frmWidgetsTags.SetPack(ui.Horizontal)
	_ = ui.CreateLabel(frmWidgetsTags, 7, 1, "  Tags ", ui.Fixed)
	txtTags = ui.CreateEditField(frmWidgetsTags, widgetWidth, "", ui.AutoSize)

	l := "       <t:cyan>SAVE=<f:>Enter <t:cyan>CLEAR FIELD=<t:>Ctrl+R <t:cyan>CLEAR ALL=<t:>Ctrl+A"
	l += " <t:cyan>COPY TEXT=<f:>Ctrl+C <t:cyan>PASTE TEXT=<f:>Ctrl+P <t:cyan>GET WORDS FROM URL=<f:>Ctrl+U"
	_ = ui.CreateLabel(frmShowInsertData, widgetWidth, 1, l, ui.Fixed)

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

	radioUrl := ui.CreateRadio(frmRadios, 10, "Url", ui.Fixed)
	radioTag := ui.CreateRadio(frmRadios, 10, "Tag", ui.Fixed)
	radioGroup := ui.CreateRadioGroup()
	radioGroup.AddItem(radioName)
	radioGroup.AddItem(radioUrl)
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
	_ = ui.CreateLabel(frmResults, widgetWidth-1, 1, "<t:cyan>OPEN=<f:>Enter <t:cyan>EDIT=<t:>Ctrl+E <t:cyan>DELETE=<t:>Ctrl+D", ui.Fixed)

	// Status label
	frmStatus := ui.CreateFrame(mainWindow, widgetWidth, 1, ui.BorderThin, ui.AutoSize)
	frmStatus.SetTitle("[STATUS]")
	lblStatus = ui.CreateLabel(frmStatus, 100, 1, "Hello user!", ui.AutoSize)

	ui.ActivateControl(mainWindow, txtUrl)

	txtName.OnKeyPress(func(key termbox.Key, r rune) bool {
		if key == termbox.KeyEnter {
			save()
		} else if key == termbox.KeyCtrlA {
			clearAll()
		} else if key == termbox.KeyCtrlU {
			txtName.SetTitle(getWords(txtUrl.Title()))
		}
		return false
	})

	txtUrl.OnKeyPress(func(key termbox.Key, r rune) bool {
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
		} else if key == termbox.KeyCtrlU {
			txtTags.SetTitle(getWords(txtUrl.Title()))
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
		} else if key == termbox.KeyCtrlD {
			remove()
		}

		return false
	})

	radioName.OnChange(func(b bool) {
		if b {
			searchCol = "name"
		}
	})
	radioUrl.OnChange(func(b bool) {
		if b {
			searchCol = "Url"
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
	txtUrl.SetTitle(strings.Split(l, "◄►")[2])
	txtTags.SetTitle(strings.Split(l, "◄►")[3])
}

func open() {
	cmd := exec.Command(config.Configuration.WebBrowser, strings.Split(listboxResults.SelectedItemText(), "◄►")[2])
	err := cmd.Run()
	if err != nil {
		log.Fatal(err)
	}
}

func save() {
	n := txtName.Title()
	p := txtUrl.Title()
	t := txtTags.Title()

	name := strings.Replace(n, "\n", "", -1)
	Url := strings.Replace(p, "\n", "", -1)
	tags := strings.Replace(t, "\n", "", -1)

	if tags == "" || name == "" || Url == "" {
		updateStatusLabel("Please fill all the bookmark input fields! >:(", 1)
	} else {
		if id == "" {
			dup := db.CheckIfDuplicated(txtUrl.Title())

			if dup {
				updateStatusLabel("Item already exist", 0)
				clearAll()
			} else {
				db.InsertData(name, Url, tags)
				updateStatusLabel("Bookmark saved successfully :)", 0)
			}
		} else {
			db.UpdateData(id, name, Url, tags)
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
	txtUrl.SetTitle("")
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

func getWords(u string) string {
	w := strings.ReplaceAll(u, "http://", "")
	w = strings.ReplaceAll(w, "https://", "")
	w = strings.ReplaceAll(w, "/", " ")
	w = strings.ReplaceAll(w, "-", " ")
	w = strings.ReplaceAll(w, "_", " ")

	return w
}

func remove() {
	l := listboxResults.SelectedItemText()
	id = strings.Split(l, "◄►")[0]

	db.RemoveData(id)

	updateStatusLabel("Item succesfully deleted.", 0)
}
