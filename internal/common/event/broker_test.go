package event

import (
	"sync"
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestChannelBroker_shouldPublish_whenNoSubscribers(t *testing.T) {
	assert := assert.New(t)
	// given
	broker := NewChannelBroker()
	defer broker.Close()

	// when
	event, err := New("data", "routingKey")
	assert.NoError(err)
	err = broker.Publish(event)

	// then
	assert.NoError(err)
}

func TestChannelBroker_shouldPublish_whenOneSubscriber(t *testing.T) {
	assert := assert.New(t)
	// given
	broker := NewChannelBroker()
	defer broker.Close()

	// and when
	sub := broker.subscribe("x.test.event")
	event, err := New("data", "x.test.event")
	assert.NoError(err)

	// and given it is synchronised
	var wg sync.WaitGroup
	wg.Add(1)

	go func() {
		defer wg.Done()
		for e := range sub {
			assert.Equal(event, e)
			return
		}
	}()

	// when
	err = broker.Publish(event)
	assert.NoError(err)

	wg.Wait()
	broker.Wait()
}

func TestChannelBroker_shouldPublish_whenMultipleSubscribers(t *testing.T) {
	assert := assert.New(t)
	// given
	broker := NewChannelBroker()
	defer broker.Close()

	// and when
	sub1 := broker.subscribe("x.test.event")
	sub2 := broker.subscribe("x.test.event")
	event, err := New("data", "x.test.event")
	assert.NoError(err)

	// and given it is synchronised
	var wg sync.WaitGroup
	wg.Add(2)

	go func() {
		defer wg.Done()
		for e := range sub1 {
			assert.Equal(event, e)
			return
		}
	}()

	go func() {
		defer wg.Done()
		for e := range sub2 {
			assert.Equal(event, e)
			data, ok := e.Data().(string)
			if !ok {
				assert.Fail("Data is not a string")
			}
			assert.Equal("data", data)
			return
		}
	}()

	// when
	err = broker.Publish(event)
	assert.NoError(err)

	wg.Wait()
	broker.Wait()
}

func TestRouting(t *testing.T) {
	assert := assert.New(t)
	// given
	broker := NewChannelBroker()
	defer broker.Close()

	// and when
	// sub1 := broker.Subscribe("x.test.event")
	// and
	var wg sync.WaitGroup
	wg.Add(1)
	called := 0
	sub1Handler := func(e Event) error {
		defer wg.Done()
		called++
		assert.Equal("x.test.event", e.RoutingKey())
		return nil
	}

	eventToHandlerMap := map[string]func(Event) error{
		"x.test.event": sub1Handler,
	}
	broker.RouteAll(eventToHandlerMap)

	event, err := New("data", "x.test.event")
	assert.NoError(err)

	// when
	err = broker.Publish(event)
	assert.NoError(err)

	wg.Wait()

	// then
	assert.Equal(1, called)
}
