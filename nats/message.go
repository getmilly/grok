package nats

import (
	"encoding/json"
	"time"

	uuid "github.com/satori/go.uuid"
)

//Message wraps all data between pubs/subs.
type Message struct {
	ID        string                 `json:"id"`
	CreatedAt time.Time              `json:"created_at"`
	Data      []byte                 `json:"data"`
	Metadata  map[string]interface{} `json:"metadata"`
}

//NewMessage creates a new message with data.
func NewMessage(data interface{}) (*Message, error) {
	bdata, err := json.Marshal(data)

	if err != nil {
		return nil, err
	}

	return &Message{
		Data:      bdata,
		CreatedAt: time.Now(),
		ID:        uuid.NewV4().String(),
		Metadata:  make(map[string]interface{}),
	}, nil
}

//SetMetadata sets message metadata.
func (message *Message) SetMetadata(key string, value interface{}) {
	message.Metadata[key] = value
}
