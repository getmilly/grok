package nats

import "time"

//Message wraps all data between pubs/subs.
type Message struct {
	CreatedAt time.Time
	Data      interface{}
	Metadata  map[string]interface{}
}

//NewMessage creates a new message with data.
func NewMessage(data interface{}) *Message {
	return &Message{
		CreatedAt: time.Now(),
		Data:      data,
		Metadata:  make(map[string]interface{}),
	}
}

func (message *Message) SetMetadata(key string, value interface{}) {
	message.Metadata[key] = value
}
