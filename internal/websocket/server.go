package websocket

import (
	"encoding/json"

	"tavinder/chess-server/internal/app/game"
	"tavinder/chess-server/internal/app/match"
	"tavinder/chess-server/internal/models"

	"github.com/gofiber/contrib/websocket"
)

type WebSocketServer struct {
	MatchManager *match.MatchManager
	GameManager  *game.GameManager
}

func NewWebSocketServer(
	matchManager *match.MatchManager,
	gameManager *game.GameManager,
) *WebSocketServer {
	return &WebSocketServer{
		MatchManager: matchManager,
		GameManager:  gameManager,
	}
}

func (ws *WebSocketServer) HandleWebSocketConnection(conn *websocket.Conn) {
	playerId := conn.Params("id")

	if playerId == "" {
		conn.WriteJSON(map[string]interface{}{
			"success": false,
			"error":   "player ID is required",
		})
		conn.Close()
		return
	}

	done := make(chan bool)
	go func() {
		ws.HandleWebSocketEvents(conn, playerId)
		done <- true
	}()

	<-done
}

func (ws *WebSocketServer) HandleWebSocketEvents(conn *websocket.Conn, playerId string) {
	defer conn.Close()

	for {
		_, p, err := conn.ReadMessage()
		if err != nil {
			break
		}

		var message *models.Message
		if err := json.Unmarshal(p, &message); err != nil {
			conn.WriteJSON(map[string]interface{}{
				"error": "invalid message format",
			})
			continue
		}

		// if err := json.Unmarshal(message, &payload); err != nil {
		// 	conn.WriteJSON(map[string]interface{}{
		// 		"error": "invalid message format",
		// 	})
		// 	continue
		// }

		event := message.Event
		if event == "" {
			conn.WriteJSON(map[string]interface{}{
				"error": "missing event field",
			})
			continue
		}

		switch event {
		case "join_queue":
			ws.MatchManager.AddPlayerToQueue(playerId, conn)
		case "move":
			ws.GameManager.MakeMove(playerId, message)
		}
	}
}
