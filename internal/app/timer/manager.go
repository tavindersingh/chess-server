package timer

import (
	"sync"
	"tavinder/chess-server/internal/events"
	"tavinder/chess-server/internal/models"
	"time"
)

type TimerManager struct {
	repository TimerRepository
	eventBus   *events.EventBus
	mutex      sync.Mutex
}

func NewTimerManager(
	repository TimerRepository,
	eventBus *events.EventBus,
) *TimerManager {
	return &TimerManager{
		repository: repository,
		eventBus:   eventBus,
	}
}

func (tm *TimerManager) CreateGameTimers(
	gameID string,
	player1ID, player2ID string,
	// player1Timeout, player2Timeout func(),
) {
	player1Timeout := func() {
		tm.eventBus.Publish(events.Event{
			Type: events.EventTimeout,
			Payload: map[string]string{
				"gameID":   gameID,
				"playerID": player1ID,
			},
		})
	}

	player2Timeout := func() {
		tm.eventBus.Publish(events.Event{
			Type: events.EventTimeout,
			Payload: map[string]string{
				"gameID":   gameID,
				"playerID": player2ID,
			},
		})
	}

	player1TimeUpdate := func(
		gameId, playerId string,
		remainingTime time.Duration,
	) {
		tm.EmitTimeUpdates(gameId)
	}

	player2TimeUpdate := func(
		gameId, playerId string,
		remainingTime time.Duration,
	) {
		tm.EmitTimeUpdates(gameId)
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

	tm.repository.AddGameTimers(gameID, gameTimers)

	gameTimers.Player1Timer.Start()
}

func (tm *TimerManager) SwitchTurns(gameID, playerID string) {
	tm.mutex.Lock()
	defer tm.mutex.Unlock()

	gameTimers, exists := tm.repository.GetGameTimers(gameID)
	if !exists {
		return
	}

	if gameTimers.Player1Timer.PlayerId == playerID {
		gameTimers.Player2Timer.Pause()
		gameTimers.Player1Timer.Start()
	} else if gameTimers.Player2Timer.PlayerId == playerID {
		gameTimers.Player1Timer.Pause()
		gameTimers.Player2Timer.Start()
	}
}

func (tm *TimerManager) StopGameTimers(gameID string) {
	tm.mutex.Lock()
	defer tm.mutex.Unlock()

	gameTimers, exists := tm.repository.GetGameTimers(gameID)
	if !exists {
		return
	}

	gameTimers.Player1Timer.Stop()
	gameTimers.Player2Timer.Stop()
	tm.repository.RemoveGameTimers(gameID)
}

func (tm *TimerManager) EmitTimeUpdates(gameId string) {
	timers, exists := tm.repository.GetGameTimers(gameId)

	if !exists {
		return
	}

	player1Id := timers.Player1Timer.PlayerId
	player2Id := timers.Player2Timer.PlayerId

	tm.eventBus.Publish(events.Event{
		Type: events.EventTimerUpdate,
		Payload: &models.TimerUpdatePayload{
			gameId:          gameId,
			player1Id:       player1Id,
			player2Id:       player2Id,
			player1TimeLeft: timers.Player1Timer.RemainingTime.Milliseconds(),
			player2TimeLeft: timers.Player2Timer.RemainingTime.Seconds(),
		},
	})
}
