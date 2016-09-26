package elasticsearch

import (
	"net/http"
	"io/ioutil"
	"encoding/json"
	"fmt"
	"bytes"
	"io"
)

func GetDecoder(reader io.Reader) (*json.Decoder) {
	decoder := json.NewDecoder(reader)
	decoder.UseNumber()
	return decoder
}

func ToJSON(value interface{}) (string) {
	asJson, err := json.MarshalIndent(value, "", "  ")
	if err != nil {
		return fmt.Sprintf("{failed to render as JSON: %v", value)
	}
	return string(asJson)
}

type ElasticSearch struct {
	baseUrl    string
	httpClient *http.Client
	index      string
}

func New(url string) *ElasticSearch {
	return &ElasticSearch{
		baseUrl:url,
		httpClient:&http.Client{
		},
		index: "logstash",
	}
}

func (es *ElasticSearch) Search(query interface{}) (*SearchResponse, error) {

	body, err := json.Marshal(query)
	if err != nil {
		return nil, err
	}

	request, err := http.NewRequest("POST",
		fmt.Sprintf("%s/%s/_search", es.baseUrl, es.index),
		bytes.NewBuffer(body))
	if err != nil {
		return nil, err
	}
	request.Header.Set("Content-Type", "application/json")

	response, err := es.httpClient.Do(request)
	if err != nil {
		return nil, err
	}

	// Read in the body before decoding, so we can keep a copy of the
	// raw response.
	responseBody, err := ioutil.ReadAll(response.Body)
	if err != nil {
		return nil, err
	}
	decoder := GetDecoder(bytes.NewBuffer(responseBody))
	var searchResponse SearchResponse
	err = decoder.Decode(&searchResponse)
	if err != nil {
		return nil, err
	}
	searchResponse.raw = string(responseBody)

	return &searchResponse, nil
}

func (es *ElasticSearch) Ping() (*PingResponse, error) {

	req, err := http.NewRequest("GET", es.baseUrl, nil)
	if err != nil {
		return nil, err
	}
	response, err := es.httpClient.Do(req)
	if err != nil {
		return nil, err
	}

	if response.StatusCode != 200 {
		return nil, NewElasticSearchError(response)
	}

	decoder := json.NewDecoder(response.Body)
	decoder.UseNumber()
	var body PingResponse
	if err := decoder.Decode(&body); err != nil {
		return nil, err
	}

	return &body, nil
}

type ElasticSearchError struct {
	// The raw error body as returned from the server.
	Raw string
}

func (e ElasticSearchError) Error() string {
	return e.Raw
}

func NewElasticSearchError(response *http.Response) ElasticSearchError {

	error := ElasticSearchError{}

	raw, _ := ioutil.ReadAll(response.Body)
	if raw != nil {
		error.Raw = string(raw)
	}

	return error
}