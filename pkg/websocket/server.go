package websocket

import (
	"log"
	"net/http"

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
		clients: make(map[string]*Client),
	}
}

type Server struct {
	upgrader websocket.Upgrader

	clients map[string]*Client
}

func (srv *Server) Handler(fn func(*Client)) func(http.ResponseWriter, *http.Request) {
	return func(w http.ResponseWriter, r *http.Request) {
		c, err := srv.upgrader.Upgrade(w, r, nil)
		if err != nil {
			log.Print("upgrade:", err)
			return
		}
		uuid4 := uuid.Must(uuid.NewV4())
		client := &Client{id: uuid4.String(), srv: srv, conn: c}
		srv.clients[uuid4.String()] = client
		defer client.Close()
		fn(client)
		for {
			mt, message, err := c.ReadMessage()
			if err != nil {
				client.receiveErrorHandler(err)
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
