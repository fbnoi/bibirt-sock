package biz

import (
	"flynoob/bibirt-sock/internal/message"
	"flynoob/bibirt-sock/pkg/websocket"
)

func (*ConnUseCase) registerMessage() {
	websocket.RegisterMessage(&message.Ping{})
	websocket.RegisterMessage(&message.Connected{})
	websocket.RegisterMessage(&message.Disconnected{})
}

func (useCase *ConnUseCase) registerMessageHandler(c *websocket.Client) {
	useCase.MonitorHealth(c)
}
