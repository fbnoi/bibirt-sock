package biz

import (
	"flynoob/bibirt-sock/internal/message"
	"flynoob/bibirt-sock/pkg/websocket"
	"sync"
	"time"

	"github.com/go-co-op/gocron"
	"github.com/go-kratos/kratos/v2/log"
)

type AuthService interface {
	ConnUUID(string) (string, error)
}

type ClientHandler struct {
	clients   map[string]*websocket.Client
	scheduler *gocron.Scheduler

	authService AuthService

	mux sync.RWMutex

	log *log.Helper
}

func NewClientHandler(authService AuthService, logger log.Logger) *ClientHandler {
	handler := &ClientHandler{
		clients:     make(map[string]*websocket.Client),
		scheduler:   gocron.NewScheduler(time.UTC),
		mux:         sync.RWMutex{},
		authService: authService,
		log:         log.NewHelper(logger),
	}
	handler.scheduler.TagsUnique()
	handler.scheduler.StartAsync()
	handler.registerMessage()

	return handler
}

func (handler *ClientHandler) HandleClient(c *websocket.Client) {
	handler.Auth(c)
	if c.Upgrade() {
		handler.handleClientConnected(c)
		handler.handleClientDisconnected(c)
		handler.registerMessageHandler(c)

		defer c.Close()
		c.Loop()
	}
}

func (handler *ClientHandler) handleClientDisconnected(c *websocket.Client) {
	c.OnClose(func(c *websocket.Client) {
		handler.mux.Lock()
		defer handler.mux.Unlock()

		c.Publish(&message.Disconnected{Id: c.ID()})
		delete(handler.clients, c.ID())
	})
}

func (handler *ClientHandler) handleClientConnected(c *websocket.Client) {
	handler.mux.Lock()
	defer handler.mux.Unlock()

	handler.clients[c.ID()] = c
	c.Publish(&message.Connected{Id: c.ID()})
}
