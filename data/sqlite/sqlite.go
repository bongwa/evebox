// +build cgo

package sqlite

import (
	"database/sql"
	"log"
	_ "github.com/mattn/go-sqlite3"
	"os"
)

func Init() *sql.DB {
	os.Remove("./sqlite.db")
	db, err := sql.Open("sqlite3", "./sqlite.db")
	if err != nil {
		log.Fatal(err);
	}
	return db
}
