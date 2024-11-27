package match

import (
	"sync"
	"tavinder/chess-server/internal/app/game"
	"tavinder/chess-server/internal/utils"

	"github.com/gofiber/contrib/websocket"
)

type MatchManager struct {
	Queue       *MatchQueue
	GameManager *game.GameManager
	Mutex       sync.Mutex
}

func NewMatchManager(gameManager *game.GameManager) *MatchManager {
	return &MatchManager{
		Queue:       NewMatchQueue(),
		GameManager: gameManager,
		Mutex:       sync.Mutex{},
	}
}

func (mm *MatchManager) AddPlayerToQueue(playerId string, conn *websocket.Conn) {
	mm.Mutex.Lock()
	defer mm.Mutex.Unlock()

	mm.Queue.Enqueue(playerId, conn)

	if mm.Queue.Size() >= 2 {
		player1, conn1 := mm.Queue.Dequeue()
		player2, conn2 := mm.Queue.Dequeue()

		gameId := utils.GenerateRandomID()

		newGame, err := mm.GameManager.AddGame(gameId, player1, player2, conn1, conn2)

		if err != nil {
			conn1.WriteJSON(map[string]interface{}{
				"type":    "error",
				"payload": "Failed to create game",
			})
			conn2.WriteJSON(map[string]interface{}{
				"type":    "error",
				"payload": "Failed to create game",
			})
			return
		}

		conn1.WriteJSON(map[string]interface{}{
			"type":    "match_found",
			"payload": newGame,
		})
		conn2.WriteJSON(map[string]interface{}{
			"type":    "match_found",
			"payload": newGame,
		})
	}
}

func (mm *MatchManager) NotifyAllPlayersByGame(
	gameId string,
	message map[string]interface{},
) {
	mm.Mutex.Lock()
	defer mm.Mutex.Unlock()

	game, exists := mm.GameManager.Games[gameId]

	if !exists {
		return
	}

	player1Id := game.Player1Id
	player2Id := game.Player2Id

	player1Conn := mm.Queue.Conns[player1Id]
	player2Conn := mm.Queue.Conns[player2Id]

	player1Conn.WriteJSON(message)
	player2Conn.WriteJSON(message)
}
