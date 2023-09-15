package biz

import (
	"flynoob/bibirt-sock/internal/message"
	"flynoob/bibirt-sock/pkg/websocket"
	"time"

	"google.golang.org/protobuf/proto"
	"google.golang.org/protobuf/types/known/timestamppb"
)

func (handler *ClientHandler) monitorHealth(client *websocket.Client) {
	handler.handlePing(client)
	handler.scheduler.Every(1).Second().Tag(client.ID()).Do(func() {
		if client.LastPingAt.Add(2 * time.Second).Before(time.Now()) {
			client.Color += 1
		}
		if client.Color == websocket.Red {
			client.Close()
			handler.scheduler.RemoveByTag(client.ID())
		}
	})
}

func (handler *ClientHandler) handlePing(client *websocket.Client) {
	ping := &message.Ping{}
	client.Subscribe(ping, func(m proto.Message) {
		client.Color = websocket.Green
		ping, ok := m.(*message.Ping)
		if !ok {
			return
		}
		client.LastPingAt = time.Now()
		ping.DownTimestamp = timestamppb.New(client.LastPingAt)
		if err := client.Send(ping); err != nil {
			handler.log.Errorf("biz.HandlePing error: %s", err)
		}
	})
}
