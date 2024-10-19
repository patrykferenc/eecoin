package eventtest

import "github.com/patrykferenc/eecoin/internal/common/event"

type MockedPublisher struct {
	published []event.Event
}

func NewMockedPublisher() *MockedPublisher {
	return &MockedPublisher{}
}

func (p *MockedPublisher) Publish(event event.Event) error {
	p.published = append(p.published, event)
	return nil
}

func (p *MockedPublisher) Published() []event.Event {
	return p.published
}
