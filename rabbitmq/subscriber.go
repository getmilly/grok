package rabbitmq

import (
	"os"
	"sync"
	"time"

	"github.com/pborman/uuid"
	"github.com/streadway/amqp"
)

//RabbitSubscriber ...
type RabbitSubscriber struct {
	address      string
	exchange     string
	exchangeType string
	queue        string
	routingKey   string
	args         amqp.Table
	fallback     FallbackHandler

	mu sync.Mutex
}

//MessageHandler ...
type MessageHandler func(body []byte) (bool, error)

var (
	rabbitCloseError chan *amqp.Error
)

//NewRabbitSubscriber instance of subscriber.
func NewRabbitSubscriber(
	address,
	exchange,
	exchangeType,
	queue,
	routingKey string,
	args map[string]interface{},
) *RabbitSubscriber {
	return &RabbitSubscriber{
		address:      address,
		exchange:     exchange,
		exchangeType: exchangeType,
		queue:        queue,
		routingKey:   routingKey,
		args:         args,
	}
}

//Fallback adds a fallback handler
func (subscriber *RabbitSubscriber) Fallback(handler FallbackHandler) {
	subscriber.mu.Lock()
	defer subscriber.mu.Unlock()
	subscriber.fallback = handler
}

//Subcribe to a queue.
func (subscriber *RabbitSubscriber) Subcribe(handler MessageHandler) error {
	cnn := subscriber.connect()
	err := subscriber.subcribe(cnn, handler)
	return err
}

//Subcribe to a queue.
func (subscriber *RabbitSubscriber) subcribe(cnn *amqp.Connection, handler MessageHandler) error {
	var rabbitErr *amqp.Error

	ch, err := cnn.Channel()
	defer ch.Close()

	if err != nil {
		panic(err)
	}

	err = ch.ExchangeDeclare(subscriber.exchange, subscriber.exchangeType, true, false, false, false, nil)

	if err != nil {
		panic(err)
	}

	q, err := ch.QueueDeclare(subscriber.queue, true, false, false, false, subscriber.args)

	if err != nil {
		panic(err)
	}

	err = ch.QueueBind(q.Name, subscriber.routingKey, subscriber.exchange, false, nil)

	if err != nil {
		panic(err)
	}

	hn, _ := os.Hostname()
	msgs, err := ch.Consume(q.Name, "amq.ctag-"+hn+"-"+uuid.New(), false, false, false, false, nil)

	if err != nil {
		panic(err)
	}

	forever := make(chan bool)

	go func() {
		for d := range msgs {
			subscriber.onMessage(d, handler)
		}
	}()

	rabbitErr = <-rabbitCloseError

	if rabbitErr != nil {
		cnn := subscriber.connect()
		subscriber.subcribe(cnn, handler)
		forever <- true
	}

	<-forever

	return nil
}

func (subscriber *RabbitSubscriber) connect() *amqp.Connection {
	for {
		conn, err := amqp.Dial(subscriber.address)

		if err == nil {
			rabbitCloseError = make(chan *amqp.Error)
			conn.NotifyClose(rabbitCloseError)
			return conn
		}

		time.Sleep(1000 * time.Millisecond)
	}
}

func (subscriber *RabbitSubscriber) onMessage(delivery amqp.Delivery, handler MessageHandler) {
	ack, _ := handler(delivery.Body)

	if ack {
		delivery.Ack(true)
		return
	}

	if !ack && !delivery.Redelivered {
		delivery.Nack(false, true)
		return
	}

	if subscriber.fallback != nil {
		fallback := NewFallback(
			subscriber.exchange,
			subscriber.exchangeType,
			subscriber.queue,
			delivery.Body,
		)
		subscriber.fallback(fallback)
	}
}
