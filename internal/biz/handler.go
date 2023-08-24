package biz

import (
	"bibirt-sock/pkg/websocket"
	"sync"
	"time"

	"github.com/go-co-op/gocron"
)

var (
	healthyScanScheduler = gocron.NewScheduler(time.UTC)
	schedulerOnce        = sync.Once{}
)

type Handler struct {
}

func (h *Handler) Handle(client *websocket.Client) {
	// client.OnClosed()
	// client.OnClosing()
	client.OnConnected(func() { (monitorHealth(client)) })
	// client.OnError()
	client.OnMessage(func(b []byte) {})
}

func monitorHealth(client *websocket.Client) {
	healthyScanScheduler.Every(1).Second().Do(func() {
		checkHealth(client)
	})
}

func checkHealth(client *websocket.Client) {
	if client.LastPingAt.Add(2 * time.Second).Before(time.Now()) {
		client.Color += 1
	}

	if client.Color == Red {
		client.Close()
	}
}
