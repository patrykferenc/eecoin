package mock

import "github.com/patrykferenc/eecoin/internal/common/event"

type Publisher struct {
	events []event.Event
	Called int
}

func (p *Publisher) Publish(event event.Event) error {
	p.Called++
	p.events = append(p.events, event)
	return nil
}

func (p *Publisher) EventWasPublished(event event.Event) bool {
	for _, e := range p.events {
		if e == event {
			return true
		}
	}
	return false
}
