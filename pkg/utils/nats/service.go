package nats

import (
	"github.com/nats-io/nats-server/v2/server"
	"github.com/nats-io/nats.go"
)

func NewNatsServer(port int) (*server.Server, error) {
	opts := &server.Options{
		Port: port,
	}
	return server.NewServer(opts)
}
func NewNatsClient(server string) (*nats.Conn, error) {
	return nats.Connect(server)
}
