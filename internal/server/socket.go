package server

import (
	"flynoob/bibirt-sock/internal/biz"
	"flynoob/bibirt-sock/internal/conf"
	"flynoob/bibirt-sock/pkg/websocket"

	"github.com/go-kratos/kratos/v2/transport/http"
)

func NewServer(c *conf.Server, useCase *biz.ClientHandler) *http.Server {
	httpSrv := http.NewServer(http.Address(c.Addr))
	httpSrv.Handle("/", websocket.NewConnHandler(useCase))
	return httpSrv
}
