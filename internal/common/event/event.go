package event

import (
	"time"

	"github.com/google/uuid"
)

type Event interface {
	ID() string
	Timestamp() time.Time
}

type simpleEvent[T any] struct {
	id        string
	timestamp time.Time
	data      T
}

func (e simpleEvent[T]) ID() string {
	return e.id
}

func (e simpleEvent[T]) Timestamp() time.Time {
	return e.timestamp
}

func New[T any](data T) (Event, error) {
	id, err := uuid.NewV7()
	if err != nil {
		return nil, err
	}

	return simpleEvent[T]{
		data:      data,
		timestamp: time.Now(),
		id:        id.String(),
	}, nil
}
