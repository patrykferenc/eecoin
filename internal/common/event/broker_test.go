package event_test

import (
	"sync"
	"testing"

	"github.com/patrykferenc/eecoin/internal/common/event"
	"github.com/stretchr/testify/assert"
)

func TestChannelBroker_shouldPublish_whenNoSubscribers(t *testing.T) {
	assert := assert.New(t)
	// given
	broker := event.NewChannelBroker()

	// when
	event, err := event.New("data", "routingKey")
	assert.NoError(err)
	err = broker.Publish(event)

	// then
	assert.NoError(err)
}

func TestChannelBroker_shouldPublish_whenOneSubscriber(t *testing.T) {
	assert := assert.New(t)
	// given
	broker := event.NewChannelBroker()

	// and when
	sub := broker.Subscribe("x.test.event")
	event, err := event.New("data", "x.test.event")
	assert.NoError(err)

	// and given it is synchronised
	var wg sync.WaitGroup
	wg.Add(1)

	go func() {
		defer wg.Done()
		for {
			select {
			case e := <-sub:
				assert.Equal(event, e)
				return
			}
		}
	}()

	// when
	err = broker.Publish(event)
	assert.NoError(err)

	wg.Wait()
	broker.Wait()
}
