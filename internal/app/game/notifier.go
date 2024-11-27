package game

import "tavinder/chess-server/internal/models"

type WebSocketNotifier interface {
	NotifyAllPlayersByGame(gameId string, message map[string]interface{})
	NotifyTimerUpdate(payload *models.TimerUpdatePayload)
	NotifyTimeout(payload *models.TimeoutPayload)
}
