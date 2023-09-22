package biz

import (
	"flynoob/bibirt-sock/internal/conf"
	"flynoob/bibirt-sock/internal/message"
	"flynoob/bibirt-sock/pkg/websocket"
	"net/http"
	"sync"
	"time"

	"github.com/go-co-op/gocron"
	"github.com/go-kratos/kratos/v2/log"
	gws "github.com/gorilla/websocket"
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

	upgrader *gws.Upgrader
}

func NewClientHandler(c *conf.Server, authService AuthService, logger log.Logger) *ClientHandler {
	handler := &ClientHandler{
		clients:     make(map[string]*websocket.Client),
		scheduler:   gocron.NewScheduler(time.UTC),
		mux:         sync.RWMutex{},
		authService: authService,
		log:         log.NewHelper(logger),
		upgrader: &gws.Upgrader{
			HandshakeTimeout: time.Millisecond * time.Duration(c.Websocket.HandshakeTimeout),
			ReadBufferSize:   int(c.Websocket.ReadBufferSize),
			WriteBufferSize:  int(c.Websocket.WriteBufferSize),
			CheckOrigin: func(r *http.Request) bool {
				return true
			},
		},
	}
	handler.scheduler.TagsUnique()
	handler.scheduler.StartAsync()
	websocket.RegisterMessage(&message.Ping{})
	websocket.RegisterMessage(&message.Connected{})
	websocket.RegisterMessage(&message.Disconnected{})
	handler.registerMessage()
	return handler
}

func (handler *ClientHandler) HandleClient(c *websocket.Client) {
	if err := handler.auth(c); err != nil {
		handler.log.Errorf("auth error: %s", err)
		return
	}
	if err := c.Upgrade(handler.upgrader); err != nil {
		handler.log.Errorf("upgrade handler error: %s", err)
		return
	}
	handler.monitorHealth(c)
	handler.registerMessageHandler(c)
	handler.handleClientConnected(c)

	defer handler.closeClient(c)
	c.Listen()
}

func (handler *ClientHandler) closeClient(c *websocket.Client) {
	c.Close()
	handler.mux.Lock()
	defer handler.mux.Unlock()
	c.Publish(&message.Disconnected{Id: c.ID()})
	delete(handler.clients, c.ID())
}

func (handler *ClientHandler) handleClientConnected(c *websocket.Client) {
	handler.mux.Lock()
	defer handler.mux.Unlock()
	if dc, ok := handler.clients[c.ID()]; ok {
		dc.Send(&message.ConnectionDuplicated{})
		handler.closeClient(dc)
	}

	handler.clients[c.ID()] = c
	c.Publish(&message.Connected{Id: c.ID()})
}
