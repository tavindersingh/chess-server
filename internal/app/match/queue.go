package match

import (
	"sync"

	"github.com/gofiber/contrib/websocket"
)

type MatchQueue struct {
	Player []string
	Conns  map[string]*websocket.Conn
	Mutex  sync.Mutex
}

func NewMatchQueue() *MatchQueue {
	return &MatchQueue{
		Player: []string{},
		Conns:  make(map[string]*websocket.Conn),
		Mutex:  sync.Mutex{},
	}
}

func (mq *MatchQueue) Enqueue(playerId string, conn *websocket.Conn) {
	mq.Mutex.Lock()
	defer mq.Mutex.Unlock()

	mq.Player = append(mq.Player, playerId)
	mq.Conns[playerId] = conn
}

func (mq *MatchQueue) Dequeue() (string, *websocket.Conn) {
	mq.Mutex.Lock()
	defer mq.Mutex.Unlock()

	if len(mq.Player) == 0 {
		return "", nil
	}

	playerId := mq.Player[0]
	mq.Player = mq.Player[1:]

	conn := mq.Conns[playerId]
	delete(mq.Conns, playerId)

	return playerId, conn
}

func (mq *MatchQueue) Size() int {
	mq.Mutex.Lock()
	defer mq.Mutex.Unlock()

	return len(mq.Player)
}
