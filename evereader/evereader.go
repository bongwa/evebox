package evereader

import (
	"os"
	"bufio"
	"gopkg.in/square/go-jose.v1/json"
)

type EveReader struct {
	decoder *json.Decoder
}

func New(path string) (*EveReader, error) {

	file, err := os.Open(path)
	if err != nil {
		return nil, err
	}

	reader := bufio.NewReader(file)
	decoder := json.NewDecoder(reader)
	decoder.UseNumber()

	return &EveReader{decoder:decoder}, nil
}

func (er *EveReader) Next() (map[string]interface{}, error) {

	var event map[string]interface{}

	if err := er.decoder.Decode(&event); err != nil {
		return nil, err
	}

	return event, nil
}