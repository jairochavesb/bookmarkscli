package main

import (
	"bookmarkscli/modules/config"
	"bookmarkscli/modules/db"
	"bookmarkscli/modules/tui"
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
