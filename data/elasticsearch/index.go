package elasticsearch

import (
	"time"
	"math/rand"
	"evebox/eve"
	"fmt"
	"bytes"
	"io/ioutil"
	"encoding/json"
	"net/http"
)

var chars = []rune("abcdefghijklmnopqrstuvwxyzABCDEFGHIJKLMNOPQRSTUVWXYZ")

func init() {
	rand.Seed(time.Now().UnixNano())
}

func GenerateId() string {
	id := make([]rune, 20)
	for i := range id {
		id[i] = chars[rand.Intn(len(chars))]
	}
	return string(id)
}

func (es *ElasticSearch) IndexRawEveEvent(event eve.RawEveEvent) (error) {
	id := GenerateId()

	timestamp, err := event.GetTimestamp()
	if err != nil {
		return err
	}
	index := fmt.Sprintf("%s-%s", es.index, timestamp.UTC().Format("2006.01.02"))

	var buf bytes.Buffer
	encoder := json.NewEncoder(&buf)
	encoder.Encode(event)
	request, err := http.NewRequest("POST",
		fmt.Sprintf("%s/%s/log/%s", es.baseUrl, index, id), &buf)
	if err != nil {
		return err
	}
	response, err := es.httpClient.Do(request)
	if err != nil {
		return err
	}

	// Required for connection re-use.
	ioutil.ReadAll(response.Body)
	response.Body.Close()

	return nil
}

