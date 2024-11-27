package game

import (
	"errors"
	"sync"
	"tavinder/chess-server/internal/app/timer"
	"tavinder/chess-server/internal/models"

	"github.com/gofiber/contrib/websocket"
)

type GameManager struct {
	Games map[string]*Game
	mutex sync.RWMutex

	timerManager *timer.TimerManager
}

func NewGameManager(
	timerManager *timer.TimerManager,
) *GameManager {
	gm := &GameManager{
		Games:        make(map[string]*Game),
		timerManager: timerManager,
	}

	return gm
}

func (gm *GameManager) AddGame(
	gameId, player1Id, player2Id string,
	player1Conn *websocket.Conn,
	player2Conn *websocket.Conn,
) (*Game, error) {
	gm.mutex.Lock()
	defer gm.mutex.Unlock()

	if _, exists := gm.Games[gameId]; exists {
		return nil, errors.New("game already exists")
	}

	game := &Game{
		GameId:          gameId,
		Player1Id:       player1Id,
		Player2Id:       player2Id,
		Player1Conn:     player1Conn,
		Player2Conn:     player2Conn,
		CurrentPlayerId: player1Id,
	}

	gm.Games[gameId] = game

	gm.timerManager.CreateGameTimers(
		gameId,
		player1Id,
		player2Id,
		gm.EmitTimeUpdates,
		func() {
			gm.EndGame(gameId, player1Id)
		},
		func() {
			gm.EndGame(gameId, player2Id)
		},
	)

	return game, nil
}

func (gm *GameManager) EndGame(gameId, losingPlayerId string) {
	gm.mutex.Lock()
	defer gm.mutex.Unlock()

	// Retrieve the game
	game, exists := gm.Games[gameId]
	if !exists {
		return
	}

	notification := map[string]string{
		"event":        "game_over",
		"losingPlayer": losingPlayerId,
	}

	game.Player1Conn.WriteJSON(notification)
	game.Player2Conn.WriteJSON(notification)

	delete(gm.Games, gameId)

	gm.timerManager.StopGameTimers(gameId)
}

func (gm *GameManager) EmitTimeUpdates(gameId string) {
	timers, exists := gm.timerManager.Repository.GetGameTimers(gameId)

	if !exists {
		return
	}

	player1Id := timers.Player1Timer.PlayerId
	player2Id := timers.Player2Timer.PlayerId

	game := gm.Games[gameId]

	payload := &models.TimerUpdatePayload{
		GameId:          gameId,
		Player1Id:       player1Id,
		Player2Id:       player2Id,
		Player1TimeLeft: timers.Player1Timer.RemainingTime.Milliseconds(),
		Player2TimeLeft: timers.Player2Timer.RemainingTime.Milliseconds(),
	}

	game.Player1Conn.WriteJSON(payload)
	game.Player2Conn.WriteJSON(payload)
}
