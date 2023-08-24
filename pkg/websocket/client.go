package websocket

import (
	"github.com/gorilla/websocket"
)

type Client struct {
	id   string
	srv  *Server
	conn *websocket.Conn

	connectedHandler    func()
	messageHandler      func(bs []byte)
	closingHandler      func()
	closedHandler       func()
	receiveErrorHandler func(error)
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

func (c *Client) OnReceiveError(fn func(error)) {
	c.receiveErrorHandler = fn
}

func (c *Client) Close() {
	c.closingHandler()
	c.conn.Close()
	c.closedHandler()
}
