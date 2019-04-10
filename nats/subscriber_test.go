package nats_test

import (
	"os"
	"reflect"
	"sync"
	"testing"

	uuid "github.com/satori/go.uuid"
	"github.com/stretchr/testify/assert"

	gnats "github.com/getmilly/grok/nats"
	"github.com/nats-io/go-nats-streaming"
)

type Testing struct {
	Value interface{}
}

func TestSubscriber_PubSub(t *testing.T) {
	conn := getConn()

	queue := uuid.NewV4().String()
	subject := uuid.NewV4().String()

	messageCount := 5
	wg := &sync.WaitGroup{}
	wg.Add(messageCount)

	go func() {
		gnats.NewSubscriber(conn).
			WithQueue(queue).
			WithSubject(subject).
			WithMessageType(reflect.TypeOf(Testing{})).
			WithHandler(func(m interface{}) error {
				_, ok := m.(*Testing)
				assert.True(t, ok)
				wg.Done()
				return nil
			}).
			Run()
	}()

	producer := gnats.NewProducer(conn)

	for i := 0; i < 5; i++ {
		message, _ := gnats.NewMessage(&Testing{
			Value: i,
		})

		err := producer.Publish(subject, message)

		assert.NoError(t, err)
	}

	wg.Wait()
}

func getConn() stan.Conn {
	conn, err := stan.Connect(os.Getenv("NATS_CLUSTER"), uuid.NewV4().String(), stan.NatsURL(os.Getenv("NATS_URL")))

	if err != nil {
		panic(err)
	}

	return conn
}
