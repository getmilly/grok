package nats

import (
	"encoding/json"

	"github.com/nats-io/go-nats-streaming"
)

//Producer wraps the publish of messages to NATS.
type Producer struct {
	conn stan.Conn
}

//NewProducer creates a new producer.
func NewProducer(conn stan.Conn) *Producer {
	return &Producer{conn}
}

//Publish sends a message to a subject.
func (producer *Producer) Publish(subject string, message *Message) error {
	m, err := json.Marshal(message)

	if err != nil {
		return err
	}

	return producer.conn.Publish(subject, m)
}
