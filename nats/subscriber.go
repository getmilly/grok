package nats

import (
	"encoding/json"
	"errors"
	"os"
	"os/signal"
	"reflect"
	"time"

	"github.com/nats-io/go-nats-streaming"

	"github.com/getmilly/grok/logging"
)

//Subscriber is an asynchronous subscriber.
type Subscriber struct {
	conn         stan.Conn
	subject      string
	queue        string
	messageType  reflect.Type
	subscription stan.Subscription
	handler      MessageHandler
	shutdown     chan os.Signal
	done         chan bool
}

//MessageHandler handles incoming subject messages.
type MessageHandler func(interface{}) error

//NewSubscriber creates a new subscriber.
func NewSubscriber(conn stan.Conn) *Subscriber {
	return &Subscriber{
		conn: conn,
	}
}

//WithSubject sets the subscription subject.
func (subscriber *Subscriber) WithSubject(subject string) *Subscriber {
	subscriber.subject = subject
	return subscriber
}

//WithHandler sets the subscription handler.
func (subscriber *Subscriber) WithHandler(handler MessageHandler) *Subscriber {
	subscriber.handler = handler
	return subscriber
}

//WithQueue sets the subscription queue.
func (subscriber *Subscriber) WithQueue(queue string) *Subscriber {
	subscriber.queue = queue
	return subscriber
}

//WithMessageType sets the subscription message type.
func (subscriber *Subscriber) WithMessageType(t reflect.Type) *Subscriber {
	subscriber.messageType = t
	return subscriber
}

//Run starts the subject subscription.
func (subscriber *Subscriber) Run() error {
	if err := subscriber.validate(); err != nil {
		return err
	}

	subscription, err := subscriber.conn.QueueSubscribe(
		subscriber.subject,
		subscriber.queue,
		subscriber.messageHandler,
		stan.SetManualAckMode(),
		stan.StartWithLastReceived(),
		stan.DurableName(subscriber.queue),
	)

	if err != nil {
		return err
	}

	subscriber.subscription = subscription

	subscriber.handleShutdown()

	<-subscriber.done

	return nil
}

func (subscriber *Subscriber) messageHandler(msg *stan.Msg) {
	message := &Message{}
	v := reflect.New(subscriber.messageType).Interface()

	if err := json.Unmarshal(msg.Data, &message); err != nil {
		logging.LogWith(err).Warn("payload wasn't type of `Message`")
		return
	}

	if err := json.Unmarshal(message.Data, v); err != nil {
		logging.LogWith(err).Warn("message data wasn't type of `%s`", subscriber.messageType.Name())
		return
	}

	defer func() {
		if err := recover(); err != nil {
			logging.LogWith(err).Error("handler panics")
		}
	}()

	logging.LogWith(v).Info("incoming message")

	if err := subscriber.handler(v); err != nil {
		logging.LogWith(err).Error("handle error")
		return
	}

	msg.Ack()
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

	if subscriber.messageType == nil {
		return errors.New("`messageType` must be set")
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

		subscriber.subscription.Unsubscribe()
		subscriber.conn.NatsConn().Drain()
		subscriber.conn.NatsConn().FlushTimeout(5 * time.Second)
		subscriber.conn.Close()

		subscriber.done <- true
	}()
}
