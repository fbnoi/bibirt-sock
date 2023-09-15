package server

import (
	"flynoob/bibirt-sock/internal/biz"
	"flynoob/bibirt-sock/pkg/websocket"

	"github.com/go-kratos/kratos/v2/transport/http"
)

func NewServer(useCase *biz.ClientHandler) *http.Server {
	httpSrv := http.NewServer(http.Address(":8080"))
	httpSrv.Handle("/", websocket.NewConnHandler(useCase))
	return httpSrv
}
