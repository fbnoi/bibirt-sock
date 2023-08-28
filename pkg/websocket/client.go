package websocket

import (
	"log"
	"net/http"
	"sync"
	"time"

	"github.com/gorilla/websocket"
	"github.com/pkg/errors"
	"google.golang.org/protobuf/proto"
	"google.golang.org/protobuf/types/known/anypb"
)

type Color int
type SocketStatus int

const (
	Green  = Color(1)
	Blue   = Color(2)
	Yellow = Color(3)
	Red    = Color(4)
)

const (
	Init      = SocketStatus(0)
	Connected = SocketStatus(1)
	Closed    = SocketStatus(2)
)

var anyFactory = sync.Pool{
	New: func() any {
		return &anypb.Any{}
	},
}

func getAny() *anypb.Any {
	a := anyFactory.Get().(*anypb.Any)
	a.Reset()
	return a
}

func putAny(a *anypb.Any) {
	anyFactory.Put(a)
}

func NewClient(server *Server, w http.ResponseWriter, r *http.Request) *Client {
	return &Client{
		Bus:    NewBus(),
		conn:   nil,
		Color:  Red,
		server: server,
		status: Init,
		Writer: w,
		Req:    r,
	}
}

type Client struct {
	Bus
	id      string
	conn    *websocket.Conn
	session map[string]any
	status  SocketStatus
	server  *Server

	LastPingAt time.Time
	Color      Color

	sync.RWMutex

	beforeUpgradeHandleFunc func(*Client) error

	Writer http.ResponseWriter
	Req    *http.Request
}

func (c *Client) Upgrade() bool {
	var err error
	if err = c.beforeUpgradeHandleFunc(c); err != nil {
		return false
	}
	c.conn, err = c.server.upgrader.Upgrade(c.Writer, c.Req, nil)

	return err == nil
}

func (c *Client) OnUpgrade(fn func(*Client) error) {
	c.beforeUpgradeHandleFunc = fn
}

func (c *Client) Loop() {
	for {
		if c.Status() != Connected {
			return
		}
		mt, message, err := c.conn.ReadMessage()
		if err != nil {
			log.Println("read:", err)
			break
		}
		go c.onReceive(mt, message)
	}
}

func (c *Client) ID() string {
	return c.id
}

func (c *Client) Send(m proto.Message) error {
	a := getAny()
	defer putAny(a)
	err := a.MarshalFrom(m)
	if err != nil {
		return err
	}
	bs, err := proto.Marshal(a)
	if err != nil {
		return err
	}

	return c.doSend(2, bs)
}

func (c *Client) Close() error {
	c.Lock()
	defer c.Unlock()
	if c.status == Closed {
		return nil
	}
	c.status = Closed
	c.conn.WriteControl(websocket.CloseMessage, nil, time.Now().Add(time.Second))
	return c.conn.Close()
}

func (c *Client) Status() SocketStatus {
	c.RLock()
	defer c.RUnlock()
	return c.status
}

func (c *Client) Set(name string, val any) {
	c.Lock()
	defer c.Unlock()
	c.session[name] = val
}

func (c *Client) Get(name string) (val any, ok bool) {
	c.RLock()
	defer c.RUnlock()
	val, ok = c.session[name]
	return
}

func (c *Client) Delete(name string) {
	c.Lock()
	defer c.Unlock()
	delete(c.session, name)
}

func (c *Client) doSend(mt int, message []byte) error {
	c.RLock()
	defer c.RUnlock()
	if c.status != Connected {
		return errors.New("Websocket is disconnected")
	}
	return c.conn.WriteMessage(mt, message)
}

func (c *Client) onReceive(mt int, message []byte) error {
	switch mt {
	case websocket.PingMessage, websocket.PongMessage, websocket.TextMessage:
	case websocket.BinaryMessage:
		a := getAny()
		defer putAny(a)
		if err := proto.Unmarshal(message, a); err != nil {
			return err
		}
		m, err := GetMessage(a.TypeUrl)
		if err != nil {
			return err
		}
		c.Publish(m)

	case websocket.CloseMessage:
		c.Close()
	}

	return nil
}
