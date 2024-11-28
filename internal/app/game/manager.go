package game

import (
	"errors"
	"tavinder/chess-server/internal/app/timer"
	"tavinder/chess-server/internal/models"
	"time"

	"github.com/charmbracelet/log"
	"github.com/gofiber/contrib/websocket"
)

type GameManager struct {
	GameRepository
	timerManager *timer.TimerManager
}

func NewGameManager(
	timerManager *timer.TimerManager,
	gameRepository GameRepository,
) *GameManager {
	gm := &GameManager{
		GameRepository: gameRepository,
		timerManager:   timerManager,
	}

	return gm
}

func (gm *GameManager) AddGame(
	gameId, player1Id, player2Id string,
	player1Conn *websocket.Conn,
	player2Conn *websocket.Conn,
) (*Game, error) {
	_, exists := gm.GameRepository.GetGame(gameId)

	if exists {
		return nil, errors.New("game already exists")
	}

	newGame := &Game{
		GameId:          gameId,
		Player1Id:       player1Id,
		Player2Id:       player2Id,
		Player1Conn:     player1Conn,
		Player2Conn:     player2Conn,
		CurrentPlayerId: player1Id,
	}

	gm.GameRepository.AddGame(newGame)

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

	return newGame, nil
}

func (gm *GameManager) EndGame(gameId, losingPlayerId string) {

	// Retrieve the game
	game, exists := gm.GameRepository.GetGame(gameId)
	if !exists {
		return
	}

	notification := map[string]string{
		"event":        "game_over",
		"losingPlayer": losingPlayerId,
	}

	game.Player1Conn.WriteJSON(notification)
	game.Player2Conn.WriteJSON(notification)

	// delete(gm.Games, gameId)

	gm.timerManager.StopGameTimers(gameId)
}

func (gm *GameManager) EmitTimeUpdates(gameId string) {
	timers, exists := gm.timerManager.Repository.GetGameTimers(gameId)

	if !exists {
		return
	}

	player1Id := timers.Player1Timer.PlayerId
	player2Id := timers.Player2Timer.PlayerId

	game, exists := gm.GameRepository.GetGame(gameId)
	if !exists {
		return
	}

	payload := &models.TimerUpdateMessage{
		Event: "timer_update",
		Payload: &models.TimerUpdatePayload{
			GameId:          gameId,
			Player1Id:       player1Id,
			Player2Id:       player2Id,
			Player1TimeLeft: timers.Player1Timer.RemainingTime.Milliseconds(),
			Player2TimeLeft: timers.Player2Timer.RemainingTime.Milliseconds(),
		},
	}

	// log.Debug(game)
	// log.Debug(payload)

	game.Player1Conn.WriteJSON(payload)
	game.Player2Conn.WriteJSON(payload)
}

func (gm *GameManager) MakeMove(playerId string, message *models.Message) {
	game, exists := gm.GameRepository.GetGame(message.Payload.GameId)
	if !exists {
		log.Error("Game not found")
		return
	}

	if game.CurrentPlayerId != playerId {
		log.Error("Not your turn")
		return
	}

	if playerId == game.Player1Id {
		moveMessage := &models.MoveMessage{
			Event: "move",
			Payload: &models.MovePayload{
				GameId:    message.Payload.GameId,
				Player1Id: game.Player1Id,
				Player2Id: game.Player2Id,
				LastMove: &models.Move{
					Move:      message.Payload.Move,
					PlayerId:  playerId,
					TimeTaken: 10 * time.Second.Milliseconds(),
				},
				CurrentPlayerId: game.Player2Id,
			},
		}

		game.CurrentPlayerId = game.Player2Id
		gm.GameRepository.UpdateGame(game.GameId, game)

		game.Player1Conn.WriteJSON(moveMessage)
		game.Player2Conn.WriteJSON(moveMessage)
	} else {
		moveMessage := &models.MoveMessage{
			Event: "move",
			Payload: &models.MovePayload{
				GameId:    message.Payload.GameId,
				Player1Id: game.Player1Id,
				Player2Id: game.Player2Id,
				LastMove: &models.Move{
					Move:     message.Payload.Move,
					PlayerId: playerId,
				},
				CurrentPlayerId: game.Player1Id,
			},
		}

		game.CurrentPlayerId = game.Player1Id
		gm.GameRepository.UpdateGame(game.GameId, game)

		game.Player1Conn.WriteJSON(moveMessage)
		game.Player2Conn.WriteJSON(moveMessage)
	}

	gm.timerManager.SwitchTurns(game.GameId, playerId)
}
