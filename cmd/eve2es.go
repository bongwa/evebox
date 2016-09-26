package main

import (
"evebox/data/sqlite"
"log"
"evebox/evereader"
	"evebox/eve"
	"evebox/data/elasticsearch"
)

func main() {
	db, err := sqlite.Init()
	if err != nil {
		log.Fatal(err)
	}
	log.Println(db)

	es := elasticsearch.New("http://10.16.1.10:9200")
	ping, err := es.Ping()
	if err != nil {
		log.Fatal("Failed to connect to Elastic Search:", err)
	}
	log.Println("Connected to Elastic Search", ping.Version.Number)

	eveReader, err := evereader.New("/var/log/suricata/eve.json")
	if err != nil {
		log.Fatal(err)
	}

	for {
		event, err := eveReader.Next()
		if err != nil {
			log.Fatal(err)
		}
		es.IndexRawEveEvent(event)
		timestamp, _ := event.GetTimestamp()
		log.Println(timestamp.UTC().Format(eve.EveTimestampFormat))
	}
}
