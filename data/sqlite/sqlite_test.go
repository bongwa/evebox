package sqlite

import (
	"testing"
	"log"
	"io/ioutil"
	"os"
	"encoding/json"
	"github.com/satori/go.uuid"
	"time"
	"evebox/evereader"
	"io"
)

func TestInit(t *testing.T) {
	log.Println("TestInit")

	db := Init()

	v1, err := os.Open("./v1.sql")
	if err != nil {
		log.Fatal(err);
	}
	buf, err := ioutil.ReadAll(v1)
	if err != nil {
		log.Fatal(err)
	}
	log.Println(buf)

	res, err := db.Exec(string(buf))
	if err != nil {
		log.Fatal(err)
	}
	log.Println(res)

	eveReader, err := evereader.New("/var/log/suricata/eve.json")
	if err != nil {
		log.Fatal(err)
	}

	var count uint64 = 0

	tx, err := db.Begin()
	if err != nil {
		log.Fatal(err)
	}

	for {

		event, err := eveReader.Next()
		if err != nil {
			if err == io.EOF {
				break
			}
			log.Fatal(err)
		}

		buf, err := json.Marshal(event)
		if err != nil {
			log.Println(err)
		}

		id := uuid.NewV4()

		timestamp := event["timestamp"].(string)
		timestamp = FormatTimestamp(timestamp)

		_, err = tx.Exec("insert into events values ($1, $2, $3)", id, timestamp, buf)
		if err != nil {
			log.Fatal(err)
		}

		_, err = tx.Exec("insert into events_fts values ($1, $2)", id, buf)
		if err != nil {
			log.Fatal(err)
		}

		count++

		if count % 1000 == 0 {
			log.Println(count)
			pos := eveReader.Pos()
			log.Printf("File position: %d", pos)

			tx.Commit()
			tx, err = db.Begin()
			if err != nil {
				log.Fatal(err)
			}
		}
	}

	log.Println("Events read", count)
}

func FormatTimestamp(timestamp string) string {
	var RFC3339Nano_Modified string = "2006-01-02T15:04:05.999999999Z0700"
	result, err := time.Parse(RFC3339Nano_Modified, timestamp)
	if err != nil {
		log.Fatal(err)
	}
	return result.UTC().Format("2006-01-02T15:04:05.999999999")
}