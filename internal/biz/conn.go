package biz

import (
	"flynoob/bibirt-sock/internal/message"
	"flynoob/bibirt-sock/pkg/websocket"
	"sync"
	"time"

	"github.com/go-co-op/gocron"
)

type AuthServiceInterface interface {
	ConnUUID(tokStr string) (string, error)
}

type ConnUseCase struct {
	clients   map[string]*websocket.Client
	scheduler *gocron.Scheduler

	authService AuthServiceInterface

	mux sync.RWMutex
}

func NewConnUseCase(authService AuthServiceInterface) *ConnUseCase {
	useCase := &ConnUseCase{
		clients:     make(map[string]*websocket.Client),
		scheduler:   gocron.NewScheduler(time.UTC),
		mux:         sync.RWMutex{},
		authService: authService,
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
		useCase.mux.Lock()
		defer useCase.mux.Unlock()

		c.Publish(&message.Disconnected{Id: c.ID()})
		delete(useCase.clients, c.ID())
	})
}

func (useCase *ConnUseCase) handleClientConnected(c *websocket.Client) {
	useCase.mux.Lock()
	defer useCase.mux.Unlock()

	useCase.clients[c.ID()] = c
	c.Publish(&message.Connected{Id: c.ID()})
}
