package rabbitmq

import (
	"encoding/json"
	"time"

	"github.com/globalsign/mgo/bson"

	"github.com/globalsign/mgo"
)

//Fallback ...
type Fallback struct {
	ID              bson.ObjectId
	Exchange        string
	ExchangeType    string
	Queue           string
	Message         interface{}
	OriginalMessage []byte
	CreatedAt       time.Time
}

//NewFallback ...
func NewFallback(exchange, exchangeType, queue string, message []byte) *Fallback {
	var m interface{}
	json.Unmarshal(message, &m)

	return &Fallback{
		Exchange:        exchange,
		ExchangeType:    exchangeType,
		Queue:           queue,
		OriginalMessage: message,
		Message:         m,
	}
}

//FallbackHandler handles messages that cannot be processed.
type FallbackHandler func(fallback *Fallback) error

//MongoFallback save messages to mongodb.
func MongoFallback(session *mgo.Session) FallbackHandler {
	return func(fallback *Fallback) error {
		s := session.Clone()
		defer s.Close()

		err := s.DB("MessageBroker").C("Fallback").Insert(fallback)

		if err != nil {
			panic(err)
		}

		return err
	}
}
