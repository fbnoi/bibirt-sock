package biz

import (
	"flynoob/bibirt-sock/internal/message"
	"flynoob/bibirt-sock/pkg/websocket"
	"log"
	"time"

	"google.golang.org/protobuf/proto"
	"google.golang.org/protobuf/types/known/timestamppb"
)

func (useCase *ConnUseCase) MonitorHealth(client *websocket.Client) {
	useCase.handlePing(client)
	useCase.scheduler.Every(1).Second().Tag(client.ID()).Do(func() {
		if client.LastPingAt.Add(2 * time.Second).Before(time.Now()) {
			client.Color += 1
		}
		if client.Color == websocket.Red {
			client.Close()
			useCase.scheduler.RemoveByTag(client.ID())
		}
	})
}

func (*ConnUseCase) handlePing(client *websocket.Client) {
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
			log.Printf("biz.HandlePing error: %s", err)
		}
	})
}
