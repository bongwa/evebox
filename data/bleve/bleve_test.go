package bleve

import (
	"testing"
	"os"
	"log"
	"evebox/evereader"
	"github.com/satori/go.uuid"
)

func TestStuff(t *testing.T) {

	os.RemoveAll("bleve.db")

	index := Init()
	log.Println(index)

	eveReader, err := evereader.New("/var/log/suricata/eve.json")
	if err != nil {
		log.Fatal(err)
	}

	var count uint64 = 0;

	for {
		event, err := eveReader.Next()
		if err != nil {
			log.Fatal(err)
		}

		id := uuid.NewV4()
		index.Index(id.String(), event)

		count++
		if count % 100 == 0 {
			log.Println(count)
		}
	}
}