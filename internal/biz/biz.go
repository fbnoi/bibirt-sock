package biz

import (
	"flynoob/bibirt-sock/internal/message"
	"flynoob/bibirt-sock/pkg/websocket"
	"time"

	"github.com/go-co-op/gocron"
)

var (
	scheduler = gocron.NewScheduler(time.UTC)
)

func Bootstrap() {
	registerMessage()
	startScheduler()
}

func registerMessage() {
	ping := &message.Ping{}
	websocket.RegisterMessage(ping)
}

func startScheduler() {
	scheduler.TagsUnique()
	scheduler.StartAsync()
}
