package timer

type TimerService interface {
	CreateGameTimers(gameID string,
		player1ID, player2ID string,
		gameTimerUpdateCallback func(gameId string),
		player1Timeout, player2Timeout func(),
	)
	StopGameTimers(gameID string)
}
