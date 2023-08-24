package server

import (
	"bibirt-sock/internal/conf"
	"bibirt-sock/pkg/websocket"

	"github.com/go-kratos/kratos/v2/log"
	"github.com/go-kratos/kratos/v2/middleware/logging"
	"github.com/go-kratos/kratos/v2/middleware/recovery"
	"github.com/go-kratos/kratos/v2/middleware/tracing"
	"github.com/go-kratos/kratos/v2/middleware/validate"
	"github.com/go-kratos/kratos/v2/transport/http"
	"github.com/gorilla/mux"
)

// NewHTTPServer new a HTTP server.
func NewHTTPServer(c *conf.Server, logger log.Logger) *http.Server {
	router := mux.NewRouter()
	handleFunc = websocket.DefaultServer.Handler(handler.WsHandler)
	router.HandleFunc("/ws", handleFunc)
	opts := []http.ServerOption{
		http.Middleware(
			recovery.Recovery(),
			tracing.Server(),
			logging.Server(logger),
			validate.Validator(),
		),
		http.Address(":8080"),
	}

	httpSrv := http.NewServer(opts...)
	httpSrv.HandlePrefix("/", router)

	return httpSrv
}
