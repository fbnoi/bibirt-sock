package biz

import (
	"flynoob/bibirt-sock/internal/message"
	"flynoob/bibirt-sock/pkg/websocket"
	"sync"
	"time"

	"github.com/go-co-op/gocron"
)

type ConnUseCase struct {
	clients   map[string]*websocket.Client
	scheduler *gocron.Scheduler

	mux sync.Mutex
}

func NewConnUseCase() *ConnUseCase {
	useCase := &ConnUseCase{
		clients:   make(map[string]*websocket.Client),
		scheduler: gocron.NewScheduler(time.UTC),
		mux:       sync.Mutex{},
	}
	useCase.scheduler.TagsUnique()
	useCase.scheduler.StartAsync()
	useCase.registerMessage()

	return useCase
}

func (useCase *ConnUseCase) HandleClient(c *websocket.Client) {
	useCase.Auth(c)
	if c.Upgrade() {
		useCase.handleClientConnected(c)
		useCase.handleClientDisconnected(c)
		useCase.registerMessageHandler(c)

		defer c.Close()
		c.Loop()
	}
}

func (useCase *ConnUseCase) handleClientDisconnected(c *websocket.Client) {
	c.OnClose(func(c *websocket.Client) {
		c.Publish(&message.Disconnected{Id: c.ID()})
		delete(useCase.clients, c.ID())
	})
}

func (useCase *ConnUseCase) handleClientConnected(c *websocket.Client) {
	useCase.clients[c.ID()] = c
	c.Publish(&message.Connected{Id: c.ID()})
}
