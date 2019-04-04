package nats

import nats "github.com/nats-io/go-nats"

//Producer wraps the publish of messages to NATS.
type Producer struct {
	conn *nats.EncodedConn
}

//NewProducer creates a new producer.
func NewProducer(conn *nats.EncodedConn) *Producer {
	return &Producer{conn}
}

//Publish sends a message to a subject.
func (producer *Producer) Publish(subject string, message interface{}) error {
	return producer.conn.Publish(subject, message)
}
