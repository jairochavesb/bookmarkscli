package main

import (
	"bookmarksV2/modules/db"
	"bookmarksV2/modules/tui"
	"os"
)

const dbName = "bookmarks.db"

func main() {
	_, err := os.Stat(dbName)
	if err != nil {
		db.InitDB()
	}

	tui.MainLoop()
}
