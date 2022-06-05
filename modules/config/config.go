package config

import (
	"bufio"
	"fmt"
	"log"
	"os"
	"strings"
	"time"
)

const configFile = ".bookmarksConfig.txt"

type Config struct {
	WebBrowser string
	DBFile     string
}

var Configuration Config

func LoadConfig() {
	file, err := os.Open(configFile)
	if err != nil {
		fmt.Println("Error reading the config file.")
		time.Sleep(3 * time.Second)
	}

	scanner := bufio.NewScanner(file)

	for scanner.Scan() {

		v := strings.Split(scanner.Text(), "=")

		if v[0] == "WEB_BROWSER" {
			Configuration.WebBrowser = v[1]
		}

		if v[0] == "DB_FILE" {
			Configuration.DBFile = v[1]
		}
	}
}

func SetConfig() {
	var webBrowser, dbFile string

	fmt.Printf("Favorite web browser to open bookmarks: ")
	fmt.Scanf("%s", &webBrowser)

	fmt.Printf("Full path and name for the database file: ")
	fmt.Scanf("%s", &dbFile)

	configString := "WEB_BROWSER=" + webBrowser + "\n"
	configString += "DB_FILE=" + dbFile

	file, err := os.Create(".bookmarksConfig.txt")
	if err != nil {
		log.Fatal("Unable to create config file.")
	}
	defer file.Close()

	_, err = file.Write([]byte(configString))
	if err != nil {
		log.Fatal("Unable to write to the config file.")
	}

}
