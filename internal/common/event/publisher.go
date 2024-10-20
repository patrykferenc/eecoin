package event

type Publisher interface {
	Publish(event Event) error
}

type Subscriber interface {
	Subscribe(routingKey string) <-chan Event
}
