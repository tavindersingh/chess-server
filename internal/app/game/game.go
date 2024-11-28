package game

import (
	"time"

	"github.com/gofiber/contrib/websocket"
)

type Game struct {
	GameId            string        `json:"gameId"`
	Player1Id         string        `json:"player1Id"`
	Player2Id         string        `json:"player2Id"`
	CurrentPlayerId   string        `json:"currentPlayerId"`
	StartTime         time.Time     `json:"startTime"`
	TotalTimeEachSide time.Duration `json:"totalTimeEachSide"`
	Player1TimeLeft   time.Duration `json:"player1TimeLeft"`
	Player2TimeLeft   time.Duration `json:"player2TimeLeft"`
	Player1Conn       *websocket.Conn
	Player2Conn       *websocket.Conn
}

func NewGame(
	gameId string,
	player1Id string,
	player2Id string,
	player1Conn *websocket.Conn,
	player2Conn *websocket.Conn,
) *Game {
	return &Game{
		GameId:      gameId,
		Player1Id:   player1Id,
		Player2Id:   player2Id,
		Player1Conn: player1Conn,
		Player2Conn: player2Conn,
		StartTime:   time.Now(),
	}
}
