package websocket

import (
	"sync"
	"time"

	"github.com/golang-jwt/jwt/v5"
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
	Connected    = SocketStatus(1)
	DisConnected = SocketStatus(2)
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

func NewClient(serve *Server, conn *websocket.Conn) *Client {
	return &Client{
		conn:      conn,
		Color:     Green,
		bus:       NewBus(),
		serve:     serve,
		status:    Connected,
		authToken: serve.GenClientToken(),
	}
}

type Client struct {
	conn    *websocket.Conn
	session map[string]any
	status  SocketStatus
	serve   *Server

	LastPingAt time.Time
	Color      Color
	bus        Bus

	authToken *jwt.Token

	sync.RWMutex
}

func (c *Client) Subscribe(m proto.Message, handleFunc func(f proto.Message)) {
	c.bus.SubscribeAsync(m, handleFunc, false)
}

func (c *Client) SubscribeOnce(m proto.Message, handleFunc func(f proto.Message)) {
	c.bus.SubscribeOnceAsync(m, handleFunc)
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
	if c.status == DisConnected {
		return nil
	}
	c.status = DisConnected
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

func (c *Client) GetToken() *jwt.Token {
	return c.authToken
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

func (c *Client) doPublish(m proto.Message) {
	c.bus.Publish(m)
	// fix me: will wait other publish call
	c.bus.WaitAsync()
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
		go func() {
			c.doPublish(m)
			PutMessage(m)
		}()

	case websocket.CloseMessage:
		c.Close()
	}

	return nil
}
