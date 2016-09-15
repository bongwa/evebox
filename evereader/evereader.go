package evereader

import (
	"os"
	"gopkg.in/square/go-jose.v1/json"
	"bufio"
)

type EveReader struct {
	file    *os.File
	reader  *bufio.Reader
	decoder *json.Decoder
	lineno  uint64
}

func New(path string) (*EveReader, error) {

	file, err := os.Open(path)
	if err != nil {
		return nil, err
	}

	reader := bufio.NewReader(file)
	decoder := json.NewDecoder(reader)
	decoder.UseNumber()

	return &EveReader{file: file, reader:reader, decoder:decoder}, nil
}

// Skip to a line number in the file. Must be called before any reading is
// done.
func (er *EveReader) SkipTo(lineno uint64) error {
	if er.lineno != 0 {
		return nil
	}
	for lineno > 0 {
		_, err := er.reader.ReadBytes('\n')
		if err != nil {
			return err
		}
		lineno--
	}
	return nil
}

// Get the current position in the file. For EveReaders this is the line number
// as the actual file offset is not useful due to buffering in the json
// decoder as well as bufio.
func (er *EveReader) Pos() (uint64) {
	return er.lineno
}


func (er *EveReader) Next() (map[string]interface{}, error) {

	var event map[string]interface{}

	if err := er.decoder.Decode(&event); err != nil {
		return nil, err
	}

	er.lineno++

	return event, nil
}