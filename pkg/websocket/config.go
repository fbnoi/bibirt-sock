package websocket

import "time"

var default_config = &Config{
	HandshakeTimeout:  time.Second,
	ReadBufferSize:    1024,
	WriteBufferSize:   1024,
	EnableCompression: false,
	PingInterval:      0,
}

type Config struct {
	HandshakeTimeout  time.Duration
	ReadBufferSize    int
	WriteBufferSize   int
	EnableCompression bool
	PingInterval      time.Duration
}
