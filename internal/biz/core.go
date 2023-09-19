package biz

import (
	"flynoob/bibirt-sock/api"
	"flynoob/bibirt-sock/internal/message"
	"flynoob/bibirt-sock/pkg/websocket"
	"net/http"
	"time"

	"google.golang.org/protobuf/proto"
)

func (handler *ClientHandler) auth(client *websocket.Client) error {
	tok := client.Req.URL.Query().Get("token")
	uuid, err := handler.authService.ConnUUID(tok)
	if err != nil {
		if api.IsTokenInvalid(err) {
			client.Writer.WriteHeader(http.StatusBadRequest)
		} else {
			client.Writer.WriteHeader(http.StatusInternalServerError)
		}
		return err
	}
	client.Set("uuid", uuid)
	return nil
}

func (handler *ClientHandler) monitorHealth(client *websocket.Client) {
	handler.handlePing(client)
	handler.scheduler.Every(1).Second().Tag(client.ID()).Do(func() {
		if client.LastPingAt.Add(10 * time.Second).Before(time.Now()) {
			client.Color += 1
		}
		if client.Color >= websocket.Red {
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
		ping.DownTimestamp = client.LastPingAt.UnixMilli()
		if err := client.Send(ping); err != nil {
			handler.log.Errorf("biz.HandlePing error: %s", err)
		}
	})
}
