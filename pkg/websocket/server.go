package websocket

import (
	"flag"
	"fmt"
	"log"
	"net/http"
	"time"

	"github.com/golang-jwt/jwt/v5"
	"github.com/gorilla/websocket"
	"github.com/pkg/errors"
)

var (
	endpoint           string
	addr               string
	handleShakeTimeout int64
	readBufferSize     int
	writeBufferSize    int
	enableCompression  bool
)

func addFlag() {
	flag.StringVar(&endpoint, "endpoint", "HZ-WH-1", "set endpoint name for server")
	flag.StringVar(&addr, "addr", ":5678", "set addr for server")
	flag.Int64Var(&handleShakeTimeout, "timeout", 3000, "set connection timeout")
	flag.IntVar(&readBufferSize, "read_buffer_size", 0, "set read buffer size")
	flag.IntVar(&writeBufferSize, "write_buffer_size", 0, "set write buffer size")
	flag.BoolVar(&enableCompression, "enable_compression", false, "enable message compression")
}

func NewServer() *Server {
	addFlag()
	return &Server{
		upgrader: websocket.Upgrader{
			HandshakeTimeout:  time.Millisecond * time.Duration(handleShakeTimeout),
			ReadBufferSize:    readBufferSize,
			WriteBufferSize:   writeBufferSize,
			EnableCompression: enableCompression,
			CheckOrigin: func(*http.Request) bool {
				return true
			},
		},
		endpoint:      endpoint,
		addr:          addr,
		bus:           NewBus(),
		beforeUpgrade: func(w http.ResponseWriter, r *http.Request) bool { return true },
		afterUpgrade:  func(*Client) {},
		handleError:   func(c *Client, b []byte, err error) {},
	}
}

type Server struct {
	bus      Bus
	addr     string
	upgrader websocket.Upgrader
	endpoint string

	beforeUpgrade func(w http.ResponseWriter, r *http.Request) bool
	afterUpgrade  func(*Client)
	handleError   func(*Client, []byte, error)
}

func (s *Server) Endpoint() string {
	return s.endpoint
}

func (s *Server) GenClientToken() *jwt.Token {
	return jwt.NewWithClaims(jwt.SigningMethodES256, jwt.MapClaims{
		"iss": s.endpoint,
		"sub": "client_auth",
		"nbf": time.Now().UnixMilli(),
		"iat": time.Now().UnixMilli(),
	})
}

func (s *Server) Run() error {
	startSchedule()
	http.HandleFunc("/", func(w http.ResponseWriter, r *http.Request) {
		if !s.beforeUpgrade(w, r) {
			return
		}
		c, err := s.upgrader.Upgrade(w, r, nil)
		if err != nil {
			w.Write([]byte(fmt.Sprintf("upgrade error: %s", err)))
			return
		}
		defer c.Close()

		client := NewClient(s, c)
		defer s.recovery(client)
		s.afterUpgrade(client)
		monitorHealth(client)
		depositMessage(client)
	})
	return http.ListenAndServe(s.addr, nil)
}

func (s *Server) Listen(addr string) error {
	s.addr = addr
	return s.Run()
}

func (s *Server) SetBeforeUpgradeHandler(handler func(w http.ResponseWriter, r *http.Request) bool) {
	s.beforeUpgrade = handler
}

func (s *Server) SetAfterUpgradeHandler(handler func(*Client)) {
	s.afterUpgrade = handler
}

func (s *Server) recovery(client *Client) {
	if message := recover(); message != nil {
		s.handleError(client, nil, errors.Errorf("Websocket error: %v", message))
	}
}

func depositMessage(client *Client) {
	for {
		if client.Status() != Connected {
			return
		}
		mt, message, err := client.conn.ReadMessage()
		if err != nil {
			log.Println("read:", err)
			break
		}
		client.onReceive(mt, message)
	}
}
