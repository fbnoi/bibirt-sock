package server

import (
	"flynoob/bibirt-sock/internal/biz"
	"flynoob/bibirt-sock/pkg/websocket"
)

func NewServer(useCase *biz.ConnUseCase) *websocket.Server {
	srv := websocket.NewServer()
	srv.OnNewConnection(useCase.HandleClient)
	return srv
}
