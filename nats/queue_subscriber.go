package nats

import (
	"os"
	"os/signal"
	"time"

	"github.com/myheartz/grok/logging"
	nats "github.com/nats-io/go-nats"
)

//Subscriber is an asynchronous subscriber.
type Subscriber struct {
	conn     *nats.EncodedConn
	subject  string
	queue    string
	handler  MessageHandler
	messages chan *Message
	shutdown chan os.Signal
	done     chan bool
}

//MessageHandler handles incoming subject messages.
type MessageHandler func(*Message) error

//NewSubscriber creates a new subscriber.
func NewSubscriber(conn *nats.EncodedConn) *Subscriber {
	return &Subscriber{
		conn:     conn,
		messages: make(chan *Message),
	}
}

//SetSubject sets the subscription subject.
func (subscriber *Subscriber) SetSubject(subject string) *Subscriber {
	subscriber.subject = subject
	return subscriber
}

//SetHandler sets the subscription handler.
func (subscriber *Subscriber) SetHandler(handler MessageHandler) *Subscriber {
	subscriber.handler = handler
	return subscriber
}

//SetQueue sets the subscription handler.
func (subscriber *Subscriber) SetQueue(queue string) *Subscriber {
	subscriber.queue = queue
	return subscriber
}

//Run starts the subject subscription.
func (subscriber *Subscriber) Run() error {
	_, err := subscriber.conn.BindRecvQueueChan(subscriber.subject, subscriber.queue, subscriber.messages)

	if err != nil {
		return err
	}

	go func() {
		for {
			subscriber.handler(
				<-subscriber.messages,
			)
		}
	}()

	<-subscriber.done

	return nil
}

func (subscriber *Subscriber) handleShutdown() {
	subscriber.shutdown = make(chan os.Signal)
	subscriber.done = make(chan bool)
	signal.Notify(subscriber.shutdown, os.Interrupt)

	go func() {
		sig := <-subscriber.shutdown
		logging.LogInfo("caught sig: %+v", sig)
		logging.LogInfo("waiting 5 seconds to finish processing")

		subscriber.conn.Drain()
		subscriber.conn.FlushTimeout(5 * time.Second)

		subscriber.done <- true
	}()
}
