package server

import (
	"flynoob/bibirt-sock/internal/biz"
	"flynoob/bibirt-sock/pkg/websocket"
)

func NewServer() *websocket.Server {
	srv := websocket.NewServer()
	biz.Bootstrap()
	srv.OnNewConnection(func(c *websocket.Client) {
		biz.Auth(c)
		biz.HandlePing(c)
		biz.MonitorHealth(c)
		if c.Upgrade() {
			defer c.Close()
			c.Loop()
		}
	})
	return srv
}
