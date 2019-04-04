package nats

import (
	"errors"
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
	stop     bool
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
		stop:     false,
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
	if err := subscriber.validate(); err != nil {
		return err
	}

	_, err := subscriber.conn.BindRecvQueueChan(subscriber.subject, subscriber.queue, subscriber.messages)

	if err != nil {
		return err
	}

	go func() {
		for !subscriber.stop {
			subscriber.handler(<-subscriber.messages)
		}
	}()

	subscriber.handleShutdown()

	<-subscriber.done

	return nil
}

func (subscriber *Subscriber) validate() error {
	if subscriber.subject == "" {
		return errors.New("`subject` must be set")
	}

	if subscriber.queue == "" {
		return errors.New("`queue` must be set")
	}

	if subscriber.handler == nil {
		return errors.New("`handler` must be set")
	}

	return nil
}

func (subscriber *Subscriber) handleShutdown() {
	subscriber.done = make(chan bool)
	subscriber.shutdown = make(chan os.Signal)
	signal.Notify(subscriber.shutdown, os.Interrupt)

	go func() {
		sig := <-subscriber.shutdown
		logging.LogInfo("caught sig: %+v", sig)
		logging.LogInfo("waiting 5 seconds to finish processing")

		subscriber.stop = true
		subscriber.conn.Drain()
		subscriber.conn.FlushTimeout(5 * time.Second)

		subscriber.done <- true
	}()
}
