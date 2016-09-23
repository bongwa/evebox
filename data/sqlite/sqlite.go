// +build cgo

package sqlite

import (
	"database/sql"
	_ "github.com/mattn/go-sqlite3"
	"github.com/satori/go.uuid"
	"time"
	"log"
	"encoding/json"
	"os"
)

var db *sql.DB

type SqliteDatastore struct {
	db *sql.DB
}

func Init() (*sql.DB, error) {
	os.Remove("./sqlite.db")
	db, err := sql.Open("sqlite3", "./sqlite.db")
	return db, err
}

func AddEvent(event map[string]interface{}) {
	id := uuid.NewV4()

	timestamp := event["timestamp"].(string)
	timestamp = FormatTimestamp(timestamp)

	buf, err := json.Marshal(event)
	if err != nil {
		log.Println(err)
	}

	_, err = db.Exec("insert into events values ($1, $2, $3)", id, timestamp, buf)
	if err != nil {
		log.Fatal(err)
	}

	_, err = db.Exec("insert into events_fts values ($1, $2)", id, buf)
	if err != nil {
		log.Fatal(err)
	}
}

// Format an event timestamp for use in a SQLite column. The format is
// already correct, just needs to be converted to UTC.
func FormatTimestamp(timestamp string) string {
	var RFC3339Nano_Modified string = "2006-01-02T15:04:05.999999999Z0700"
	result, err := time.Parse(RFC3339Nano_Modified, timestamp)
	if err != nil {
		log.Fatal(err)
	}
	return result.UTC().Format("2006-01-02T15:04:05.999999999")
}

