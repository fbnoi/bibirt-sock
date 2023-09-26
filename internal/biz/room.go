package biz

import (
	"flynoob/bibirt-sock/internal/message"
	"flynoob/bibirt-sock/pkg/websocket"
	"sync"
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

	current *Player

	Watchers []*websocket.Client

	paly     chan int
	confirm  chan int
	finished chan int

	// how mach player comfirmed
	counter int

	// if game is suspended
	suspend bool

	mux sync.Mutex
}

func NewRoom(node *MatchListNode) *Room {
	return &Room{}
}

func NewPlayer(client *websocket.Client) *Player {
	return &Player{}
}

// init: add player client subscribers
func (r *Room) init() {
	r.addSubscriber(r.WhitePlayer)
	r.addSubscriber(r.BlackPlayer)
}

// start send matched message and wait for confirming message
// from player until timeout
func (r *Room) start() {
	select {
	case <-time.After(10 * time.Second):
		r.timeout()
	case <-r.finished:
		r.end()
	case <-r.confirm:
		r.play()
	}
}

// play start to play the game
func (r *Room) play() {
	timer := time.NewTimer(10 * time.Second)
	for {
		if !timer.Stop() {
			<-timer.C
		}
		timer.Reset(10 * time.Second)
		select {
		case <-timer.C:
			r.switchPlayer()
		case <-r.finished:
			r.end()
			break
		case <-r.paly:
			continue
		}
	}
}

func (r *Room) switchPlayer() {
	r.mux.Lock()
	defer r.mux.Unlock()
	if r.current == r.WhitePlayer {
		r.current = r.BlackPlayer
	} else {
		r.current = r.WhitePlayer
	}
}

func (r *Room) timeout() {
	r.finished <- 1
}

func (r *Room) end() {
	r.broadcast(&message.GameEnd{})
}

func (r *Room) addSubscriber(player *Player) {
	r.onConfirm(player)
}

func (r *Room) broadcast(mesage proto.Message) {
	r.BlackPlayer.Client.Send(mesage)
	r.WhitePlayer.Client.Send(mesage)
}

func (r *Room) checkBoardFinished() bool {
	return false
}

func (r *Room) onConfirm(player *Player) {
	player.Client.Subscribe(&message.ConfirmedRequest{}, func(m proto.Message) {
		r.mux.Lock()
		defer r.mux.Unlock()
		r.counter++
		r.broadcast(&message.ConfirmedReply{
			Uuid: m.(*message.ConfirmedRequest).Uuid,
		})
		if r.counter >= 2 {
			r.broadcast(&message.GameStart{})
		}
	})
}

func (r *Room) onClose(player *Player) {
	player.Client.Subscribe(&message.Disconnected{}, func(m proto.Message) {
		r.mux.Lock()
		defer r.mux.Unlock()
		r.counter--
		r.broadcast(&message.Disconnected{
			Uuid: m.(*message.Disconnected).Uuid,
		})

		if r.counter <= 0 {
			r.finished <- 1
		}
	})
}

func (r *Room) onDrop(player *Player) {
	player.Client.Subscribe(&message.Disconnected{}, func(m proto.Message) {
		dm := m.(*message.DropRequest)
		if dm.Uuid == player.Uuid {
			r.broadcast(&message.DropReply{
				Uuid: dm.Uuid,
				X:    dm.X,
				Y:    dm.Y,
			})
			if !r.checkBoardFinished() {
				r.switchPlayer()
			}
		}
	})
}
