package websocket

import (
	"flag"
	"log"
	"net/http"
	"time"

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
		endpoint:    endpoint,
		addr:        addr,
		bus:         NewBus(),
		handleError: func(c *Client, b []byte, err error) {},
	}
}

type Server struct {
	bus      Bus
	addr     string
	upgrader websocket.Upgrader
	endpoint string

	handleError         func(*Client, []byte, error)
	newClientHandleFunc func(*Client)
}

func (s *Server) Endpoint() string {
	return s.endpoint
}

func (s *Server) Run() error {
	http.HandleFunc("/", func(w http.ResponseWriter, r *http.Request) {
		defer recovery()
		client := NewClient(s, w, r)
		s.newClientHandleFunc(client)
	})
	return http.ListenAndServe(s.addr, nil)
}

func (s *Server) OnNewConnection(fn func(*Client)) {
	s.newClientHandleFunc = fn
}

func (s *Server) Listen(addr string) error {
	s.addr = addr
	return s.Run()
}

func recovery() {
	if message := recover(); message != nil {
		log.Println(errors.Errorf("Websocket error: %v", message))
	}
}
