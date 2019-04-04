package nats

import (
	"time"

	uuid "github.com/satori/go.uuid"
)

//Message wraps all data between pubs/subs.
type Message struct {
	ID        string
	CreatedAt time.Time
	Data      interface{}
	Metadata  map[string]interface{}
}

//NewMessage creates a new message with data.
func NewMessage(data interface{}) *Message {
	return &Message{
		Data:      data,
		CreatedAt: time.Now(),
		ID:        uuid.NewV4().String(),
		Metadata:  make(map[string]interface{}),
	}
}

//SetMetadata sets message metadata.
func (message *Message) SetMetadata(key string, value interface{}) {
	message.Metadata[key] = value
}
