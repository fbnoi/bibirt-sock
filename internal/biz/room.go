package biz

import (
	"flynoob/bibirt-sock/pkg/websocket"
	"time"

	"google.golang.org/protobuf/proto"
)

type Player struct {
	Client *websocket.Client
	Uuid   string
	Name   string
}

type Room struct {
	Name        string
	BlackPlayer *Player
	WhitePlayer *Player

	Watchers     []*websocket.Client
	palyerChan   chan proto.Message
	proccessChan chan int
}

func NewRoom(node *MatchListNode) *Room {
	return &Room{}
}

func NewPlayer(client *websocket.Client) *Player {
	return &Player{}
}

func (r *Room) start() {
	select {
	case <-time.After(10 * time.Second):
		r.timeout()
		break
	case message := <-r.palyerChan:
		r.broadcast(message)
		// case c
	}

	r.clear()
}

func (r *Room) gameProccess() {
	tk := time.NewTicker(10 * time.Second)
	for {
		select {
		case <-tk.C:
		}
	}
}

func (r *Room) timeout() {

}

func (r *Room) broadcast(proto.Message) {

}

func (r *Room) suspend()

func (r *Room) clear() {

}

func (p *Player) confirm() error {
	return nil
}
