package event

type Publisher interface {
	Publish(event Event) error
}
