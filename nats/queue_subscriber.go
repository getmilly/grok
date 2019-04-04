package nats

import nats "github.com/nats-io/go-nats"

//Subscriber is an asynchronous subscriber.
type Subscriber struct {
	conn     *nats.EncodedConn
	subject  string
	handler  MessageHandler
	messages chan *Message
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

func shutdown() {

}
