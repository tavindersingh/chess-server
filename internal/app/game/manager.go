package game

import (
	"errors"
	"sync"
	"tavinder/chess-server/internal/app/timer"
	"tavinder/chess-server/internal/events"

	"github.com/charmbracelet/log"
	"github.com/gofiber/contrib/websocket"
)

type GameManager struct {
	Games map[string]*Game
	mutex sync.RWMutex

	timerService timer.TimerService
	eventBus     *events.EventBus
	wsNotifier   WebSocketNotifier
}

func NewGameManager(
	timerService timer.TimerService,
	eventBus *events.EventBus,
) *GameManager {
	gm := &GameManager{
		Games:        make(map[string]*Game),
		timerService: timerService,
		eventBus:     eventBus,
	}

	gm.eventBus.Subscribe(events.EventTimeout, gm.handleTimeout)
	gm.eventBus.Subscribe(events.EventTimerUpdate, gm.handleTimerUpdate)

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

	gm.timerService.CreateGameTimers(
		gameId,
		player1Id,
		player2Id,
		// func() {
		// 	gm.EndGame(gameId, player1Id)
		// },
		// func() {
		// 	gm.EndGame(gameId, player2Id)
		// },
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

	gm.timerService.StopGameTimers(gameId)
}

// func (gm *GameManager) NotifyPlayers(
// 	gameID, playerID string,
// 	remainingTime float64,
// ) {
// 	game, exists := gm.Games[gameID]
// 	if !exists {
// 		return
// 	}

// 	message := map[string]interface{}{
// 		"event": "timer_update",
// 		"data": map[string]interface{}{
// 			"player_id": playerID,
// 			"time":      remainingTime,
// 		},
// 	}

// 	// Notify Player 1
// 	if err := game.Player1Conn.WriteJSON(message); err != nil {
// 		// handle error (e.g., log it)
// 	}

// 	// Notify Player 2
// 	if err := game.Player2Conn.WriteJSON(message); err != nil {
// 		// handle error (e.g., log it)
// 	}
// }

func (gm *GameManager) handleTimeout(event events.Event) {
	log.Info("handleTimeout")
	log.Info(event)
	// payload := event.Payload.(map[string]interface{})

	// gm.wsNotifier.NotifyAllPlayersByGame(
	// 	payload["gameId"].(string),
	// 	map[string]interface{}{
	// 		"type":    "timer_update",
	// 		"payload": payload,
	// 	},
	// )
}

func (gm *GameManager) handleTimerUpdate(event events.Event) {
	log.Info("handle timer update")
	log.Info(event)
	payload := event.Payload

	gm.wsNotifier.NotifyAllPlayersByGame(
		payload["gameId"].(string),
		map[string]interface{}{
			"type":    "timer_update",
			"payload": payload,
		},
	)
}
