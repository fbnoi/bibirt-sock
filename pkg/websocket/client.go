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

func NewClient(w http.ResponseWriter, r *http.Request) *Client {
	return &Client{
		Bus:     NewBus(),
		conn:    nil,
		Color:   Red,
		status:  Init,
		Writer:  w,
		Req:     r,
		session: make(map[string]any),
	}
}

type Client struct {
	Bus
	sync.RWMutex

	id      string
	conn    *websocket.Conn
	session map[string]any
	status  SocketStatus

	LastPingAt time.Time
	Color      Color

	Writer http.ResponseWriter
	Req    *http.Request
}

func (c *Client) Upgrade(upgrader *websocket.Upgrader) (err error) {
	c.conn, err = upgrader.Upgrade(c.Writer, c.Req, nil)
	if err == nil {
		c.status = Connected
		c.Color = Green
	}

	return
}

func (c *Client) Listen() {
	for {
		if c.Status() != Connected {
			return
		}
		mt, message, err := c.conn.ReadMessage()
		if err != nil {
			log.Println("read:", err)
			break
		}
		go func() {
			if err := c.onReceive(mt, message); err != nil {
				log.Println("onReceive:", err)
			}
		}()
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
		if err := a.UnmarshalTo(m); err != nil {
			log.Printf("onReceive: %s\n", err)
			return err
		}
		c.Publish(m)

	case websocket.CloseMessage:
		c.Close()
	}

	return nil
}
