package db

import (
	"database/sql"
	"log"
	"os"
	"strconv"
	"strings"

	_ "github.com/mattn/go-sqlite3"
)

var id int
var name, path, tags string

func InitDB() {
	file, err := os.Create("bookmarks.db")
	if err != nil {
		log.Fatal(err)
	}
	file.Close()

	sqlDB, err := sql.Open("sqlite3", "bookmarks.db")
	if err != nil {
		log.Fatal(err)
	}

	createItemsTable := `CREATE TABLE items (
      "id" integer NOT NULL PRIMARY KEY AUTOINCREMENT,      
      "name" TEXT,
      "path" TEXT,
      "tags" TEXT    
     );`

	statement, err := sqlDB.Prepare(createItemsTable)
	if err != nil {
		log.Fatal(err)
	}

	statement.Exec()
	defer sqlDB.Close()
}

func InsertData(n, p, t string) {
	q := "INSERT INTO items(name, path, tags) VALUES (\"" + n + "\",\"" + p + "\",\"" + t + "\");"
	execStatement(q)
}

func UpdateData(i, n, p, t string) {
	q := "UPDATE items set tags=\"" + t + "\",name=\"" + n + "\",path=\"" + p + "\" where id=\"" + i + "\";"
	execStatement(q)
}

func RemoveData(i string) {
	q := "delete from items where id=" + i + ";"
	execStatement(q)
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

	sqlDB, err := sql.Open("sqlite3", "bookmarks.db")
	checkErr(err)
	defer sqlDB.Close()

	row, err := sqlDB.Query(q)
	checkErr(err)

	defer row.Close()
	var results []string
	var r string

	for row.Next() {
		row.Scan(&id, &name, &path, &tags)
		i := strconv.Itoa(id)
		r = i + "◄►" + name + "◄►" + path + "◄►" + tags
		results = append(results, r)
	}

	return results
}

func execStatement(q string) {
	sqlDB, err := sql.Open("sqlite3", "bookmarks.db")
	checkErr(err)
	defer sqlDB.Close()

	statement, err := sqlDB.Prepare(q)
	checkErr(err)

	_, err = statement.Exec()
	checkErr(err)
}

func checkErr(err error) {
	if err != nil {
		log.Fatal(err)
	}
}
