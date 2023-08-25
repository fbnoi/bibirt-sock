package biz

import (
	"bibirt-sock/internal/message"
	"bibirt-sock/pkg/websocket"
	"time"

	"github.com/go-kratos/kratos/v2/log"
	"google.golang.org/protobuf/proto"
	timestamppb "google.golang.org/protobuf/types/known/timestamppb"
)

var ping_type_url = message.TypeUrl(&message.Ping{})

func (h *Handler) handlePing(client *websocket.Client, m proto.Message) {
	ping, ok := m.(*message.Ping)
	if !ok {
		return
	}
	client.LastPingAt = time.Now()
	ping.DownTimestamp = timestamppb.New(client.LastPingAt)
	if err := h.SendToClient(client, ping); err != nil {
		log.Errorf("biz.handlePing error: %s", err)
	}
}
