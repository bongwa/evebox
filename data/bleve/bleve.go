package bleve

import (
	"github.com/blevesearch/bleve"
	"log"
)

func Init() bleve.Index {

	mapping := bleve.NewIndexMapping()
	index, err := bleve.New("bleve.db", mapping)
	if err != nil {
		log.Fatal(err)
	}

	return index

}