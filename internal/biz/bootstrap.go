package biz

import (
	"flynoob/bibirt-sock/internal/message"
	"flynoob/bibirt-sock/pkg/websocket"
)

func (*ClientHandler) registerMessage() {
	websocket.RegisterMessage(&message.Ping{})
	websocket.RegisterMessage(&message.Connected{})
	websocket.RegisterMessage(&message.Disconnected{})
}

func (handler *ClientHandler) registerMessageHandler(c *websocket.Client) {
	handler.monitorHealth(c)
}
