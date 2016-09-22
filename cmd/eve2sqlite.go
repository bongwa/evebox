package main

import (
	"evebox/data/sqlite"
	"log"
)

func main() {
	db := sqlite.Init()
	log.Println(db)
}
