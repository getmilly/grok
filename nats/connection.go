package nats

import (
	"time"

	"github.com/getmilly/grok/logging"
	nats "github.com/nats-io/go-nats"
)

// Connect creates a new encoded connection to a NATS server.
// JSON is used as encoder.
// Multiple servers can be specified comma-separated at `address`.
func Connect(address string, token string) (*nats.EncodedConn, error) {
	nc, _ := nats.Connect(
		address,
		nats.Token(token),
		nats.MaxReconnects(10),
		nats.ReconnectWait(2*time.Second),
		nats.DontRandomize(),
		nats.ClosedHandler(closedHandler),
	)
	return nats.NewEncodedConn(nc, nats.JSON_ENCODER)
}

func closedHandler(conn *nats.Conn) {
	err := conn.LastError()

	if err != nil {
		logging.LogWith(err).Error("subscription error")
	}
}
