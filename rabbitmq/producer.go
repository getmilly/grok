package rabbitmq

import (
	"encoding/json"

	"github.com/streadway/amqp"
)

//RabbitProducer ...
type RabbitProducer struct {
	address string
}

//NewRabbitProducer create a new RabbitMQ producer.
func NewRabbitProducer(address string) *RabbitProducer {
	return &RabbitProducer{
		address: address,
	}
}

//Publish a message to RabbitMQ
func (broker RabbitProducer) Publish(exchange, exchangeType string, body ...interface{}) error {
	cnn, err := amqp.Dial(broker.address)
	defer cnn.Close()

	if err != nil {
		return err
	}

	ch, err := cnn.Channel()
	defer ch.Close()

	if err != nil {
		return err
	}

	err = ch.ExchangeDeclare(exchange, exchangeType, true, false, false, false, nil)

	if err != nil {
		return err
	}

	for _, b := range body {
		messageBody, err := json.Marshal(b)
		if err != nil {
			return err
		}

		err = ch.Publish(exchange, exchangeType, false, false, amqp.Publishing{ContentType: "application/json", Body: messageBody})

		if err != nil {
			return err
		}
	}

	return nil
}
