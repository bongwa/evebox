package elasticsearch

import (
	"testing"
	"log"
	"encoding/json"
	"io"
)

func _TestElasticSearch(t *testing.T) {

	es := New("http://10.16.1.10:9200")
	response, err := es.Ping()
	if err != nil {
		log.Println("Failed to ping elastic search:", err)
	}

	log.Println("Name:", response.Name)
	log.Println("Cluster Name:", response.ClusterName)
	log.Println("Version:", response.Version.Number)

	//log.Println(response.Get("version"))
	//log.Println(response.Get2("version").Get("number").Value())

}

func TestFindById(t *testing.T) {
	es := New("http://10.16.1.10:9200")
	service := NewService(es)
	service.GetEventById("AVcqpfbxGHWznrB0lL-T")
}

func _TestSearch(t *testing.T) {

	//es := New("http://10.16.1.10:9200")

	query := mapping{
		"query": mapping{
			"filtered": mapping{
				"filter": mapping{
					"and": list{
						mapping{
							"exists": mapping{
								"field": "event_type",
							},
						},
					},
				},
			},
		},
	}

	query["sort"] = list{
		mapping{
			"@timestamp": mapping{
				"order": SORT_DESC,
			},
		},
	}

	asJson, err := json.MarshalIndent(query, "", "    ")
	if err != nil {
		log.Fatal(err)
	}
	log.Println("JSON:", string(asJson))

	if true {
		return
	}

	//query := map[string]interface{}{
	//	"query": map[string]interface{}{
	//		"filtered": nil,
	//	},
	//}

	log.Println(query)
}

type FakeReader struct {
	count int
}

func (fr *FakeReader) Read(p []byte) (int, error) {
	if fr.count > 3 {
		return 0, io.EOF
	}
	fr.count++
	log.Println("Fake reader: array size:", len(p))
	log.Println("Capacity:", cap(p))
	p[0] = 'A'
	return 1, nil
}

func TestBulk(t *testing.T) {
	fakeReader := &FakeReader{}
	es := New("http://10.16.1.10:9200")
	es.Bulk(fakeReader)
}