package core

type EventService interface {
	GetEventById(id string) (*Event, error)
}

type Event struct {
	Id     string
	Source map[string]interface{}
}