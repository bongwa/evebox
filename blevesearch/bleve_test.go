package blevesearch

import (
	"log"
	"testing"

	"github.com/blevesearch/bleve"
	"github.com/jasonish/evebox/evereader"
	"github.com/satori/go.uuid"
	"time"
)

func TestBleve(t *testing.T) {


	//mapping := bleve.NewIndexMapping()
	index, err := bleve.Open("bleve.db")
	if err != nil {
		mapping := bleve.NewIndexMapping()
		index, err = bleve.New("bleve.db", mapping)
	}

	if err != nil {
		t.Fatal(err)
	}

	//log.Println(mapping)
	log.Println(index)

	log.Println("Opening eve.json...")
	reader, err := evereader.New("/var/log/suricata/eve.json-20161124")
	if err != nil {
		log.Fatal(err)
	}
	log.Println("Opened...")

	log.Println("Reading...")

	lastTs := time.Now()
	count := int64(0)

	for {
		next, err := reader.Next()
		if err != nil {
			log.Fatal(err)
		}

		id := uuid.NewV1()

		err = index.Index(id.String(), next)
		if err != nil {
			log.Fatal(err)
		}

		count++

		now := time.Now()
		//diff := now.Sub(lastTs).Seconds()
		//log.Println(diff)
		if now.Sub(lastTs).Seconds() > 1 {
			log.Println("Count:", count)
			lastTs = now
			count = 0
		}

	}

}
