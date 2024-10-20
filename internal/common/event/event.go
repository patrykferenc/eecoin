package event

import (
	"time"

	"github.com/google/uuid"
)

type Event interface {
	ID() string
	Timestamp() time.Time
	RoutingKey() string
}

type simpleEvent[T any] struct {
	id         string
	timestamp  time.Time
	data       T
	routingKey string
}

func (e simpleEvent[T]) ID() string {
	return e.id
}

func (e simpleEvent[T]) Timestamp() time.Time {
	return e.timestamp
}

func (e simpleEvent[T]) RoutingKey() string {
	return e.routingKey
}

func New[T any](data T, routingKey string) (Event, error) {
	id, err := uuid.NewV7()
	if err != nil {
		return nil, err
	}

	return simpleEvent[T]{
		data:       data,
		timestamp:  time.Now(),
		id:         id.String(),
		routingKey: routingKey,
	}, nil
}
