package timer

import (
	"sync"
	"time"
)

type TimerManager struct {
	Repository TimerRepository
	mutex      sync.Mutex
}

func NewTimerManager(
	repository TimerRepository,
) *TimerManager {
	return &TimerManager{
		Repository: repository,
	}
}

func (tm *TimerManager) CreateGameTimers(
	gameID string,
	player1ID, player2ID string,
	gameTimerUpdateCallback func(gameId string),
	player1Timeout, player2Timeout func(),
) {
	player1TimeUpdate := func(
		gameId, playerId string,
		remainingTime time.Duration,
	) {
		gameTimerUpdateCallback(gameId)
	}

	player2TimeUpdate := func(
		gameId, playerId string,
		remainingTime time.Duration,
	) {
		gameTimerUpdateCallback(gameId)
	}

	player1Timer := NewTimer(gameID,
		player1ID,
		10*time.Minute,
		player1Timeout,
		player1TimeUpdate,
	)

	player2Timer := NewTimer(
		gameID,
		player2ID,
		10*time.Minute,
		player2Timeout,
		player2TimeUpdate,
	)

	gameTimers := &GameTimers{
		Player1Timer: player1Timer,
		Player2Timer: player2Timer,
	}

	tm.Repository.AddGameTimers(gameID, gameTimers)

	gameTimers.Player1Timer.Start()
}

func (tm *TimerManager) SwitchTurns(gameID, playerID string) {
	tm.mutex.Lock()
	defer tm.mutex.Unlock()

	gameTimers, exists := tm.Repository.GetGameTimers(gameID)
	if !exists {
		return
	}

	if gameTimers.Player1Timer.PlayerId == playerID {
		gameTimers.Player1Timer.Pause()
		gameTimers.Player2Timer.Start()
	} else if gameTimers.Player2Timer.PlayerId == playerID {
		gameTimers.Player2Timer.Pause()
		gameTimers.Player1Timer.Start()
	}
}

func (tm *TimerManager) StopGameTimers(gameID string) {
	tm.mutex.Lock()
	defer tm.mutex.Unlock()

	gameTimers, exists := tm.Repository.GetGameTimers(gameID)
	if !exists {
		return
	}

	gameTimers.Player1Timer.Stop()
	gameTimers.Player2Timer.Stop()
	tm.Repository.RemoveGameTimers(gameID)
}
