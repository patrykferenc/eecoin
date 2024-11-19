package event

import (
	"time"

	"github.com/google/uuid"
)

type Event interface {
	ID() string
	Timestamp() time.Time
	RoutingKey() string
	Data() interface{}
}

type SimpleEvent struct {
	id         string
	timestamp  time.Time
	data       interface{}
	routingKey string
}

func (e SimpleEvent) ID() string {
	return e.id
}

func (e SimpleEvent) Timestamp() time.Time {
	return e.timestamp
}

func (e SimpleEvent) RoutingKey() string {
	return e.routingKey
}

func (e SimpleEvent) Data() interface{} {
	return e.data
}

func New[T any](data T, routingKey string) (SimpleEvent, error) {
	id, err := uuid.NewV7()
	if err != nil {
		return SimpleEvent{}, err
	}

	return SimpleEvent{
		data:       data,
		timestamp:  time.Now(),
		id:         id.String(),
		routingKey: routingKey,
	}, nil
}
