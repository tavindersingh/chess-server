package timer

type TimerService interface {
	CreateGameTimers(gameID string,
		player1ID, player2ID string,
		// player1Timeout, player2Timeout func(),
	)
	StopGameTimers(gameID string)
}
