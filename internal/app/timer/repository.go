package timer

type GameTimers struct {
	Player1Timer *Timer
	Player2Timer *Timer
	// repository    game.GameRepository
}

type TimerRepository interface {
	AddGameTimers(gameID string, timers *GameTimers)
	GetGameTimers(gameID string) (*GameTimers, bool)
	RemoveGameTimers(gameID string)
}

// func NewGameTimer(game game.Game) *GameTimers {
// 	gameTimer := &GameTimers{
// 		Player1Timer:  NewTimer(game.Player1Id, game.Player1TimeLeft, game.OnUpdatePlayer1TimeLeft),
// 		Player2Timer:  NewTimer(game.Player2Id, game.Player2TimeLeft, game.OnUpdatePlayer2TimeLeft),
// 		isPlayer1Turn: true,
// 		// repository:    repository,
// 	}

// 	return gameTimer
// }

// func (gt *GameTimers) createTimeUpdateCallback(playerId string) TimeUpdateCallback {
// 	return func(timeLeft time.Duration) {
// 		game, err := gt.repository.GetGame(gt.gameId)
// 		if err != nil {
// 			return
// 		}

// 	}
// }
