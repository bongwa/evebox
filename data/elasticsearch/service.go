package elasticsearch

import (
	"evebox/core"
	"log"
)

// The Elastic Search event service.
//
// Implements the core.EventService interface.
type Service struct {
	elasticsearch *ElasticSearch
}

func NewService(elasticsearch *ElasticSearch) core.EventService {
	return &Service{elasticsearch:elasticsearch}
}

func (s *Service) GetEventById(id string) (*core.Event, error) {
	query := mapping{
		"query": mapping{
			"filtered": mapping{
				"filter": mapping{
					"term": mapping{
						"_id": id,
					},
				},
			},
		},
	}
	response, err := s.elasticsearch.Search(query)
	if err != nil {
		log.Println("error: %v", err)
	} else {
		log.Println(ToJSON(response))
	}

	return nil, nil
}