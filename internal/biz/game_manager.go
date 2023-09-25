package biz

import (
	"flynoob/bibirt-sock/pkg/websocket"
	"math"
	"sync"
)

const GAP = 50

type GamemManager struct {
	list *MatchList

	rooms map[string]*Room
}

func (gm *GamemManager) StartGame(c *websocket.Client) {
	if node := gm.list.FindAndRermoveOrJoin(c); node != nil {
		room := NewRoom(node)
		gm.rooms[room.Name] = room
		room.bootstrap()
	}
}

type MatchList struct {
	header *MatchListNode
	end    *MatchListNode

	mux sync.Mutex
}

func (m *MatchList) FindAndRermoveOrJoin(c *websocket.Client) *MatchListNode {
	m.mux.Lock()
	defer m.mux.Unlock()
	score, _ := c.Get("score")
	current := m.header
	for current != nil {
		if current.isMatch(score.(int)) {
			m.remove(current)
			return current
		}
	}
	m.append(NewMatchListNode(c))
	return nil
}

func (m *MatchList) append(node *MatchListNode) {
	if m.header == nil {
		m.header = node
		m.end = node
	} else {
		m.end.next = node
		node.prev = m.end
		m.end = node
	}
}

func (m *MatchList) remove(node *MatchListNode) {
	if m.header == node {
		m.header = node.next
	}
	if m.end == node {
		m.end = node.prev
	}

	if node.prev != nil {
		node.prev.next = node.next
	}
	if node.next != nil {
		node.next.prev = node.prev
	}
	node.prev = nil
	node.next = nil
}

type MatchListNode struct {
	CP    *MatchCP
	Score int
	next  *MatchListNode
	prev  *MatchListNode
}

func NewMatchListNode(c *websocket.Client) *MatchListNode {
	return &MatchListNode{
		CP: &MatchCP{},
	}
}

func (m *MatchListNode) isMatch(score int) bool {
	return math.Abs(float64(score-m.Score)) < 50
}

type MatchCP struct {
	Client1 *websocket.Client
	Client2 *websocket.Client
}
