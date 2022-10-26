package db

import (
	"bookmarkscli/modules/config"
	"database/sql"
	"log"
	"os"
	"strconv"
	"strings"

	_ "github.com/mattn/go-sqlite3"
)

var id int
var name, url, tags string

func InitDB() {
	file, err := os.Create(config.Configuration.DBFile)
	if err != nil {
		log.Fatal(err)
	}
	file.Close()

	sqlDB, err := sql.Open("sqlite3", config.Configuration.DBFile)
	if err != nil {
		log.Fatal(err)
	}

	createItemsTable := `CREATE TABLE items (
      "id" integer NOT NULL PRIMARY KEY AUTOINCREMENT,      
      "name" TEXT,
      "url" TEXT,
      "tags" TEXT    
     );`

	statement, err := sqlDB.Prepare(createItemsTable)
	if err != nil {
		log.Fatal(err)
	}

	statement.Exec()
	defer sqlDB.Close()
}

func InsertData(n, u, t string) {
	q := "INSERT INTO items(name, url, tags) VALUES (\"" + n + "\",\"" + u + "\",\"" + t + "\");"
	_ = execStatement(q)
}

func UpdateData(i, n, u, t string) {
	q := "UPDATE items set tags=\"" + t + "\",name=\"" + n + "\",url=\"" + u + "\" where id=\"" + i + "\";"
	_ = execStatement(q)
}

func RemoveData(i string) {
	q := "delete from items where id=" + i + ";"
	_ = execStatement(q)
}

func SearchData(k, c string) []string {

	q := "SELECT * FROM items WHERE " + c
	tmp := ""

	keywords := strings.SplitAfter(k, " ")

	if len(keywords) >= 1 {
		for i, v := range keywords {
			v = strings.Replace(v, " ", "", -1)
			tmp += " LIKE '%" + v + "%'"
			if i != len(keywords)-1 {
				tmp += " AND " + c
			} else {
				tmp += ";"
			}
		}
		q += tmp

	} else {
		q += " LIKE '%" + k + "%';"
	}

	sqlDB, err := sql.Open("sqlite3", config.Configuration.DBFile)
	checkErr(err)
	defer sqlDB.Close()

	row, err := sqlDB.Query(q)
	checkErr(err)

	defer row.Close()
	var results []string
	var r string

	for row.Next() {
		row.Scan(&id, &name, &url, &tags)
		i := strconv.Itoa(id)
		r = i + "◄►" + name + "◄►" + url + "◄►" + tags
		results = append(results, r)
	}

	return results
}

func execStatement(q string) sql.Result {
	sqlDB, err := sql.Open("sqlite3", config.Configuration.DBFile)
	checkErr(err)
	defer sqlDB.Close()

	statement, err := sqlDB.Prepare(q)
	checkErr(err)

	result, err := statement.Exec()
	checkErr(err)

	return result
}

func CheckIfDuplicated(s string) bool {
	q := "SELECT COUNT(url) FROM items WHERE url=\"" + s + "\";"

	sqlDB, err := sql.Open("sqlite3", config.Configuration.DBFile)
	checkErr(err)
	defer sqlDB.Close()

	row, err := sqlDB.Query(q)
	checkErr(err)

	defer row.Close()
	var count int

	for row.Next() {
		row.Scan(&count)
	}

	return count > 0

}

func checkErr(err error) {
	if err != nil {
		log.Fatal(err)
	}
}
