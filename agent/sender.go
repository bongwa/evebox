package agent

type EventSender struct {

	queue []interface{}
}

func NewEventSender() *EventSender {

	queue := make([]interface{}, 0, 1000)

	return &EventSender{
		queue: queue,
	}

}

func (s *EventSender) Submit(event interface{}) {
	s.queue = append(s.queue, event)
}

func (s *EventSender) Reset() {
	s.queue = s.queue[:0]
}