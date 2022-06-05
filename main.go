package main

import (
	"bookmarksV2/modules/config"
	"bookmarksV2/modules/db"
	"bookmarksV2/modules/tui"
	"os"
)

const dbFile = ".bookmarksConfig.txt"

func main() {
	_, err := os.Stat(dbFile)
	if err != nil {
		config.SetConfig()
	}
	config.LoadConfig()

	_, err = os.Stat(config.Configuration.DBFile)
	if err != nil {
		db.InitDB()
	}

	tui.MainLoop()
}
