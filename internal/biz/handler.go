package biz

import (
	"bibirt-sock/pkg/websocket"
	"fmt"
	"sync"

	"github.com/go-kratos/kratos/v2/errors"
	"github.com/go-kratos/kratos/v2/log"
	"google.golang.org/protobuf/proto"
	"google.golang.org/protobuf/types/known/anypb"
)

const BS_MESSAGE_TYPE = 2

var anyPool = &sync.Pool{
	New: func() any {
		return &anypb.Any{}
	},
}

type Handler struct {
	clients map[string]*websocket.Client
}

func (h *Handler) RegisterHandler(client *websocket.Client) {
	client.OnMessage(h.dispatch(client))
}

func (h *Handler) SendToClient(client *websocket.Client, m proto.Message) error {
	return h.doSend(client, m)
}

func (h *Handler) SendToClientByUUID(id string, m proto.Message) error {
	if client, ok := h.clients[id]; ok {
		return h.doSend(client, m)
	}
	return errors.NotFound(fmt.Sprintf("biz.Handler.SendToClientByUUID: client %s not found", id), "client not found")
}

func (h *Handler) SendToClients(m proto.Message, ids ...string) {
	var clients []*websocket.Client
	for _, id := range ids {
		if client, ok := h.clients[id]; ok {
			clients = append(clients, client)
		}
	}
	for _, client := range clients {
		h.doSend(client, m)
	}
}

func (h *Handler) dispatch(client *websocket.Client) func(bs []byte) {
	return func(bs []byte) {
		a := getAny()
		defer putAny(a)
		if err := proto.Unmarshal(bs, a); err != nil {
			log.Error(err)
			return
		}
		switch a.TypeUrl {
		case ping_type_url:
			h.handlePing(client, a)
		}
	}
}

func (h *Handler) doSend(client *websocket.Client, m proto.Message) error {
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

	return client.Send(BS_MESSAGE_TYPE, bs)
}

func getAny() *anypb.Any {
	a := anyPool.Get().(*anypb.Any)
	a.Reset()
	return a
}

func putAny(a *anypb.Any) {
	anyPool.Put(a)
}
