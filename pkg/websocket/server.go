package websocket

import (
	"net/http"
	"time"

	"github.com/go-co-op/gocron"
	"github.com/gofrs/uuid"
	"github.com/gorilla/websocket"
	"github.com/pkg/errors"
)

var DefaultServer = &Server{
	upgrader: websocket.Upgrader{
		HandshakeTimeout:  default_config.HandshakeTimeout,
		ReadBufferSize:    default_config.ReadBufferSize,
		WriteBufferSize:   default_config.WriteBufferSize,
		EnableCompression: default_config.EnableCompression,
	},
	clients: make(map[string]*Client),
}

func NewServer(conf *Config) *Server {
	return &Server{
		upgrader: websocket.Upgrader{
			HandshakeTimeout:  conf.HandshakeTimeout,
			ReadBufferSize:    conf.ReadBufferSize,
			WriteBufferSize:   conf.WriteBufferSize,
			EnableCompression: conf.EnableCompression,
		},
		clients:   make(map[string]*Client),
		scheduler: gocron.NewScheduler(time.UTC),
	}
}

type Server struct {
	upgrader websocket.Upgrader

	clients   map[string]*Client
	scheduler *gocron.Scheduler
}

func (srv *Server) Handler(fn func(*Client)) func(http.ResponseWriter, *http.Request) {
	return func(w http.ResponseWriter, r *http.Request) {
		c, err := srv.upgrader.Upgrade(w, r, nil)
		client := srv.newClient(c)
		fn(client)
		if err != nil {
			client.handleError(err)
			delete(srv.clients, client.id)
			return
		}
		defer client.Close()
		client.connectedHandler()
		srv.monitorHealth(client)
		for {
			mt, message, err := c.ReadMessage()
			if err != nil {
				client.handleError(err)
			}
			switch mt {
			case websocket.PingMessage, websocket.PongMessage:
			case websocket.TextMessage, websocket.BinaryMessage:
				client.Receive(message)
			case websocket.CloseMessage:
				return
			}
		}
	}
}

func (srv *Server) SendToClient(id string, mt int, message []byte) error {
	if client, ok := srv.clients[id]; ok {
		return client.Send(mt, message)
	}

	return errors.Errorf("websocket.SendToClient: client %s not found", id)
}

func (srv *Server) monitorHealth(client *Client) {
	srv.scheduler.Every(1).Second().Do(func() {
		if client.LastPingAt.Add(2 * time.Second).Before(time.Now()) {
			client.Color += 1
		}

		if client.Color == Red {
			client.Close()
		}
	})
}

func (srv *Server) newClient(c *websocket.Conn) *Client {
	uuid4 := uuid.Must(uuid.NewV4())
	client := &Client{
		id:             uuid4.String(),
		srv:            srv,
		conn:           c,
		Color:          Green,
		messageHandler: func(bs []byte) {},
		closingHandler: func() {},
		closedHandler:  func() {},
		handleError:    func(error) {},
	}
	srv.clients[uuid4.String()] = client
	return client
}
