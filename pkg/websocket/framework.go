package websocket

import (
	"sync"
	"time"

	"bibirt-sock/pkg/websocket/pb"

	"github.com/go-co-op/gocron"
	"google.golang.org/protobuf/proto"
)

func WrapHandler(c *Client, fn func(*Client, proto.Message)) HandleFunc {
	return func(m proto.Message) { fn(c, m) }
}

var (
	healthyScanScheduler = gocron.NewScheduler(time.UTC)
	schedulerOnce        = sync.Once{}
)

func startSchedule() {
	schedulerOnce.Do(func() {
		healthyScanScheduler.StartAsync()
	})
}

func monitorHealth(client *Client) {
	startSchedule()
	client.Subscribe(&pb.Heartbeat{}, WrapHandler(client, func(c *Client, m proto.Message) {
		heartbeat, ok := m.(*pb.Heartbeat)
		if !ok {
			return
		}
		client.LastPingAt = time.Now()
		heartbeat.DownTimestamp = client.LastPingAt.UnixMilli()
		client.Send(heartbeat)
	}))
	healthyScanScheduler.Every(1).Second().Do(func() {
		checkHealth(client)
	})
}

func checkHealth(client *Client) {
	if client.LastPingAt.Add(2 * time.Second).Before(time.Now()) {
		client.Color += 1
	}

	if client.Color == Red {
		client.Close()
	}
}
