package biz

import (
	"flynoob/bibirt-sock/internal/message"
	"flynoob/bibirt-sock/pkg/websocket"
	"fmt"
	"sync"
	"time"

	"google.golang.org/protobuf/proto"
)

const (
	MAX_PIECE_COUNT = 255

	PIECE_COLOR_BLACK = 1
	PIECE_COLOR_WHITE = 1
)

// []int{ x, y, color}
type Piece struct {
	x, y, color int
}

type Player struct {
	Client *websocket.Client
	Uuid   string
	Name   string
}

type Room struct {
	Name        string
	BlackPlayer *Player
	WhitePlayer *Player
	playerMux   sync.Mutex

	current *Player

	Watchers []*websocket.Client

	confirm  chan int
	finished chan int

	// how mach player comfirmed
	counter int

	// if game is suspended
	// suspend bool

	boardMux  sync.RWMutex
	board     map[string]*Piece
	lastPiece *Piece
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
		}
	}
}

func (r *Room) switchPlayer() {
	r.playerMux.Lock()
	defer r.playerMux.Unlock()
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
	r.onClose(player)
	r.onConfirm(player)
	r.onDrop(player)
}

func (r *Room) broadcast(message proto.Message) {
	r.BlackPlayer.Client.Send(message)
	r.WhitePlayer.Client.Send(message)
}

func (r *Room) onConfirm(player *Player) {
	player.Client.Subscribe(&message.ConfirmedRequest{}, func(m proto.Message) {
		r.playerMux.Lock()
		defer r.playerMux.Unlock()
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
		r.playerMux.Lock()
		defer r.playerMux.Unlock()
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
			color := PIECE_COLOR_WHITE
			if player == r.BlackPlayer {
				color = PIECE_COLOR_BLACK
			}
			if r.dropPiece(color, int(dm.X), int(dm.Y)) {
				if !r.judgeBoardFinished() {
					r.switchPlayer()
				} else {
				}
			}
		}
	})
}

func (r *Room) dropPiece(color, x, y int) (success bool) {
	r.boardMux.Lock()
	defer r.boardMux.Unlock()
	msg := &message.DropReply{Uuid: r.current.Uuid}
	if (color == PIECE_COLOR_WHITE && r.current == r.WhitePlayer) ||
		(color == PIECE_COLOR_BLACK && r.current == r.BlackPlayer) {
		key := fmt.Sprintf("%d-%d-%d", x, y, color)
		if _, ok := r.board[key]; !ok {
			r.board[key] = &Piece{x, y, color}
			r.lastPiece = r.board[key]
			x64, y64 := int64(x), int64(y)
			msg.X, msg.Y, msg.Success = &x64, &y64, 1
			success = true
		} else {
			m := "con't drop here"
			msg.Success, msg.Message = 0, &m
			success = false
		}
	} else {
		m := "wrong color"
		msg.Success, msg.Message = 0, &m
		success = false
	}
	r.broadcast(msg)
	return
}

func (r *Room) judgeBoardFinished() bool {
	var x, y, score int
	var key string
	if r.lastPiece != nil {
		score = 0
		x, y = inclineSectionU(r.lastPiece.x, r.lastPiece.y)
		for i := 0; i < 15; i++ {
			key = fmt.Sprintf("%d-%d-%d", x, y, r.lastPiece.color)
			if _, ok := r.board[key]; ok {
				score++
				if score >= 5 {
					return true
				}
			} else {
				score = 0
			}
			x++
			y++
		}

		x, y = inclineSectionD(r.lastPiece.x, r.lastPiece.y)
		for x < 15 && y < 15 {
			key = fmt.Sprintf("%d-%d-%d", x, y, r.lastPiece.color)
			if _, ok := r.board[key]; ok {
				score++
				if score >= 5 {
					return true
				}
			} else {
				score = 0
			}
			x++
			y--
		}
	}

	return false
}

func inclineSectionU(x, y int) (int, int) {
	for x > 0 && y > 0 {
		x--
		y--
	}

	return x, y
}

func inclineSectionD(x, y int) (int, int) {
	for x > 0 && y > 0 {
		x--
		y++
	}

	return x, y
}
