package event

import (
	"log/slog"
	"sync"
)

type ChannelBroker struct {
	channels         map[string][]chan Event
	wg               sync.WaitGroup
	subscribtionLock sync.Mutex
}

func NewChannelBroker() *ChannelBroker {
	return &ChannelBroker{
		channels: make(map[string][]chan Event),
	}
}

func (b *ChannelBroker) Publish(event Event) error {
	b.subscribtionLock.Lock()
	defer b.subscribtionLock.Unlock()

	subscribers, ok := b.channels[event.RoutingKey()]
	if !ok {
		slog.Warn("No subscribers", "routingKey", event.RoutingKey())
		return nil
	}

	slog.Debug("Publishing event", "event", event, "subscribers", len(subscribers))
	for _, subscriber := range subscribers {
		if subscriber == nil {
			slog.Warn("Subscriber is nil, will skip")
			continue
		}
		b.wg.Add(1)
		go func(sub chan Event, e Event) {
			defer b.wg.Done()
			if e == nil {
				slog.Warn("Event is nil, will skip")
				return
			}
			sub <- e
		}(subscriber, event)
	}

	return nil
}

func (b *ChannelBroker) Subscribe(routingKey string) <-chan Event {
	b.subscribtionLock.Lock()
	defer b.subscribtionLock.Unlock()

	if _, ok := b.channels[routingKey]; !ok {
		b.channels[routingKey] = make([]chan Event, 0)
	}

	channel := make(chan Event)
	b.channels[routingKey] = append(b.channels[routingKey], channel)

	return channel
}

func (b *ChannelBroker) Wait() {
	b.wg.Wait()
}
