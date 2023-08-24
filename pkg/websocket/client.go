package websocket

import (
	"time"

	"github.com/gorilla/websocket"
)

type Color int

const (
	Green  = Color(1)
	Blue   = Color(2)
	Yellow = Color(3)
	Red    = Color(4)
)

type Client struct {
	id   string
	srv  *Server
	conn *websocket.Conn

	LastPingAt time.Time
	Color      Color

	connectedHandler func()
	messageHandler   func(bs []byte)
	closingHandler   func()
	closedHandler    func()
	handleError      func(error)
}

func (c *Client) Send(mt int, message []byte) error {
	return c.conn.WriteMessage(mt, message)
}

func (c *Client) SendToClient(id string, mt int, message []byte) error {
	return c.srv.SendToClient(id, mt, message)
}

func (c *Client) ID() string {
	return c.id
}

func (c *Client) OnMessage(fn func([]byte)) {
	c.messageHandler = fn
}

func (c *Client) Receive(bs []byte) {
	c.messageHandler(bs)
}

func (c *Client) OnConnected(fn func()) {
	c.connectedHandler = fn
}

func (c *Client) OnClosing(fn func()) {
	c.closingHandler = fn
}

func (c *Client) OnClosed(fn func()) {
	c.closedHandler = fn
}

func (c *Client) OnError(fn func(error)) {
	c.handleError = fn
}

func (c *Client) Close() {
	c.closingHandler()
	c.conn.Close()
	delete(c.srv.clients, c.id)
	c.closedHandler()
}
