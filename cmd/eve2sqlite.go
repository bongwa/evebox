package main

import (
	"evebox/data/sqlite"
	"log"
	"evebox/evereader"
)

func main() {
	db, err := sqlite.Init()
	if err != nil {
		log.Fatal(err)
	}
	log.Println(db)

	eveReader, err := evereader.New("/var/log/suricata/eve.json")
	if err != nil {
		log.Fatal(err)
	}

	for {
		event, err := eveReader.Next()
		if err != nil {
			log.Fatal(err)
		}
		log.Println(event)
	}
}
