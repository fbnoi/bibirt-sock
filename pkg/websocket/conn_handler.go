package websocket

import (
	"net/http"

	gs "github.com/gorilla/websocket"
)

type ClientHandler interface {
	HandleClient(c *Client)
}

type ConnHandler struct {
	clientHandler ClientHandler
	upgrader      gs.Upgrader
}

func NewConnHandler(handler ClientHandler) http.Handler {
	cHandler := &ConnHandler{clientHandler: handler}
	return cHandler
}

func (c *ConnHandler) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	client := NewClient(w, r)
	c.clientHandler.HandleClient(client)
}
